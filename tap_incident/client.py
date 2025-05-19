"""Incident.io REST API client."""

import backoff
import logging
import requests
from typing import Dict, Any, Optional, List, Tuple, Union, Iterator

from tap_incident import __version__

LOGGER = logging.getLogger(__name__)


class IncidentClientError(Exception):
    """Custom exception for Incident API errors."""
    pass


class IncidentClient:
    """Incident.io API client."""

    def __init__(self, api_key: str, api_endpoint: str = "https://api.incident.io"):
        """Initialize client.
        
        Args:
            api_key: Incident.io API key
            api_endpoint: API endpoint URL
        """
        self.api_key = api_key
        self.api_endpoint = api_endpoint.rstrip("/")
        self.session = requests.Session()
        self.session.headers.update({
            "Authorization": f"Bearer {api_key}",
            "Content-Type": "application/json",
            "User-Agent": f"tap-incident/{__version__}"
        })

    @backoff.on_exception(
        backoff.expo,
        (requests.exceptions.RequestException),
        max_tries=5,
        factor=2,
    )
    def _make_request(
        self, 
        method: str, 
        endpoint: str, 
        params: Optional[Dict[str, Any]] = None, 
        data: Optional[Dict[str, Any]] = None
    ) -> Dict[str, Any]:
        """Make an API request with retry logic.
        
        Args:
            method: HTTP method
            endpoint: API endpoint path
            params: Query parameters
            data: Request body
            
        Returns:
            Response JSON data
            
        Raises:
            IncidentClientError: If the request fails
        """
        url = f"{self.api_endpoint}/{endpoint.lstrip('/')}"
        
        LOGGER.debug(f"Making {method} request to {url}")
        
        try:
            response = self.session.request(method, url, params=params, json=data)
            response.raise_for_status()
            return response.json()
        except requests.exceptions.HTTPError as e:
            error_message = str(e)
            try:
                error_details = response.json()
                error_message = f"{error_message}: {error_details}"
            except (ValueError, KeyError):
                pass
            
            LOGGER.error(f"HTTP error: {error_message}")
            raise IncidentClientError(error_message) from e
        except requests.exceptions.RequestException as e:
            LOGGER.error(f"Request error: {str(e)}")
            raise IncidentClientError(str(e)) from e

    def paginate(
        self, 
        endpoint: str, 
        params: Optional[Dict[str, Any]] = None
    ) -> Iterator[Dict[str, Any]]:
        """Paginate through API results.
        
        Args:
            endpoint: API endpoint path
            params: Query parameters
            
        Yields:
            Each response page
        """
        params = params or {}
        
        # Default page size if not specified
        if "page_size" not in params:
            params["page_size"] = 50
            
        after = None
        
        while True:
            # Add after cursor if we have one
            if after:
                params["after"] = after
                
            response = self._make_request("GET", endpoint, params=params)
            
            yield response
            
            # Check if we have more pages
            # Look for the last item's ID in the response
            # The structure of response depends on the endpoint, but typically it's a list
            # inside a key matching the endpoint name
            
            # Try to find the key that contains a list of records
            items = None
            for key, value in response.items():
                if isinstance(value, list) and len(value) > 0:
                    items = value
                    break
            
            if not items or len(items) == 0:
                # No more results
                break
                
            # Use the ID of the last item as the 'after' cursor
            after = items[-1].get("id")
            
            if not after:
                # No valid cursor found
                break

    # API endpoint methods
    def get_actions(self) -> List[Dict[str, Any]]:
        """Get all actions."""
        response = self._make_request("GET", "v2/actions")
        return response.get("actions", [])

    def get_alerts(self) -> List[Dict[str, Any]]:
        """Get all alerts."""
        alerts = []
        for page in self.paginate("v2/alerts"):
            alerts.extend(page.get("alerts", []))
        return alerts

    def get_alert_attributes(self) -> List[Dict[str, Any]]:
        """Get all alert attributes."""
        response = self._make_request("GET", "v2/alert_attributes")
        return response.get("alert_attributes", [])

    def get_alert_sources(self) -> List[Dict[str, Any]]:
        """Get all alert sources."""
        response = self._make_request("GET", "v2/alert_sources")
        return response.get("alert_sources", [])

    def get_custom_fields(self) -> List[Dict[str, Any]]:
        """Get all custom fields."""
        response = self._make_request("GET", "v2/custom_fields")
        return response.get("custom_fields", [])

    def get_custom_field_options(self, custom_field_id: str) -> List[Dict[str, Any]]:
        """Get options for a specific custom field."""
        options = []
        for page in self.paginate(f"v1/custom_field_options", {"custom_field_id": custom_field_id}):
            options.extend(page.get("custom_field_options", []))
        return options

    def get_follow_ups(self) -> List[Dict[str, Any]]:
        """Get all follow-ups."""
        response = self._make_request("GET", "v2/follow_ups")
        return response.get("follow_ups", [])

    def get_incident_roles(self) -> List[Dict[str, Any]]:
        """Get all incident roles."""
        response = self._make_request("GET", "v2/incident_roles")
        return response.get("incident_roles", [])

    def get_incident_statuses(self) -> List[Dict[str, Any]]:
        """Get all incident statuses."""
        response = self._make_request("GET", "v1/incident_statuses")
        return response.get("incident_statuses", [])

    def get_incident_timestamps(self) -> List[Dict[str, Any]]:
        """Get all incident timestamps."""
        response = self._make_request("GET", "v2/incident_timestamps")
        return response.get("incident_timestamps", [])

    def get_incident_types(self) -> List[Dict[str, Any]]:
        """Get all incident types."""
        response = self._make_request("GET", "v1/incident_types")
        return response.get("incident_types", [])

    def get_incident_updates(self, incident_id: Optional[str] = None) -> List[Dict[str, Any]]:
        """Get all incident updates, optionally filtered by incident ID."""
        params = {}
        if incident_id:
            params["incident_id"] = incident_id
            
        updates = []
        for page in self.paginate("v2/incident_updates", params):
            updates.extend(page.get("incident_updates", []))
        return updates

    def get_incidents(self) -> List[Dict[str, Any]]:
        """Get all incidents."""
        incidents = []
        for page in self.paginate("v2/incidents"):
            incidents.extend(page.get("incidents", []))
        return incidents

    def get_incident_attachments(self, incident_id: str) -> List[Dict[str, Any]]:
        """Get attachments for a specific incident."""
        response = self._make_request("GET", "v1/incident_attachments", {"incident_id": incident_id})
        return response.get("incident_attachments", [])

    def get_severities(self) -> List[Dict[str, Any]]:
        """Get all severities."""
        response = self._make_request("GET", "v1/severities")
        return response.get("severities", [])

    def get_users(self) -> List[Dict[str, Any]]:
        """Get all users."""
        users = []
        for page in self.paginate("v2/users"):
            users.extend(page.get("users", []))
        return users