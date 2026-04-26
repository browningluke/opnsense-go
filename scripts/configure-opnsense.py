#!/usr/bin/env python3
"""
Apply OPNsense pre-test configuration defined in opnsense-setup.json.

This script is called by CI before running acceptance tests. It reads a list
of API calls from scripts/opnsense-setup.json and executes them in order.

To add or remove pre-conditions, edit opnsense-setup.json — no Python
changes required.

Reads credentials from environment variables:
  OPNSENSE_URI        - Base URL, e.g. https://localhost:8443
  OPNSENSE_API_KEY    - API key
  OPNSENSE_API_SECRET - API secret
"""

import json
import os
import pathlib
import ssl
import sys
import urllib.request
from base64 import b64encode

SETUP_FILE = pathlib.Path(__file__).parent / "opnsense-setup.json"


def api_post(base_url, key, secret, path, body):
    url = base_url.rstrip("/") + path
    credentials = b64encode(f"{key}:{secret}".encode()).decode()
    headers = {
        "Authorization": f"Basic {credentials}",
        "Content-Type": "application/json",
    }
    # Allow self-signed certificates (mirrors AllowInsecure=true in tests)
    ctx = ssl.create_default_context()
    ctx.check_hostname = False
    ctx.verify_mode = ssl.CERT_NONE

    req = urllib.request.Request(
        url, data=json.dumps(body).encode(), headers=headers, method="POST"
    )
    try:
        with urllib.request.urlopen(req, context=ctx, timeout=30) as resp:
            return json.loads(resp.read().decode())
    except urllib.error.HTTPError as e:
        raise RuntimeError(f"HTTP {e.code}: {e.read().decode()}") from e


def main():
    base_url = os.environ.get("OPNSENSE_URI", "").rstrip("/")
    key = os.environ.get("OPNSENSE_API_KEY", "")
    secret = os.environ.get("OPNSENSE_API_SECRET", "")

    if not base_url or not key or not secret:
        print(
            "ERROR: OPNSENSE_URI, OPNSENSE_API_KEY, and OPNSENSE_API_SECRET must be set",
            file=sys.stderr,
        )
        sys.exit(1)

    steps = json.loads(SETUP_FILE.read_text())
    print(f"Applying {len(steps)} setup step(s) to {base_url} ...")

    errors = []
    for step in steps:
        desc = step["description"]
        try:
            result = api_post(base_url, key, secret, step["endpoint"], step["body"])
            if result.get("result") not in ("saved", "ok"):
                raise RuntimeError(f"unexpected result: {result}")
            print(f"  [OK] {desc}")
        except Exception as e:
            print(f"  [FAIL] {desc}: {e}", file=sys.stderr)
            errors.append(desc)

    if errors:
        print(f"\n{len(errors)} step(s) failed.", file=sys.stderr)
        sys.exit(1)

    print("Setup complete.")


if __name__ == "__main__":
    main()
