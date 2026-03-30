"""Application entrypoint for ForgeNet."""

from __future__ import annotations

import argparse
import asyncio

from sqlalchemy import select

from forgenet.config import get_settings
from forgenet.storage import create_all, create_session_factory
from forgenet.storage.seed import seed_demo_data
from forgenet.storage.tables import Incident, Job
from forgenet.transport import (
    PyTAKRuntime,
    build_incident_cot,
    build_job_cot,
    parse_cot_event,
    record_published_cot,
    record_received_cot,
)


def build_parser() -> argparse.ArgumentParser:
    """Create the top-level CLI parser."""

    parser = argparse.ArgumentParser(
        description="ForgeNet local operator tooling"
    )
    subparsers = parser.add_subparsers(dest="command", required=True)

    subparsers.add_parser("init-db", help="Create the local SQLite schema")
    subparsers.add_parser(
        "seed-demo", help="Seed a tiny demo dataset if the DB is empty"
    )
    subparsers.add_parser(
        "bootstrap",
        help="Create the local SQLite schema and seed a demo dataset if needed",
    )
    publish_incident = subparsers.add_parser(
        "publish-incident",
        help="Publish an incident record as a CoT event via PyTAK",
    )
    publish_incident.add_argument(
        "--incident-id", help="Incident ID to publish"
    )
    publish_incident.add_argument(
        "--cot-url", help="Override the default CoT URL"
    )

    publish_job = subparsers.add_parser(
        "publish-job",
        help="Publish a job record as a CoT event via PyTAK",
    )
    publish_job.add_argument("--job-id", help="Job ID to publish")
    publish_job.add_argument("--cot-url", help="Override the default CoT URL")

    receive_once = subparsers.add_parser(
        "receive-once",
        help="Wait for one inbound CoT event and persist it to the event log",
    )
    receive_once.add_argument("--cot-url", help="Override the default CoT URL")
    receive_once.add_argument(
        "--timeout",
        type=float,
        default=10.0,
        help="Seconds to wait for one inbound CoT event",
    )

    return parser


def cmd_init_db() -> None:
    """Create the local database schema."""

    settings = get_settings()
    create_all(settings.db_path)
    print(f"Initialized ForgeNet database at {settings.db_path}")


def cmd_seed_demo() -> None:
    """Seed demo data into the local database."""

    settings = get_settings()
    create_all(settings.db_path)
    session_factory = create_session_factory(settings.db_path)
    with session_factory() as session:
        seed_demo_data(session)
    print(f"Seeded ForgeNet demo data in {settings.db_path}")


async def _publish_incident(
    incident_id: str | None, cot_url: str | None
) -> None:
    settings = get_settings()
    create_all(settings.db_path)
    session_factory = create_session_factory(settings.db_path)

    with session_factory() as session:
        query = select(Incident).order_by(Incident.created_at.asc())
        if incident_id:
            query = query.where(Incident.id == incident_id)
        incident = session.scalar(query.limit(1))
        if incident is None:
            raise SystemExit("No incident found to publish")

        payload = build_incident_cot(incident)
        runtime = PyTAKRuntime(cot_url or settings.cot_url)
        await runtime.start()
        try:
            await runtime.publish(payload)
        finally:
            await runtime.stop()

        record_published_cot(
            session,
            summary=f"Published incident {incident.id} as CoT",
            incident_id=incident.id,
            payload={
                "uid": incident.external_uid
                or f"forgenet-incident-{incident.id}",
                "cot_url": cot_url or settings.cot_url,
                "kind": "incident",
            },
        )
        print(
            f"Published incident {incident.id} to {cot_url or settings.cot_url}"
        )


async def _publish_job(job_id: str | None, cot_url: str | None) -> None:
    settings = get_settings()
    create_all(settings.db_path)
    session_factory = create_session_factory(settings.db_path)

    with session_factory() as session:
        query = select(Job).order_by(Job.created_at.asc())
        if job_id:
            query = query.where(Job.id == job_id)
        job = session.scalar(query.limit(1))
        if job is None:
            raise SystemExit("No job found to publish")

        payload = build_job_cot(job)
        runtime = PyTAKRuntime(cot_url or settings.cot_url)
        await runtime.start()
        try:
            await runtime.publish(payload)
        finally:
            await runtime.stop()

        record_published_cot(
            session,
            summary=f"Published job {job.id} as CoT",
            incident_id=job.incident_id,
            job_id=job.id,
            payload={
                "uid": f"forgenet-job-{job.id}",
                "cot_url": cot_url or settings.cot_url,
                "kind": "job",
            },
        )
        print(f"Published job {job.id} to {cot_url or settings.cot_url}")


async def _receive_once(cot_url: str | None, timeout: float) -> None:
    settings = get_settings()
    create_all(settings.db_path)
    session_factory = create_session_factory(settings.db_path)
    runtime = PyTAKRuntime(cot_url or settings.cot_url)
    await runtime.start()
    try:
        data = await runtime.receive_once(timeout=timeout)
    finally:
        await runtime.stop()

    parsed = parse_cot_event(data)
    with session_factory() as session:
        event = record_received_cot(session, parsed)

    print(f"Persisted inbound CoT event as audit event {event.id}")


def main() -> None:
    """Run the ForgeNet command-line entrypoint."""

    args = build_parser().parse_args()

    if args.command == "init-db":
        cmd_init_db()
        return

    if args.command == "seed-demo":
        cmd_seed_demo()
        return

    if args.command == "bootstrap":
        cmd_init_db()
        cmd_seed_demo()
        return

    if args.command == "publish-incident":
        asyncio.run(_publish_incident(args.incident_id, args.cot_url))
        return

    if args.command == "publish-job":
        asyncio.run(_publish_job(args.job_id, args.cot_url))
        return

    if args.command == "receive-once":
        asyncio.run(_receive_once(args.cot_url, args.timeout))
        return

    raise SystemExit(f"Unknown command: {args.command}")


if __name__ == "__main__":
    main()
