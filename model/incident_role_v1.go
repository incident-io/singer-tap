package model

import "github.com/incident-io/singer-tap/client"

type incidentRoleV1 struct{}

var IncidentRoleV1 incidentRoleV1

func (incidentRoleV1) Schema() Property {
	return Property{
		Types: []string{"object"},
		Properties: map[string]Property{
			"id": {
				Types: []string{"string"},
			},
			"name": {
				Types: []string{"string"},
			},
			"description": {
				Types: []string{"string"},
			},
			"instructions": {
				Types: []string{"string"},
			},
			"required": {
				Types: []string{"boolean"},
			},
			"role_type": {
				Types: []string{"string"},
			},
			"short_form": {
				Types: []string{"string"},
			},
			"created_at": DateTime.Schema(),
			"updated_at": DateTime.Schema(),
		},
	}
}

func (incidentRoleV1) Serialize(input client.IncidentRoleV1) map[string]any {
	// Just flat convert everything into a map[string]any
	return DumpToMap(input)
}
