package model

import (
	"github.com/incident-io/singer-tap/client"
)

type alertAttributeValueV2 struct{}

var AlertAttributeValueV2 alertAttributeValueV2

func (alertAttributeValueV2) Schema() Property {
	return Property{
		Types: []string{"object"},
		Properties: map[string]Property{
			"catalog_entry": Optional(Property{
				Types: []string{"object", "null"},
				Properties: map[string]Property{
					"id": {
						Types: []string{"string"},
					},
					"name": {
						Types: []string{"string"},
					},
				},
			}),
			"label": {
				Types: []string{"string", "null"},
			},
			"literal": {
				Types: []string{"string", "null"},
			},
		},
	}
}

func (alertAttributeValueV2) Serialize(input client.AlertAttributeValueV2) map[string]any {
	result := map[string]any{}
	
	if input.CatalogEntry != nil {
		result["catalog_entry"] = map[string]any{
			"id":   input.CatalogEntry.Id,
			"name": input.CatalogEntry.Name,
		}
	}
	
	if input.Label != nil {
		result["label"] = *input.Label
	}
	
	if input.Literal != nil {
		result["literal"] = *input.Literal
	}
	
	return result
}