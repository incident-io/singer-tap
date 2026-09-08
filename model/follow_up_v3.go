package model

import incident "github.com/incident-io/sdk-go"

type followUpV3 struct{}

var FollowUpV3 followUpV3

func (followUpV3) Schema() Property {
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
			"category": Optional(FollowUpCategoryV3.Schema()),
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

func (followUpV3) Serialize(input incident.FollowUpV3) map[string]any {
	var external_issue_reference map[string]any
	if input.ExternalIssueReference != nil {
		external_issue_reference = ExternalIssueReferenceV2.Serialize(input.ExternalIssueReference)
	}

	var assignee map[string]any
	if input.Assignee != nil {
		assignee = UserV2.Serialize(*input.Assignee)
	}

	return map[string]any{
		"assignee":                 assignee,
		"id":                       input.Id,
		"incident_id":              input.IncidentId,
		"priority":                 FollowUpPriorityV2.Serialize(input.Priority),
		"category":                 FollowUpCategoryV3.Serialize(input.Category),
		"status":                   input.Status,
		"title":                    input.Title,
		"description":              input.Description,
		"external_issue_reference": external_issue_reference,
		"completed_at":             input.CompletedAt,
		"created_at":               input.CreatedAt,
		"updated_at":               input.UpdatedAt,
	}
}
