package model

import "github.com/incident-io/singer-tap/client"

type incidentUpdateV2 struct{}

var IncidentUpdateV2 incidentUpdateV2

func (incidentUpdateV2) Schema() Property {
	return Property{
		Types: []string{"object"},
		Properties: map[string]Property{
			"id": {
				Types: []string{"string"},
			},
			"incident_id": {
				Types: []string{"string"},
			},
			"message": {
				Types: []string{"string", "null"},
			},
			"new_incident_status": IncidentStatusV1.Schema(),
			"new_severity":        Optional(SeverityV2.Schema()),
			"updater":             ActorV2.Schema(),
			"created_at":          DateTime.Schema(),
		},
	}
}

func (incidentUpdateV2) Serialize(input client.IncidentUpdateV2) map[string]any {
	var severity map[string]any
	if input.NewSeverity != nil {
		severity = SeverityV2.Serialize(input.NewSeverity)
	}

	return map[string]any{
		"id":                  input.Id,
		"incident_id":         input.IncidentId,
		"message":             input.Message,
		"new_incident_status": IncidentStatusV1.Serialize(input.NewIncidentStatus),
		"new_severity":        severity,
		"updater":             ActorV2.Serialize(input.Updater),
		"created_at":          input.CreatedAt,
	}
}
