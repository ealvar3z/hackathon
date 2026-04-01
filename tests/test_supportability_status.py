"""Tests for logistics status reporting and supportability scoring."""

from __future__ import annotations

from pathlib import Path

from section4.storage import create_all, create_session_factory
from section4.storage.tables import (
    Capability,
    LogisticsStatusReport,
    LXDRRequestRecord,
)
from section4.supportability import (
    SupportabilityAssessment,
    assess_request_supportability,
)


def build_session_factory(tmp_path: Path):
    """Create a clean test database for supportability tests."""

    db_path = tmp_path / "section4-supportability.db"
    create_all(db_path)
    return create_session_factory(db_path)


def test_supportability_prefers_on_hand_stock_over_fabrication(
    tmp_path: Path,
) -> None:
    """Available stock should outrank slower fabrication when sufficient."""

    session_factory = build_session_factory(tmp_path)
    with session_factory() as session:
        request = LXDRRequestRecord(
            request_unique_identification_local="3838JBNM5X",
            request_type="SR",
            primary_segment_name="supply_request",
            request_priority="02",
            element_unit_identification_callsign="KL9K",
            physical_location_of_requestor="4QFJ1234567890",
            created_at_local="2027OCT13T15470352",
            segment_count=1,
            canonical_text="REQ",
            payload_json={},
            latest_frame_state="QUEUED",
        )
        capability = Capability(
            node_id="atak-node-2",
            callsign="NODE2",
            capability_type="fabrication",
            title="Fabrication cell",
            availability_status="available",
            throughput_per_day=4,
            lead_time_minutes=90,
        )
        stock = LogisticsStatusReport(
            node_id="atak-node-3",
            callsign="NODE3",
            report_type="on_hand_usage",
            item_reference="BRKT-4421",
            on_hand_quantity=3,
            usage_rate_per_day=1.0,
            required_quantity=1,
            required_by_local="2027OCT15",
            delivery_method="direct",
            transport_mode="ground",
            payload_json={},
        )
        session.add_all([request, capability, stock])
        session.commit()

        assessment = assess_request_supportability(
            session=session,
            request=request,
        )

    assert isinstance(assessment, SupportabilityAssessment)
    assert assessment.recommended_course_of_action == "issue_from_stock"
    assert assessment.supportable is True
    assert assessment.supporting_node_id == "atak-node-3"


def test_supportability_prefers_fabrication_when_stock_is_zero(
    tmp_path: Path,
) -> None:
    """Fabrication should be selected when no stock is on hand."""

    session_factory = build_session_factory(tmp_path)
    with session_factory() as session:
        request = LXDRRequestRecord(
            request_unique_identification_local="3838JBNM5X",
            request_type="PM",
            primary_segment_name="mobility_pax",
            request_priority="02",
            element_unit_identification_callsign="KL9K",
            physical_location_of_requestor="4QFJ1234567890",
            created_at_local="2027OCT13T15470352",
            segment_count=1,
            canonical_text="REQ",
            payload_json={},
            latest_frame_state="QUEUED",
        )
        capability = Capability(
            node_id="atak-node-2",
            callsign="NODE2",
            capability_type="fabrication",
            title="Fabrication cell",
            availability_status="available",
            throughput_per_day=4,
            lead_time_minutes=90,
        )
        stock = LogisticsStatusReport(
            node_id="atak-node-3",
            callsign="NODE3",
            report_type="on_hand_usage",
            item_reference="BRKT-4421",
            on_hand_quantity=0,
            usage_rate_per_day=1.0,
            required_quantity=1,
            required_by_local="2027OCT15",
            delivery_method="direct",
            transport_mode="ground",
            payload_json={},
        )
        session.add_all([request, capability, stock])
        session.commit()

        assessment = assess_request_supportability(
            session=session,
            request=request,
        )

    assert assessment.recommended_course_of_action == "fabricate"
    assert assessment.supporting_node_id == "atak-node-2"


def test_supportability_marks_request_blocked_without_stock_or_capability(
    tmp_path: Path,
) -> None:
    """Requests without stock or support capability should be blocked."""

    session_factory = build_session_factory(tmp_path)
    with session_factory() as session:
        request = LXDRRequestRecord(
            request_unique_identification_local="3838JBNM5X",
            request_type="SR",
            primary_segment_name="supply_request",
            request_priority="02",
            element_unit_identification_callsign="KL9K",
            physical_location_of_requestor="4QFJ1234567890",
            created_at_local="2027OCT13T15470352",
            segment_count=1,
            canonical_text="REQ",
            payload_json={},
            latest_frame_state="QUEUED",
        )
        session.add(request)
        session.commit()

        assessment = assess_request_supportability(
            session=session,
            request=request,
        )

    assert assessment.supportable is False
    assert assessment.recommended_course_of_action == "blocked"
