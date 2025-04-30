package model

import "github.com/incident-io/singer-tap/client"

type incidentTypeV2 struct{}

var IncidentTypeV2 incidentTypeV2

func (incidentTypeV2) Schema() Property {
	return Property{
		Types: []string{"object"},
		Properties: map[string]Property{
			"id": {
				Types: []string{"string"},
			},
			"name": {
				Types: []string{"string"},
			},
			"is_default": {
				Types: []string{"boolean"},
			},
			"description": {
				Types: []string{"string"},
			},
			"private_incidents_only": {
				Types: []string{"boolean"},
			},
			"create_in_triage": {
				Types: []string{"string"},
			},
			"created_at": DateTime.Schema(),
			"updated_at": DateTime.Schema(),
		},
	}
}

func (incidentTypeV2) Serialize(input *client.IncidentTypeV2) map[string]any {
	if input == nil {
		return nil
	}

	// Just flat convert everything into a map[string]any
	return DumpToMap(input)
}