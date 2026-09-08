package tap

import (
	"context"

	kitlog "github.com/go-kit/log"
	incident "github.com/incident-io/sdk-go"
	"github.com/incident-io/singer-tap/model"
	"github.com/pkg/errors"
)

func init() {
	register(&StreamCustomFields{})
}

type StreamCustomFields struct {
}

func (s *StreamCustomFields) Output() *Output {
	return &Output{
		Type:   OutputTypeSchema,
		Stream: "custom_fields",
		Schema: &model.Schema{
			HasAdditionalProperties: false,
			Type:                    []string{"object"},
			Properties:              model.CustomFieldV2.Schema().Properties,
		},
		KeyProperties:      []string{"id"},
		BookmarkProperties: []string{},
	}
}

func (s *StreamCustomFields) GetRecords(ctx context.Context, logger kitlog.Logger, cl *incident.ClientWithResponses) ([]map[string]any, error) {
	var (
		results = []map[string]any{}
	)

	response, err := cl.CustomFieldsV2ListWithResponse(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "listing custom fields")
	}

	for _, element := range response.JSON200.CustomFields {
		results = append(results, model.CustomFieldV2.Serialize(element))
	}

	return results, nil
}
