package tap

import (
	"context"

	kitlog "github.com/go-kit/log"
	"github.com/incident-io/singer-tap/client"
	"github.com/incident-io/singer-tap/model"
	"github.com/pkg/errors"
)

func init() {
	register(&StreamAlertSources{})
}

type StreamAlertSources struct{}

func (s *StreamAlertSources) Output() *Output {
	return &Output{
		Type:   OutputTypeSchema,
		Stream: "alert_sources",
		Schema: &model.Schema{
			HasAdditionalProperties: false,
			Type:                    []string{"object"},
			Properties:              model.AlertSourceV2.Schema().Properties,
		},
		KeyProperties:      []string{"id"},
		BookmarkProperties: []string{},
	}
}

func (s *StreamAlertSources) GetRecords(ctx context.Context, logger kitlog.Logger, cl *client.ClientWithResponses) ([]map[string]any, error) {
	var (
		results = []map[string]any{}
	)

	logger.Log("msg", "loading alert sources")
	page, err := cl.AlertSourcesV2ListWithResponse(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "listing alert sources")
	}

	for _, element := range page.JSON200.AlertSources {
		results = append(results, model.AlertSourceV2.Serialize(element))
	}

	return results, nil
}