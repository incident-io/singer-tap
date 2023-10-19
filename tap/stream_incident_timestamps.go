package tap

import (
	"context"

	kitlog "github.com/go-kit/log"
	"github.com/incident-io/singer-tap/client"
	"github.com/incident-io/singer-tap/model"
	"github.com/pkg/errors"
)

func init() {
	register(&StreamIncidentTimestamps{})
}

type StreamIncidentTimestamps struct {
}

func (s *StreamIncidentTimestamps) Output() *Output {
	return &Output{
		Type:   OutputTypeSchema,
		Stream: "incident_timestamps",
		Schema: &model.Schema{
			HasAdditionalProperties: false,
			Type:                    []string{"object"},
			Properties:              model.IncidentTimestampV2.Schema().Properties,
		},
		KeyProperties:      []string{"id"},
		BookmarkProperties: []string{},
	}
}

func (s *StreamIncidentTimestamps) GetRecords(ctx context.Context, logger kitlog.Logger, cl *client.ClientWithResponses) ([]map[string]any, error) {
	var (
		results = []map[string]any{}
	)

	response, err := cl.IncidentTimestampsV2ListWithResponse(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "listing incident timestamps")
	}

	for _, element := range response.JSON200.IncidentTimestamps {
		results = append(results, model.IncidentTimestampV2.Serialize(element))
	}

	return results, nil
}
