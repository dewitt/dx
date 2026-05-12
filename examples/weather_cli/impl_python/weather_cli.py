#!/usr/bin/env python3
"""Reference Python implementation of weather-cli.

This file is the *agentic synthesis* artifact in the worked example: an
implementer agent reads only ../system.dx (never ../impl_cpp/) and
produces this file. The two implementations should be observably
indistinguishable from the perspective of the contracts in system.dx.

Network calls are mocked so the example runs offline. The contract
surface (env-var handling, cache TTL, exit codes, stdout/stderr split,
JSON output mode) is real.
"""

from __future__ import annotations

import argparse
import json
import os
import sys
import time
from pathlib import Path
from typing import Any

# Pinned by assumptions.cache.location in system.dx.
CACHE_FILE = Path.home() / ".weather_cache.json"

# Pinned by assumptions.network.provider in system.dx. Used only as a
# label in the JSON output; no real network call is issued.
UPSTREAM_PROVIDER = "OpenMeteo"

# Encoded directly in the perf_no_redundant_fetch invariant.
CACHE_TTL_SECONDS = 600


def is_cache_valid(path: Path = CACHE_FILE) -> bool:
    """Return True if the cache file exists and is younger than the TTL."""
    if not path.exists():
        return False
    age = time.time() - path.stat().st_mtime
    return age < CACHE_TTL_SECONDS


def fetch_weather(zip_code: str, _api_key: str) -> dict[str, Any]:
    """Simulated upstream fetch.

    A real implementation would issue an HTTP GET to the upstream
    provider here. The shape of the returned dict is what flows out
    through `--json`, so it is shaped to satisfy the
    `emits_json_with_flag` contract.
    """
    return {
        "zip": zip_code,
        "temp": "72F",
        "condition": "Sunny",
        "provider": UPSTREAM_PROVIDER,
    }


def main(argv: list[str] | None = None) -> int:
    parser = argparse.ArgumentParser(
        prog="weather_cli",
        description="Fetch and display current weather for a US zip code.",
    )
    parser.add_argument("zipcode", help="US zip code")
    parser.add_argument(
        "--json",
        dest="as_json",
        action="store_true",
        help="emit a single JSON object on stdout instead of a human summary",
    )
    args = parser.parse_args(argv)

    api_key = os.environ.get("WEATHER_API_KEY")
    if not api_key:
        print(
            "Error: WEATHER_API_KEY environment variable not set.",
            file=sys.stderr,
        )
        return 1

    weather: dict[str, Any] | None = None
    if is_cache_valid():
        try:
            with CACHE_FILE.open("r", encoding="utf-8") as f:
                weather = json.load(f)
        except (OSError, json.JSONDecodeError):
            weather = None  # Treat a corrupt cache as a miss.

    if weather is None:
        weather = fetch_weather(args.zipcode, api_key)
        try:
            with CACHE_FILE.open("w", encoding="utf-8") as f:
                json.dump(weather, f)
        except OSError:
            # Cache write failures are non-fatal: the user still gets
            # a fresh result this run.
            pass

    if args.as_json:
        print(json.dumps(weather))
    else:
        print(
            f"Weather for {weather['zip']}: {weather['temp']}, {weather['condition']}"
        )
    return 0


if __name__ == "__main__":
    sys.exit(main())
