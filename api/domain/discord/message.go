package discord

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/samber/oops"
)

func (svc *DiscordService) SendMessage(channelID string, content string) error {
	sfChannelID, err := snowflake.Parse(channelID)
	if err != nil {
		return oops.With("channelID", channelID).Wrapf(err, "Failed to parse channel ID")
	}

	msgCreate := discord.MessageCreate{
		Content: content,
	}

	_, err = svc.botClient.Rest().CreateMessage(sfChannelID, msgCreate)
	if err != nil {
		return oops.With("channelID", channelID).Wrapf(err, "Failed to send message")
	}

	return nil
}
