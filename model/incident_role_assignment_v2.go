package model

import "github.com/incident-io/singer-tap/client"

type incidentRoleAssignmentV2 struct{}

var IncidentRoleAssignmentV2 incidentRoleAssignmentV2

func (incidentRoleAssignmentV2) Schema() Property {
	return Property{
		Types: []string{"object"},
		Properties: map[string]Property{
			"assignee": Optional(UserV2.Schema()),
			"role":     IncidentRoleV2.Schema(),
		},
	}
}

func (incidentRoleAssignmentV2) Serialize(input client.IncidentRoleAssignmentV2) map[string]any {
	var assignee map[string]any
	if input.Assignee != nil {
		assignee = UserV2.Serialize(*input.Assignee)
	}

	return map[string]any{
		"assignee": assignee,
		"role":     EmbeddedIncidentRoleV2.Serialize(input.Role),
	}
}