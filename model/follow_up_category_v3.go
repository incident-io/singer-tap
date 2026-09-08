package model

import incident "github.com/incident-io/sdk-go"

type followUpCategoryV3 struct{}

var FollowUpCategoryV3 followUpCategoryV3

func (followUpCategoryV3) Schema() Property {
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

func (followUpCategoryV3) Serialize(input *incident.FollowUpCategoryV3) map[string]any {
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
