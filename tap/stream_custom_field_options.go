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
	register(&StreamCustomFieldOptions{})
}

type StreamCustomFieldOptions struct {
}

func (s *StreamCustomFieldOptions) Output() *Output {
	return &Output{
		Type:   OutputTypeSchema,
		Stream: "custom_field_options",
		Schema: &model.Schema{
			HasAdditionalProperties: false,
			Type:                    []string{"object"},
			Properties:              model.CustomFieldOptionV1.Schema().Properties,
		},
		KeyProperties:      []string{"id"},
		BookmarkProperties: []string{},
	}
}

func (s *StreamCustomFieldOptions) GetRecords(ctx context.Context, logger kitlog.Logger, cl *client.ClientWithResponses) ([]map[string]any, error) {
	var (
		results = []map[string]any{}
	)

	// We need to go over all custom fields to build the options
	response, err := cl.CustomFieldsV2ListWithResponse(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "listing incidents")
	}

	for _, element := range response.JSON200.CustomFields {
		options, err := s.GetOptions(ctx, logger, cl, element.Id)
		if err != nil {
			return nil, errors.Wrap(err, "listing custom field options")
		}

		for _, option := range options {
			results = append(results, model.CustomFieldOptionV1.Serialize(&option))
		}
	}

	return results, nil
}

func (s *StreamCustomFieldOptions) GetOptions(
	ctx context.Context,
	logger kitlog.Logger,
	cl *client.ClientWithResponses,
	customFieldId string,
) ([]client.CustomFieldOptionV1, error) {
	var (
		after    *string
		pageSize = int64(250)
		results  = []client.CustomFieldOptionV1{}
	)

	for {
		page, err := cl.CustomFieldOptionsV1ListWithResponse(ctx, &client.CustomFieldOptionsV1ListParams{
			CustomFieldId: customFieldId,
			PageSize:      &pageSize,
			After:         after,
		})
		if err != nil {
			return nil, errors.Wrap(err, "listing custom field options")
		}

		results = append(results, page.JSON200.CustomFieldOptions...)

		if count := len(page.JSON200.CustomFieldOptions); count == 0 {
			return results, nil // end pagination
		} else {
			after = lo.ToPtr(page.JSON200.CustomFieldOptions[count-1].Id)
		}
	}
}
