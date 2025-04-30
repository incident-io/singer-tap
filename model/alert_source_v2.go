package model

import (
	"github.com/incident-io/singer-tap/client"
)

type alertSourceV2 struct{}

var AlertSourceV2 alertSourceV2

func (alertSourceV2) Schema() Property {
	return Property{
		Types: []string{"object"},
		Properties: map[string]Property{
			"id": {
				Types: []string{"string"},
			},
			"name": {
				Types: []string{"string"},
			},
			"source_type": {
				Types: []string{"string"},
			},
			"email_options": Optional(Property{
				Types: []string{"object", "null"},
				Properties: map[string]Property{
					"email_address": {
						Types: []string{"string"},
					},
				},
			}),
			"jira_options": Optional(Property{
				Types: []string{"object", "null"},
				Properties: map[string]Property{
					"project_ids": ArrayOf(Property{
						Types: []string{"string"},
					}),
				},
			}),
		},
	}
}

func (alertSourceV2) Serialize(input client.AlertSourceV2) map[string]any {
	result := map[string]any{
		"id":          input.Id,
		"name":        input.Name,
		"source_type": input.SourceType,
	}

	if input.EmailOptions != nil {
		result["email_options"] = map[string]any{
			"email_address": input.EmailOptions.EmailAddress,
		}
	}

	if input.JiraOptions != nil {
		result["jira_options"] = map[string]any{
			"project_ids": input.JiraOptions.ProjectIds,
		}
	}

	return result
}
