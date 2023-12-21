package model

import "github.com/incident-io/singer-tap/client"

type incidentTimestampValueV2 struct{}

var IncidentTimestampValueV2 incidentTimestampValueV2

func (incidentTimestampValueV2) Schema() Property {
	return Property{
		Types: []string{"object"},
		Properties: map[string]Property{
			"value": Optional(DateTime.Schema()),
		},
	}
}

func (incidentTimestampValueV2) Serialize(input *client.IncidentTimestampValueV2) map[string]any {
	if input == nil {
		return nil
	}

	return map[string]any{
		"value": input.Value,
	}
}
