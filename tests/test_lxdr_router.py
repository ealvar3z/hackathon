"""Tests for the first in-memory LXDR router behavior."""

from __future__ import annotations

import json

from section4.lxdr.header import LXDRHeader
from section4.lxdr.link import LXDRLinkFrame
from section4.lxdr.request import LXDRRequestContainer
from section4.lxdr.router import LXDRRouter
from section4.lxdr.segments import PAX_MOVEMENT, LXDRSegment
from section4.lxdr.types import DeliveryMethod, LinkRepresentation, LinkState


def build_sample_header(local_id: str = "3838JBNM5X") -> LXDRHeader:
    """Build a valid ADRIAN-style header for router tests."""

    return LXDRHeader(
        date_request_created_local="2027OCT13",
        time_request_created_local="15470352",
        physical_location_of_requestor="4QFJ1234567890",
        request_unique_identification_local=local_id,
        request_priority="02",
        element_unit_identification_callsign="KL9K",
        request_segments=1,
    )


def build_sample_request(local_id: str = "3838JBNM5X") -> LXDRRequestContainer:
    """Build a minimal valid PAX request container."""

    segment = LXDRSegment(
        spec=PAX_MOVEMENT,
        values={
            "segment_number": "1",
            "request_type": "PM",
            "request_priority": "02",
            "zap_or_edipi": "1010919789",
            "earliest_departure_date": "2027OCT15",
            "latest_departure_date": "2027OCT20",
            "departure_location": "4QFJ1234567890",
            "destination_location": "4QFJ4567890123",
            "estimated_baggage_weight_lbs": "075",
            "hazardous_material_type": "X",
        },
    )
    return LXDRRequestContainer(
        header=build_sample_header(local_id),
        segments=[segment],
    )


def build_sync_response_frame(
    local_id: str,
    sync_id: str,
) -> LXDRLinkFrame:
    """Build a synchronization response frame for a single request."""

    return LXDRLinkFrame(
        link_message_id="sync-msg-1",
        sender_id="ALOC",
        recipient_id="NODE1",
        created_at_local="2027OCT13T15470352",
        delivery_method=DeliveryMethod.SYNCHRONIZATION,
        representation=LinkRepresentation.INLINE_STRUCTURED,
        payload_count=1,
        payloads=[
            json.dumps(
                {
                    "local_request_id": local_id,
                    "sync_request_id": sync_id,
                },
                sort_keys=True,
                separators=(",", ":"),
            )
        ],
        state=LinkState.DELIVERED,
        sync_of=local_id,
    )


def test_router_queues_outbound_request_as_created_frame() -> None:
    """Queuing a request should create an outbox entry and queued frame."""

    router = LXDRRouter(sender_id="NODE1")
    request = build_sample_request()

    frame = router.queue_request(
        request=request,
        recipient_id="ALOC",
        created_at_local="2027OCT13T15470352",
    )

    assert frame.sender_id == "NODE1"
    assert frame.recipient_id == "ALOC"
    assert frame.state is LinkState.QUEUED
    assert router.outbox_size == 1
    assert (
        router.pending_sync_requests()[0].header.to_dict()
        == request.header.to_dict()
    )


def test_router_builds_sync_bundle_from_unsynchronized_requests() -> None:
    """The router should bundle all unsynchronized queued requests."""

    router = LXDRRouter(sender_id="NODE1")
    router.queue_request(
        request=build_sample_request("3838JBNM5X"),
        recipient_id="ALOC",
        created_at_local="2027OCT13T15470352",
    )
    router.queue_request(
        request=build_sample_request("3838JBNM6X"),
        recipient_id="ALOC",
        created_at_local="2027OCT13T15480352",
    )

    bundle = router.build_sync_bundle(
        recipient_id="ALOC",
        created_at_local="2027OCT13T15490352",
        correlation_id="sync-pass-1",
    )

    assert bundle.representation is LinkRepresentation.BUNDLE
    assert bundle.delivery_method is DeliveryMethod.SYNCHRONIZATION
    assert bundle.payload_count == 1
    bundled_requests = bundle.embedded_requests()
    assert len(bundled_requests) == 2
    assert bundled_requests[0].header.request_unique_identification_local == (
        "3838JBNM5X"
    )
    assert bundled_requests[1].header.request_unique_identification_local == (
        "3838JBNM6X"
    )


def test_router_dedupes_inbound_link_frames_by_message_id() -> None:
    """The same inbound frame should only be accepted once."""

    router = LXDRRouter(sender_id="NODE1")
    frame = build_sync_response_frame("3838JBNM5X", "ABCD1234EFGH")

    assert router.accept_inbound_frame(frame) is True
    assert router.accept_inbound_frame(frame) is False


def test_router_applies_sync_identifier_to_matching_request() -> None:
    """A sync response should update the matching request header."""

    router = LXDRRouter(sender_id="NODE1")
    router.queue_request(
        request=build_sample_request("3838JBNM5X"),
        recipient_id="ALOC",
        created_at_local="2027OCT13T15470352",
    )

    frame = build_sync_response_frame("3838JBNM5X", "ABCD1234EFGH")
    updated = router.apply_sync_response(frame)

    assert updated == 1
    request = router.request_by_local_id("3838JBNM5X")
    assert request is not None
    assert request.header.request_unique_identification_sync == (
        "ABCD1234EFGH"
    )
    assert router.pending_sync_requests() == []


def test_router_rejects_sync_response_for_unknown_request() -> None:
    """A sync response for an unknown local request ID should fail cleanly."""

    router = LXDRRouter(sender_id="NODE1")
    frame = build_sync_response_frame("MISSING0001", "ABCD1234EFGH")

    assert router.apply_sync_response(frame) == 0
