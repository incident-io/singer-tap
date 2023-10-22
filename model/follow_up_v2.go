package model

import "github.com/incident-io/singer-tap/client"

type followUpV2 struct{}

var FollowUpV2 followUpV2

func (followUpV2) Schema() Property {
	return Property{
		Types: []string{"object"},
		Properties: map[string]Property{
			"assignee": Optional(UserV1.Schema()),
			"id": {
				Types: []string{"string"},
			},
			"incident_id": {
				Types: []string{"string"},
			},
			"priority": Optional(FollowUpPriorityV2.Schema()),
			"status": {
				Types: []string{"string"},
			},
			"title": {
				Types: []string{"string"},
			},
			"description": {
				Types: []string{"string", "null"},
			},
			"external_issue_reference": Optional(ExternalIssueReferenceV2.Schema()),
			"completed_at":             Optional(DateTime.Schema()),
			"created_at":               DateTime.Schema(),
			"updated_at":               DateTime.Schema(),
		},
	}
}

func (followUpV2) Serialize(input client.FollowUpV2) map[string]any {
	var external_issue_reference map[string]any
	if input.ExternalIssueReference != nil {
		external_issue_reference = ExternalIssueReferenceV2.Serialize(input.ExternalIssueReference)
	}

	var assignee map[string]any
	if input.Assignee != nil {
		assignee = UserV1.Serialize(*input.Assignee)
	}

	return map[string]any{
		"assignee":                 assignee,
		"id":                       input.Id,
		"incident_id":              input.IncidentId,
		"priority":                 FollowUpPriorityV2.Serialize(input.Priority),
		"status":                   input.Status,
		"title":                    input.Title,
		"description":              input.Description,
		"external_issue_reference": external_issue_reference,
		"completed_at":             input.CompletedAt,
		"created_at":               input.CreatedAt,
		"updated_at":               input.UpdatedAt,
	}
}
