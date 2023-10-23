# Development

You'll need [Go installed][go] to be able to contribute to `singer-tap`.

You can run all of the tests using `go test`:

```
go test ./...
```

To build the binary, run `make`: this will place a built version of the tool in
the `bin` directory. If you work for incident.io and have a local instance of
the app running, then you can point it to your local environment using the
`--api-key` flag, or an environment variable:

```
export INCIDENT_ENDPOINT="http://localhost:3001/api/public"
```

[go]: https://go.dev/doc/install

## Adding new streams

Each table the tap imports is implemented as a `Stream`. If you want to export a new table then you can just add a new `stream_<tablename>` and implement the stream interface. See the tap folder for more examples.

## Adding new fields and schemas

The translation / interface layer between incident.io types and singer types are all contained within the `model` folder. To implement a new type or schema add a new file and implement the following methods:

- `Schema()` - singer schema representing the full type, may have `Optional()` and `ArrayOf` types
- `Serialize()` - takes the incident.io client type and returns data in the format specified by the `Schema()` method
