package model

import (
	incident "github.com/incident-io/sdk-go"
	"github.com/samber/lo"
)

type customFieldEntryV1 struct{}

var CustomFieldEntryV1 customFieldEntryV1

func (customFieldEntryV1) Schema() Property {
	return Property{
		Types: []string{"object"},
		Properties: map[string]Property{
			"custom_field": CustomFieldTypeInfoV1.Schema(),
			"values":       ArrayOf(CustomFieldValueV1.Schema()),
		},
	}
}

func (customFieldEntryV1) Serialize(input incident.CustomFieldEntryV1) map[string]any {
	return map[string]any{
		"custom_field": CustomFieldTypeInfoV1.Serialize(input.CustomField),
		"values": lo.Map(input.Values, func(value incident.CustomFieldValueV1, _ int) map[string]any {
			return CustomFieldValueV1.Serialize(value)
		}),
	}
}
