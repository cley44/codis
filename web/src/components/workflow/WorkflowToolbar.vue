<script setup lang="ts">
import { ref } from 'vue'

const emit = defineEmits<{
  save: []
  addNode: [nodeType: string]
  addTrigger: [eventType: string]
  validate: []
}>()

const showAddMenu = ref(false)
const showTriggerMenu = ref(false)

const nodeTypes = [
  { type: 'discord_node_type_add_member_role', label: 'Add Role', icon: '➕' },
  { type: 'discord_node_type_remove_member_role', label: 'Remove Role', icon: '➖' },
  { type: 'discord_node_type_send_message', label: 'Send Message', icon: '💬' },
]

const triggerTypes = [
  { type: 'discord_event_type_message_create', label: 'Message Created', icon: '📝' },
  { type: 'discord_event_type_message_reaction_add', label: 'Reaction Added', icon: '👍' },
]
</script>

<template>
  <div class="toolbar">
    <div class="toolbar-section">
      <button class="btn btn-primary" @click="emit('save')">Save Workflow</button>
      <button class="btn btn-secondary" @click="emit('validate')">Validate</button>
    </div>

    <div class="toolbar-section">
      <div class="dropdown">
        <button class="btn btn-trigger" @click="showTriggerMenu = !showTriggerMenu">
          + Add Trigger
        </button>
        <div v-if="showTriggerMenu" class="dropdown-menu">
          <button
            v-for="triggerType in triggerTypes"
            :key="triggerType.type"
            class="dropdown-item"
            @click="[emit('addTrigger', triggerType.type), (showTriggerMenu = false)]"
          >
            <span class="item-icon">{{ triggerType.icon }}</span>
            {{ triggerType.label }}
          </button>
        </div>
      </div>

      <div class="dropdown">
        <button class="btn btn-accent" @click="showAddMenu = !showAddMenu">+ Add Node</button>
        <div v-if="showAddMenu" class="dropdown-menu">
          <button
            v-for="nodeType in nodeTypes"
            :key="nodeType.type"
            class="dropdown-item"
            @click="[emit('addNode', nodeType.type), (showAddMenu = false)]"
          >
            <span class="item-icon">{{ nodeType.icon }}</span>
            {{ nodeType.label }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1rem 1.5rem;
  background: rgba(15, 23, 42, 0.96);
  border-bottom: 1px solid rgba(31, 41, 55, 0.9);
}

.toolbar-section {
  display: flex;
  gap: 0.75rem;
}

.btn {
  padding: 0.5rem 1rem;
  border-radius: 0.5rem;
  border: none;
  cursor: pointer;
  font-size: 0.9rem;
  font-weight: 500;
  transition: all 0.15s ease;
}

.btn-primary {
  background: linear-gradient(135deg, #5865f2, #3b82f6);
  color: white;
}

.btn-primary:hover {
  filter: brightness(1.1);
}

.btn-secondary {
  background: rgba(55, 65, 81, 0.6);
  color: #e5e7eb;
  border: 1px solid rgba(148, 163, 184, 0.3);
}

.btn-secondary:hover {
  background: rgba(55, 65, 81, 0.8);
}

.btn-accent {
  background: rgba(34, 197, 94, 0.2);
  color: #22c55e;
  border: 1px solid rgba(34, 197, 94, 0.4);
}

.btn-accent:hover {
  background: rgba(34, 197, 94, 0.3);
}

.btn-trigger {
  background: rgba(251, 191, 36, 0.2);
  color: #fbbf24;
  border: 1px solid rgba(251, 191, 36, 0.4);
}

.btn-trigger:hover {
  background: rgba(251, 191, 36, 0.3);
}

.dropdown {
  position: relative;
}

.dropdown-menu {
  position: absolute;
  top: calc(100% + 0.5rem);
  right: 0;
  background: rgba(15, 23, 42, 0.98);
  border: 1px solid rgba(148, 163, 184, 0.3);
  border-radius: 0.5rem;
  min-width: 200px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4);
  z-index: 1000;
}

.dropdown-item {
  width: 100%;
  padding: 0.75rem 1rem;
  background: transparent;
  border: none;
  color: #e5e7eb;
  font-size: 0.9rem;
  text-align: left;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.dropdown-item:hover {
  background: rgba(55, 65, 81, 0.5);
}

.item-icon {
  font-size: 1.1rem;
}
</style>
