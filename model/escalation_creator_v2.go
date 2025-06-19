package model

import "github.com/incident-io/singer-tap/client"

type escalationCreatorV2 struct{}

var EscalationCreatorV2 escalationCreatorV2

func (escalationCreatorV2) Schema() Property {
	return Property{
		Types: []string{"object"},
		Properties: map[string]Property{
			"user": Optional(UserV2.Schema()),
			"alert":    Optional(AlertActorV2.Schema()),
			"workflow": Optional(WorkflowActorV2.Schema()),
		},
	}
}

func (escalationCreatorV2) Serialize(input client.EscalationCreatorV2) map[string]any {
	result := make(map[string]any)

	if input.User != nil {
		result["user"] = UserV2.Serialize(*input.User)
	}

	if input.Alert != nil {
		result["alert"] = AlertActorV2.Serialize(*input.Alert)
	}

	if input.Workflow != nil {
		result["workflow"] = WorkflowActorV2.Serialize(*input.Workflow)
	}

	return result
}