package model

import (
	"github.com/incident-io/singer-tap/client"
	"github.com/samber/lo"
)

type alertAttributeEntryV2 struct{}

var AlertAttributeEntryV2 alertAttributeEntryV2

func (alertAttributeEntryV2) Schema() Property {
	return Property{
		Types: []string{"object"},
		Properties: map[string]Property{
			"attribute": AlertAttributeV2.Schema(),
			"value": Optional(AlertAttributeValueV2.Schema()),
			"array_value": Optional(ArrayOf(AlertAttributeValueV2.Schema())),
		},
	}
}

func (alertAttributeEntryV2) Serialize(input client.AlertAttributeEntryV2) map[string]any {
	result := map[string]any{
		"attribute": AlertAttributeV2.Serialize(input.Attribute),
	}

	if input.Value != nil {
		result["value"] = AlertAttributeValueV2.Serialize(*input.Value)
	}

	if input.ArrayValue != nil && len(*input.ArrayValue) > 0 {
		result["array_value"] = lo.Map(*input.ArrayValue, func(value client.AlertAttributeValueV2, _ int) map[string]any {
			return AlertAttributeValueV2.Serialize(value)
		})
	}

	return result
}