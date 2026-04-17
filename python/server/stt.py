from __future__ import annotations

from fastapi import APIRouter

router = APIRouter(prefix="/stt", tags=["stt"])


@router.get("/health")
def health() -> dict[str, str]:
    return {"status": "ok"}
