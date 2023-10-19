package model

import "github.com/incident-io/singer-tap/client"

type incidentAttachmentV1 struct{}

var IncidentAttachmentV1 incidentAttachmentV1

func (incidentAttachmentV1) Schema() Property {
	return Property{
		Types: []string{"object"},
		Properties: map[string]Property{
			"id": {
				Types: []string{"string"},
			},
			"incident_id": {
				Types: []string{"string"},
			},
			"resource": ExternalResourceV1.Schema(),
		},
	}
}

func (incidentAttachmentV1) Serialize(input client.IncidentAttachmentV1) map[string]any {
	return map[string]any{
		"id":          input.Id,
		"incident_id": input.IncidentId,
		"resource":    ExternalResourceV1.Serialize(input.Resource),
	}
}
