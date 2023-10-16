FROM alpine:20230329 AS runtime

# goreleaser supplies this for us
COPY tap-incident /usr/local/bin

ENTRYPOINT ["/usr/local/bin/tap-incident"]
