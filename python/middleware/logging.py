from __future__ import annotations

import time

from fastapi import Request


async def logging_middleware(request: Request, call_next):
    started = time.perf_counter()
    response = await call_next(request)
    response.headers["x-process-time"] = f"{time.perf_counter() - started:.4f}"
    return response
