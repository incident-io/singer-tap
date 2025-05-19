# incident.io Singer tap

This is the official Singer tap for [incident.io](https://incident.io/), used to
extract incident.io data that can be used alongside Singer targets to load into
a data warehouse.

## Getting started

### Installation

#### Using pip

```console
pip install tap-incident
```

#### From source

```console
pip install git+https://github.com/Bilanc/bilanc-incident-io-tap.git
```

### Configuration

Create a `config.json` file that looks like this:

```json
{
  "api_key": "<your-api-key>",
  "endpoint": "https://api.incident.io"
}
```

You'll need an incident.io API key with permission to:
- View data, like public incidents and organisation settings
- View catalog types and entries

If you want this tap to have access to private incident data, also include the
following scope:
- View all incident data, including private incidents

### Running the tap

```console
tap-incident --config config.json
```

### Discovery mode

To use the discovery mode, which outputs a catalog of streams:

```console
tap-incident --config config.json --discover > catalog.json
```

### Configuring exports

By default, the tap will export all data it can.

If you want to control what fields and tables you wish to export, you will need a catalog file. You can use the discovery command to create a default catalog file which will contain all streams and columns along with their properties.

```console
$ tap-incident --discover --config=config.json > catalog.json
```

The catalog file follows the Singer specification to enable or disable a stream. Add the field `selected: true|false` inside the stream's top-level metadata. For example, to disable exporting severities:

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
        ...
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

## Running with a catalog file

```console
tap-incident --config config.json --catalog catalog.json
```

## Available streams

The tap-incident package supports the following streams:

- actions
- alerts
- alert_attributes
- alert_sources
- custom_field_options
- custom_fields
- follow_ups
- incident_roles
- incident_statuses
- incident_timestamps
- incident_types
- incident_updates
- incidents
- severities
- users

## For developers

If you want to contribute to this Singer tap, please see the [development guide](https://github.com/incident-io/tap-incident/blob/main/development.md) for more details.