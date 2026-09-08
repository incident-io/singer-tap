package model

import (
	incident "github.com/incident-io/sdk-go"
	"github.com/samber/lo"
)

type customFieldTypeInfoV2 struct{}

var CustomFieldTypeInfoV2 customFieldTypeInfoV2

func (customFieldTypeInfoV2) Schema() Property {
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

func (customFieldTypeInfoV2) Serialize(input incident.CustomFieldTypeInfoV2) map[string]any {
	return map[string]any{
		"id":          input.Id,
		"name":        input.Name,
		"description": input.Description,
		"field_type":  input.FieldType,
		"options": lo.Map(input.Options, func(option incident.CustomFieldOptionV2, _ int) map[string]any {
			return CustomFieldOptionV2.Serialize(&option)
		}),
	}
}