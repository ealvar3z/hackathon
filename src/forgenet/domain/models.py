"""Core ForgeNet domain enums and shared types."""

from __future__ import annotations

from enum import StrEnum


class IncidentStatus(StrEnum):
    """Lifecycle states for a reported incident."""

    NEW = "new"
    TRIAGED = "triaged"
    ASSIGNED = "assigned"
    IN_PROGRESS = "in_progress"
    RESOLVED = "resolved"
    CANCELLED = "cancelled"


class JobStatus(StrEnum):
    """Lifecycle states for a logistics job."""

    QUEUED = "queued"
    ASSIGNED = "assigned"
    ACKNOWLEDGED = "acknowledged"
    IN_PROGRESS = "in_progress"
    BLOCKED = "blocked"
    COMPLETED = "completed"
    CANCELLED = "cancelled"


class CapabilityType(StrEnum):
    """Capability categories advertised by a node."""

    REPAIR = "repair"
    FABRICATION = "fabrication"
    SUPPLY = "supply"
    TRANSPORT = "transport"
    DIAGNOSTICS = "diagnostics"


class ArtifactKind(StrEnum):
    """Types of artifact files or references attached to ForgeNet records."""

    IMAGE = "image"
    DOCUMENT = "document"
    SOP = "sop"
    PART_MODEL = "part_model"
    DATA_PACKAGE = "data_package"
    OTHER = "other"


class EventKind(StrEnum):
    """Audit and workflow events emitted by ForgeNet."""

    INCIDENT_REPORTED = "incident_reported"
    INCIDENT_UPDATED = "incident_updated"
    CAPABILITY_ADVERTISED = "capability_advertised"
    JOB_CREATED = "job_created"
    JOB_ASSIGNED = "job_assigned"
    JOB_ACKNOWLEDGED = "job_acknowledged"
    JOB_UPDATED = "job_updated"
    JOB_COMPLETED = "job_completed"
    ARTIFACT_REGISTERED = "artifact_registered"
    COT_RECEIVED = "cot_received"
    COT_PUBLISHED = "cot_published"
    SYSTEM = "system"


class EventActorType(StrEnum):
    """Origin of an event in the audit log."""

    SYSTEM = "system"
    ALOC = "aloc"
    EUD = "eud"
    EXTERNAL = "external"
