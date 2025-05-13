#!/usr/bin/env python3
"""Command line interface for tap-incident."""

import argparse
import json
import logging
import sys
from typing import Dict, Any

from tap_incident import __version__
from tap_incident.client import IncidentClient
from tap_incident.discover import discover as discover_streams
from tap_incident.sync import sync as sync_streams
from tap_incident.state import load_state

LOGGER = logging.getLogger(__name__)
REQUIRED_CONFIG_KEYS = ["api_key"]


def load_json(path: str) -> Dict[str, Any]:
    """Load a JSON file from disk."""
    try:
        with open(path) as f:
            return json.load(f)
    except (IOError, json.JSONDecodeError) as e:
        LOGGER.error(f"Failed to load JSON file {path}: {e}")
        sys.exit(1)


def parse_args() -> argparse.Namespace:
    """Parse command line arguments."""
    parser = argparse.ArgumentParser(description="Singer.io tap for extracting data from incident.io API")
    parser.add_argument("-c", "--config", required=True, help="Config file")
    parser.add_argument("-s", "--state", help="State file")
    parser.add_argument("-p", "--properties", help="Property selections (deprecated)")
    parser.add_argument("--catalog", help="Catalog file")
    parser.add_argument("-d", "--discover", action="store_true", help="Run in discovery mode")
    parser.add_argument("-v", "--version", action="version", version=f"tap-incident {__version__}")
    parser.add_argument("--debug", action="store_true", default=False, help="Enable debug logging")
    parser.add_argument(
        "--full-refresh", 
        action="store_true", 
        default=False, 
        help="Ignore state and perform a full refresh"
    )

    return parser.parse_args()


def main():
    """Main entry point."""
    args = parse_args()

    # Setup logging
    log_level = logging.DEBUG if args.debug else logging.INFO
    logging.basicConfig(level=log_level, format="%(asctime)s - %(name)s - %(levelname)s - %(message)s")

    # Parse config
    try:
        config = load_json(args.config)
    except Exception as e:
        LOGGER.error(f"Failed to load config: {e}")
        LOGGER.error(
            """
            We expect a config file in JSON format that looks like:
            {
              "api_key": "<your-api-key>",
              "endpoint": "https://api.incident.io"
            }
            """
        )
        sys.exit(1)

    # Check for required config keys
    missing_keys = [key for key in REQUIRED_CONFIG_KEYS if key not in config]
    if missing_keys:
        LOGGER.error(f"Missing required config keys: {', '.join(missing_keys)}")
        sys.exit(1)

    # If no endpoint is specified, use the default
    if "endpoint" not in config:
        config["endpoint"] = "https://api.incident.io"

    # Set up API client
    client = IncidentClient(config["api_key"], config["endpoint"])

    # Load state if provided and not doing a full refresh
    state = {}
    if args.state and not args.full_refresh:
        LOGGER.info(f"Loading state from {args.state}")
        state = load_state(args.state)
    elif args.full_refresh and args.state:
        LOGGER.info("Full refresh requested, ignoring state file")

    # Catalog handling - either from args or discovery
    catalog = None
    if args.catalog:
        LOGGER.info(f"Loading catalog from {args.catalog}")
        catalog = load_json(args.catalog)
    
    # Deprecated property selection via --properties
    if args.properties:
        LOGGER.warning(
            "The --properties parameter is deprecated. "
            "Please use --catalog instead."
        )

    # Discovery mode
    if args.discover:
        LOGGER.info("Running in discovery mode")
        catalog = discover_streams(client)
        json.dump(catalog, sys.stdout, indent=2)
        return

    # Sync
    LOGGER.info("Starting sync")
    sync_streams(client, catalog, state)
    LOGGER.info("Sync completed")


if __name__ == "__main__":
    main()
