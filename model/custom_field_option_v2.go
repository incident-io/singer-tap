package model

import "github.com/incident-io/singer-tap/client"

type customFieldOptionV2 struct{}

var CustomFieldOptionV2 customFieldOptionV2

func (customFieldOptionV2) Schema() Property {
	return Property{
		Types: []string{"object"},
		Properties: map[string]Property{
			"id": {
				Types: []string{"string"},
			},
			"custom_field_id": {
				Types: []string{"string"},
			},
			"sort_key": {
				Types: []string{"integer"},
			},
			"value": {
				Types: []string{"string"},
			},
		},
	}
}

func (customFieldOptionV2) Serialize(input *client.CustomFieldOptionV2) map[string]any {
	if input == nil {
		return nil
	}

	return DumpToMap(input)
}