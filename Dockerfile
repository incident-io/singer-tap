FROM alpine:3.24.2 AS runtime

# Upgrade the base packages so we pick up security fixes released since the
# base image was built.
RUN apk upgrade --no-cache && apk add --no-cache ca-certificates

# goreleaser supplies this for us
COPY tap-incident /usr/local/bin

ENTRYPOINT ["/usr/local/bin/tap-incident"]
