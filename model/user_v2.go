package model

import incident "github.com/incident-io/sdk-go"

type userV2 struct{}

var UserV2 userV2

func (userV2) Schema() Property {
	return Property{
		Types: []string{"object"},
		Properties: map[string]Property{
			"id": {
				Types: []string{"string"},
			},
			"name": {
				Types: []string{"string"},
			},
			"email": {
				Types: []string{"null", "string"},
			},
			"slack_user_id": {
				Types: []string{"null", "string"},
			},
		},
	}
}

func (userV2) Serialize(input incident.UserV2) map[string]any {
	return map[string]any{
		"id":            input.Id,
		"name":          input.Name,
		"email":         input.Email,
		"slack_user_id": input.SlackUserId,
	}
}
