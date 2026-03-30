"""SQLAlchemy table definitions for ForgeNet persistence."""

from __future__ import annotations

from datetime import datetime
from typing import Any
from uuid import uuid4

from sqlalchemy import JSON, DateTime, Enum, Float, ForeignKey, Index, Integer, String, Text
from sqlalchemy.orm import Mapped, mapped_column, relationship

from forgenet.domain.models import (
    ArtifactKind,
    CapabilityType,
    EventActorType,
    EventKind,
    IncidentStatus,
    JobStatus,
)
from forgenet.storage.db import Base


def _utcnow() -> datetime:
    """Return a naive UTC timestamp for SQLite storage."""

    return datetime.utcnow()


def _new_id() -> str:
    """Generate a compact UUID string for primary keys."""

    return str(uuid4())


class Incident(Base):
    """Reported maintenance or logistics incident."""

    __tablename__ = "incidents"

    id: Mapped[str] = mapped_column(String(36), primary_key=True, default=_new_id)
    external_uid: Mapped[str | None] = mapped_column(String(255), index=True)
    title: Mapped[str] = mapped_column(String(255))
    description: Mapped[str] = mapped_column(Text)
    failed_component: Mapped[str | None] = mapped_column(String(255))
    part_number: Mapped[str | None] = mapped_column(String(255), index=True)
    unit_name: Mapped[str | None] = mapped_column(String(255))
    reporting_node_id: Mapped[str | None] = mapped_column(String(255), index=True)
    reporting_callsign: Mapped[str | None] = mapped_column(String(255))
    location_label: Mapped[str | None] = mapped_column(String(255))
    latitude: Mapped[float | None] = mapped_column(Float)
    longitude: Mapped[float | None] = mapped_column(Float)
    priority: Mapped[int] = mapped_column(Integer, default=3)
    urgency: Mapped[int] = mapped_column(Integer, default=3)
    mission_impact: Mapped[str | None] = mapped_column(Text)
    local_stock_on_hand: Mapped[int | None] = mapped_column(Integer)
    requested_quantity: Mapped[int | None] = mapped_column(Integer)
    recommended_coa: Mapped[str | None] = mapped_column(String(64))
    recommended_coa_confidence: Mapped[float | None] = mapped_column(Float)
    recommended_coa_rationale: Mapped[str | None] = mapped_column(Text)
    readiness_delta: Mapped[float | None] = mapped_column(Float)
    eta_minutes: Mapped[int | None] = mapped_column(Integer)
    status: Mapped[IncidentStatus] = mapped_column(
        Enum(IncidentStatus, native_enum=False), default=IncidentStatus.NEW, index=True
    )
    created_at: Mapped[datetime] = mapped_column(DateTime, default=_utcnow)
    updated_at: Mapped[datetime] = mapped_column(DateTime, default=_utcnow, onupdate=_utcnow)

    jobs: Mapped[list["Job"]] = relationship(back_populates="incident")
    artifacts: Mapped[list["Artifact"]] = relationship(back_populates="incident")
    events: Mapped[list["Event"]] = relationship(back_populates="incident")


class Capability(Base):
    """Capability advertised by a node on the network."""

    __tablename__ = "capabilities"

    id: Mapped[str] = mapped_column(String(36), primary_key=True, default=_new_id)
    node_id: Mapped[str] = mapped_column(String(255), index=True)
    callsign: Mapped[str | None] = mapped_column(String(255))
    capability_type: Mapped[CapabilityType] = mapped_column(
        Enum(CapabilityType, native_enum=False), index=True
    )
    title: Mapped[str] = mapped_column(String(255))
    description: Mapped[str | None] = mapped_column(Text)
    materials: Mapped[dict[str, Any] | None] = mapped_column(JSON)
    equipment: Mapped[dict[str, Any] | None] = mapped_column(JSON)
    skills: Mapped[dict[str, Any] | None] = mapped_column(JSON)
    throughput_per_day: Mapped[int | None] = mapped_column(Integer)
    lead_time_minutes: Mapped[int | None] = mapped_column(Integer)
    availability_status: Mapped[str] = mapped_column(String(64), default="available")
    active: Mapped[bool] = mapped_column(default=True, index=True)
    last_reported_at: Mapped[datetime] = mapped_column(DateTime, default=_utcnow, index=True)
    created_at: Mapped[datetime] = mapped_column(DateTime, default=_utcnow)
    updated_at: Mapped[datetime] = mapped_column(DateTime, default=_utcnow, onupdate=_utcnow)

    jobs: Mapped[list["Job"]] = relationship(back_populates="assigned_capability")
    events: Mapped[list["Event"]] = relationship(back_populates="capability")


class Job(Base):
    """Actionable workflow derived from an incident."""

    __tablename__ = "jobs"

    id: Mapped[str] = mapped_column(String(36), primary_key=True, default=_new_id)
    incident_id: Mapped[str] = mapped_column(
        String(36), ForeignKey("incidents.id", ondelete="CASCADE"), index=True
    )
    assigned_capability_id: Mapped[str | None] = mapped_column(
        String(36), ForeignKey("capabilities.id", ondelete="SET NULL"), index=True
    )
    assigned_node_id: Mapped[str | None] = mapped_column(String(255), index=True)
    assigned_callsign: Mapped[str | None] = mapped_column(String(255))
    job_type: Mapped[str] = mapped_column(String(64))
    title: Mapped[str] = mapped_column(String(255))
    description: Mapped[str | None] = mapped_column(Text)
    course_of_action: Mapped[str | None] = mapped_column(String(64))
    status: Mapped[JobStatus] = mapped_column(
        Enum(JobStatus, native_enum=False), default=JobStatus.QUEUED, index=True
    )
    priority: Mapped[int] = mapped_column(Integer, default=3)
    estimated_eta_minutes: Mapped[int | None] = mapped_column(Integer)
    actual_completion_at: Mapped[datetime | None] = mapped_column(DateTime)
    assigned_at: Mapped[datetime | None] = mapped_column(DateTime)
    acknowledged_at: Mapped[datetime | None] = mapped_column(DateTime)
    started_at: Mapped[datetime | None] = mapped_column(DateTime)
    completed_at: Mapped[datetime | None] = mapped_column(DateTime)
    created_at: Mapped[datetime] = mapped_column(DateTime, default=_utcnow)
    updated_at: Mapped[datetime] = mapped_column(DateTime, default=_utcnow, onupdate=_utcnow)

    incident: Mapped["Incident"] = relationship(back_populates="jobs")
    assigned_capability: Mapped["Capability | None"] = relationship(back_populates="jobs")
    artifacts: Mapped[list["Artifact"]] = relationship(back_populates="job")
    events: Mapped[list["Event"]] = relationship(back_populates="job")


class Artifact(Base):
    """Artifact attached to an incident or job."""

    __tablename__ = "artifacts"

    id: Mapped[str] = mapped_column(String(36), primary_key=True, default=_new_id)
    incident_id: Mapped[str | None] = mapped_column(
        String(36), ForeignKey("incidents.id", ondelete="CASCADE"), index=True
    )
    job_id: Mapped[str | None] = mapped_column(
        String(36), ForeignKey("jobs.id", ondelete="CASCADE"), index=True
    )
    external_uid: Mapped[str | None] = mapped_column(String(255), index=True)
    kind: Mapped[ArtifactKind] = mapped_column(Enum(ArtifactKind, native_enum=False), index=True)
    title: Mapped[str] = mapped_column(String(255))
    description: Mapped[str | None] = mapped_column(Text)
    file_name: Mapped[str | None] = mapped_column(String(255))
    file_path: Mapped[str | None] = mapped_column(String(1024))
    media_type: Mapped[str | None] = mapped_column(String(255))
    sha256: Mapped[str | None] = mapped_column(String(64), index=True)
    size_bytes: Mapped[int | None] = mapped_column(Integer)
    source: Mapped[str | None] = mapped_column(String(255))
    metadata_json: Mapped[dict[str, Any] | None] = mapped_column(JSON)
    created_at: Mapped[datetime] = mapped_column(DateTime, default=_utcnow)
    updated_at: Mapped[datetime] = mapped_column(DateTime, default=_utcnow, onupdate=_utcnow)

    incident: Mapped["Incident | None"] = relationship(back_populates="artifacts")
    job: Mapped["Job | None"] = relationship(back_populates="artifacts")
    events: Mapped[list["Event"]] = relationship(back_populates="artifact")


class Event(Base):
    """Immutable audit and workflow event."""

    __tablename__ = "events"

    id: Mapped[str] = mapped_column(String(36), primary_key=True, default=_new_id)
    incident_id: Mapped[str | None] = mapped_column(
        String(36), ForeignKey("incidents.id", ondelete="CASCADE"), index=True
    )
    job_id: Mapped[str | None] = mapped_column(
        String(36), ForeignKey("jobs.id", ondelete="CASCADE"), index=True
    )
    capability_id: Mapped[str | None] = mapped_column(
        String(36), ForeignKey("capabilities.id", ondelete="SET NULL"), index=True
    )
    artifact_id: Mapped[str | None] = mapped_column(
        String(36), ForeignKey("artifacts.id", ondelete="SET NULL"), index=True
    )
    kind: Mapped[EventKind] = mapped_column(Enum(EventKind, native_enum=False), index=True)
    actor_type: Mapped[EventActorType] = mapped_column(
        Enum(EventActorType, native_enum=False), default=EventActorType.SYSTEM, index=True
    )
    actor_id: Mapped[str | None] = mapped_column(String(255), index=True)
    actor_callsign: Mapped[str | None] = mapped_column(String(255))
    summary: Mapped[str] = mapped_column(String(255))
    detail: Mapped[str | None] = mapped_column(Text)
    payload_json: Mapped[dict[str, Any] | None] = mapped_column(JSON)
    occurred_at: Mapped[datetime] = mapped_column(DateTime, default=_utcnow, index=True)
    created_at: Mapped[datetime] = mapped_column(DateTime, default=_utcnow)

    incident: Mapped["Incident | None"] = relationship(back_populates="events")
    job: Mapped["Job | None"] = relationship(back_populates="events")
    capability: Mapped["Capability | None"] = relationship(back_populates="events")
    artifact: Mapped["Artifact | None"] = relationship(back_populates="events")


Index("ix_capabilities_node_type_active", Capability.node_id, Capability.capability_type, Capability.active)
Index("ix_jobs_incident_status", Job.incident_id, Job.status)
Index("ix_artifacts_incident_job", Artifact.incident_id, Artifact.job_id)
Index("ix_events_kind_occurred", Event.kind, Event.occurred_at)
