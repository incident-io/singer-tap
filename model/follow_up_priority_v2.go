package model

import "github.com/incident-io/singer-tap/client"

type followUpPriorityV2 struct{}

var FollowUpPriorityV2 followUpPriorityV2

func (followUpPriorityV2) Schema() Property {
	return Property{
		Types: []string{"object"},
		Properties: map[string]Property{
			"id": {
				Types: []string{"string"},
			},
			"name": {
				Types: []string{"string"},
			},
			"rank": {
				Types: []string{"integer"},
			},
			"description": {
				Types: []string{"string", "null"},
			},
		},
	}
}

func (followUpPriorityV2) Serialize(input *client.FollowUpPriorityV2) map[string]any {
	if input == nil {
		return nil
	}

	return map[string]any{
		"id":          input.Id,
		"name":        input.Name,
		"rank":        input.Rank,
		"description": input.Description,
	}
}
