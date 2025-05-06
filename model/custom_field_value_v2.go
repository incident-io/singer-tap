package model

import "github.com/incident-io/singer-tap/client"

type customFieldValueV2 struct{}

var CustomFieldValueV2 customFieldValueV2

func (customFieldValueV2) Schema() Property {
	return Property{
		Types: []string{"object"},
		Properties: map[string]Property{
			"value_link": {
				Types: []string{"null", "string"},
			},
			"value_numeric": {
				Types: []string{"null", "number"},
			},
			"value_text": {
				Types: []string{"null", "string"},
			},
			"value_catalog_entry": Optional(EmbeddedCatalogEntryV1.Schema()),
			"value_option":        Optional(CustomFieldOptionV1.Schema()),
		},
	}
}

func (customFieldValueV2) Serialize(input client.CustomFieldValueV2) map[string]any {
	return map[string]any{
		"value_link":          input.ValueLink,
		"value_numeric":       input.ValueNumeric,
		"value_text":          input.ValueText,
		"value_catalog_entry": EmbeddedCatalogEntryV2.Serialize(input.ValueCatalogEntry),
		"value_option":        CustomFieldOptionV2.Serialize(input.ValueOption),
	}
}