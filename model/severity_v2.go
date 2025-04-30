package model

import "github.com/incident-io/singer-tap/client"

type severityV1 struct{}

var SeverityV1 severityV1

func (severityV1) Schema() Property {
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
			"rank": {
				Types: []string{"integer"},
			},
			"created_at": DateTime.Schema(),
			"updated_at": DateTime.Schema(),
		},
	}
}

func (severityV1) Serialize(input *client.SeverityV1) map[string]any {
	if input == nil {
		return nil
	}

	return DumpToMap(input)
}

type severityV2 struct{}

var SeverityV2 severityV2

func (severityV2) Schema() Property {
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
			"rank": {
				Types: []string{"integer"},
			},
			"created_at": DateTime.Schema(),
			"updated_at": DateTime.Schema(),
		},
	}
}

func (severityV2) Serialize(input *client.SeverityV2) map[string]any {
	if input == nil {
		return nil
	}

	return DumpToMap(input)
}
