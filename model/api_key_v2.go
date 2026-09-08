package model

import incident "github.com/incident-io/sdk-go"

type apiKeyV2 struct{}

var APIKey apiKeyV2

func (apiKeyV2) Schema() Property {
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

func (apiKeyV2) Serialize(input incident.APIKeyActorV2) map[string]any {
	return DumpToMap(input)
}
