<script setup lang="ts">
import { Handle, Position, useNode } from '@vue-flow/core'
import { computed } from 'vue'

const { node } = useNode()

const channelId = computed({
  get: () => node.data.channelId || '',
  set: (value: string) => {
    node.data.channelId = value
  },
})

const messageContent = computed({
  get: () => node.data.messageContent || '',
  set: (value: string) => {
    node.data.messageContent = value
  },
})
</script>

<template>
  <div class="action-node">
    <Handle type="target" :position="Position.Top" />
    <div class="node-header">
      <span class="node-icon">💬</span>
      <span class="node-title">Send Message</span>
    </div>
    <div class="node-body">
      <div class="field">
        <label for="channel-id">Channel ID:</label>
        <input
          id="channel-id"
          v-model="channelId"
          type="text"
          placeholder="Enter channel ID..."
          class="input"
        />
      </div>
      <div class="field">
        <label for="message">Message:</label>
        <textarea
          id="message"
          v-model="messageContent"
          placeholder="Enter message..."
          class="textarea"
          rows="3"
        />
      </div>
    </div>
    <Handle type="source" :position="Position.Bottom" />
  </div>
</template>

<style scoped>
.action-node {
  background: rgba(15, 23, 42, 0.98);
  border: 1px solid rgba(148, 163, 184, 0.3);
  border-radius: 0.75rem;
  min-width: 280px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4);
  color: #e5e7eb;
}

.node-header {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1rem;
  border-bottom: 1px solid rgba(55, 65, 81, 0.5);
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
  gap: 0.35rem;
}

.field + .field {
  margin-top: 0.5rem;
}

label {
  font-size: 0.75rem;
  font-weight: 500;
  text-transform: uppercase;
  color: #9ca3af;
}

.input,
.textarea {
  background: rgba(31, 41, 55, 0.8);
  border: 1px solid rgba(75, 85, 99, 0.6);
  border-radius: 0.5rem;
  padding: 0.5rem 0.75rem;
  color: #e5e7eb;
  font-size: 0.9rem;
  outline: none;
  transition: border-color 0.15s ease;
  font-family: inherit;
}

.input:focus,
.textarea:focus {
  border-color: #3b82f6;
}

.input::placeholder,
.textarea::placeholder {
  color: #6b7280;
}

.textarea {
  resize: vertical;
  min-height: 60px;
}
</style>
