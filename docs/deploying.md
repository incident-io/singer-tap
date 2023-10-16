# Deploying using CI

Most people run the tap in a CI pipeline that triggers periodically depending on
their needs.

## CircleCI

If you run on CircleCI, an example config is below.

> You can configure a [scheduled pipeline](https://circleci.com/docs/scheduled-pipelines/)
> to run the sync on a regular cadence. This is recommended if your importer
> config uses sources other than local files.

```yaml
# .circleci/config.yml
---
version: 2.1

jobs:
  sync:
    docker:
      - image: cimg/base:2023.04
    working_directory: ~/app
    steps:
      - checkout
      - run:
          name: Install incident-tap
          command: |
            VERSION="0.1.0"

            echo "Installing importer v${VERSION}..."
            curl -L \
              -o "/tmp/incident-tap_${VERSION}_linux_amd64.tar.gz" \
              "https://github.com/incident-io/singer-tap/releases/download/v${VERSION}/incident-tap_${VERSION}_linux_amd64.tar.gz"
            tar zxf "/tmp/incident-tap_${VERSION}_linux_amd64.tar.gz" -C /tmp
      - run:
          name: Sync
          command: |
            /tmp/incident-tap --config config.json

workflows:
  version: 2
  sync:
    jobs:
      - sync
```

## GitHub Actions

If you run on GitHub Actions, an example config is:

```yaml
name: Sync

# Run hourly:
on:
  schedule:
    - cron: "55 * * * *" # hourly, on the 55th minute

jobs:
  sync:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v3
        with:
          fetch-depth: 0
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: "1.21"
      - name: Install incident-tap
        run: |
          go install github.com/incident-io/singer-tap/cmd/incident-tap@latest
      - name: Sync
        run: |
          incident-tap --config config.json
```
