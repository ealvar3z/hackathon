"""Configuration models and settings for section4."""

from __future__ import annotations

from dataclasses import dataclass
from pathlib import Path

from section4.storage.db import DEFAULT_DB_PATH


@dataclass(slots=True)
class Settings:
    """Runtime settings for the local section4 node."""

    db_path: Path = DEFAULT_DB_PATH
    artifacts_dir: Path = Path("data/artifacts")
    cot_url: str = "udp://239.2.3.1:6969"


def get_settings() -> Settings:
    """Return default settings for the current local environment."""

    settings = Settings()
    settings.artifacts_dir.mkdir(parents=True, exist_ok=True)
    settings.db_path.parent.mkdir(parents=True, exist_ok=True)
    return settings
