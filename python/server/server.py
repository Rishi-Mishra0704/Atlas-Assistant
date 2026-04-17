from __future__ import annotations

from fastapi import FastAPI

from middleware.logging import logging_middleware
from server import embed, rag, rerank, stt, tts, wake


ROUTERS = {
    "wake": wake.router,
    "stt": stt.router,
    "tts": tts.router,
    "embed": embed.router,
    "rerank": rerank.router,
    "rag": rag.router,
}


def create_app(sidecar: str = "rag") -> FastAPI:
    app = FastAPI(title=f"Atlas {sidecar} sidecar")
    app.middleware("http")(logging_middleware)
    app.include_router(ROUTERS.get(sidecar, rag.router))
    return app
