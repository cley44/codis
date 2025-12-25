package models

type DiscordNodeType string

const (
	DiscordNodeTypeAddMemberRole    DiscordNodeType = "discord_node_type_add_member_role"
	DiscordNodeTypeRemoveMemberRole DiscordNodeType = "discord_node_type_remove_member_role"
	DiscordNodeTypeSendMessage      DiscordNodeType = "discord_node_type_send_message"
)
