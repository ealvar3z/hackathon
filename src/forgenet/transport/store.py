"""Persistence helpers for transport-layer events."""

from __future__ import annotations

from datetime import datetime

from sqlalchemy.orm import Session

from forgenet.domain.models import CapabilityType, EventActorType, EventKind
from forgenet.storage.tables import Capability, Event
from forgenet.transport.cot import ParsedCoTEvent


def record_received_cot(session: Session, parsed: ParsedCoTEvent) -> Event:
    """Persist a received CoT event into the ForgeNet audit log."""

    incident_id = parsed.detail_attributes.get("incident_id")
    job_id = parsed.detail_attributes.get("job_id")
    capability_id = parsed.detail_attributes.get("capability_id")
    actor_type = EventActorType.EXTERNAL
    if parsed.callsign:
        actor_type = EventActorType.EUD

    event = Event(
        incident_id=incident_id,
        job_id=job_id,
        capability_id=capability_id,
        kind=EventKind.COT_RECEIVED,
        actor_type=actor_type,
        actor_id=parsed.uid,
        actor_callsign=parsed.callsign,
        summary=f"Received CoT event {parsed.cot_type or 'unknown'}",
        detail=parsed.remarks,
        payload_json={
            "uid": parsed.uid,
            "cot_type": parsed.cot_type,
            "callsign": parsed.callsign,
            "lat": parsed.lat,
            "lon": parsed.lon,
            "how": parsed.how,
            "forgenet": parsed.detail_attributes,
            "raw_xml": parsed.raw_xml,
        },
    )
    session.add(event)
    session.commit()
    return event


def record_published_cot(
    session: Session,
    *,
    summary: str,
    payload: dict,
    incident_id: str | None = None,
    job_id: str | None = None,
    capability_id: str | None = None,
) -> Event:
    """Persist a locally published CoT event into the ForgeNet audit log."""

    event = Event(
        incident_id=incident_id,
        job_id=job_id,
        capability_id=capability_id,
        kind=EventKind.COT_PUBLISHED,
        actor_type=EventActorType.ALOC,
        actor_id="forgenet-aloc",
        actor_callsign="ALOC",
        summary=summary,
        payload_json=payload,
    )
    session.add(event)
    session.commit()
    return event


def upsert_capability_from_cot(
    session: Session,
    parsed: ParsedCoTEvent,
) -> Capability | None:
    """Upsert a capability from a ForgeNet capability CoT event."""

    attrs = parsed.detail_attributes
    if attrs.get("object") != "capability":
        return None

    capability_id = attrs.get("capability_id")
    node_id = attrs.get("node_id")
    capability_type = attrs.get("capability_type")
    title = attrs.get("title")
    availability_status = attrs.get("availability_status")

    if not all([capability_id, node_id, capability_type, title]):
        return None

    capability = session.get(Capability, capability_id)
    if capability is None:
        capability = Capability(
            id=capability_id,
            node_id=node_id,
            capability_type=CapabilityType(capability_type),
            title=title,
        )
        session.add(capability)

    capability.callsign = parsed.callsign
    capability.description = parsed.remarks
    capability.availability_status = availability_status or "available"
    capability.throughput_per_day = _parse_optional_int(
        attrs.get("throughput_per_day")
    )
    capability.lead_time_minutes = _parse_optional_int(
        attrs.get("lead_time_minutes")
    )
    capability.last_reported_at = datetime.utcnow()
    capability.active = True

    session.commit()
    return capability


def _parse_optional_int(value: str | None) -> int | None:
    """Parse an optional integer value from CoT detail attributes."""

    if value is None:
        return None
    return int(value)
