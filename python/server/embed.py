from __future__ import annotations

from fastapi import APIRouter

router = APIRouter(prefix="/embed", tags=["embed"])


@router.get("/health")
def health() -> dict[str, str]:
    return {"status": "ok"}
