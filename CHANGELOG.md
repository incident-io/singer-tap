# Changelog

## v0.7.0

- Actions and follow-ups are now read from the paginated `/v3/actions` and
  `/v3/follow_ups` endpoints. The V2 endpoints they replace returned an
  organisation's whole history in a single unpaginated response.
- Follow-up records gain a `category` field.
- The API client now comes from the official Go SDK,
  [`incident-io/sdk-go`](https://github.com/incident-io/sdk-go), rather than
  being generated into this repository.
- Building the tap now needs Go 1.24.

## v0.6.0

- Added support for escalations stream

## v0.5.0

- Added support for alerts, alert attributes, and alert sources.
