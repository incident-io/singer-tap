package model

import (
	"github.com/incident-io/singer-tap/client"
)

type alertAttributeV2 struct{}

var AlertAttributeV2 alertAttributeV2

func (alertAttributeV2) Schema() Property {
	return Property{
		Types: []string{"object"},
		Properties: map[string]Property{
			"id": {
				Types: []string{"string"},
			},
			"name": {
				Types: []string{"string"},
			},
			"array": {
				Types: []string{"boolean"},
			},
			"type": {
				Types: []string{"string"},
			},
		},
	}
}

func (alertAttributeV2) Serialize(input client.AlertAttributeV2) map[string]any {
	return map[string]any{
		"id":    input.Id,
		"name":  input.Name,
		"array": input.Array,
		"type":  input.Type,
	}
}