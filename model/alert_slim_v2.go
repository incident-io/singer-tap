package model

import "github.com/incident-io/singer-tap/client"

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

func (alertSlimV2) Serialize(input client.AlertSlimV2) map[string]any {
	return map[string]any{
		"id":    input.Id,
		"title": input.Title,
	}
}