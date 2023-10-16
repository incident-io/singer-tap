# incident.io Singer tap

Jump to [Getting started](#getting-started) if you want to begin from an example
configuration or just want to experiment.

Otherwise check-out the rest of our documentation for details on how to:

- [Run the importer from CI tools like CircleCI or GitHub Actions](deploying.md)

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
