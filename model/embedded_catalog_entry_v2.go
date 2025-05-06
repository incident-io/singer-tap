package model

import (
	"github.com/incident-io/singer-tap/client"
)

type embeddedCatalogEntryV2 struct{}

var EmbeddedCatalogEntryV2 embeddedCatalogEntryV2

func (embeddedCatalogEntryV2) Schema() Property {
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

func (embeddedCatalogEntryV2) Serialize(input *client.EmbeddedCatalogEntryV2) map[string]any {
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