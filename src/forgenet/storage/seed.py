"""Bootstrap and seed helpers for ForgeNet."""

from __future__ import annotations

from sqlalchemy import select
from sqlalchemy.orm import Session

from forgenet.domain.models import (
    ArtifactKind,
    CapabilityType,
    EventActorType,
    EventKind,
    IncidentStatus,
    JobStatus,
)
from forgenet.storage.tables import Artifact, Capability, Event, Incident, Job


def seed_demo_data(session: Session) -> None:
    """Insert a small, deterministic demo dataset if the database is empty."""

    existing_incident = session.scalar(select(Incident.id).limit(1))
    if existing_incident is not None:
        return

    requester_incident = Incident(
        external_uid="incident-alpha-001",
        title="MRZR suspension bracket failure",
        description=(
            "Forward team reports a failed suspension bracket and "
            "depleted local spares. "
            "Vehicle is mission critical for casualty movement."
        ),
        failed_component="Suspension bracket",
        part_number="MRZR-BRKT-4421",
        unit_name="Node 1 / Forward Maintenance Team",
        reporting_node_id="atak-node-1",
        reporting_callsign="NODE1",
        location_label="OBJ HARBOR",
        latitude=32.7157,
        longitude=-117.1611,
        priority=1,
        urgency=1,
        mission_impact=(
            "Vehicle unavailable for casualty movement until repaired "
            "or rerouted."
        ),
        local_stock_on_hand=0,
        requested_quantity=1,
        recommended_coa="fabricate",
        recommended_coa_confidence=0.82,
        recommended_coa_rationale=(
            "Nearest fabrication node has matching material stock and "
            "lower ETA than reroute."
        ),
        readiness_delta=-0.15,
        eta_minutes=95,
        status=IncidentStatus.TRIAGED,
    )
    session.add(requester_incident)
    session.flush()

    fab_capability = Capability(
        node_id="atak-node-2",
        callsign="NODE2",
        capability_type=CapabilityType.FABRICATION,
        title="Expeditionary additive fabrication cell",
        description=(
            "Polymer and aluminum bracket fabrication with basic "
            "post-processing."
        ),
        materials={"pla_cf": True, "aluminum_plate_mm": [3, 5, 8]},
        equipment={"printers": 2, "cnc_router": 1},
        skills={"cad_repair": True, "print_prep": True},
        throughput_per_day=6,
        lead_time_minutes=60,
        availability_status="available",
    )
    repair_capability = Capability(
        node_id="atak-node-3",
        callsign="NODE3",
        capability_type=CapabilityType.REPAIR,
        title="Forward repair detachment",
        description=(
            "Field expedient repair and install team with mobile welder."
        ),
        materials={"steel_stock": True, "fasteners": True},
        equipment={"welder": 1, "vehicle_lift": 1},
        skills={"field_repair": True, "installation": True},
        throughput_per_day=4,
        lead_time_minutes=140,
        availability_status="available",
    )
    session.add_all([fab_capability, repair_capability])
    session.flush()

    job = Job(
        incident_id=requester_incident.id,
        assigned_capability_id=fab_capability.id,
        assigned_node_id=fab_capability.node_id,
        assigned_callsign=fab_capability.callsign,
        job_type="fabrication",
        title="Fabricate replacement suspension bracket",
        description=(
            "Produce one replacement bracket and push install "
            "instructions to requester."
        ),
        course_of_action="fabricate",
        status=JobStatus.ASSIGNED,
        priority=1,
        estimated_eta_minutes=95,
    )
    session.add(job)
    session.flush()

    artifact = Artifact(
        incident_id=requester_incident.id,
        job_id=job.id,
        external_uid="artifact-stl-001",
        kind=ArtifactKind.PART_MODEL,
        title="Suspension bracket STL",
        description=(
            "Printable replacement bracket model approved for demo workflow."
        ),
        file_name="mrzr-suspension-bracket-v1.stl",
        file_path="data/artifacts/mrzr-suspension-bracket-v1.stl",
        media_type="model/stl",
        source="forgenet-demo-seed",
        metadata_json={"revision": "v1", "unit": "mm"},
    )
    session.add(artifact)

    session.add_all(
        [
            Event(
                incident_id=requester_incident.id,
                kind=EventKind.INCIDENT_REPORTED,
                actor_type=EventActorType.EUD,
                actor_id="atak-node-1",
                actor_callsign="NODE1",
                summary="Incident reported from forward node.",
                detail="Initial failure report received and stored locally.",
                payload_json={"external_uid": requester_incident.external_uid},
            ),
            Event(
                incident_id=requester_incident.id,
                capability_id=fab_capability.id,
                kind=EventKind.CAPABILITY_ADVERTISED,
                actor_type=EventActorType.EUD,
                actor_id=fab_capability.node_id,
                actor_callsign=fab_capability.callsign,
                summary="Fabrication capability advertised.",
                detail="Node 2 reported additive manufacturing availability.",
            ),
            Event(
                incident_id=requester_incident.id,
                job_id=job.id,
                capability_id=fab_capability.id,
                kind=EventKind.JOB_ASSIGNED,
                actor_type=EventActorType.ALOC,
                actor_id="forgenet-aloc",
                actor_callsign="ALOC",
                summary="Fabrication job assigned.",
                detail=(
                    "ALOC selected fabrication as the recommended "
                    "course of action."
                ),
            ),
            Event(
                incident_id=requester_incident.id,
                job_id=job.id,
                artifact_id=artifact.id,
                kind=EventKind.ARTIFACT_REGISTERED,
                actor_type=EventActorType.SYSTEM,
                actor_id="forgenet-bootstrap",
                summary="Seed artifact registered.",
                detail="Demo part model attached to job and incident.",
            ),
        ]
    )

    session.commit()
