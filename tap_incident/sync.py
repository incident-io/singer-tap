"""Synchronization logic for tap-incident."""

import datetime
import time
import singer
from typing import Dict, Any, List, Optional, Set

from tap_incident.client import IncidentClient
from tap_incident.streams import STREAMS

LOGGER = singer.get_logger()


def get_selected_streams(catalog: Dict[str, Any]) -> Dict[str, Dict[str, Any]]:
    """Get selected streams from catalog.
    
    Args:
        catalog: Singer catalog
        
    Returns:
        Dictionary of selected streams with their schemas
    """
    selected_streams = {}
    
    for stream in catalog.get("streams", []):
        stream_name = stream["tap_stream_id"]
        
        # Get the top-level metadata
        mdata = singer.metadata.to_map(stream["metadata"])
        top_level_md = mdata.get(())
        
        # If the stream is selected or selected by default and not explicitly unselected
        if (
            top_level_md.get("selected", False) or
            (top_level_md.get("selected-by-default", False) and "selected" not in top_level_md)
        ):
            selected_streams[stream_name] = stream
            
    return selected_streams


def get_selected_fields(stream: Dict[str, Any]) -> Set[str]:
    """Get selected fields for a stream.
    
    Args:
        stream: Stream definition from catalog
        
    Returns:
        Set of selected field names
    """
    mdata = singer.metadata.to_map(stream["metadata"])
    properties = stream["schema"]["properties"]
    
    selected_fields = set()
    
    for field_name in properties:
        field_md = mdata.get(("properties", field_name))
        
        # Field is selected if explicitly selected or selected by default and not explicitly unselected
        if (
            field_md.get("selected", False) or
            (field_md.get("selected-by-default", False) and "selected" not in field_md)
        ):
            selected_fields.add(field_name)
            
    return selected_fields


def sync_stream(
    client: IncidentClient,
    stream_name: str,
    stream_schema: Dict[str, Any],
    selected_fields: Set[str],
) -> None:
    """Sync a single stream.
    
    Args:
        client: IncidentClient instance
        stream_name: Name of the stream to sync
        stream_schema: Schema for the stream
        selected_fields: Set of selected field names
    """
    LOGGER.info(f"Syncing stream: {stream_name}")
    
    extraction_time = datetime.datetime.now(datetime.timezone.utc)
    
    # Write schema
    singer.write_schema(
        stream_name=stream_name,
        schema=stream_schema,
        key_properties=STREAMS[stream_name].key_properties,
    )
    
    # Get stream instance from registry
    stream_object = STREAMS[stream_name]
    
    # Sync records
    record_count = 0
    for record in stream_object.sync(client):
        # Filter record to only include selected fields
        filtered_record = {k: v for k, v in record.items() if k in selected_fields}
        
        # Write record with extraction time
        singer.write_record(
            stream_name=stream_name,
            record=filtered_record,
            time_extracted=extraction_time,
        )
        record_count += 1
        
    LOGGER.info(f"Synced {record_count} records from {stream_name}")


def sync(
    client: IncidentClient,
    catalog: Optional[Dict[str, Any]] = None,
    state: Optional[Dict[str, Any]] = None,
) -> None:
    """Sync all selected streams.
    
    Args:
        client: IncidentClient instance
        catalog: Singer catalog (optional)
        state: Singer state (optional)
    """
    if catalog:
        # Use provided catalog
        selected_streams = get_selected_streams(catalog)
    else:
        # If no catalog is provided, sync all streams
        from tap_incident.discover import discover
        catalog = discover(client)
        selected_streams = {
            stream["tap_stream_id"]: stream for stream in catalog["streams"]
        }
    
    for stream_name, stream in selected_streams.items():
        if stream_name not in STREAMS:
            LOGGER.warning(f"Stream '{stream_name}' not found in available streams, skipping")
            continue
            
        selected_fields = get_selected_fields(stream)
        sync_stream(
            client,
            stream_name,
            stream["schema"],
            selected_fields,
        )
