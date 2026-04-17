from __future__ import annotations

from fastapi import APIRouter

router = APIRouter(prefix="/tts", tags=["tts"])


@router.get("/health")
def health() -> dict[str, str]:
    return {"status": "ok"}
