# Development

You'll need [Go installed][go] to be able to contribute to `singer-tap`.

You can run all of the tests using `go test`:

```
go test ./...
```

To build the binary, run `make`: this will place a built version of the tool in
the `bin` directory.

## Using Environment Variables

You can provide configuration via environment variables instead of a config file:

```bash
# Required: Set your API key
export INCIDENT_API_KEY="your-api-key"

# Optional: Set a custom endpoint (defaults to https://api.incident.io)
export INCIDENT_ENDPOINT="http://localhost:3001/api/public"

# Run discovery mode
./bin/tap-incident --discover

# Extract data
./bin/tap-incident
```

Note: If you provide both environment variables and a config file, the config file values take precedence.

[go]: https://go.dev/doc/install

## Updating the API client

The API client comes from [`github.com/incident-io/sdk-go`][sdk], the official Go
SDK, which is generated from our published OpenAPI schema and released whenever
that schema changes. There is no client to regenerate here — to pick up new
endpoints or fields, bump the dependency:

```
go get github.com/incident-io/sdk-go@latest
go mod tidy
```

`client/client.go` wraps the SDK's constructor to add retries and to turn non-2xx
responses into errors.

[sdk]: https://github.com/incident-io/sdk-go

## Adding new streams

Each table the tap imports is implemented as a `Stream`. If you want to export a new table then you can just add a new `stream_<tablename>` and implement the stream interface. See the tap folder for more examples.

## Adding new fields and schemas

The translation / interface layer between incident.io types and singer types are all contained within the `model` folder. To implement a new type or schema add a new file and implement the following methods:

- `Schema()` - singer schema representing the full type, may have `Optional()` and `ArrayOf` types
- `Serialize()` - takes the incident.io client type and returns data in the format specified by the `Schema()` method
