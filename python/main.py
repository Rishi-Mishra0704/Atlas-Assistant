from __future__ import annotations

import argparse

import uvicorn

from server.server import create_app


def main() -> None:
    parser = argparse.ArgumentParser(description="Run an Atlas Python sidecar.")
    parser.add_argument("sidecar", nargs="?", default="rag")
    parser.add_argument("--host", default="127.0.0.1")
    parser.add_argument("--port", type=int, default=8000)
    args = parser.parse_args()

    uvicorn.run(create_app(args.sidecar), host=args.host, port=args.port)


if __name__ == "__main__":
    main()
