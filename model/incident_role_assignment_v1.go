package model

import "github.com/incident-io/singer-tap/client"

type incidentRoleAssignmentV1 struct{}

var IncidentRoleAssignmentV1 incidentRoleAssignmentV1

func (incidentRoleAssignmentV1) Schema() Property {
	return Property{
		Types: []string{"object"},
		Properties: map[string]Property{
			"assignee": Optional(UserV1.Schema()),
			"role":     IncidentRoleV1.Schema(),
		},
	}
}

func (incidentRoleAssignmentV1) Serialize(input client.IncidentRoleAssignmentV1) map[string]any {
	var assignee map[string]any
	if input.Assignee != nil {
		assignee = UserV1.Serialize(*input.Assignee)
	}

	return map[string]any{
		"assignee": assignee,
		"role":     IncidentRoleV1.Serialize(input.Role),
	}
}
