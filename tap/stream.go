package tap

import (
	"context"
	"fmt"
	"time"

	kitlog "github.com/go-kit/log"
	incident "github.com/incident-io/sdk-go"
)

var streams = map[string]Stream{}

func register(s Stream) {
	op := s.Output()
	if _, ok := streams[op.Stream]; ok {
		panic(fmt.Sprintf("stream already registered: %s", op.Stream))
	}

	streams[op.Stream] = s
}

// Stream is a data model from the incident.io API that we want to represent as a Singer
// tap stream.
type Stream interface {
	// Output is the schema of the stream, in JSON schema format.
	Output() *Output
	// GetRecords returns a slice of entries in the stream. People will eventually ask for
	// this to be a channel, but we're going simple and loading everything for now.
	GetRecords(ctx context.Context, logger kitlog.Logger, cl *incident.ClientWithResponses) ([]map[string]any, error)
}

// IncrementalStream is a Stream that can fetch only the records updated since a bookmark,
// rather than the whole history on every run.
type IncrementalStream interface {
	Stream
	// GetRecordsSince returns the records updated at or after since. If since is nil, it
	// returns every record.
	GetRecordsSince(ctx context.Context, logger kitlog.Logger, cl *incident.ClientWithResponses, since *time.Time) ([]map[string]any, error)
}
