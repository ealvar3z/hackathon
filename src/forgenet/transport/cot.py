"""CoT encoding and decoding helpers for ForgeNet."""

from __future__ import annotations

import xml.etree.ElementTree as ET
from dataclasses import dataclass
from typing import Any

import pytak

from forgenet.storage.tables import Incident, Job

FORGENET_DETAIL_TAG = "forgenet"
FORGENET_COT_TYPE_INCIDENT = "b-m-r"
FORGENET_COT_TYPE_JOB = "b-m-t"


@dataclass(slots=True)
class ParsedCoTEvent:
    """Normalized view of a received CoT event."""

    uid: str | None
    cot_type: str | None
    callsign: str | None
    lat: float | None
    lon: float | None
    how: str | None
    remarks: str | None
    detail_attributes: dict[str, Any]
    raw_xml: str


def _base_event(
    uid: str, cot_type: str, stale_seconds: int = 300
) -> ET.Element:
    """Create a minimal CoT event skeleton."""

    event = ET.Element("event")
    event.set("version", "2.0")
    event.set("uid", uid)
    event.set("type", cot_type)
    event.set("how", "h-g-i-g-o")
    event.set("time", pytak.cot_time())
    event.set("start", pytak.cot_time())
    event.set("stale", pytak.cot_time(stale_seconds))
    return event


def build_incident_cot(incident: Incident) -> bytes:
    """Encode an incident as a ForgeNet CoT event."""

    uid = incident.external_uid or f"forgenet-incident-{incident.id}"
    event = _base_event(uid=uid, cot_type=FORGENET_COT_TYPE_INCIDENT)

    point = ET.SubElement(event, "point")
    point.set(
        "lat", str(incident.latitude if incident.latitude is not None else 0.0)
    )
    point.set(
        "lon",
        str(incident.longitude if incident.longitude is not None else 0.0),
    )
    point.set("hae", "9999999.0")
    point.set("ce", "9999999.0")
    point.set("le", "9999999.0")

    detail = ET.SubElement(event, "detail")

    if incident.reporting_callsign:
        contact = ET.SubElement(detail, "contact")
        contact.set("callsign", incident.reporting_callsign)

    remarks = ET.SubElement(detail, "remarks")
    remarks.text = incident.description

    forge = ET.SubElement(detail, FORGENET_DETAIL_TAG)
    forge.set("object", "incident")
    forge.set("incident_id", incident.id)
    forge.set("status", incident.status.value)
    if incident.part_number:
        forge.set("part_number", incident.part_number)
    if incident.failed_component:
        forge.set("failed_component", incident.failed_component)
    if incident.recommended_coa:
        forge.set("recommended_coa", incident.recommended_coa)

    return pytak.DEFAULT_XML_DECLARATION + b"\n" + ET.tostring(event)


def build_job_cot(job: Job) -> bytes:
    """Encode a job as a ForgeNet CoT event."""

    uid = f"forgenet-job-{job.id}"
    event = _base_event(uid=uid, cot_type=FORGENET_COT_TYPE_JOB)

    point = ET.SubElement(event, "point")
    incident = job.incident
    lat = (
        incident.latitude if incident and incident.latitude is not None else 0.0
    )
    lon = (
        incident.longitude
        if incident and incident.longitude is not None
        else 0.0
    )
    point.set("lat", str(lat))
    point.set("lon", str(lon))
    point.set("hae", "9999999.0")
    point.set("ce", "9999999.0")
    point.set("le", "9999999.0")

    detail = ET.SubElement(event, "detail")
    if job.assigned_callsign:
        contact = ET.SubElement(detail, "contact")
        contact.set("callsign", job.assigned_callsign)

    remarks = ET.SubElement(detail, "remarks")
    remarks.text = job.description or job.title

    forge = ET.SubElement(detail, FORGENET_DETAIL_TAG)
    forge.set("object", "job")
    forge.set("job_id", job.id)
    forge.set("incident_id", job.incident_id)
    forge.set("status", job.status.value)
    forge.set("job_type", job.job_type)
    if job.course_of_action:
        forge.set("course_of_action", job.course_of_action)

    return pytak.DEFAULT_XML_DECLARATION + b"\n" + ET.tostring(event)


def parse_cot_event(data: bytes) -> ParsedCoTEvent:
    """Parse incoming CoT bytes into a normalized structure."""

    raw_xml = data.decode("utf-8", errors="replace")
    root = ET.fromstring(data)
    point = root.find("point")
    detail = root.find("detail")
    contact = detail.find("contact") if detail is not None else None
    remarks = detail.find("remarks") if detail is not None else None
    forge = detail.find(FORGENET_DETAIL_TAG) if detail is not None else None

    lat = None
    lon = None
    if point is not None:
        lat_value = point.attrib.get("lat")
        lon_value = point.attrib.get("lon")
        lat = float(lat_value) if lat_value is not None else None
        lon = float(lon_value) if lon_value is not None else None

    detail_attributes: dict[str, Any] = {}
    if forge is not None:
        detail_attributes = dict(forge.attrib)

    return ParsedCoTEvent(
        uid=root.attrib.get("uid"),
        cot_type=root.attrib.get("type"),
        callsign=contact.attrib.get("callsign")
        if contact is not None
        else None,
        lat=lat,
        lon=lon,
        how=root.attrib.get("how"),
        remarks=remarks.text if remarks is not None else None,
        detail_attributes=detail_attributes,
        raw_xml=raw_xml,
    )
