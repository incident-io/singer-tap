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

## Contributing

We're happy to accept open-source contributions or feedback. Just open a
PR/issue and we'll get back to you. This repo contains details on
[how to get started with local development](./development.md).
