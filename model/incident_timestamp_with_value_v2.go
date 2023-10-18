package model

import "github.com/incident-io/singer-tap/client"

type incidentTimestampWithValueV2 struct{}

var IncidentTimestampWithValueV2 incidentTimestampWithValueV2

func (incidentTimestampWithValueV2) Schema() Property {
	return Property{
		Types: []string{"object"},
		Properties: map[string]Property{
			"incident_timestamp": IncidentTimestampV2.Schema(),
			"value":              Optional(IncidentTimestampValueV2.Schema()),
		},
	}
}

func (incidentTimestampWithValueV2) Serialize(input client.IncidentTimestampWithValueV2) map[string]any {
	return map[string]any{
		"incident_timestamp": IncidentTimestampV2.Serialize(input.IncidentTimestamp),
		"value":              IncidentTimestampValueV2.Serialize(input.Value),
	}
}
