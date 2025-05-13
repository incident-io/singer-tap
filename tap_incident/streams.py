"""Stream definitions for tap-incident."""

import logging
from abc import ABC, abstractmethod
from datetime import datetime
from typing import Dict, Any, List, Iterator, Optional, Set

from tap_incident.client import IncidentClient
from tap_incident.state import get_bookmark_date, update_bookmark, write_state

LOGGER = logging.getLogger(__name__)


class Stream(ABC):
    """Base stream class."""
    
    # Class variables that should be defined by each stream
    name = None  # Stream name
    key_properties = None  # Primary key fields
    replication_method = "FULL_TABLE"  # Default to full table replication
    replication_key = None  # Field to use for bookmarking/incremental (if supported)
    valid_replication_keys = []  # List of fields that could be used for incremental
    
    @abstractmethod
    def sync(self, client: IncidentClient, state: Dict[str, Any] = None) -> Iterator[Dict[str, Any]]:
        """Sync data from the stream.
        
        Args:
            client: IncidentClient instance
            state: Current state
            
        Yields:
            Records from the stream
        """
        pass
    
    def get_starting_time(self, state: Dict[str, Any]) -> Optional[datetime]:
        """Get the starting timestamp for incremental replication.
        
        Args:
            state: Current state
            
        Returns:
            Timestamp to start from, or None for all data
        """
        if not self.replication_key:
            return None
            
        return get_bookmark_date(state, self.name, self.replication_key)
    
    def update_state(self, state: Dict[str, Any], record: Dict[str, Any]) -> Dict[str, Any]:
        """Update state with record's replication key value.
        
        Args:
            state: Current state
            record: Current record
            
        Returns:
            Updated state
        """
        if not self.replication_key or self.replication_key not in record:
            return state
            
        current_bookmark = self.get_starting_time(state)
        record_time = datetime.fromisoformat(record[self.replication_key].replace("Z", "+00:00"))
        
        if current_bookmark is None or record_time > current_bookmark:
            return update_bookmark(state, self.name, self.replication_key, record[self.replication_key])
            
        return state
    
    def update_state_for_full_table(self, state: Dict[str, Any], records: List[Dict[str, Any]]) -> Dict[str, Any]:
        """Update state for FULL_TABLE replication mode.
        
        This processes all records at once after they've been synced to find the latest
        value of the replication key for a future incremental sync.
        
        Args:
            state: Current state
            records: All synced records
            
        Returns:
            Updated state
        """
        if not self.replication_key or not self.valid_replication_keys:
            return state
            
        # Find the record with the latest value of the replication key
        latest_record = None
        latest_time = None
        
        for record in records:
            if self.replication_key not in record:
                continue
                
            try:
                record_time = datetime.fromisoformat(record[self.replication_key].replace("Z", "+00:00"))
                if latest_time is None or record_time > latest_time:
                    latest_time = record_time
                    latest_record = record
            except (ValueError, TypeError):
                LOGGER.warning(f"Invalid date format for {self.replication_key} in record: {record.get('id', 'unknown')}")
                continue
        
        # Update the state with the latest record
        if latest_record:
            state = update_bookmark(state, self.name, self.replication_key, latest_record[self.replication_key])
        
        return state


class ActionsStream(Stream):
    """Actions stream."""
    
    name = "actions"
    key_properties = ["id", "incident_id"]
    valid_replication_keys = ["updated_at"]
    
    def __init__(self, replication_method=None, replication_key=None):
        self.replication_method = replication_method or "FULL_TABLE"
        self.replication_key = replication_key or "updated_at" 
    
    def sync(self, client: IncidentClient, state: Dict[str, Any] = None) -> Iterator[Dict[str, Any]]:
        """Sync actions."""
        state = state or {}
        starting_date = self.get_starting_time(state)
        
        actions = client.get_actions()
        record_count = 0
        all_records = []  # Keep track of all records for FULL_TABLE
        
        for action in actions:
            all_records.append(action)
            
            # If we're doing incremental and have a bookmark date, filter records
            if self.replication_method == "INCREMENTAL" and starting_date and self.replication_key:
                record_date = datetime.fromisoformat(action[self.replication_key].replace("Z", "+00:00"))
                if record_date <= starting_date:
                    continue
            
            record_count += 1
            
            # For incremental, update state record by record
            if self.replication_method == "INCREMENTAL":
                state = self.update_state(state, action)
            
                # Every 100 records, write the state
                if record_count % 100 == 0:
                    write_state(state)
                
            yield action
        
        # For FULL_TABLE, update state once at the end with the latest record
        if self.replication_method == "FULL_TABLE":
            state = self.update_state_for_full_table(state, all_records)
        
        # Write the state one final time
        write_state(state)


class AlertsStream(Stream):
    """Alerts stream."""
    
    name = "alerts"
    key_properties = ["id"]
    valid_replication_keys = ["created_at"]
    
    def __init__(self, replication_method=None, replication_key=None):
        self.replication_method = replication_method or "FULL_TABLE"
        self.replication_key = replication_key or "created_at"
    
    def sync(self, client: IncidentClient, state: Dict[str, Any] = None) -> Iterator[Dict[str, Any]]:
        """Sync alerts."""
        state = state or {}
        starting_date = self.get_starting_time(state)
        
        alerts = client.get_alerts()
        record_count = 0
        all_records = []  # Keep track of all records for FULL_TABLE
        
        for alert in alerts:
            all_records.append(alert)
            
            # Filter by bookmark date for incremental replication
            if self.replication_method == "INCREMENTAL" and starting_date and self.replication_key:
                record_date = datetime.fromisoformat(alert[self.replication_key].replace("Z", "+00:00"))
                if record_date <= starting_date:
                    continue
            
            record_count += 1
            
            # For incremental, update state record by record
            if self.replication_method == "INCREMENTAL":
                state = self.update_state(state, alert)
            
                # Every 100 records, write the state
                if record_count % 100 == 0:
                    write_state(state)
                
            yield alert
        
        # For FULL_TABLE, update state once at the end with the latest record
        if self.replication_method == "FULL_TABLE":
            state = self.update_state_for_full_table(state, all_records)
        
        # Write the state one final time
        write_state(state)


class AlertAttributesStream(Stream):
    """Alert attributes stream."""
    
    name = "alert_attributes"
    key_properties = ["id"]
    
    def sync(self, client: IncidentClient, state: Dict[str, Any] = None) -> Iterator[Dict[str, Any]]:
        """Sync alert attributes."""
        for attribute in client.get_alert_attributes():
            yield attribute


class AlertSourcesStream(Stream):
    """Alert sources stream."""
    
    name = "alert_sources"
    key_properties = ["id"]
    
    def sync(self, client: IncidentClient, state: Dict[str, Any] = None) -> Iterator[Dict[str, Any]]:
        """Sync alert sources."""
        for source in client.get_alert_sources():
            yield source


class CustomFieldsStream(Stream):
    """Custom fields stream."""
    
    name = "custom_fields"
    key_properties = ["id"]
    valid_replication_keys = ["updated_at"]
    
    def __init__(self, replication_method=None, replication_key=None):
        self.replication_method = replication_method or "FULL_TABLE"
        self.replication_key = replication_key or "updated_at"
    
    def sync(self, client: IncidentClient, state: Dict[str, Any] = None) -> Iterator[Dict[str, Any]]:
        """Sync custom fields."""
        state = state or {}
        starting_date = self.get_starting_time(state)
        
        custom_fields = client.get_custom_fields()
        record_count = 0
        all_records = []
        
        for custom_field in custom_fields:
            all_records.append(custom_field)
            
            # Filter by bookmark date for incremental replication
            if self.replication_method == "INCREMENTAL" and starting_date and self.replication_key:
                record_date = datetime.fromisoformat(custom_field[self.replication_key].replace("Z", "+00:00"))
                if record_date <= starting_date:
                    continue
            
            record_count += 1
            
            # For incremental, update state record by record
            if self.replication_method == "INCREMENTAL":
                state = self.update_state(state, custom_field)
            
                # Every 100 records, write the state
                if record_count % 100 == 0:
                    write_state(state)
                
            yield custom_field
        
        # For FULL_TABLE, update state once at the end with the latest record
        if self.replication_method == "FULL_TABLE":
            state = self.update_state_for_full_table(state, all_records)
            
        # Write the state one final time
        write_state(state)


class CustomFieldOptionsStream(Stream):
    """Custom field options stream."""
    
    name = "custom_field_options"
    key_properties = ["id", "custom_field_id"]
    
    def sync(self, client: IncidentClient, state: Dict[str, Any] = None) -> Iterator[Dict[str, Any]]:
        """Sync custom field options."""
        custom_fields = client.get_custom_fields()
        
        for custom_field in custom_fields:
            custom_field_id = custom_field["id"]
            
            LOGGER.info(f"Syncing options for custom field {custom_field_id}")
            
            options = client.get_custom_field_options(custom_field_id)
            for option in options:
                yield option


class FollowUpsStream(Stream):
    """Follow ups stream."""
    
    name = "follow_ups"
    key_properties = ["id", "incident_id"]
    valid_replication_keys = ["updated_at"]
    
    def __init__(self, replication_method=None, replication_key=None):
        self.replication_method = replication_method or "FULL_TABLE"
        self.replication_key = replication_key or "updated_at"
    
    def sync(self, client: IncidentClient, state: Dict[str, Any] = None) -> Iterator[Dict[str, Any]]:
        """Sync follow ups."""
        state = state or {}
        starting_date = self.get_starting_time(state)
        
        follow_ups = client.get_follow_ups()
        record_count = 0
        all_records = []
        
        for follow_up in follow_ups:
            all_records.append(follow_up)
            
            # Filter by bookmark date for incremental replication
            if self.replication_method == "INCREMENTAL" and starting_date and self.replication_key:
                record_date = datetime.fromisoformat(follow_up[self.replication_key].replace("Z", "+00:00"))
                if record_date <= starting_date:
                    continue
            
            record_count += 1
            
            # For incremental, update state record by record
            if self.replication_method == "INCREMENTAL":
                state = self.update_state(state, follow_up)
            
                # Every 100 records, write the state
                if record_count % 100 == 0:
                    write_state(state)
                
            yield follow_up
        
        # For FULL_TABLE, update state once at the end with the latest record
        if self.replication_method == "FULL_TABLE":
            state = self.update_state_for_full_table(state, all_records)
            
        # Write the state one final time
        write_state(state)


class IncidentRolesStream(Stream):
    """Incident roles stream."""
    
    name = "incident_roles"
    key_properties = ["id"]
    valid_replication_keys = ["updated_at"]
    
    def __init__(self, replication_method=None, replication_key=None):
        self.replication_method = replication_method or "FULL_TABLE"
        self.replication_key = replication_key or "updated_at"
    
    def sync(self, client: IncidentClient, state: Dict[str, Any] = None) -> Iterator[Dict[str, Any]]:
        """Sync incident roles."""
        state = state or {}
        starting_date = self.get_starting_time(state)
        
        roles = client.get_incident_roles()
        record_count = 0
        all_records = []
        
        for role in roles:
            all_records.append(role)
            
            # Filter by bookmark date for incremental replication
            if self.replication_method == "INCREMENTAL" and starting_date and self.replication_key:
                record_date = datetime.fromisoformat(role[self.replication_key].replace("Z", "+00:00"))
                if record_date <= starting_date:
                    continue
            
            record_count += 1
            
            # For incremental, update state record by record
            if self.replication_method == "INCREMENTAL":
                state = self.update_state(state, role)
            
                # Every 100 records, write the state
                if record_count % 100 == 0:
                    write_state(state)
                
            yield role
        
        # For FULL_TABLE, update state once at the end with the latest record
        if self.replication_method == "FULL_TABLE":
            state = self.update_state_for_full_table(state, all_records)
            
        # Write the state one final time
        write_state(state)


class IncidentStatusesStream(Stream):
    """Incident statuses stream."""
    
    name = "incident_statuses"
    key_properties = ["id"]
    valid_replication_keys = ["updated_at"]
    
    def __init__(self, replication_method=None, replication_key=None):
        self.replication_method = replication_method or "FULL_TABLE"
        self.replication_key = replication_key or "updated_at"
    
    def sync(self, client: IncidentClient, state: Dict[str, Any] = None) -> Iterator[Dict[str, Any]]:
        """Sync incident statuses."""
        state = state or {}
        starting_date = self.get_starting_time(state)
        
        statuses = client.get_incident_statuses()
        record_count = 0
        all_records = []
        
        for status in statuses:
            all_records.append(status)
            
            # Filter by bookmark date for incremental replication
            if self.replication_method == "INCREMENTAL" and starting_date and self.replication_key:
                record_date = datetime.fromisoformat(status[self.replication_key].replace("Z", "+00:00"))
                if record_date <= starting_date:
                    continue
            
            record_count += 1
            
            # For incremental, update state record by record
            if self.replication_method == "INCREMENTAL":
                state = self.update_state(state, status)
            
                # Every 100 records, write the state
                if record_count % 100 == 0:
                    write_state(state)
                
            yield status
        
        # For FULL_TABLE, update state once at the end with the latest record
        if self.replication_method == "FULL_TABLE":
            state = self.update_state_for_full_table(state, all_records)
            
        # Write the state one final time
        write_state(state)


class IncidentTimestampsStream(Stream):
    """Incident timestamps stream."""
    
    name = "incident_timestamps"
    key_properties = ["id"]
    
    def sync(self, client: IncidentClient, state: Dict[str, Any] = None) -> Iterator[Dict[str, Any]]:
        """Sync incident timestamps."""
        for timestamp in client.get_incident_timestamps():
            yield timestamp


class IncidentTypesStream(Stream):
    """Incident types stream."""
    
    name = "incident_types"
    key_properties = ["id"]
    valid_replication_keys = ["updated_at"]
    
    def __init__(self, replication_method=None, replication_key=None):
        self.replication_method = replication_method or "FULL_TABLE"
        self.replication_key = replication_key or "updated_at"
    
    def sync(self, client: IncidentClient, state: Dict[str, Any] = None) -> Iterator[Dict[str, Any]]:
        """Sync incident types."""
        state = state or {}
        starting_date = self.get_starting_time(state)
        
        incident_types = client.get_incident_types()
        record_count = 0
        all_records = []
        
        for incident_type in incident_types:
            all_records.append(incident_type)
            
            # Filter by bookmark date for incremental replication
            if self.replication_method == "INCREMENTAL" and starting_date and self.replication_key:
                record_date = datetime.fromisoformat(incident_type[self.replication_key].replace("Z", "+00:00"))
                if record_date <= starting_date:
                    continue
            
            record_count += 1
            
            # For incremental, update state record by record
            if self.replication_method == "INCREMENTAL":
                state = self.update_state(state, incident_type)
            
                # Every 100 records, write the state
                if record_count % 100 == 0:
                    write_state(state)
                
            yield incident_type
        
        # For FULL_TABLE, update state once at the end with the latest record
        if self.replication_method == "FULL_TABLE":
            state = self.update_state_for_full_table(state, all_records)
            
        # Write the state one final time
        write_state(state)


class IncidentUpdatesStream(Stream):
    """Incident updates stream."""
    
    name = "incident_updates"
    key_properties = ["id", "incident_id"]
    valid_replication_keys = ["created_at"]
    
    def __init__(self, replication_method=None, replication_key=None):
        self.replication_method = replication_method or "FULL_TABLE"
        self.replication_key = replication_key or "created_at"
    
    def sync(self, client: IncidentClient, state: Dict[str, Any] = None) -> Iterator[Dict[str, Any]]:
        """Sync incident updates."""
        state = state or {}
        starting_date = self.get_starting_time(state)
        
        updates = client.get_incident_updates()
        record_count = 0
        all_records = []
        
        for update in updates:
            all_records.append(update)
            
            # Filter by bookmark date for incremental replication
            if self.replication_method == "INCREMENTAL" and starting_date and self.replication_key:
                record_date = datetime.fromisoformat(update[self.replication_key].replace("Z", "+00:00"))
                if record_date <= starting_date:
                    continue
            
            record_count += 1
            
            # For incremental, update state record by record
            if self.replication_method == "INCREMENTAL":
                state = self.update_state(state, update)
            
                # Every 100 records, write the state
                if record_count % 100 == 0:
                    write_state(state)
                
            yield update
        
        # For FULL_TABLE, update state once at the end with the latest record
        if self.replication_method == "FULL_TABLE":
            state = self.update_state_for_full_table(state, all_records)
            
        # Write the state one final time
        write_state(state)


class IncidentsStream(Stream):
    """Incidents stream."""
    
    name = "incidents"
    key_properties = ["id"]
    valid_replication_keys = ["updated_at"]
    
    def __init__(self, replication_method=None, replication_key=None):
        self.replication_method = replication_method or "FULL_TABLE"
        self.replication_key = replication_key or "updated_at"
    
    def sync(self, client: IncidentClient, state: Dict[str, Any] = None) -> Iterator[Dict[str, Any]]:
        """Sync incidents."""
        state = state or {}
        starting_date = self.get_starting_time(state)
        
        incidents = client.get_incidents()
        record_count = 0
        all_records = []
        
        for incident in incidents:
            all_records.append(incident)
            
            # Filter by bookmark date for incremental replication
            if self.replication_method == "INCREMENTAL" and starting_date and self.replication_key:
                record_date = datetime.fromisoformat(incident[self.replication_key].replace("Z", "+00:00"))
                if record_date <= starting_date:
                    continue
            
            incident_id = incident["id"]
            
            # Fetch related data for each incident
            attachments = client.get_incident_attachments(incident_id)
            incident["attachments"] = attachments
            
            updates = client.get_incident_updates(incident_id)
            incident["updates"] = updates
            
            record_count += 1
            
            # For incremental, update state record by record
            if self.replication_method == "INCREMENTAL":
                state = self.update_state(state, incident)
            
                # Every 100 records, write the state
                if record_count % 100 == 0:
                    write_state(state)
                
            yield incident
        
        # For FULL_TABLE, update state once at the end with the latest record
        if self.replication_method == "FULL_TABLE":
            state = self.update_state_for_full_table(state, all_records)
            
        # Write the state one final time
        write_state(state)


class SeveritiesStream(Stream):
    """Severities stream."""
    
    name = "severities"
    key_properties = ["id"]
    valid_replication_keys = ["updated_at"]
    
    def __init__(self, replication_method=None, replication_key=None):
        self.replication_method = replication_method or "FULL_TABLE"
        self.replication_key = replication_key or "updated_at"
    
    def sync(self, client: IncidentClient, state: Dict[str, Any] = None) -> Iterator[Dict[str, Any]]:
        """Sync severities."""
        state = state or {}
        starting_date = self.get_starting_time(state)
        
        severities = client.get_severities()
        record_count = 0
        all_records = []
        
        for severity in severities:
            all_records.append(severity)
            
            # Filter by bookmark date for incremental replication
            if self.replication_method == "INCREMENTAL" and starting_date and self.replication_key:
                record_date = datetime.fromisoformat(severity[self.replication_key].replace("Z", "+00:00"))
                if record_date <= starting_date:
                    continue
            
            record_count += 1
            
            # For incremental, update state record by record
            if self.replication_method == "INCREMENTAL":
                state = self.update_state(state, severity)
            
                # Every 100 records, write the state
                if record_count % 100 == 0:
                    write_state(state)
                
            yield severity
        
        # For FULL_TABLE, update state once at the end with the latest record
        if self.replication_method == "FULL_TABLE":
            state = self.update_state_for_full_table(state, all_records)
            
        # Write the state one final time
        write_state(state)


class UsersStream(Stream):
    """Users stream."""
    
    name = "users"
    key_properties = ["id"]
    
    def sync(self, client: IncidentClient, state: Dict[str, Any] = None) -> Iterator[Dict[str, Any]]:
        """Sync users."""
        for user in client.get_users():
            yield user


# Dictionary of all available streams
STREAMS = {
    "actions": ActionsStream(),
    "alerts": AlertsStream(),
    "alert_attributes": AlertAttributesStream(),
    "alert_sources": AlertSourcesStream(),
    "custom_fields": CustomFieldsStream(),
    "custom_field_options": CustomFieldOptionsStream(),
    "follow_ups": FollowUpsStream(),
    "incident_roles": IncidentRolesStream(),
    "incident_statuses": IncidentStatusesStream(),
    "incident_timestamps": IncidentTimestampsStream(),
    "incident_types": IncidentTypesStream(),
    "incident_updates": IncidentUpdatesStream(),
    "incidents": IncidentsStream(),
    "severities": SeveritiesStream(),
    "users": UsersStream(),
}