# Changelog

## v0.8.0

- The actions and follow_ups streams now sync incrementally on `updated_at`.
  The tap accepts a `--state` file and emits `STATE` messages with a bookmark
  per stream. With a bookmark, only actions and follow-ups changed since that
  time are fetched. Without one, the tap fetches everything, as before.
- Discovery marks these two streams `INCREMENTAL` with `updated_at` as the
  replication key, and always includes their `id` and `updated_at` fields.
- Targets that replaced these tables on each run now upsert into them instead,
  so actions and follow-ups deleted in incident.io are no longer removed from
  the destination.
- Upgraded `sdk-go` from 1.9.0 to 1.22.0.

## v0.7.1

- Building the tap now needs Go 1.27. Release binaries are built with Go 1.27.1.
- Upgraded all Go dependencies, including `sdk-go` from 1.0.98 to 1.9.0.
- The Docker image now uses `alpine:3.24.2` and upgrades its packages at
  build time.

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
