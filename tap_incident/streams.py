"""Stream definitions for tap-incident."""

import logging
from abc import ABC, abstractmethod
from typing import Dict, Any, List, Iterator, Optional

from tap_incident.client import IncidentClient

LOGGER = logging.getLogger(__name__)


class Stream(ABC):
    """Base stream class."""
    
    # Class variables that should be defined by each stream
    name = None
    key_properties = None
    
    @abstractmethod
    def sync(self, client: IncidentClient) -> Iterator[Dict[str, Any]]:
        """Sync data from the stream.
        
        Args:
            client: IncidentClient instance
            
        Yields:
            Records from the stream
        """
        pass


class ActionsStream(Stream):
    """Actions stream."""
    
    name = "actions"
    key_properties = ["id", "incident_id"]
    
    def sync(self, client: IncidentClient) -> Iterator[Dict[str, Any]]:
        """Sync actions."""
        for action in client.get_actions():
            yield action


class AlertsStream(Stream):
    """Alerts stream."""
    
    name = "alerts"
    key_properties = ["id"]
    
    def sync(self, client: IncidentClient) -> Iterator[Dict[str, Any]]:
        """Sync alerts."""
        for alert in client.get_alerts():
            yield alert


class AlertAttributesStream(Stream):
    """Alert attributes stream."""
    
    name = "alert_attributes"
    key_properties = ["id"]
    
    def sync(self, client: IncidentClient) -> Iterator[Dict[str, Any]]:
        """Sync alert attributes."""
        for attribute in client.get_alert_attributes():
            yield attribute


class AlertSourcesStream(Stream):
    """Alert sources stream."""
    
    name = "alert_sources"
    key_properties = ["id"]
    
    def sync(self, client: IncidentClient) -> Iterator[Dict[str, Any]]:
        """Sync alert sources."""
        for source in client.get_alert_sources():
            yield source


class CustomFieldsStream(Stream):
    """Custom fields stream."""
    
    name = "custom_fields"
    key_properties = ["id"]
    
    def sync(self, client: IncidentClient) -> Iterator[Dict[str, Any]]:
        """Sync custom fields."""
        for custom_field in client.get_custom_fields():
            yield custom_field


class CustomFieldOptionsStream(Stream):
    """Custom field options stream."""
    
    name = "custom_field_options"
    key_properties = ["id", "custom_field_id"]
    
    def sync(self, client: IncidentClient) -> Iterator[Dict[str, Any]]:
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
    
    def sync(self, client: IncidentClient) -> Iterator[Dict[str, Any]]:
        """Sync follow ups."""
        for follow_up in client.get_follow_ups():
            yield follow_up


class IncidentRolesStream(Stream):
    """Incident roles stream."""
    
    name = "incident_roles"
    key_properties = ["id"]
    
    def sync(self, client: IncidentClient) -> Iterator[Dict[str, Any]]:
        """Sync incident roles."""
        for role in client.get_incident_roles():
            yield role


class IncidentStatusesStream(Stream):
    """Incident statuses stream."""
    
    name = "incident_statuses"
    key_properties = ["id"]
    
    def sync(self, client: IncidentClient) -> Iterator[Dict[str, Any]]:
        """Sync incident statuses."""
        for status in client.get_incident_statuses():
            yield status


class IncidentTimestampsStream(Stream):
    """Incident timestamps stream."""
    
    name = "incident_timestamps"
    key_properties = ["id"]
    
    def sync(self, client: IncidentClient) -> Iterator[Dict[str, Any]]:
        """Sync incident timestamps."""
        for timestamp in client.get_incident_timestamps():
            yield timestamp


class IncidentTypesStream(Stream):
    """Incident types stream."""
    
    name = "incident_types"
    key_properties = ["id"]
    
    def sync(self, client: IncidentClient) -> Iterator[Dict[str, Any]]:
        """Sync incident types."""
        for incident_type in client.get_incident_types():
            yield incident_type


class IncidentUpdatesStream(Stream):
    """Incident updates stream."""
    
    name = "incident_updates"
    key_properties = ["id", "incident_id"]
    
    def sync(self, client: IncidentClient) -> Iterator[Dict[str, Any]]:
        """Sync incident updates."""
        for update in client.get_incident_updates():
            yield update


class IncidentsStream(Stream):
    """Incidents stream."""
    
    name = "incidents"
    key_properties = ["id"]
    
    def sync(self, client: IncidentClient) -> Iterator[Dict[str, Any]]:
        """Sync incidents."""
        incidents = client.get_incidents()
        
        for incident in incidents:
            incident_id = incident["id"]
            
            # Fetch related data for each incident
            attachments = client.get_incident_attachments(incident_id)
            incident["attachments"] = attachments
            
            updates = client.get_incident_updates(incident_id)
            incident["updates"] = updates
            
            yield incident


class SeveritiesStream(Stream):
    """Severities stream."""
    
    name = "severities"
    key_properties = ["id"]
    
    def sync(self, client: IncidentClient) -> Iterator[Dict[str, Any]]:
        """Sync severities."""
        for severity in client.get_severities():
            yield severity


class UsersStream(Stream):
    """Users stream."""
    
    name = "users"
    key_properties = ["id"]
    
    def sync(self, client: IncidentClient) -> Iterator[Dict[str, Any]]:
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

STREAMS_back = {
    "actions": ActionsStream(),
    "incidents": IncidentsStream(),
}