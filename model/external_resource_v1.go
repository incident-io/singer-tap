package model

import "github.com/incident-io/singer-tap/client"

type externalResourceV1 struct{}

var ExternalResourceV1 externalResourceV1

func (externalResourceV1) Schema() Property {
	return Property{
		Types: []string{"object"},
		Properties: map[string]Property{
			"external_id": {
				Types: []string{"string"},
			},
			"permalink": {
				Types: []string{"string"},
			},
			"resource_type": {
				Types: []string{"string"},
			},
			"title": {
				Types: []string{"string"},
			},
		},
	}
}

func (externalResourceV1) Serialize(input client.ExternalResourceV1) map[string]any {
	return DumpToMap(input)
}
