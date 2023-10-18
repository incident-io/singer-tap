package model

import "github.com/incident-io/singer-tap/client"

type customFieldOptionV1 struct{}

var CustomFieldOptionV1 customFieldOptionV1

func (customFieldOptionV1) Schema() Property {
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

func (customFieldOptionV1) Serialize(input *client.CustomFieldOptionV1) map[string]any {
	if input == nil {
		return nil
	}

	return DumpToMap(input)
}
