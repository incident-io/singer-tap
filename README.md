# incident.io Singer tap

This is the official Singer tap for [incident.io](https://incident.io/), used to
extract incident.io data that can be used alongside Singer targets to load into
a data warehouse.

## Getting started

macOS users can install the catalog using brew:

```console
brew tap incident-io/homebrew-taps
brew install tap-incident
```

Otherwise, ensure that the go runtime is installed and then:

```console
go install github.com/incident-io/singer-tap/cmd/tap-incident@latest
```

Once installed, see [documentation](docs) for example configurations.

## Using Docker

[hub]: https://hub.docker.com/r/incidentio/singer-tap/tags

A Docker image is available for containerised environments; see [Docker
Hub][hub] for more details of the image and available tags.

## Integration snapshots

The `integration` package contains snapshots of data from the integration test
incident.io account. These snapshots are not intended to be managed by hand, and
should be updated automatically instead.

Whenever you need to update the snapshots, run:

```console
$ export TEST_INCIDENT_API_KEY="<test-key>"
$ export TAP_SNAPSHOT_UPDATE='true' ginkgo -tags=integration -r ./integration
Writing snapshot file testdata/sync/follow_ups.json
...
```
[test-key]: https://start.1password.com/open/i?a=5P4EMQNU2RGBHC6TS4KB3PUFO4&v=vdgfxfm7m46sq2gzodftwi2yle&i=gayjefie3pcngvbv5tz2vsbndi&h=incident-io.1password.com

For incident.io employees, you can find the test key in 1Password
[here][test-key].

## Contributing

We're happy to accept open-source contributions or feedback. Just open a
PR/issue and we'll get back to you. This repo contains details on
[how to get started with local development](./development.md).
