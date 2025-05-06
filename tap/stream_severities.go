package tap

import (
	"context"

	kitlog "github.com/go-kit/log"
	"github.com/incident-io/singer-tap/client"
	"github.com/incident-io/singer-tap/model"
	"github.com/pkg/errors"
)

func init() {
	register(&StreamSeverities{})
}

type StreamSeverities struct {
}

func (s *StreamSeverities) Output() *Output {
	return &Output{
		Type:   OutputTypeSchema,
		Stream: "severities",
		Schema: &model.Schema{
			HasAdditionalProperties: false,
			Type:                    []string{"object"},
			Properties:              model.SeverityV2.Schema().Properties,
		},
		KeyProperties:      []string{"id"},
		BookmarkProperties: []string{},
	}
}

func (s *StreamSeverities) GetRecords(ctx context.Context, logger kitlog.Logger, cl *client.ClientWithResponses) ([]map[string]any, error) {
	var (
		results = []map[string]any{}
	)

	response, err := cl.SeveritiesV1ListWithResponse(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "listing severities")
	}

	for _, element := range response.JSON200.Severities {
		results = append(results, model.SeverityV1.Serialize(&element))
	}

	return results, nil
}
