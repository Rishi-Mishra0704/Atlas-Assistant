from __future__ import annotations

from fastapi import APIRouter

router = APIRouter(prefix="/rerank", tags=["rerank"])


@router.get("/health")
def health() -> dict[str, str]:
    return {"status": "ok"}
