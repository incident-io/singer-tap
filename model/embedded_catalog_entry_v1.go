package model

import (
	incident "github.com/incident-io/sdk-go"
)

type embeddedCatalogEntryV1 struct{}

var EmbeddedCatalogEntryV1 embeddedCatalogEntryV1

func (embeddedCatalogEntryV1) Schema() Property {
	return Property{
		Types: []string{"object"},
		Properties: map[string]Property{
			"id": {
				Types: []string{"string"},
			},
			"name": {
				Types: []string{"string"},
			},
			"aliases": {
				Types: []string{"array", "null"},
				Items: &ArrayItem{
					Type: "string",
				},
			},
			"external_id": {
				Types: []string{"string", "null"},
			},
		},
	}
}

func (embeddedCatalogEntryV1) Serialize(input *incident.EmbeddedCatalogEntryV1) map[string]any {
	if input == nil {
		return nil
	}

	return map[string]any{
		"id":          input.Id,
		"name":        input.Name,
		"aliases":     input.Aliases,
		"external_id": input.ExternalId,
	}
}
