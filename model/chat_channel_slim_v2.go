package model

import "github.com/incident-io/singer-tap/client"

type chatChannelSlimV2 struct{}

var ChatChannelSlimV2 chatChannelSlimV2

func (chatChannelSlimV2) Schema() Property {
	return Property{
		Types: []string{"object"},
		Properties: map[string]Property{
			"slack_channel_id": Optional(Property{
				Types: []string{"string"},
			}),
			"slack_team_id": Optional(Property{
				Types: []string{"string"},
			}),
			"microsoft_teams_channel_id": Optional(Property{
				Types: []string{"string"},
			}),
			"microsoft_teams_team_id": Optional(Property{
				Types: []string{"string"},
			}),
		},
	}
}

func (chatChannelSlimV2) Serialize(input client.ChatChannelSlimV2) map[string]any {
	result := make(map[string]any)
	
	if input.SlackChannelId != nil {
		result["slack_channel_id"] = *input.SlackChannelId
	}
	if input.SlackTeamId != nil {
		result["slack_team_id"] = *input.SlackTeamId
	}
	if input.MicrosoftTeamsChannelId != nil {
		result["microsoft_teams_channel_id"] = *input.MicrosoftTeamsChannelId
	}
	if input.MicrosoftTeamsTeamId != nil {
		result["microsoft_teams_team_id"] = *input.MicrosoftTeamsTeamId
	}
	
	return result
}