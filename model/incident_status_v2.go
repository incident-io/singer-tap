package model

import "github.com/incident-io/singer-tap/client"

type incidentStatusV2 struct{}

var IncidentStatusV2 incidentStatusV2

func (incidentStatusV2) Schema() Property {
	return Property{
		Types: []string{"object"},
		Properties: map[string]Property{
			"id": {
				Types: []string{"string"},
			},
			"name": {
				Types: []string{"string"},
			},
			"description": {
				Types: []string{"string"},
			},
			"category": {
				Types: []string{"string"},
			},
			"rank": {
				Types: []string{"integer"},
			},
			"created_at": DateTime.Schema(),
			"updated_at": DateTime.Schema(),
		},
	}
}

func (incidentStatusV2) Serialize(input client.IncidentStatusV2) map[string]any {
	return map[string]any{
		"id":          input.Id,
		"name":        input.Name,
		"description": input.Description,
		"category":    input.Category,
		"rank":        input.Rank,
		"created_at":  input.CreatedAt,
		"updated_at":  input.UpdatedAt,
	}
}