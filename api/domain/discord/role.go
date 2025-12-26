package discord

import (
	"github.com/disgoorg/snowflake/v2"
	"github.com/samber/oops"
)

func (svc *DiscordService) AddMemberRole(guildID string, userID string, roleID string) error {
	sfGuildID, err := snowflake.Parse(guildID)
	if err != nil {
		return oops.With("guildID", guildID).Wrapf(err, "Failed to parse guild ID")
	}

	sfUserID, err := snowflake.Parse(userID)
	if err != nil {
		return oops.With("guildID", guildID).With("userID", userID).Wrapf(err, "Failed to parse user ID")
	}

	sfRoleID, err := snowflake.Parse(roleID)
	if err != nil {
		return oops.With("guildID", guildID).With("userID", userID).With("roleID", roleID).Wrapf(err, "Failed to parse role ID")
	}

	err = svc.botClient.Rest().AddMemberRole(sfGuildID, sfUserID, sfRoleID)
	if err != nil {
		return oops.With("guildID", guildID).With("userID", userID).With("roleID", roleID).Wrapf(err, "Failed to add role to member")
	}

	return nil
}

func (svc *DiscordService) RemoveRoleFromMember(guildID string, userID string, roleID string) error {
	sfGuildID, err := snowflake.Parse(guildID)
	if err != nil {
		return oops.With("guildID", guildID).Wrapf(err, "Failed to parse guild ID")
	}

	sfUserID, err := snowflake.Parse(userID)
	if err != nil {
		return oops.With("guildID", guildID).With("userID", userID).Wrapf(err, "Failed to parse user ID")
	}

	sfRoleID, err := snowflake.Parse(roleID)
	if err != nil {
		return oops.With("guildID", guildID).With("userID", userID).With("roleID", roleID).Wrapf(err, "Failed to parse role ID")
	}

	err = svc.botClient.Rest().RemoveMemberRole(sfGuildID, sfUserID, sfRoleID)
	if err != nil {
		return oops.With("guildID", guildID).With("userID", userID).With("roleID", roleID).Wrapf(err, "Failed to remove role from member")
	}
	return nil
}
