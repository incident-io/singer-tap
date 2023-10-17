package model

import (
	"github.com/incident-io/singer-tap/client"
	"github.com/samber/lo"
)

type customFieldTypeInfoV1 struct{}

var CustomFieldTypeInfoV1 customFieldTypeInfoV1

func (customFieldTypeInfoV1) Schema() Property {
	return Property{
		Types: []string{"object"},
		Properties: map[string]Property{
			"id": {
				Types: []string{"string"},
			},
			"name": {
				Types: []string{"string"},
			},
			"description": {
				Types: []string{"string"},
			},
			"field_type": {
				Types: []string{"string"},
			},
			"options": ArrayOf(CustomFieldOptionV1.Schema()),
		},
	}
}

func (customFieldTypeInfoV1) Serialize(input client.CustomFieldTypeInfoV1) map[string]any {
	return map[string]any{
		"id":          input.Id,
		"name":        input.Name,
		"description": input.Description,
		"field_type":  input.FieldType,
		"options": lo.Map(input.Options, func(option client.CustomFieldOptionV1, _ int) map[string]any {
			return CustomFieldOptionV1.Serialize(&option)
		}),
	}
}
