def parse_api_response(resp: dict) -> dict:
    if not resp.get("success"):
        raise ValueError(f"API Error: {resp}")

    payload = resp.get("data")
    if payload is None:
        raise ValueError(f"Missing data: {resp}")

    return payload