"""Persistence helpers for transport-layer events."""

from __future__ import annotations

from sqlalchemy.orm import Session

from forgenet.domain.models import EventActorType, EventKind
from forgenet.storage.tables import Event
from forgenet.transport.cot import ParsedCoTEvent


def record_received_cot(session: Session, parsed: ParsedCoTEvent) -> Event:
    """Persist a received CoT event into the ForgeNet audit log."""

    incident_id = parsed.detail_attributes.get("incident_id")
    job_id = parsed.detail_attributes.get("job_id")
    actor_type = EventActorType.EXTERNAL
    if parsed.callsign:
        actor_type = EventActorType.EUD

    event = Event(
        incident_id=incident_id,
        job_id=job_id,
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
) -> Event:
    """Persist a locally published CoT event into the ForgeNet audit log."""

    event = Event(
        incident_id=incident_id,
        job_id=job_id,
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
