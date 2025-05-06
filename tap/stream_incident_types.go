package tap

import (
	"context"

	kitlog "github.com/go-kit/log"
	"github.com/incident-io/singer-tap/client"
	"github.com/incident-io/singer-tap/model"
	"github.com/pkg/errors"
)

func init() {
	register(&StreamIncidentTypes{})
}

type StreamIncidentTypes struct {
}

func (s *StreamIncidentTypes) Output() *Output {
	return &Output{
		Type:   OutputTypeSchema,
		Stream: "incident_types",
		Schema: &model.Schema{
			HasAdditionalProperties: false,
			Type:                    []string{"object"},
			Properties:              model.IncidentTypeV2.Schema().Properties,
		},
		KeyProperties:      []string{"id"},
		BookmarkProperties: []string{},
	}
}

func (s *StreamIncidentTypes) GetRecords(ctx context.Context, logger kitlog.Logger, cl *client.ClientWithResponses) ([]map[string]any, error) {
	var (
		results = []map[string]any{}
	)

	response, err := cl.IncidentTypesV1ListWithResponse(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "listing incident types")
	}

	for _, element := range response.JSON200.IncidentTypes {
		results = append(results, model.IncidentTypeV1.Serialize(&element))
	}

	return results, nil
}
