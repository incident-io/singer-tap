package model

import incident "github.com/incident-io/sdk-go"

type escalationPriorityV2 struct{}

var EscalationPriorityV2 escalationPriorityV2

func (escalationPriorityV2) Schema() Property {
	return Property{
		Types: []string{"object"},
		Properties: map[string]Property{
			"name": {
				Types: []string{"string"},
			},
		},
	}
}

func (escalationPriorityV2) Serialize(input incident.EscalationPriorityV2) map[string]any {
	return map[string]any{
		"name": input.Name,
	}
}