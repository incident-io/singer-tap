"""Discovery mode for tap-incident."""

import json
import os
import singer
from typing import Dict, Any, List

from tap_incident.client import IncidentClient
from tap_incident.streams import STREAMS

LOGGER = singer.get_logger()


def get_abs_path(path: str) -> str:
    """Get absolute path to file."""
    return os.path.join(os.path.dirname(os.path.realpath(__file__)), path)


def load_schema(stream_name: str) -> Dict[str, Any]:
    """Load stream schema from schemas directory."""
    path = get_abs_path(f"schemas/{stream_name}.json")
    with open(path, "r") as f:
        schema = json.load(f)
    return schema


def discover(client: IncidentClient) -> Dict[str, Any]:
    """Discover available streams.
    
    Args:
        client: IncidentClient instance
        
    Returns:
        Catalog dictionary for Singer
    """
    catalog = {"streams": []}
    
    for stream_name, stream_object in STREAMS.items():
        LOGGER.info(f"Discovering schema for {stream_name}")
        
        schema = load_schema(stream_name)
        
        stream_metadata = [
            {
                "breadcrumb": [],
                "metadata": {
                    "inclusion": "available",
                    "selected-by-default": True,
                    "forced-replication-method": "FULL_TABLE"
                }
            }
        ]
        
        # Add metadata for each property in the schema
        for prop_name in schema["properties"].keys():
            stream_metadata.append({
                "breadcrumb": ["properties", prop_name],
                "metadata": {
                    "inclusion": "available",
                    "selected-by-default": True,
                }
            })
            
        catalog_entry = {
            "stream": stream_name,
            "tap_stream_id": stream_name,
            "schema": schema,
            "metadata": stream_metadata
        }
        
        catalog["streams"].append(catalog_entry)
        
    return catalog