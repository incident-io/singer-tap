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


def get_replication_metadata(stream_obj: Any) -> Dict[str, Any]:
    """Get replication metadata for a stream.
    
    Args:
        stream_obj: Stream object
        
    Returns:
        Dictionary with replication metadata
    """
    metadata = {
        "forced-replication-method": stream_obj.replication_method
    }
    
    # Only include replication key for INCREMENTAL method
    if (
        hasattr(stream_obj, "valid_replication_keys") and
        len(stream_obj.valid_replication_keys) > 0
    ):
        metadata["replication-key"] = stream_obj.replication_key or stream_obj.valid_replication_keys[0]
        metadata["valid-replication-keys"] = stream_obj.valid_replication_keys
    
    return metadata


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
        
        # Get replication metadata
        replication_metadata = get_replication_metadata(stream_object)
        
        stream_metadata = [
            {
                "breadcrumb": [],
                "metadata": {
                    "inclusion": "available",
                    "selected-by-default": True,
                    **replication_metadata
                }
            }
        ]
        
        # Add metadata for each property in the schema
        for prop_name in schema["properties"].keys():
            # Determine if this property is one of the primary keys
            is_key = prop_name in stream_object.key_properties
            
            # Determine if this property is the replication key
            is_replication_key = (
                hasattr(stream_object, "replication_key") and 
                stream_object.replication_key == prop_name
            )
            
            # Build property metadata
            prop_metadata = {
                "inclusion": "automatic" if is_key else "available",
                "selected-by-default": True
            }
            
            # If this is a replication key, mark it
            if is_replication_key:
                prop_metadata["is-replication-key"] = True
            
            stream_metadata.append({
                "breadcrumb": ["properties", prop_name],
                "metadata": prop_metadata
            })
            
        catalog_entry = {
            "stream": stream_name,
            "tap_stream_id": stream_name,
            "schema": schema,
            "metadata": stream_metadata
        }
        
        # Add key properties to the catalog entry
        if hasattr(stream_object, "key_properties") and stream_object.key_properties:
            catalog_entry["key_properties"] = stream_object.key_properties
        
        catalog["streams"].append(catalog_entry)
        
    return catalog
