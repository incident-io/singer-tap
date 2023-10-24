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

## Regenerating the CI test files

For this you will need to work for incident.io and have access to the required API key.

First you need to get the python target-csv working locally. Annoyingly the official version only works on super old pythons so you can install the forked version. Recommend setting up a virtual env for this first:

```
pyenv virtualenv 3.10.13 csv-target
pyenv activate csv-target
python -m pip install 'target-csv @ git+https://github.com/rliddler/target-csv@master'
```

You will also need to create a config for both the tap and the target. For the target you will pretty much need the following:

```
{
  "destination_path": "integration/testdata/sync",
  "timestamp_override": "2023",
  "override_file": true
}
```

You then need to delete the existing integration/testdata/sync/\* files (annoying I know but the target only appends to each file) and regenerate them using the following command:

```
tap-incident --config tmp/config.json --catalog integration/catalog.json | target-csv --config tmp/target-config.json
```

Check any new changes in the git diff - and if you approve commit them. The github workflow will basically do the same as this and check there are no differences in the generated rows.
