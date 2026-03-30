"""Database helpers for ForgeNet."""

from __future__ import annotations

from pathlib import Path

from sqlalchemy import create_engine
from sqlalchemy.orm import DeclarativeBase, Session, sessionmaker


DEFAULT_DB_PATH = Path("data/sqlite/forgenet.db")


class Base(DeclarativeBase):
    """Base declarative model for ForgeNet tables."""


def get_database_url(path: Path | str = DEFAULT_DB_PATH) -> str:
    """Build a SQLite connection URL for the given path."""

    db_path = Path(path)
    db_path.parent.mkdir(parents=True, exist_ok=True)
    return f"sqlite:///{db_path}"


def create_engine_for_path(path: Path | str = DEFAULT_DB_PATH):
    """Create a SQLAlchemy engine for the configured SQLite file."""

    return create_engine(get_database_url(path), future=True)


def create_session_factory(path: Path | str = DEFAULT_DB_PATH) -> sessionmaker[Session]:
    """Create a session factory bound to the SQLite engine."""

    engine = create_engine_for_path(path)
    return sessionmaker(bind=engine, expire_on_commit=False, class_=Session)


def create_all(path: Path | str = DEFAULT_DB_PATH) -> None:
    """Create all registered tables."""

    engine = create_engine_for_path(path)
    Base.metadata.create_all(engine)
