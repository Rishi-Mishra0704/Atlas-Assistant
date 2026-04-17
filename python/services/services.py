from __future__ import annotations

from dataclasses import dataclass


@dataclass
class ServiceResult:
    ok: bool
    message: str = ""
