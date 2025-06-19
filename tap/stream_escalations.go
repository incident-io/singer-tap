package tap

import (
	"context"

	kitlog "github.com/go-kit/log"
	"github.com/incident-io/singer-tap/client"
	"github.com/incident-io/singer-tap/model"
	"github.com/pkg/errors"
)

func init() {
	register(&StreamEscalations{})
}

type StreamEscalations struct {
}

func (s *StreamEscalations) Output() *Output {
	return &Output{
		Type:   OutputTypeSchema,
		Stream: "escalations",
		Schema: &model.Schema{
			HasAdditionalProperties: false,
			Type:                    []string{"object"},
			Properties:              model.EscalationV2.Schema().Properties,
		},
		KeyProperties:      []string{"id"},
		BookmarkProperties: []string{},
	}
}

func (s *StreamEscalations) GetRecords(ctx context.Context, logger kitlog.Logger, cl *client.ClientWithResponses) ([]map[string]any, error) {
	var (
		results = []map[string]any{}
		after   *string
	)

	for {
		params := &client.EscalationsV2ListParams{
			PageSize: Ptr(int64(50)), // Max page size
			After:    after,
		}

		response, err := cl.EscalationsV2ListWithResponse(ctx, params)
		if err != nil {
			return nil, errors.Wrap(err, "listing escalations")
		}

		if response.StatusCode() != 200 {
			return nil, errors.Errorf("unexpected status code: %d", response.StatusCode())
		}

		for _, element := range response.JSON200.Escalations {
			results = append(results, model.EscalationV2.Serialize(element))
		}

		// Check if there are more pages
		if response.JSON200.PaginationMeta.After == nil {
			break
		}

		after = response.JSON200.PaginationMeta.After
	}

	return results, nil
}

// Ptr is a helper function to return a pointer to a value
func Ptr[T any](v T) *T {
	return &v
}