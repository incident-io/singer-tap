package tap

import (
	"context"

	kitlog "github.com/go-kit/log"
	"github.com/incident-io/singer-tap/client"
	"github.com/incident-io/singer-tap/model"
	"github.com/pkg/errors"
)

func init() {
	register(&StreamActions{})
}

type StreamActions struct {
}

func (s *StreamActions) Output() *Output {
	return &Output{
		Type:   OutputTypeSchema,
		Stream: "actions",
		Schema: &model.Schema{
			HasAdditionalProperties: false,
			Type:                    []string{"object"},
			Properties:              model.ActionV2.Schema().Properties,
		},
		KeyProperties:      []string{"id"},
		BookmarkProperties: []string{},
	}
}

func (s *StreamActions) GetRecords(ctx context.Context, logger kitlog.Logger, cl *client.ClientWithResponses) ([]map[string]any, error) {
	var (
		results = []map[string]any{}
	)

	response, err := cl.ActionsV2ListWithResponse(ctx, &client.ActionsV2ListParams{})
	if err != nil {
		return nil, errors.Wrap(err, "listing incidents for actions stream")
	}

	for _, element := range response.JSON200.Actions {
		results = append(results, model.ActionV2.Serialize(element))
	}

	return results, nil
}
