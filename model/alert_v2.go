package model

import (
	"github.com/incident-io/singer-tap/client"
	"github.com/samber/lo"
)

type alertV2 struct{}

var AlertV2 alertV2

func (alertV2) Schema() Property {
	return Property{
		Types: []string{"object"},
		Properties: map[string]Property{
			"id": {
				Types: []string{"string"},
			},
			"alert_source_id": {
				Types: []string{"string"},
			},
			"attributes": ArrayOf(AlertAttributeEntryV2.Schema()),
			"created_at": DateTime.Schema(),
			"deduplication_key": {
				Types: []string{"string"},
			},
			"description": {
				Types: []string{"string", "null"},
			},
			"resolved_at": Optional(DateTime.Schema()),
			"source_url": {
				Types: []string{"string", "null"},
			},
			"status": {
				Types: []string{"string"},
			},
			"title": {
				Types: []string{"string"},
			},
		},
	}
}

func (alertV2) Serialize(input client.AlertV2) map[string]any {
	attributes := []map[string]any{}
	if len(input.Attributes) > 0 {
		attributes = lo.Map(input.Attributes, func(entry client.AlertAttributeEntryV2, _ int) map[string]any {
			return AlertAttributeEntryV2.Serialize(entry)
		})
	}

	result := map[string]any{
		"id":                input.Id,
		"alert_source_id":   input.AlertSourceId,
		"attributes":        attributes,
		"created_at":        input.CreatedAt,
		"deduplication_key": input.DeduplicationKey,
		"status":            input.Status,
		"title":             input.Title,
	}

	if input.Description != nil {
		result["description"] = *input.Description
	}

	if input.ResolvedAt != nil {
		result["resolved_at"] = *input.ResolvedAt
	}

	if input.SourceUrl != nil {
		result["source_url"] = *input.SourceUrl
	}

	return result
}