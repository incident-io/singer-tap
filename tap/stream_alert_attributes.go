package tap

import (
	"context"

	kitlog "github.com/go-kit/log"
	"github.com/incident-io/singer-tap/client"
	"github.com/incident-io/singer-tap/model"
	"github.com/pkg/errors"
)

func init() {
	register(&StreamAlertAttributes{})
}

type StreamAlertAttributes struct{}

func (s *StreamAlertAttributes) Output() *Output {
	return &Output{
		Type:   OutputTypeSchema,
		Stream: "alert_attributes",
		Schema: &model.Schema{
			HasAdditionalProperties: false,
			Type:                    []string{"object"},
			Properties:              model.AlertAttributeV2.Schema().Properties,
		},
		KeyProperties:      []string{"id"},
		BookmarkProperties: []string{},
	}
}

func (s *StreamAlertAttributes) GetRecords(ctx context.Context, logger kitlog.Logger, cl *client.ClientWithResponses) ([]map[string]any, error) {
	var (
		results = []map[string]any{}
	)

	logger.Log("msg", "loading alert attributes")
	page, err := cl.AlertAttributesV2ListWithResponse(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "listing alert attributes")
	}

	for _, element := range page.JSON200.AlertAttributes {
		results = append(results, model.AlertAttributeV2.Serialize(element))
	}

	return results, nil
}
