package model

import (
	"github.com/incident-io/singer-tap/client"
	"github.com/samber/lo"
)

type customFieldEntryV2 struct{}

var CustomFieldEntryV2 customFieldEntryV2

func (customFieldEntryV2) Schema() Property {
	return Property{
		Types: []string{"object"},
		Properties: map[string]Property{
			"custom_field": CustomFieldTypeInfoV1.Schema(),
			"values":       ArrayOf(CustomFieldValueV1.Schema()),
		},
	}
}

func (customFieldEntryV2) Serialize(input client.CustomFieldEntryV2) map[string]any {
	return map[string]any{
		"custom_field": CustomFieldTypeInfoV2.Serialize(input.CustomField),
		"values": lo.Map(input.Values, func(value client.CustomFieldValueV2, _ int) map[string]any {
			return CustomFieldValueV2.Serialize(value)
		}),
	}
}
