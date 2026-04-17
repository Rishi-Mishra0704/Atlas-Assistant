from __future__ import annotations

import os
from dataclasses import dataclass


@dataclass(frozen=True)
class Config:
    postgres_url: str
    qdrant_url: str


def load_config() -> Config:
    return Config(
        postgres_url=os.getenv(
            "POSTGRES_URL",
            "postgresql://atlas:atlas@localhost:5432/atlas",
        ),
        qdrant_url=os.getenv("QDRANT_URL", "http://localhost:6333"),
    )
