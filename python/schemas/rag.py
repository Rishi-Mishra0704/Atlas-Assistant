from __future__ import annotations

from pydantic import BaseModel


class RetrievalRequest(BaseModel):
    query: str
    limit: int = 10


class RetrievalResponse(BaseModel):
    results: list[str]
