"""Tests for persisted LXDR router state."""

from __future__ import annotations

import json
from pathlib import Path

from section4.lxdr.header import LXDRHeader
from section4.lxdr.link import LXDRLinkFrame
from section4.lxdr.request import LXDRRequestContainer
from section4.lxdr.router import LXDRRouter
from section4.lxdr.router_store import PersistentLXDRRouter
from section4.lxdr.segments import PAX_MOVEMENT, LXDRSegment
from section4.lxdr.types import DeliveryMethod, LinkRepresentation, LinkState
from section4.storage import create_all, create_session_factory
from section4.storage.tables import (
    LXDRInboundFrame,
    LXDROutboundFrame,
    LXDRRequestRecord,
)


def build_sample_header(local_id: str = "3838JBNM5X") -> LXDRHeader:
    """Build a valid ADRIAN-style header for router-store tests."""

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


def build_router(tmp_path: Path) -> tuple[PersistentLXDRRouter, Path]:
    """Create a persistent router backed by a temporary SQLite database."""

    db_path = tmp_path / "section4-router.db"
    create_all(db_path)
    session_factory = create_session_factory(db_path)
    memory_router = LXDRRouter(sender_id="NODE1")
    return (
        PersistentLXDRRouter(
            session_factory=session_factory,
            memory_router=memory_router,
        ),
        db_path,
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


def build_inbound_request_frame(
    local_id: str = "REMOTE0001",
) -> LXDRLinkFrame:
    """Build a direct inbound frame carrying one ADRIAN request."""

    request = build_sample_request(local_id)
    return LXDRLinkFrame(
        link_message_id="inbound-req-1",
        sender_id="ALOC",
        recipient_id="NODE1",
        created_at_local="2027OCT13T16010352",
        delivery_method=DeliveryMethod.DIRECT,
        representation=LinkRepresentation.INLINE_STRUCTURED,
        payload_count=1,
        payloads=[
            json.dumps(
                request.to_dict(),
                sort_keys=True,
                separators=(",", ":"),
            )
        ],
        state=LinkState.DELIVERED,
    )


def test_persistent_router_stores_outbound_frames(tmp_path: Path) -> None:
    """Queued outbound requests should also be persisted."""

    router, db_path = build_router(tmp_path)
    router.queue_request(
        request=build_sample_request(),
        recipient_id="ALOC",
        created_at_local="2027OCT13T15470352",
    )

    session_factory = create_session_factory(db_path)
    with session_factory() as session:
        rows = session.query(LXDROutboundFrame).all()

    assert len(rows) == 1
    assert rows[0].sender_id == "NODE1"
    assert rows[0].recipient_id == "ALOC"
    assert rows[0].state == "QUEUED"


def test_persistent_router_stores_adrian_request_record(
    tmp_path: Path,
) -> None:
    """Queued requests should persist a first-class ADRIAN request record."""

    router, db_path = build_router(tmp_path)
    request = build_sample_request("3838JBNM5X")
    router.queue_request(
        request=request,
        recipient_id="ALOC",
        created_at_local="2027OCT13T15470352",
    )

    session_factory = create_session_factory(db_path)
    with session_factory() as session:
        row = session.query(LXDRRequestRecord).one()

    assert row.request_unique_identification_local == "3838JBNM5X"
    assert row.request_direction == "OUTBOUND"
    assert row.source_sender_id == "NODE1"
    assert row.source_recipient_id == "ALOC"
    assert row.request_type == "PM"
    assert row.primary_segment_name == "mobility_pax"
    assert row.latest_frame_state == "QUEUED"
    assert row.canonical_text.startswith("2027OCT13-15470352-")


def test_persistent_router_stores_inbound_frame_dedupe_record(
    tmp_path: Path,
) -> None:
    """Accepted inbound frames should be recorded for dedupe persistence."""

    router, db_path = build_router(tmp_path)

    router.record_inbound_frame(
        LXDRLinkFrame(
            link_message_id="inbound-001",
            sender_id="ALOC",
            recipient_id="NODE1",
            created_at_local="2027OCT13T15470352",
            delivery_method=DeliveryMethod.DIRECT,
            representation=LinkRepresentation.INLINE_STRUCTURED,
            payload_count=1,
            payloads=["{}"],
            state=LinkState.DELIVERED,
        )
    )

    session_factory = create_session_factory(db_path)
    with session_factory() as session:
        rows = session.query(LXDRInboundFrame).all()

    assert len(rows) == 1
    assert rows[0].link_message_id == "inbound-001"
    assert rows[0].sender_id == "ALOC"
    assert rows[0].payload_json["sender_id"] == "ALOC"


def test_persistent_router_updates_request_sync_identifier(
    tmp_path: Path,
) -> None:
    """Applying sync responses should update persisted outbound records."""

    router, db_path = build_router(tmp_path)
    router.queue_request(
        request=build_sample_request("3838JBNM5X"),
        recipient_id="ALOC",
        created_at_local="2027OCT13T15470352",
    )
    router.apply_sync_update(
        local_request_id="3838JBNM5X",
        sync_request_id="ABCD1234EFGH",
    )

    session_factory = create_session_factory(db_path)
    with session_factory() as session:
        row = session.query(LXDROutboundFrame).one()
        request_row = session.query(LXDRRequestRecord).one()

    assert row.request_unique_identification_local == "3838JBNM5X"
    assert row.request_unique_identification_sync == "ABCD1234EFGH"
    assert row.state == "SYNCED"
    assert request_row.request_unique_identification_sync == "ABCD1234EFGH"
    assert request_row.latest_frame_state == "SYNCED"


def test_persistent_router_processes_sync_response_frame(
    tmp_path: Path,
) -> None:
    """Inbound sync frames should update request and inbound frame state."""

    router, db_path = build_router(tmp_path)
    router.queue_request(
        request=build_sample_request("3838JBNM5X"),
        recipient_id="ALOC",
        created_at_local="2027OCT13T15470352",
    )

    updates = router.receive_inbound_frame(
        build_sync_response_frame("3838JBNM5X", "ABCD1234EFGH")
    )

    session_factory = create_session_factory(db_path)
    with session_factory() as session:
        inbound_row = session.query(LXDRInboundFrame).one()
        request_row = session.query(LXDRRequestRecord).one()

    assert updates == 1
    assert inbound_row.link_message_id == "sync-msg-1"
    assert inbound_row.payload_json["sync_of"] == "3838JBNM5X"
    assert request_row.request_unique_identification_sync == "ABCD1234EFGH"
    assert request_row.latest_frame_state == "SYNCED"


def test_persistent_router_persists_inbound_request_frame(
    tmp_path: Path,
) -> None:
    """Inbound request frames should persist received ADRIAN requests."""

    router, db_path = build_router(tmp_path)

    updates = router.receive_inbound_frame(
        build_inbound_request_frame("REMOTE0001")
    )

    session_factory = create_session_factory(db_path)
    with session_factory() as session:
        inbound_row = session.query(LXDRInboundFrame).one()
        request_row = session.query(LXDRRequestRecord).one()

    assert updates == 1
    assert inbound_row.link_message_id == "inbound-req-1"
    assert request_row.request_unique_identification_local == "REMOTE0001"
    assert request_row.request_direction == "RECEIVED"
    assert request_row.source_sender_id == "ALOC"
    assert request_row.source_recipient_id == "NODE1"
    assert request_row.latest_frame_state == "DELIVERED"


def test_persistent_router_tracks_send_attempts(tmp_path: Path) -> None:
    """Marking a frame as sending should increment attempt state."""

    router, db_path = build_router(tmp_path)
    frame = router.queue_request(
        request=build_sample_request(),
        recipient_id="ALOC",
        created_at_local="2027OCT13T15470352",
    )

    router.mark_frame_sending(
        link_message_id=frame.link_message_id,
        attempted_at="2027OCT13T15480352",
    )

    session_factory = create_session_factory(db_path)
    with session_factory() as session:
        row = session.query(LXDROutboundFrame).one()
        request_row = session.query(LXDRRequestRecord).one()

    assert row.state == "SENDING"
    assert row.attempt_count == 1
    assert row.last_attempt_at == "2027OCT13T15480352"
    assert request_row.latest_frame_state == "SENDING"


def test_persistent_router_records_failed_frame_state(tmp_path: Path) -> None:
    """A failed frame should persist its failure reason and state."""

    router, db_path = build_router(tmp_path)
    frame = router.queue_request(
        request=build_sample_request(),
        recipient_id="ALOC",
        created_at_local="2027OCT13T15470352",
    )

    router.mark_frame_failed(
        link_message_id=frame.link_message_id,
        error_message="No route available",
    )

    session_factory = create_session_factory(db_path)
    with session_factory() as session:
        row = session.query(LXDROutboundFrame).one()
        request_row = session.query(LXDRRequestRecord).one()

    assert row.state == "FAILED"
    assert row.last_error == "No route available"
    assert request_row.latest_frame_state == "FAILED"


def test_persistent_router_lists_retryable_frames(tmp_path: Path) -> None:
    """Failed and queued frames should be visible for retry scheduling."""

    router, _ = build_router(tmp_path)
    first = router.queue_request(
        request=build_sample_request("3838JBNM5X"),
        recipient_id="ALOC",
        created_at_local="2027OCT13T15470352",
    )
    second = router.queue_request(
        request=build_sample_request("3838JBNM6X"),
        recipient_id="ALOC",
        created_at_local="2027OCT13T15480352",
    )

    router.mark_frame_failed(
        link_message_id=first.link_message_id,
        error_message="Transient link failure",
    )
    router.mark_frame_sent(link_message_id=second.link_message_id)

    retryable = router.retryable_outbound_frames()

    assert [frame.link_message_id for frame in retryable] == [
        first.link_message_id,
    ]
