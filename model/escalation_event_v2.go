package model

import "github.com/incident-io/singer-tap/client"

type escalationEventV2 struct{}

var EscalationEventV2 escalationEventV2

func (escalationEventV2) Schema() Property {
	return Property{
		Types: []string{"object"},
		Properties: map[string]Property{
			"id": {
				Types: []string{"string"},
			},
			"event": {
				Types: []string{"string"},
			},
			"occurred_at": DateTime.Schema(),
			"urgency": Optional(Property{
				Types: []string{"string"},
			}),
			"channels": Optional(ArrayOf(ChatChannelSlimV2.Schema())),
			"users": Optional(ArrayOf(UserV2.Schema())),
		},
	}
}

func (escalationEventV2) Serialize(input client.EscalationEventV2) map[string]any {
	result := map[string]any{
		"id":          input.Id,
		"event":       input.Event,
		"occurred_at": input.OccurredAt,
	}

	if input.Urgency != nil {
		result["urgency"] = *input.Urgency
	}

	if input.Channels != nil {
		channels := make([]map[string]any, 0, len(*input.Channels))
		for _, channel := range *input.Channels {
			channels = append(channels, ChatChannelSlimV2.Serialize(channel))
		}
		result["channels"] = channels
	}

	if input.Users != nil {
		users := make([]map[string]any, 0, len(*input.Users))
		for _, user := range *input.Users {
			users = append(users, UserV2.Serialize(user))
		}
		result["users"] = users
	}

	return result
}