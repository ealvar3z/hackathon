"""Tests for the initial LXDR core and link serialization contract."""

from __future__ import annotations

import pytest

from section4.lxdr.codec.packed import (
    PACKED_HEADER_SIZE_BYTES,
    PACKED_PAX_SIZE_BYTES,
    decode_header,
    decode_pax_segment,
    encode_header,
    encode_pax_segment,
)
from section4.lxdr.header import LXDRHeader
from section4.lxdr.link import (
    LXDRLinkFrame,
    embed_bundle,
    embed_inline_structured,
    embed_inline_text,
)
from section4.lxdr.request import LXDRRequestContainer
from section4.lxdr.segments import (
    MAINTENANCE_REQUEST,
    PAX_MOVEMENT,
    LXDRSegment,
)
from section4.lxdr.types import DeliveryMethod, LinkRepresentation, LinkState


def build_sample_header() -> LXDRHeader:
    """Build a valid ADRIAN-style header for tests."""

    return LXDRHeader(
        date_request_created_local="2027OCT13",
        time_request_created_local="15470352",
        physical_location_of_requestor="4QFJ1234567890",
        request_unique_identification_local="3838JBNM5X",
        request_priority="02",
        element_unit_identification_callsign="KL9K",
        request_segments=1,
    )


def build_sample_pax_request() -> LXDRRequestContainer:
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
        header=build_sample_header(),
        segments=[segment],
    )


def build_sample_pax_request_with_local_id(
    local_request_id: str,
) -> LXDRRequestContainer:
    """Build a minimal valid PAX request with a specific local ID."""

    request = build_sample_pax_request()
    request.header.request_unique_identification_local = local_request_id
    return request


def build_base_link_frame() -> LXDRLinkFrame:
    """Build a valid base link frame without embedded payloads."""

    return LXDRLinkFrame(
        link_message_id="msg-1",
        sender_id="S4ALOC",
        recipient_id="NODE2",
        created_at_local="2027OCT13T15470352",
        delivery_method=DeliveryMethod.DIRECT,
        representation=LinkRepresentation.INLINE_TEXT,
        payload_count=1,
        state=LinkState.CREATED,
    )


def test_header_render_text_matches_adrian_shape() -> None:
    """Header text rendering should preserve ADRIAN field order."""

    header = build_sample_header()

    assert (
        header.render_text()
        == "2027OCT13-15470352-4QFJ1234567890-3838JBNM5X-02-KL9K-01"
    )


def test_header_parse_text_round_trips() -> None:
    """Header text should parse back into the structured header model."""

    header = LXDRHeader.parse_text(
        "2027OCT13-15470352-4QFJ1234567890-3838JBNM5X-02-KL9K-01"
    )

    assert header.date_request_created_local == "2027OCT13"
    assert header.request_unique_identification_local == "3838JBNM5X"
    assert header.request_segments == 1


def test_header_packed_round_trip_restores_original_values() -> None:
    """Packed header encoding should round-trip deterministically."""

    header = build_sample_header()

    packed = encode_header(header)
    restored = decode_header(packed)

    assert len(packed) == PACKED_HEADER_SIZE_BYTES
    assert restored.to_dict() == header.to_dict()


def test_header_packed_rejects_invalid_length() -> None:
    """Packed header decoding should reject truncated data."""

    with pytest.raises(ValueError, match="exactly"):
        decode_header(b"\x00" * (PACKED_HEADER_SIZE_BYTES - 1))


def test_pax_packed_round_trip_restores_edipi_segment() -> None:
    """Packed PAX encoding should round-trip for an EDI-PI payload."""

    segment = build_sample_pax_request().segments[0]

    packed = encode_pax_segment(segment)
    restored = decode_pax_segment(packed)

    assert len(packed) == PACKED_PAX_SIZE_BYTES
    assert restored.values == segment.values


def test_pax_packed_round_trip_restores_zap_with_plus() -> None:
    """Packed PAX encoding should preserve ZAP values including '+'."""

    segment = LXDRSegment(
        spec=PAX_MOVEMENT,
        values={
            "segment_number": "1",
            "request_type": "PM",
            "request_priority": "02",
            "zap_or_edipi": "SC9789AB+",
            "earliest_departure_date": "2027OCT15",
            "latest_departure_date": "2027OCT20",
            "departure_location": "4QFJ1234567890",
            "destination_location": "4QFJ4567890123",
            "estimated_baggage_weight_lbs": "075",
            "hazardous_material_type": "X",
        },
    )

    restored = decode_pax_segment(encode_pax_segment(segment))

    assert restored.values["zap_or_edipi"] == "SC9789AB+"


def test_pax_packed_rejects_invalid_length() -> None:
    """Packed PAX decoding should reject truncated data."""

    with pytest.raises(ValueError, match="exactly"):
        decode_pax_segment(b"\x00" * (PACKED_PAX_SIZE_BYTES - 1))


def test_header_rejects_invalid_segment_count() -> None:
    """ADRIAN header segment count must remain within 1..9."""

    header = build_sample_header()
    header.request_segments = 0

    with pytest.raises(ValueError, match="1..9"):
        header.validate()


def test_inline_text_link_round_trip_restores_request() -> None:
    """INLINE_TEXT frames should round-trip through link serialization."""

    request = build_sample_pax_request()
    frame = embed_inline_text(build_base_link_frame(), [request])

    parsed_frame = LXDRLinkFrame.parse(frame.serialize())
    parsed_requests = parsed_frame.embedded_requests()

    assert parsed_frame.representation is LinkRepresentation.INLINE_TEXT
    assert len(parsed_requests) == 1
    assert parsed_requests[0].header.request_priority == "02"
    assert parsed_requests[0].segments[0].values["request_type"] == "PM"


def test_inline_structured_link_round_trip_restores_request() -> None:
    """INLINE_STRUCTURED frames should round-trip through link JSON."""

    request = build_sample_pax_request()
    frame = embed_inline_structured(build_base_link_frame(), [request])

    parsed_frame = LXDRLinkFrame.parse(frame.serialize())
    parsed_requests = parsed_frame.embedded_requests()

    assert (
        parsed_frame.representation
        is LinkRepresentation.INLINE_STRUCTURED
    )
    assert len(parsed_requests) == 1
    assert parsed_requests[0].segments[0].values["zap_or_edipi"] == (
        "1010919789"
    )


def test_link_frame_rejects_payload_count_mismatch() -> None:
    """Link validation should fail if payload_count is inconsistent."""

    frame = LXDRLinkFrame(
        link_message_id="msg-2",
        sender_id="S4ALOC",
        recipient_id="NODE2",
        created_at_local="2027OCT13T15470352",
        delivery_method=DeliveryMethod.DIRECT,
        representation=LinkRepresentation.INLINE_TEXT,
        payload_count=2,
        payloads=["payload-1"],
    )

    with pytest.raises(ValueError, match="payload_count"):
        frame.validate()


def test_bundle_round_trip_restores_multiple_requests() -> None:
    """Bundle frames should round-trip multiple structured requests."""

    request_one = build_sample_pax_request_with_local_id("3838JBNM5X")
    request_two = build_sample_pax_request_with_local_id("3838JBNM6X")
    base = LXDRLinkFrame(
        link_message_id="msg-3",
        sender_id="S4ALOC",
        recipient_id="NODE2",
        created_at_local="2027OCT13T15470352",
        delivery_method=DeliveryMethod.SYNCHRONIZATION,
        representation=LinkRepresentation.BUNDLE,
        payload_count=1,
        sync_of="sync-1",
    )

    frame = embed_bundle(base, [request_one, request_two])
    parsed_requests = LXDRLinkFrame.parse(frame.serialize()).embedded_requests()

    assert frame.representation is LinkRepresentation.BUNDLE
    assert len(parsed_requests) == 2
    assert parsed_requests[0].header.request_unique_identification_local == (
        "3838JBNM5X"
    )
    assert parsed_requests[1].header.request_unique_identification_local == (
        "3838JBNM6X"
    )


def test_bundle_decode_rejects_non_object_payloads() -> None:
    """Bundle payloads must decode to structured request objects."""

    frame = LXDRLinkFrame(
        link_message_id="msg-4",
        sender_id="S4ALOC",
        recipient_id="NODE2",
        created_at_local="2027OCT13T15470352",
        delivery_method=DeliveryMethod.SYNCHRONIZATION,
        representation=LinkRepresentation.BUNDLE,
        payload_count=1,
        payloads=['"not-a-request-object"'],
    )

    with pytest.raises(ValueError, match="must be an object"):
        frame.embedded_requests()


def test_maintenance_request_type_remains_unresolved_for_text_parse() -> None:
    """Maintenance text parsing must not invent a request type code."""

    segment = LXDRSegment(
        spec=MAINTENANCE_REQUEST,
        values={
            "segment_number": "1",
            "request_type": "CM",
            "request_priority": "02",
            "serial_number": "123456789012345",
            "niin": "015519434",
            "model_of_equipment": "MODELX",
            "item_nomenclature": "RADIO",
            "number_of_pieces": "1",
            "equipment_operational_condition": "C",
            "date_maintenance_support_required": "2027OCT20",
            "location_of_equipment": "4QFJ1234567890",
            "type_of_maintenance_support": "R1",
            "type_of_repair": "D1",
            "repair_major_defect": "MD13",
            "attachment": "1",
            "narrative": "Needs bench evaluation",
        },
    )
    request = LXDRRequestContainer(
        header=build_sample_header(),
        segments=[segment],
    )
    frame = embed_inline_structured(build_base_link_frame(), [request])
    parsed = LXDRLinkFrame.parse(frame.serialize()).embedded_requests()[0]

    assert parsed.segments[0].spec.request_type is None
