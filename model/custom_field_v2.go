package model

import "github.com/incident-io/singer-tap/client"

type customFieldV2 struct{}

var CustomFieldV2 customFieldV2

func (customFieldV2) Schema() Property {
	return Property{
		Types: []string{"object"},
		Properties: map[string]Property{
			"id": {
				Types: []string{"string"},
			},
			"name": {
				Types: []string{"string"},
			},
			"field_type": {
				Types: []string{"string"},
			},
			"description": {
				Types: []string{"string"},
			},
			"created_at": DateTime.Schema(),
			"updated_at": DateTime.Schema(),
		},
	}

}

func (customFieldV2) Serialize(input client.CustomFieldV2) map[string]any {
	return map[string]any{
		"id":          input.Id,
		"name":        input.Name,
		"field_type":  input.FieldType,
		"description": input.Description,
		"created_at":  input.CreatedAt,
		"updated_at":  input.UpdatedAt,
	}
}
