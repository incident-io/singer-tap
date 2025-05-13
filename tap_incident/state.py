"""State handling utilities for tap-incident."""

import copy
import datetime
import json
import logging
import singer
from typing import Dict, Any, Optional

LOGGER = logging.getLogger(__name__)


def get_stream_state(state: Dict[str, Any], stream_name: str) -> Dict[str, Any]:
    """Get state for a specific stream.
    
    Args:
        state: Current state
        stream_name: Name of the stream
        
    Returns:
        Stream state dictionary
    """
    return state.get("bookmarks", {}).get(stream_name, {})


def write_state(state: Dict[str, Any]) -> None:
    """Write state to standard output.
    
    Args:
        state: Current state
    """
    if state is not None:
        singer.write_state(state)


def update_bookmark(state: Dict[str, Any], stream_name: str, bookmark_key: str, bookmark_value: str) -> Dict[str, Any]:
    """Update bookmark in state.
    
    Args:
        state: Current state
        stream_name: Name of the stream
        bookmark_key: Key for the bookmark (e.g., "updated_at")
        bookmark_value: Value for the bookmark (e.g., timestamp)
        
    Returns:
        Updated state dictionary
    """
    if "bookmarks" not in state:
        state["bookmarks"] = {}
    
    if stream_name not in state["bookmarks"]:
        state["bookmarks"][stream_name] = {}
        
    state["bookmarks"][stream_name][bookmark_key] = bookmark_value
    
    return state


def get_bookmark(state: Dict[str, Any], stream_name: str, bookmark_key: str, default: str = None) -> Optional[str]:
    """Get bookmark from state.
    
    Args:
        state: Current state
        stream_name: Name of the stream
        bookmark_key: Key for the bookmark (e.g., "updated_at")
        default: Default value if bookmark doesn't exist
        
    Returns:
        Bookmark value or default
    """
    return state.get("bookmarks", {}).get(stream_name, {}).get(bookmark_key, default)


def get_bookmark_date(state: Dict[str, Any], stream_name: str, bookmark_key: str) -> Optional[datetime.datetime]:
    """Get bookmark as datetime.
    
    Args:
        state: Current state
        stream_name: Name of the stream
        bookmark_key: Key for the bookmark (e.g., "updated_at")
        
    Returns:
        Bookmark datetime or None
    """
    bookmark = get_bookmark(state, stream_name, bookmark_key)
    if bookmark:
        try:
            return datetime.datetime.fromisoformat(bookmark.replace("Z", "+00:00"))
        except (ValueError, TypeError):
            LOGGER.warning(f"Invalid bookmark date for {stream_name}.{bookmark_key}: {bookmark}")
    
    return None


def reset_stream_state(state: Dict[str, Any], stream_name: str) -> Dict[str, Any]:
    """Reset state for a specific stream.
    
    Args:
        state: Current state
        stream_name: Name of the stream
        
    Returns:
        Updated state with stream reset
    """
    result = copy.deepcopy(state)
    if "bookmarks" in result and stream_name in result["bookmarks"]:
        del result["bookmarks"][stream_name]
    
    return result


def load_state(state_path: str) -> Dict[str, Any]:
    """Load state from a file.
    
    Args:
        state_path: Path to state file
        
    Returns:
        State dictionary
    """
    try:
        with open(state_path, "r") as f:
            return json.load(f)
    except (IOError, json.JSONDecodeError) as e:
        LOGGER.warning(f"Failed to load state file: {e}")
        return {}