from __future__ import annotations

from pydantic import BaseModel


class VoiceRequest(BaseModel):
    text: str


class VoiceResponse(BaseModel):
    text: str
