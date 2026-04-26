#!/usr/bin/env python3
"""
Pre-configure OPNsense state required for acceptance tests.

This script is called by CI before running tests to ensure the OPNsense
instance has the required configuration (e.g. interfaces registered in
service modules) that tests depend on.

Reads credentials from environment variables:
  OPNSENSE_URI          - Base URL, e.g. https://localhost:8443
  OPNSENSE_API_KEY      - API key
  OPNSENSE_API_SECRET   - API secret
"""

import json
import os
import ssl
import sys
import urllib.request
from base64 import b64encode


def api_post(base_url, key, secret, path, body):
    url = base_url.rstrip("/") + path
    credentials = b64encode(f"{key}:{secret}".encode()).decode()
    headers = {
        "Authorization": f"Basic {credentials}",
        "Content-Type": "application/json",
    }
    data = json.dumps(body).encode("utf-8")

    # Allow self-signed certificates (mirrors AllowInsecure=true in tests)
    ctx = ssl.create_default_context()
    ctx.check_hostname = False
    ctx.verify_mode = ssl.CERT_NONE

    req = urllib.request.Request(url, data=data, headers=headers, method="POST")
    try:
        with urllib.request.urlopen(req, context=ctx, timeout=30) as resp:
            return json.loads(resp.read().decode("utf-8"))
    except urllib.error.HTTPError as e:
        body = e.read().decode("utf-8")
        raise RuntimeError(f"HTTP {e.code} from {url}: {body}") from e


def configure_kea_dhcpv4(base_url, key, secret):
    """Enable the wan interface in Kea DHCPv4 general settings."""
    result = api_post(
        base_url, key, secret,
        "/api/kea/dhcpv4/set",
        {"dhcpv4": {"general": {"enabled": "0", "interfaces": "wan", "valid_lifetime": "4000", "fwrules": "1"}}},
    )
    if result.get("result") != "saved":
        raise RuntimeError(f"Failed to configure Kea DHCPv4 general settings: {result}")
    print("  [OK] Kea DHCPv4: wan interface enabled in general settings")


def configure_kea_dhcpv6(base_url, key, secret):
    """Enable the wan interface in Kea DHCPv6 general settings."""
    result = api_post(
        base_url, key, secret,
        "/api/kea/dhcpv6/set",
        {"dhcpv6": {"general": {"enabled": "0", "interfaces": "wan", "valid_lifetime": "4000", "fwrules": "1"}}},
    )
    if result.get("result") != "saved":
        raise RuntimeError(f"Failed to configure Kea DHCPv6 general settings: {result}")
    print("  [OK] Kea DHCPv6: wan interface enabled in general settings")


def main():
    base_url = os.environ.get("OPNSENSE_URI", "").rstrip("/")
    key = os.environ.get("OPNSENSE_API_KEY", "")
    secret = os.environ.get("OPNSENSE_API_SECRET", "")

    if not base_url or not key or not secret:
        print("ERROR: OPNSENSE_URI, OPNSENSE_API_KEY, and OPNSENSE_API_SECRET must be set", file=sys.stderr)
        sys.exit(1)

    print(f"Configuring OPNsense at {base_url} ...")

    errors = []

    for name, fn in [
        ("Kea DHCPv4", configure_kea_dhcpv4),
        ("Kea DHCPv6", configure_kea_dhcpv6),
    ]:
        try:
            fn(base_url, key, secret)
        except Exception as e:
            print(f"  [FAIL] {name}: {e}", file=sys.stderr)
            errors.append(name)

    if errors:
        print(f"\nConfiguration failed for: {', '.join(errors)}", file=sys.stderr)
        sys.exit(1)

    print("OPNsense configuration complete.")


if __name__ == "__main__":
    main()
