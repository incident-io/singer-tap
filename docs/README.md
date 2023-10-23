# incident.io Singer tap

Jump to [Getting started](#getting-started) if you want to begin from an example
configuration or just want to experiment.

Otherwise check-out the rest of our documentation for details on how to:

- [Run the exporter from CI tools like CircleCI or GitHub Actions](deploying.md)

[open-issue]: https://github.com/incident-io/singer-tap/issues/new

If you can't find an answer to your question, please [open an issue][open-issue]
with your request and we'll be happy to help out.

## Getting started

[api-keys]: https://app.incident.io/settings/api-keys

Start by visiting your [incident dashboard][api-keys] to create an API key with
permission to:

- View data, like public incidents and organisation settings
- View catalog types and entries

If you want this tap to have access to private incident data, also include the
following scope:

- View all incident data, including private incidents

Once you've created the key, create a `config.json` file that you'll use to
configure the tap that looks like this:

```json
{
  "api_key": "<your-api-key"
}
```

You can check this works by running:

```console
$ tap-incident --discover --config=config.json
```

## Configuring exports

By default the tap will export all data it can.

If you want to control what fields and tables you wish to export you will need a catalog file. You can use the discover command to create a default catalog file which will contain all streams and columns along with their properties.

```console
$ tap-incident --discover --config=config.json > catalog.json
```

The catalog file follows the Singer specification to enable or disable a stream add the field `selected: <true|false>` inside the streams top level metadata. For example to disable exporting severities:

```json
{
  "streams": [
    {
      "stream": "severities",
      "tap_stream_id": "severities",
      "schema": {
        ...
      },
      "metadata": [
        {
          "breadcrumb": [],
          "metadata": {
            "selected": false, // Add this
            "inclusion": "available",
            "selected-by-default": true,
            "forced-replication-method": "FULL_TABLE"
          }
        },
        {
          "breadcrumb": [
            "properties",
            "id"
          ],
```

You can adjust which columns within a stream you wish to export similarly, by adding the same field inside the metadata for that column, for example to not export a description field:

```json
    {
      "breadcrumb": [
        "properties",
        "description"
      ],
      "metadata": {
        "inclusion": "available",
        "selected": false, // Add this
        "selected-by-default": true
      }
    },
```

## Table Information

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
