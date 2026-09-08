package tap

import (
	"context"

	kitlog "github.com/go-kit/log"
	incident "github.com/incident-io/sdk-go"
	"github.com/incident-io/singer-tap/model"
	"github.com/pkg/errors"
)

func init() {
	register(&StreamIncidentRoles{})
}

type StreamIncidentRoles struct {
}

func (s *StreamIncidentRoles) Output() *Output {
	return &Output{
		Type:   OutputTypeSchema,
		Stream: "incident_roles",
		Schema: &model.Schema{
			HasAdditionalProperties: false,
			Type:                    []string{"object"},
			Properties:              model.IncidentRoleV2.Schema().Properties,
		},
		KeyProperties:      []string{"id"},
		BookmarkProperties: []string{},
	}
}

func (s *StreamIncidentRoles) GetRecords(ctx context.Context, logger kitlog.Logger, cl *incident.ClientWithResponses) ([]map[string]any, error) {
	var (
		results = []map[string]any{}
	)

	response, err := cl.IncidentRolesV2ListWithResponse(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "listing incident roles")
	}

	for _, element := range response.JSON200.IncidentRoles {
		results = append(results, model.IncidentRoleV2.Serialize(element))
	}

	return results, nil
}
