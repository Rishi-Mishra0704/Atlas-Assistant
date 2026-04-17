from __future__ import annotations

import uvicorn

from server.server import create_app

if __name__ == "__main__":
    uvicorn.run(create_app("tts"), host="127.0.0.1", port=8012)
