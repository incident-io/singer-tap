package model

import "github.com/incident-io/singer-tap/client"

type actionV2 struct{}

var ActionV2 actionV2

func (actionV2) Schema() Property {
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

func (actionV2) Serialize(input client.ActionV2) map[string]any {
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
