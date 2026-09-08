package model

import incident "github.com/incident-io/sdk-go"

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

func (workflowActorV2) Serialize(input incident.WorkflowActorV2) map[string]any {
	return map[string]any{
		"id":   input.Id,
		"name": input.Name,
	}
}