package tap

import (
	"context"

	kitlog "github.com/go-kit/log"
	"github.com/incident-io/singer-tap/client"
	"github.com/incident-io/singer-tap/model"
	"github.com/pkg/errors"
)

func init() {
	register(&StreamIncidentStatuses{})
}

type StreamIncidentStatuses struct {
}

func (s *StreamIncidentStatuses) Output() *Output {
	return &Output{
		Type:   OutputTypeSchema,
		Stream: "incident_statuses",
		Schema: &model.Schema{
			HasAdditionalProperties: false,
			Type:                    []string{"object"},
			Properties:              model.IncidentStatusV1.Schema().Properties,
		},
		KeyProperties:      []string{"id"},
		BookmarkProperties: []string{},
	}
}

func (s *StreamIncidentStatuses) GetRecords(ctx context.Context, logger kitlog.Logger, cl *client.ClientWithResponses) ([]map[string]any, error) {
	var (
		results = []map[string]any{}
	)

	response, err := cl.IncidentStatusesV1ListWithResponse(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "listing incident statuses")
	}

	for _, element := range response.JSON200.IncidentStatuses {
		results = append(results, model.IncidentStatusV1.Serialize(element))
	}

	return results, nil
}
