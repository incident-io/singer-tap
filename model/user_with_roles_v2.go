package model

import "github.com/incident-io/singer-tap/client"

type userWithRolesV2 struct{}

var UserWithRolesV2 userWithRolesV2

func (userWithRolesV2) Schema() Property {
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

func (userWithRolesV2) Serialize(input client.UserWithRolesV2) map[string]any {
	return map[string]any{
		"id":            input.Id,
		"name":          input.Name,
		"email":         input.Email,
		"slack_user_id": input.SlackUserId,
	}
}