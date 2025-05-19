"""Synchronization logic for tap-incident."""

import datetime
import time
import singer
from typing import Dict, Any, List, Optional, Set

from tap_incident.client import IncidentClient
from tap_incident.streams import STREAMS
from tap_incident.state import write_state

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


def get_stream_replication_method(stream: Dict[str, Any]) -> str:
    """Get replication method from catalog metadata.
    
    Args:
        stream: Stream definition from catalog
        
    Returns:
        Replication method (FULL_TABLE or INCREMENTAL)
    """
    mdata = singer.metadata.to_map(stream["metadata"])
    top_level_md = mdata.get(())
    
    # Check for forced replication method
    forced_method = top_level_md.get("forced-replication-method")
    if forced_method:
        return forced_method
    
    # Check for user-specified method
    return top_level_md.get("replication-method", "FULL_TABLE")


def get_stream_replication_key(stream: Dict[str, Any]) -> Optional[str]:
    """Get replication key from catalog metadata.
    
    Args:
        stream: Stream definition from catalog
        
    Returns:
        Replication key field name
    """
    mdata = singer.metadata.to_map(stream["metadata"])
    top_level_md = mdata.get(())
    
    return top_level_md.get("replication-key")


def sync_stream(
    client: IncidentClient,
    stream_name: str,
    stream_schema: Dict[str, Any],
    selected_fields: Set[str],
    state: Dict[str, Any],
    replication_method: str = "FULL_TABLE",
    replication_key: Optional[str] = None,
) -> None:
    """Sync a single stream.
    
    Args:
        client: IncidentClient instance
        stream_name: Name of the stream to sync
        stream_schema: Schema for the stream
        selected_fields: Set of selected field names
        state: Current state
        replication_method: Replication method
        replication_key: Replication key field
    """
    LOGGER.info(f"Syncing stream: {stream_name}")
    LOGGER.info(f"Replication method: {replication_method}")
    if replication_method == "INCREMENTAL":
        LOGGER.info(f"Replication key: {replication_key}")
    
    extraction_time = datetime.datetime.now(datetime.timezone.utc)
    
    # Write schema
    singer.write_schema(
        stream_name=stream_name,
        schema=stream_schema,
        key_properties=STREAMS[stream_name].key_properties,
    )
    
    # Get stream instance from registry
    stream_object = STREAMS[stream_name]
    
    # Update the stream's replication settings
    if hasattr(stream_object, "replication_method"):
        stream_object.replication_method = replication_method
    
    if hasattr(stream_object, "replication_key") and replication_method == "INCREMENTAL":
        if replication_key and replication_key in stream_object.valid_replication_keys:
            stream_object.replication_key = replication_key
        elif len(stream_object.valid_replication_keys) > 0:
            # Use the first valid key if user didn't specify one
            stream_object.replication_key = stream_object.valid_replication_keys[0]
            LOGGER.info(f"Using default replication key: {stream_object.replication_key}")
    
    # Sync records
    record_count = 0
    for record in stream_object.sync(client, state):
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
    state = state or {}
    
    # Write initial state to establish a baseline for incremental syncs
    write_state(state)
    
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
        replication_method = get_stream_replication_method(stream)
        replication_key = get_stream_replication_key(stream)
        
        sync_stream(
            client,
            stream_name,
            stream["schema"],
            selected_fields,
            state,
            replication_method,
            replication_key,
        )
        
    # Write final state
    write_state(state)
