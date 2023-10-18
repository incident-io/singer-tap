package model

import "github.com/incident-io/singer-tap/client"

type incidentStatusV1 struct{}

var IncidentStatusV1 incidentStatusV1

func (incidentStatusV1) Schema() Property {
	return Property{
		Types: []string{"object"},
		Properties: map[string]Property{
			"id": {
				Types: []string{"string"},
			},
			"name": {
				Types: []string{"string"},
			},
			"category": {
				Types: []string{"string"},
			},
			"description": {
				Types: []string{"string"},
			},
			"rank": {
				Types: []string{"integer"},
			},
			"created_at": DateTime.Schema(),
			"updated_at": DateTime.Schema(),
		},
	}
}

func (incidentStatusV1) Serialize(input client.IncidentStatusV1) map[string]any {
	// Just flat convert everything into a map[string]any
	return DumpToMap(input)
}
