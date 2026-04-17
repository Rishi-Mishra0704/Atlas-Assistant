from __future__ import annotations

from fastapi import APIRouter

router = APIRouter(prefix="/wake", tags=["wake"])


@router.get("/health")
def health() -> dict[str, str]:
    return {"status": "ok"}
