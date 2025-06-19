package model

import "github.com/incident-io/singer-tap/client"

type workflowActorV2 struct{}

var WorkflowActorV2 workflowActorV2

func (workflowActorV2) Schema() Property {
	return Property{
		Types: []string{"object"},
		Properties: map[string]Property{
			"id": {
				Types: []string{"string"},
			},
			"name": {
				Types: []string{"string"},
			},
		},
	}
}

func (workflowActorV2) Serialize(input client.WorkflowActorV2) map[string]any {
	return map[string]any{
		"id":   input.Id,
		"name": input.Name,
	}
}