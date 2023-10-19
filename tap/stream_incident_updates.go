package tap

import (
	"context"

	kitlog "github.com/go-kit/log"
	"github.com/incident-io/singer-tap/client"
	"github.com/incident-io/singer-tap/model"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

func init() {
	register(&StreamIncidentUpdates{})
}

type StreamIncidentUpdates struct {
}

func (s *StreamIncidentUpdates) Output() *Output {
	return &Output{
		Type:   OutputTypeSchema,
		Stream: "incident_updates",
		Schema: &model.Schema{
			HasAdditionalProperties: false,
			Type:                    []string{"object"},
			Properties:              model.IncidentUpdateV2.Schema().Properties,
		},
		KeyProperties:      []string{"id"},
		BookmarkProperties: []string{},
	}
}

func (s *StreamIncidentUpdates) GetRecords(ctx context.Context, logger kitlog.Logger, cl *client.ClientWithResponses) ([]map[string]any, error) {
	var (
		after    *string
		pageSize = 250
		results  = []map[string]any{}
	)

	for {
		logger.Log("msg", "loading page", "page_size", pageSize, "after", after)
		page, err := cl.IncidentUpdatesV2ListWithResponse(ctx, &client.IncidentUpdatesV2ListParams{
			PageSize: &pageSize,
			After:    after,
		})
		if err != nil {
			return nil, errors.Wrap(err, "listing incident updates")
		}

		for _, element := range page.JSON200.IncidentUpdates {
			results = append(results, model.IncidentUpdateV2.Serialize(element))
		}
		if count := len(page.JSON200.IncidentUpdates); count == 0 {
			return results, nil // end pagination
		} else {
			after = lo.ToPtr(page.JSON200.IncidentUpdates[count-1].Id)
		}
	}
}
