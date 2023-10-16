FROM alpine:20230329 AS runtime

# goreleaser supplies this for us
COPY incident-tap /usr/local/bin

ENTRYPOINT ["/usr/local/bin/incident-tap"]
