package model

import "github.com/incident-io/singer-tap/client"

type escalationV2 struct{}

var EscalationV2 escalationV2

func (escalationV2) Schema() Property {
	return Property{
		Types: []string{"object"},
		Properties: map[string]Property{
			"id": {
				Types: []string{"string"},
			},
			"title": {
				Types: []string{"string"},
			},
			"status": {
				Types: []string{"string"},
			},
			"escalation_path_id": Optional(Property{
				Types: []string{"string"},
			}),
			"creator":           EscalationCreatorV2.Schema(),
			"priority":          EscalationPriorityV2.Schema(),
			"events":            ArrayOf(EscalationEventV2.Schema()),
			"created_at":        DateTime.Schema(),
			"related_alerts":    ArrayOf(AlertSlimV2.Schema()),
			"related_incidents": ArrayOf(IncidentSlimV2.Schema()),
		},
	}
}

func (escalationV2) Serialize(input client.EscalationV2) map[string]any {
	relatedAlerts := make([]map[string]any, 0, len(input.RelatedAlerts))
	for _, alert := range input.RelatedAlerts {
		relatedAlerts = append(relatedAlerts, AlertSlimV2.Serialize(alert))
	}

	relatedIncidents := make([]map[string]any, 0, len(input.RelatedIncidents))
	for _, incident := range input.RelatedIncidents {
		relatedIncidents = append(relatedIncidents, IncidentSlimV2.Serialize(incident))
	}

	events := make([]map[string]any, 0, len(input.Events))
	for _, event := range input.Events {
		events = append(events, EscalationEventV2.Serialize(event))
	}

	return map[string]any{
		"id":                 input.Id,
		"title":              input.Title,
		"status":             input.Status,
		"escalation_path_id": input.EscalationPathId,
		"creator":            EscalationCreatorV2.Serialize(input.Creator),
		"priority":           EscalationPriorityV2.Serialize(input.Priority),
		"events":             events,
		"created_at":         input.CreatedAt,
		"related_alerts":     relatedAlerts,
		"related_incidents":  relatedIncidents,
	}
}
