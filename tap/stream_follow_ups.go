package tap

import (
	"context"

	kitlog "github.com/go-kit/log"
	"github.com/incident-io/singer-tap/client"
	"github.com/incident-io/singer-tap/model"
	"github.com/pkg/errors"
)

func init() {
	register(&StreamFollowUps{})
}

type StreamFollowUps struct {
}

func (s *StreamFollowUps) Output() *Output {
	return &Output{
		Type:   OutputTypeSchema,
		Stream: "follow_ups",
		Schema: &model.Schema{
			HasAdditionalProperties: false,
			Type:                    []string{"object"},
			Properties:              model.FollowUpV2.Schema().Properties,
		},
		KeyProperties:      []string{"id"},
		BookmarkProperties: []string{},
	}
}

func (s *StreamFollowUps) GetRecords(ctx context.Context, logger kitlog.Logger, cl *client.ClientWithResponses) ([]map[string]any, error) {
	var (
		results = []map[string]any{}
	)

	response, err := cl.FollowUpsV2ListWithResponse(ctx, &client.FollowUpsV2ListParams{})
	if err != nil {
		return nil, errors.Wrap(err, "listing incidents for actions stream")
	}

	for _, element := range response.JSON200.FollowUps {
		results = append(results, model.FollowUpV2.Serialize(element))
	}

	return results, nil
}
