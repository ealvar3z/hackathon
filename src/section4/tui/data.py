"""Read-only data loading helpers for the s4net console."""

from __future__ import annotations

from dataclasses import dataclass

from sqlalchemy import func, select
from sqlalchemy.orm import Session, sessionmaker

from section4.storage.tables import (
    Artifact,
    Capability,
    Event,
    Incident,
    Job,
    LXDRInboundFrame,
    LXDROutboundFrame,
)


@dataclass(frozen=True)
class DetailLine:
    """A rendered line in the right-hand detail panel."""

    text: str
    style: str = "body"


@dataclass(frozen=True)
class BrowserItem:
    """A list item and its rendered detail text."""

    item_id: str
    label: str
    detail_lines: list[DetailLine]


@dataclass(frozen=True)
class DashboardData:
    """Top-level dashboard values for the s4net COP view."""

    summary_lines: list[DetailLine]
    recent_event_lines: list[DetailLine]


def _fmt_value(value: object | None) -> str:
    """Return a UI-friendly string for optional values."""

    if value is None:
        return "-"
    return str(value)


def _join_json_keys(value: object | None) -> str:
    """Render JSON-like dictionaries as a short comma-delimited list."""

    if not isinstance(value, dict) or not value:
        return "-"
    return ", ".join(sorted(str(key) for key in value))


def _section(title: str) -> DetailLine:
    """Create a section heading for the detail panel."""

    return DetailLine(title.upper(), style="section")


def _field(label: str, value: object | None) -> DetailLine:
    """Create a simple key-value line."""

    return DetailLine(f"{label}: {_fmt_value(value)}")


def _blank() -> DetailLine:
    """Create a blank spacer line."""

    return DetailLine("", style="muted")


def load_dashboard(session: Session) -> DashboardData:
    """Load aggregate counts and recent audit events."""

    incident_count = (
        session.scalar(select(func.count()).select_from(Incident)) or 0
    )
    capability_count = (
        session.scalar(select(func.count()).select_from(Capability)) or 0
    )
    job_count = session.scalar(select(func.count()).select_from(Job)) or 0
    artifact_count = (
        session.scalar(select(func.count()).select_from(Artifact)) or 0
    )
    event_count = session.scalar(select(func.count()).select_from(Event)) or 0
    outbound_count = (
        session.scalar(select(func.count()).select_from(LXDROutboundFrame)) or 0
    )
    inbound_count = (
        session.scalar(select(func.count()).select_from(LXDRInboundFrame)) or 0
    )
    queued_count = _count_outbound_state(session, "QUEUED")
    sending_count = _count_outbound_state(session, "SENDING")
    failed_count = _count_outbound_state(session, "FAILED")
    synced_count = _count_outbound_state(session, "SYNCED")

    latest_incident = session.scalar(
        select(Incident).order_by(Incident.updated_at.desc()).limit(1)
    )
    latest_job = session.scalar(
        select(Job).order_by(Job.updated_at.desc()).limit(1)
    )

    recent_events = session.scalars(
        select(Event).order_by(Event.occurred_at.desc()).limit(5)
    ).all()

    summary_lines = [
        _section("COP Summary"),
        DetailLine(f"Requests: {incident_count}"),
        DetailLine(f"Capabilities: {capability_count}"),
        DetailLine(f"Tasks: {job_count}"),
        DetailLine(f"Artifacts: {artifact_count}"),
        DetailLine(f"Events: {event_count}"),
        _blank(),
        _section("LXDR Lifecycle"),
        DetailLine(f"Outbound frames: {outbound_count}"),
        DetailLine(f"Inbound frames: {inbound_count}"),
        DetailLine(f"Queued: {queued_count}"),
        DetailLine(f"Sending: {sending_count}"),
        DetailLine(f"Failed: {failed_count}"),
        DetailLine(f"Synced: {synced_count}"),
        _blank(),
        _section("Latest Request"),
        DetailLine(
            f"{latest_incident.title} [{latest_incident.status}]"
            if latest_incident
            else "-"
        ),
        _blank(),
        _section("Latest Task"),
        DetailLine(
            f"{latest_job.title} [{latest_job.status}]" if latest_job else "-"
        ),
    ]
    recent_event_lines = [
        DetailLine(
            f"{event.occurred_at:%Y-%m-%d %H:%M} {event.kind}: {event.summary}"
        )
        for event in recent_events
    ] or [DetailLine("No events recorded.", style="muted")]
    recent_event_lines.insert(0, _section("Latest Events"))
    return DashboardData(
        summary_lines=summary_lines,
        recent_event_lines=recent_event_lines,
    )


def load_incident_items(session: Session) -> list[BrowserItem]:
    """Load incidents for the ADRIAN request page."""

    incidents = session.scalars(
        select(Incident).order_by(Incident.updated_at.desc())
    ).all()
    items: list[BrowserItem] = []
    for incident in incidents:
        items.append(
            BrowserItem(
                item_id=incident.id,
                label=(
                    f"{incident.title} [{incident.status}] P{incident.priority}"
                ),
                detail_lines=[
                    _section("Overview"),
                    _field("Request ID", incident.id),
                    _field("Request UID", incident.external_uid),
                    _field("Status", incident.status),
                    _field("Unit", incident.unit_name),
                    _field("Reporting callsign", incident.reporting_callsign),
                    _field("Reporting node", incident.reporting_node_id),
                    _field("Location", incident.location_label),
                    _field("Urgency", incident.urgency),
                    _field("Priority", incident.priority),
                    _field("Stock on hand", incident.local_stock_on_hand),
                    _field("Requested quantity", incident.requested_quantity),
                    _field("Recommended COA", incident.recommended_coa),
                    _field("ETA minutes", incident.eta_minutes),
                    _blank(),
                    _section("Requested Support"),
                    _field("Component", incident.failed_component),
                    _field("Part number", incident.part_number),
                    _blank(),
                    _section("Description"),
                    DetailLine(incident.description),
                    _blank(),
                    _section("Mission Impact"),
                    DetailLine(_fmt_value(incident.mission_impact)),
                    _blank(),
                    _section("COA Rationale"),
                    DetailLine(_fmt_value(incident.recommended_coa_rationale)),
                ],
            )
        )
    return items


def load_capability_items(session: Session) -> list[BrowserItem]:
    """Load capabilities for the capabilities page."""

    capabilities = session.scalars(
        select(Capability).order_by(Capability.last_reported_at.desc())
    ).all()
    items: list[BrowserItem] = []
    for capability in capabilities:
        items.append(
            BrowserItem(
                item_id=capability.id,
                label=(
                    f"{capability.callsign or capability.node_id} "
                    f"[{capability.capability_type}]"
                ),
                detail_lines=[
                    _section("Overview"),
                    _field("ID", capability.id),
                    _field("Node ID", capability.node_id),
                    _field("Callsign", capability.callsign),
                    _field("Type", capability.capability_type),
                    _field("Title", capability.title),
                    _field("Availability", capability.availability_status),
                    _field("Throughput/day", capability.throughput_per_day),
                    _field("Lead time minutes", capability.lead_time_minutes),
                    _field("Active", capability.active),
                    _field(
                        "Last reported",
                        capability.last_reported_at.strftime("%Y-%m-%d %H:%M"),
                    ),
                    _blank(),
                    _section("Description"),
                    DetailLine(_fmt_value(capability.description)),
                    _blank(),
                    _section("Resources"),
                    _field("Materials", _join_json_keys(capability.materials)),
                    _field("Equipment", _join_json_keys(capability.equipment)),
                    _field("Skills", _join_json_keys(capability.skills)),
                ],
            )
        )
    return items


def load_job_items(session: Session) -> list[BrowserItem]:
    """Load jobs for the task page."""

    jobs = session.scalars(select(Job).order_by(Job.updated_at.desc())).all()
    items: list[BrowserItem] = []
    for job in jobs:
        items.append(
            BrowserItem(
                item_id=job.id,
                label=(
                    f"{job.title} "
                    f"[{job.status}] "
                    f"{_fmt_value(job.assigned_callsign)}"
                ),
                detail_lines=[
                    _section("Overview"),
                    _field("ID", job.id),
                    _field("Incident ID", job.incident_id),
                    _field("Status", job.status),
                    _field("Type", job.job_type),
                    _field("Assigned capability", job.assigned_capability_id),
                    _field("Assigned node", job.assigned_node_id),
                    _field("Assigned callsign", job.assigned_callsign),
                    _field("Course of action", job.course_of_action),
                    _field("Priority", job.priority),
                    _field("Estimated ETA", job.estimated_eta_minutes),
                    _blank(),
                    _section("Description"),
                    DetailLine(_fmt_value(job.description)),
                ],
            )
        )
    return items


def load_event_items(session: Session) -> list[BrowserItem]:
    """Load recent events for the sync and event log page."""

    events = session.scalars(
        select(Event).order_by(Event.occurred_at.desc())
    ).all()
    items: list[BrowserItem] = []
    for event in events:
        items.append(
            BrowserItem(
                item_id=event.id,
                label=(
                    f"{event.occurred_at:%Y-%m-%d %H:%M} "
                    f"{event.kind}: {event.summary}"
                ),
                detail_lines=[
                    _section("Overview"),
                    _field("ID", event.id),
                    _field("Kind", event.kind),
                    _field(
                        "Occurred",
                        event.occurred_at.strftime("%Y-%m-%d %H:%M:%S"),
                    ),
                    _field("Actor type", event.actor_type),
                    _field("Actor ID", event.actor_id),
                    _field("Actor callsign", event.actor_callsign),
                    _field("Incident ID", event.incident_id),
                    _field("Job ID", event.job_id),
                    _field("Capability ID", event.capability_id),
                    _field("Artifact ID", event.artifact_id),
                    _blank(),
                    _section("Summary"),
                    DetailLine(event.summary),
                    _blank(),
                    _section("Detail"),
                    DetailLine(_fmt_value(event.detail)),
                ],
            )
        )
    return items


def load_sync_log_items(session: Session) -> list[BrowserItem]:
    """Load persisted LXDR inbound and outbound frame state."""

    outbound_rows = session.scalars(
        select(LXDROutboundFrame).order_by(LXDROutboundFrame.updated_at.desc())
    ).all()
    inbound_rows = session.scalars(
        select(LXDRInboundFrame).order_by(LXDRInboundFrame.created_at.desc())
    ).all()
    items: list[BrowserItem] = []

    for row in outbound_rows:
        items.append(
            BrowserItem(
                item_id=row.id,
                label=(
                    f"OUT {row.link_message_id} "
                    f"[{row.state}] {row.request_unique_identification_local}"
                ),
                detail_lines=[
                    _section("Outbound Frame"),
                    _field("Message ID", row.link_message_id),
                    _field("Sender", row.sender_id),
                    _field("Recipient", row.recipient_id),
                    _field(
                        "Local request ID",
                        row.request_unique_identification_local,
                    ),
                    _field(
                        "Sync request ID",
                        row.request_unique_identification_sync,
                    ),
                    _field("Delivery method", row.delivery_method),
                    _field("Representation", row.representation),
                    _field("State", row.state),
                    _field("Attempt count", row.attempt_count),
                    _field("Last attempt", row.last_attempt_at),
                    _field("Last error", row.last_error),
                    _field("Created local", row.created_at_local),
                    _field("Correlation ID", row.correlation_id),
                ],
            )
        )

    for row in inbound_rows:
        items.append(
            BrowserItem(
                item_id=row.id,
                label=(
                    f"IN {row.link_message_id} "
                    f"[{row.payload_count} payloads]"
                ),
                detail_lines=[
                    _section("Inbound Frame"),
                    _field("Message ID", row.link_message_id),
                    _field("Sender", row.sender_id),
                    _field("Recipient", row.recipient_id),
                    _field("Payload count", row.payload_count),
                    _field(
                        "Recorded",
                        row.created_at.strftime("%Y-%m-%d %H:%M:%S"),
                    ),
                ],
            )
        )

    return items


def _count_outbound_state(session: Session, state: str) -> int:
    """Count outbound frames in a specific lifecycle state."""

    return (
        session.scalar(
            select(func.count())
            .select_from(LXDROutboundFrame)
            .where(LXDROutboundFrame.state == state)
        )
        or 0
    )


def load_artifact_items(session: Session) -> list[BrowserItem]:
    """Load artifacts for the artifact browser page."""

    artifacts = session.scalars(
        select(Artifact).order_by(Artifact.updated_at.desc())
    ).all()
    items: list[BrowserItem] = []
    for artifact in artifacts:
        items.append(
            BrowserItem(
                item_id=artifact.id,
                label=f"{artifact.title} [{artifact.kind}]",
                detail_lines=[
                    _section("Overview"),
                    _field("ID", artifact.id),
                    _field("Kind", artifact.kind),
                    _field("Incident ID", artifact.incident_id),
                    _field("Job ID", artifact.job_id),
                    _field("File name", artifact.file_name),
                    _field("File path", artifact.file_path),
                    _field("Media type", artifact.media_type),
                    _field("Source", artifact.source),
                    _blank(),
                    _section("Description"),
                    DetailLine(_fmt_value(artifact.description)),
                ],
            )
        )
    return items


def load_page_items(
    session_factory: sessionmaker[Session],
    page_name: str,
) -> tuple[str, list[BrowserItem]] | DashboardData:
    """Load a named page from the SQLite-backed section4 state."""

    with session_factory() as session:
        if page_name == "dashboard":
            return load_dashboard(session)
        if page_name == "incidents":
            return ("Requests", load_incident_items(session))
        if page_name == "capabilities":
            return ("Capabilities", load_capability_items(session))
        if page_name == "jobs":
            return ("Tasks", load_job_items(session))
        if page_name == "artifacts":
            return ("Artifacts", load_artifact_items(session))
        if page_name == "events":
            return ("Sync Log", load_sync_log_items(session))

    raise ValueError(f"Unknown page name: {page_name}")
