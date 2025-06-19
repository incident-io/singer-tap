package model

import "github.com/incident-io/singer-tap/client"

type alertActorV2 struct{}

var AlertActorV2 alertActorV2

func (alertActorV2) Schema() Property {
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

func (alertActorV2) Serialize(input client.AlertActorV2) map[string]any {
	return map[string]any{
		"id":    input.Id,
		"title": input.Title,
	}
}