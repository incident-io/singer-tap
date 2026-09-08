package model

import incident "github.com/incident-io/sdk-go"

type alertSlimV2 struct{}

var AlertSlimV2 alertSlimV2

func (alertSlimV2) Schema() Property {
	return Property{
		Types: []string{"object"},
		Properties: map[string]Property{
			"id": {
				Types: []string{"string"},
			},
			"title": {
				Types: []string{"string"},
			},
		},
	}
}

func (alertSlimV2) Serialize(input incident.AlertSlimV2) map[string]any {
	return map[string]any{
		"id":    input.Id,
		"title": input.Title,
	}
}