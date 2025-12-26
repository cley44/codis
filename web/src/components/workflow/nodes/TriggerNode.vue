<script setup lang="ts">
import { Handle, Position } from '@vue-flow/core'

const props = defineProps<{
  data: {
    event: string
  }
}>()

const eventLabels: Record<string, string> = {
  discord_event_type_message_create: 'Message Created',
  discord_event_type_message_reaction_add: 'Reaction Added',
}
</script>

<template>
  <div class="trigger-node">
    <div class="node-header">
      <span class="node-icon">⚡</span>
      <span class="node-title">Trigger</span>
    </div>
    <div class="node-body">
      <div class="field">
        <label>Event:</label>
        <span class="value">{{ eventLabels[data.event] || data.event }}</span>
      </div>
    </div>
    <Handle type="source" :position="Position.Bottom" />
  </div>
</template>

<style scoped>
.trigger-node {
  background: linear-gradient(135deg, #fbbf24, #f59e0b);
  border-radius: 0.75rem;
  min-width: 250px;
  box-shadow: 0 8px 24px rgba(251, 191, 36, 0.35);
  color: #1f2937;
}

.node-header {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1rem;
  border-bottom: 1px solid rgba(0, 0, 0, 0.1);
}

.node-icon {
  font-size: 1.25rem;
}

.node-title {
  font-weight: 600;
  font-size: 0.95rem;
}

.node-body {
  padding: 0.75rem 1rem;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

label {
  font-size: 0.75rem;
  font-weight: 500;
  text-transform: uppercase;
  opacity: 0.7;
}

.value {
  font-size: 0.9rem;
  font-weight: 500;
}
</style>
