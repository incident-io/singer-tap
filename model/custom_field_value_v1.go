package model

import "github.com/incident-io/singer-tap/client"

type customFieldValueV1 struct{}

var CustomFieldValueV1 customFieldValueV1

func (customFieldValueV1) Schema() Property {
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

func (customFieldValueV1) Serialize(input client.CustomFieldValueV1) map[string]any {
	return map[string]any{
		"value_link":          input.ValueLink,
		"value_numeric":       input.ValueNumeric,
		"value_text":          input.ValueText,
		"value_catalog_entry": EmbeddedCatalogEntryV1.Serialize(input.ValueCatalogEntry),
		"value_option":        CustomFieldOptionV1.Serialize(input.ValueOption),
	}
}
