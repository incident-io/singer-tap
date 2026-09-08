package model

import incident "github.com/incident-io/sdk-go"

type incidentTimestampV2 struct{}

var IncidentTimestampV2 incidentTimestampV2

func (incidentTimestampV2) Schema() Property {
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
		},
	}
}

func (incidentTimestampV2) Serialize(input incident.IncidentTimestampV2) map[string]any {
	// Just flat convert everything into a map[string]any
	return DumpToMap(input)
}
