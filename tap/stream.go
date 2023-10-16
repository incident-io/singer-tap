package tap

import (
	"context"
	"fmt"

	kitlog "github.com/go-kit/log"
	"github.com/incident-io/singer-tap/client"
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
	GetRecords(ctx context.Context, logger kitlog.Logger, cl *client.ClientWithResponses) ([]map[string]any, error)
}
