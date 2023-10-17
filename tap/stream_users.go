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
	register(&StreamUsers{})
}

type StreamUsers struct {
}

func (s *StreamUsers) Output() *Output {
	return &Output{
		Type:   OutputTypeSchema,
		Stream: "users",
		Schema: &model.Schema{
			HasAdditionalProperties: false,
			Type:                    []string{"object"},
			Properties:              model.UserV1.Schema().Properties,
		},
		KeyProperties:      []string{"id"},
		BookmarkProperties: []string{},
	}
}

func (s *StreamUsers) GetRecords(ctx context.Context, logger kitlog.Logger, cl *client.ClientWithResponses) ([]map[string]any, error) {
	var (
		after    *string
		pageSize = 250
		results  = []map[string]any{}
	)

	for {
		logger.Log("msg", "loading page", "page_size", pageSize, "after", after)
		page, err := cl.UsersV2ListWithResponse(ctx, &client.UsersV2ListParams{
			PageSize: &pageSize,
			After:    after,
		})
		if err != nil {
			return nil, errors.Wrap(err, "listing incidents")
		}

		for _, element := range page.JSON200.Users {
			results = append(results, model.UserV1.Serialize(element))
		}
		if count := len(page.JSON200.Users); count == 0 {
			return results, nil // end pagination
		} else {
			after = lo.ToPtr(page.JSON200.Users[count-1].Id)
		}
	}
}
