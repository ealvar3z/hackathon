"""Tests for persisted LXDR router state."""

from __future__ import annotations

from pathlib import Path

from section4.lxdr.header import LXDRHeader
from section4.lxdr.request import LXDRRequestContainer
from section4.lxdr.router import LXDRRouter
from section4.lxdr.router_store import PersistentLXDRRouter
from section4.lxdr.segments import PAX_MOVEMENT, LXDRSegment
from section4.storage import create_all, create_session_factory
from section4.storage.tables import LXDRInboundFrame, LXDROutboundFrame


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


def test_persistent_router_stores_inbound_frame_dedupe_record(
    tmp_path: Path,
) -> None:
    """Accepted inbound frames should be recorded for dedupe persistence."""

    router, db_path = build_router(tmp_path)

    router.record_inbound_frame(
        link_message_id="inbound-001",
        sender_id="ALOC",
        recipient_id="NODE1",
        payload_count=1,
    )

    session_factory = create_session_factory(db_path)
    with session_factory() as session:
        rows = session.query(LXDRInboundFrame).all()

    assert len(rows) == 1
    assert rows[0].link_message_id == "inbound-001"
    assert rows[0].sender_id == "ALOC"


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

    assert row.request_unique_identification_local == "3838JBNM5X"
    assert row.request_unique_identification_sync == "ABCD1234EFGH"
    assert row.state == "SYNCED"
