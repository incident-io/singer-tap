package model

import "github.com/incident-io/singer-tap/client"

type incidentSlimV2 struct{}

var IncidentSlimV2 incidentSlimV2

func (incidentSlimV2) Schema() Property {
	return Property{
		Types: []string{"object"},
		Properties: map[string]Property{
			"id": {
				Types: []string{"string"},
			},
			"external_id": {
				Types: []string{"integer"},
			},
			"reference": {
				Types: []string{"string"},
			},
			"name": {
				Types: []string{"string"},
			},
		},
	}
}

func (incidentSlimV2) Serialize(input client.IncidentSlimV2) map[string]any {
	return map[string]any{
		"id":          input.Id,
		"external_id": input.ExternalId,
		"reference":   input.Reference,
		"name":        input.Name,
	}
}