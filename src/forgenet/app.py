"""Application entrypoint for ForgeNet."""

from __future__ import annotations

import argparse

from forgenet.config import get_settings
from forgenet.storage import create_all, create_session_factory
from forgenet.storage.seed import seed_demo_data


def build_parser() -> argparse.ArgumentParser:
    """Create the top-level CLI parser."""

    parser = argparse.ArgumentParser(description="ForgeNet local operator tooling")
    subparsers = parser.add_subparsers(dest="command", required=True)

    subparsers.add_parser("init-db", help="Create the local SQLite schema")
    subparsers.add_parser("seed-demo", help="Seed a tiny demo dataset if the DB is empty")
    subparsers.add_parser(
        "bootstrap",
        help="Create the local SQLite schema and seed a demo dataset if needed",
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

    raise SystemExit(f"Unknown command: {args.command}")


if __name__ == "__main__":
    main()
