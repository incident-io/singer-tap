# incident.io (tap-incident)

## Connecting incident.io

### Requirements

To set up the incident.io tap in Stitch, you need an incident.io API key with permission to view incidents and org settings. To create a new API key you will need the right user permissions or to ask someone within your organisation who does.

### Setup

To create a new API key you can navigate [to your API key settings page](https://app.incident.io/~/settings/api-keys) and click "Add new".

You need to create the key with the following permissions:

- View data, like public incidents and organisation settings
- View catalog types and entries

If you want this tap to have access to private incident data, also include the
following scope:

- View all incident data, including private incidents

---

## incident.io Tap Replication

There is no partial or incremental replication available in the incident.io tap. Instead each stream will perform a full table replication each time.

The amount of data in each stream is relatively low so this should not be an issue for most customers.

---

## incident.io Table Schemas

### Incidents

- Table name: incidents
- Description: Incidents are a core resource, on which many other resources (actions, etc) are created. You will find incident specific actions, updates, timestamps etc nested on this resource.
- Primary key column(s): id
- Replication: full table
- API documentation: [Incidents V2](https://api-docs.incident.io/tag/Incidents-V2)

### Actions

- Table name: actions
- Description: Incident actions are used during an incident, to track work such as 'restart the database' or 'contact the customer'. Actions are also included in the incidents table.
- Primary key column(s): id, incident_id
- Replication: full table
- API documentation: [Actions V2](https://api-docs.incident.io/tag/Actions-V2)

### Custom Field Options

- Table name: custom_field_options
- Description:
- Primary key column(s): id, custom_field_id
- Replication: full table
- API documentation: [Custom Field Options V1](https://api-docs.incident.io/tag/Custom-Field-Options-V1)

### Custom Fields

- Table name: custom_fields
- Description: Custom fields are used to attach metadata to incidents, which you can use when searching for incidents in the dashboard, triggering workflows, building announcement rules or for your own data needs.
- Primary key column(s): id
- Replication: full table
- API documentation: [Custom Fields V2](https://api-docs.incident.io/tag/Custom-Fields-V2)

### Follow Ups

- Table name: follow_ups
- Description: Incidents can have follow-ups associated with them, which track work that should be done after an incident (e.g. improving some documentation, or upgrading a dependency). They can also be exported to external issue trackers.
- Primary key column(s): id, incident_id
- Replication: full table
- API documentation: [Follow Ups V2](https://api-docs.incident.io/tag/Follow-ups-V2)

### Incident Roles

- Table name: incident_roles
- Description: During an incident, you can assign responders to one of the incident roles that are configured in your organisation settings.
- Primary key column(s): id
- Replication: full table
- API documentation: [Incident Roles V2](https://api-docs.incident.io/tag/Incident-Roles-V2)

### Incident Statuses

- Table name: incident_statuses
- Description: Each incident has a status, picked from one of the statuses configured in your organisations settings.
- Primary key column(s): id
- Replication: full table
- API documentation: [Incident Statuses V1](https://api-docs.incident.io/tag/Incident-Statuses-V1)

### Incident Timestamps

- Table name: incident_timestamps
- Description: Each incident has a number of timestamps; some being defaults that we set on each incident for you, and other being configured for your organisation within settings.
- Primary key column(s): id
- Replication: full table
- API documentation: [Incident Timestamps V2](https://api-docs.incident.io/tag/Incident-Timestamps-V2)

### Incident Types

- Table name: incident_types
- Description: With incident types enabled, you can tailor your process to the situation you're responding to with different custom fields and roles for each incident type.
- Primary key column(s): id
- Replication: full table
- API documentation: [Incident Types V1](https://api-docs.incident.io/tag/Incident-Types-V1)

### Incident Updates

- Table name: incident_updates
- Description: Incident Updates allows you to see all the updates that have been shared against a particular incident. This will include any time that the Severity or Status of an incident changed, alongside any additional updates that were provided. Incident updates are also included in the Incidents table
- Primary key column(s): id, incident_id
- Replication: full table
- API documentation: [Incident Updates V2](https://api-docs.incident.io/tag/Incident-Updates-V2)

### Severities

- Table name: severities
- Description: Each incident has a severity, picked from one of the severities configured in your organisations settings.
- Primary key column(s): id
- Replication: full table
- API documentation: [Severities V1](https://api-docs.incident.io/tag/Severities-V1)

### Users

- Table name: users
- Description: Users all have a single base role, and can be assigned multiple custom roles. They can be managed via your Slack workspace or SAML provider.
- Primary key column(s): id, custom_field_id
- Replication: full table
- API documentation: [Users V2](https://api-docs.incident.io/tag/Users-V2)
