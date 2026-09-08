package tap

import (
	"context"

	kitlog "github.com/go-kit/log"
	incident "github.com/incident-io/sdk-go"
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
			Properties:              model.FollowUpV3.Schema().Properties,
		},
		KeyProperties:      []string{"id"},
		BookmarkProperties: []string{},
	}
}

func (s *StreamFollowUps) GetRecords(ctx context.Context, logger kitlog.Logger, cl *incident.ClientWithResponses) ([]map[string]any, error) {
	var (
		after    *string
		pageSize = int64(250)
		results  = []map[string]any{}
	)

	for {
		logger.Log("msg", "loading follow-ups page", "page_size", pageSize, "after", after)
		page, err := cl.FollowUpsV3ListWithResponse(ctx, &incident.FollowUpsV3ListParams{
			PageSize: &pageSize,
			After:    after,
		})
		if err != nil {
			return nil, errors.Wrap(err, "listing follow-ups")
		}

		for _, element := range page.JSON200.FollowUps {
			results = append(results, model.FollowUpV3.Serialize(element))
		}

		// The endpoint sets after only when the page was full, so no cursor means no more rows.
		if page.JSON200.PaginationMeta.After == nil {
			return results, nil
		}

		after = page.JSON200.PaginationMeta.After
	}
}
