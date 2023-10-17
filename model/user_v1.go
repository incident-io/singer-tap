package model

import "github.com/incident-io/singer-tap/client"

type userV1 struct{}

var UserV1 userV1

func (userV1) Schema() Property {
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

func (userV1) Serialize(input client.UserV1) map[string]any {
	// Deprecated role field needs removing - so build manually and omit it
	return map[string]any{
		"id":            input.Id,
		"name":          input.Name,
		"email":         input.Email,
		"slack_user_id": input.SlackUserId,
	}
}
