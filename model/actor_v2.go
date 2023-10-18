package model

import "github.com/incident-io/singer-tap/client"

type actorV2 struct{}

var ActorV2 actorV2

func (actorV2) Schema() Property {
	return Property{
		Types: []string{"object"},
		Properties: map[string]Property{
			"api_key": Optional(APIKey.Schema()),
			"user":    Optional(UserV1.Schema()),
		},
	}
}

func (actorV2) Serialize(input client.ActorV2) map[string]any {
	var user map[string]any
	if input.User != nil {
		user = UserV1.Serialize(*input.User)
	}

	var apiKey map[string]any
	if input.ApiKey != nil {
		apiKey = APIKey.Serialize(*input.ApiKey)
	}

	return map[string]any{
		"api_key": apiKey,
		"user":    user,
	}
}
