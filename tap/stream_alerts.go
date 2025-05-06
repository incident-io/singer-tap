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
	register(&StreamAlerts{})
}

type StreamAlerts struct{}

func (s *StreamAlerts) Output() *Output {
	return &Output{
		Type:   OutputTypeSchema,
		Stream: "alerts",
		Schema: &model.Schema{
			HasAdditionalProperties: false,
			Type:                    []string{"object"},
			Properties:              model.AlertV2.Schema().Properties,
		},
		KeyProperties:      []string{"id"},
		BookmarkProperties: []string{},
	}
}

func (s *StreamAlerts) GetRecords(ctx context.Context, logger kitlog.Logger, cl *client.ClientWithResponses) ([]map[string]any, error) {
	var (
		after    *string
		pageSize = int64(50)
		results  = []map[string]any{}
	)

	for {
		logger.Log("msg", "loading alerts page", "page_size", pageSize, "after", after)
		page, err := cl.AlertsV2ListWithResponse(ctx, &client.AlertsV2ListParams{
			PageSize: pageSize,
			After:    after,
		})
		if err != nil {
			return nil, errors.Wrap(err, "listing alerts")
		}

		for _, element := range page.JSON200.Alerts {
			results = append(results, model.AlertV2.Serialize(element))
		}
		if count := len(page.JSON200.Alerts); count == 0 {
			return results, nil // end pagination
		} else {
			after = lo.ToPtr(page.JSON200.Alerts[count-1].Id)
		}
	}
}
