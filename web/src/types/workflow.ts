// Backend API types
export interface Workflow {
  id: string
  starting_nodes_ids: string[] | null
  guild_id: string
  starting_discord_events: DiscordEventType[]
  nodes: Node[]
  created_at: string
  updated_at: string
}

export interface Node {
  id: string
  workflow_id: string
  type: DiscordNodeType
  next_node_id: string | null
  data?: {
    role_id?: string | null
    channel_id?: string | null
    message_content?: string | null
  }
  created_at: string
  updated_at: string
}

export type DiscordNodeType =
  | 'discord_node_type_add_member_role'
  | 'discord_node_type_remove_member_role'
  | 'discord_node_type_send_message'

export type DiscordEventType =
  | 'discord_event_type_message_create'
  | 'discord_event_type_message_reaction_add'

// Vue Flow types
export interface VueFlowNode {
  id: string
  type: string // 'trigger' | DiscordNodeType
  position: { x: number; y: number }
  data: NodeData
}

export interface NodeData {
  // For trigger nodes
  event?: DiscordEventType

  // For action nodes
  nodeType?: DiscordNodeType
  nodeId?: string // Backend node ID

  // Node-specific properties (for inline editing)
  roleId?: string
  channelId?: string
  messageContent?: string
}

export interface VueFlowEdge {
  id: string
  source: string
  target: string
  type?: string
}

// Validation result
export interface ValidationResult {
  valid: boolean
  errors: string[]
  warnings: string[]
}
