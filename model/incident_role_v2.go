package model

import "github.com/incident-io/singer-tap/client"

type incidentRoleV2 struct{}

var IncidentRoleV2 incidentRoleV2

// Actually identical to V1 but hey
func (incidentRoleV2) Schema() Property {
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
			"shortform": {
				Types: []string{"string"},
			},
			"created_at": DateTime.Schema(),
			"updated_at": DateTime.Schema(),
		},
	}
}

func (incidentRoleV2) Serialize(input client.IncidentRoleV2) map[string]any {
	// Just flat convert everything into a map[string]any
	return DumpToMap(input)
}
