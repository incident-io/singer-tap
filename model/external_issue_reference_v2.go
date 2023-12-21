package model

import "github.com/incident-io/singer-tap/client"

type externalIssueReferenceV2 struct{}

var ExternalIssueReferenceV2 externalIssueReferenceV2

func (externalIssueReferenceV2) Schema() Property {
	return Property{
		Types: []string{"object"},
		Properties: map[string]Property{
			"issue_name": {
				Types: []string{"string"},
			},
			"issue_permalink": {
				Types: []string{"string"},
			},
			"provider": {
				Types: []string{"string"},
			},
		},
	}
}

func (externalIssueReferenceV2) Serialize(input *client.ExternalIssueReferenceV2) map[string]any {
	if input == nil {
		return nil
	}

	return map[string]any{
		"issue_name":      input.IssueName,
		"issue_permalink": input.IssuePermalink,
		"provider":        input.Provider,
	}
}
