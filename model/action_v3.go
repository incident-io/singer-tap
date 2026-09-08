package model

import incident "github.com/incident-io/sdk-go"

type actionV3 struct{}

var ActionV3 actionV3

func (actionV3) Schema() Property {
	return Property{
		Types: []string{"object"},
		Properties: map[string]Property{
			"id": {
				Types: []string{"string"},
			},
			"incident_id": {
				Types: []string{"string"},
			},
			"status": {
				Types: []string{"string"},
			},
			"description": {
				Types: []string{"string"},
			},
			"assignee":     Optional(UserV2.Schema()),
			"completed_at": Optional(DateTime.Schema()),
			"created_at":   DateTime.Schema(),
			"updated_at":   DateTime.Schema(),
		},
	}
}

func (actionV3) Serialize(input incident.ActionV3) map[string]any {
	var assignee map[string]any
	if input.Assignee != nil {
		assignee = UserV2.Serialize(*input.Assignee)
	}

	return map[string]any{
		"id":           input.Id,
		"incident_id":  input.IncidentId,
		"status":       input.Status,
		"description":  input.Description,
		"assignee":     assignee,
		"completed_at": input.CompletedAt,
		"created_at":   input.CreatedAt,
		"updated_at":   input.UpdatedAt,
	}
}
