<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { VueFlow, useVueFlow, Panel } from '@vue-flow/core'
import { Background } from '@vue-flow/background'
import { Controls } from '@vue-flow/controls'
import { useWorkflowStore } from '@/stores/workflow'
import { useWorkflowTransform } from '@/composables/useWorkflowTransform'
import WorkflowToolbar from '@/components/workflow/WorkflowToolbar.vue'
import TriggerNode from '@/components/workflow/nodes/TriggerNode.vue'
import AddRoleNode from '@/components/workflow/nodes/AddRoleNode.vue'
import RemoveRoleNode from '@/components/workflow/nodes/RemoveRoleNode.vue'
import SendMessageNode from '@/components/workflow/nodes/SendMessageNode.vue'
import type { DiscordEventType } from '@/types/workflow'
import { uuid } from 'vue-uuid'

import '@vue-flow/core/dist/style.css'
import '@vue-flow/core/dist/theme-default.css'
import '@vue-flow/controls/dist/style.css'

const props = defineProps<{
  guildId: string
  workflowId: string
}>()

const router = useRouter()
const workflowStore = useWorkflowStore()
const { backendToVueFlow, vueFlowToBackend, validateWorkflow } = useWorkflowTransform()

const {
  nodes,
  edges,
  addNodes,
  addEdges,
  applyNodeChanges,
  applyEdgeChanges,
  setNodes,
  setEdges,
  onConnect,
} = useVueFlow()

// Handle new connections between nodes
onConnect((params) => {
  addEdges([params])
})

// Custom node types
const nodeTypes: any = {
  trigger: TriggerNode,
  discord_node_type_add_member_role: AddRoleNode,
  discord_node_type_remove_member_role: RemoveRoleNode,
  discord_node_type_send_message: SendMessageNode,
}

// Load workflow on mount
onMounted(async () => {
  try {
    await workflowStore.fetchWorkflow(props.workflowId)

    if (workflowStore.currentWorkflow) {
      console.log('Loaded workflow:', workflowStore.currentWorkflow)
      const { nodes: vfNodes, edges: vfEdges } = backendToVueFlow(workflowStore.currentWorkflow)
      console.log('Transformed nodes:', vfNodes)
      console.log('Transformed edges:', vfEdges)
      setNodes(vfNodes)
      setEdges(vfEdges)
    }
  } catch (err) {
    console.error('Failed to load workflow:', err)
    alert(
      'Failed to load workflow. Returning to list.\n\nError: ' +
        (err instanceof Error ? err.message : String(err)),
    )
    router.push({ name: 'WorkflowList', params: { guildId: props.guildId } })
  }
})

// Sync nodes/edges changes to store
watch(
  [nodes, edges],
  () => {
    workflowStore.setNodes(nodes.value)
    workflowStore.setEdges(edges.value)
  },
  { deep: true },
)

async function handleSave() {
  const validation = validateWorkflow(nodes.value, edges.value)

  if (!validation.valid) {
    alert('Validation failed:\n' + validation.errors.join('\n'))
    return
  }

  if (validation.warnings.length > 0) {
    const proceed = confirm('Warnings:\n' + validation.warnings.join('\n') + '\n\nContinue saving?')
    if (!proceed) return
  }

  try {
    // Transform to backend format
    const { workflow, nodes: backendNodes } = vueFlowToBackend(
      workflowStore.currentWorkflow!.id,
      props.guildId,
      nodes.value,
      edges.value,
    )

    // Save to backend with the transformed data
    await workflowStore.saveWorkflow(
      workflow.starting_nodes_ids || [],
      workflow.starting_discord_events || [],
      backendNodes,
    )

    alert('Workflow saved successfully!')
  } catch (err) {
    console.error('Failed to save workflow:', err)
    alert('Failed to save workflow. Please try again.')
  }
}

function handleValidate() {
  const validation = validateWorkflow(nodes.value, edges.value)

  if (validation.valid) {
    const message =
      validation.warnings.length > 0
        ? 'Workflow is valid! ✓\n\nWarnings:\n' + validation.warnings.join('\n')
        : 'Workflow is valid! ✓'
    alert(message)
  } else {
    alert('Validation errors:\n' + validation.errors.join('\n'))
  }
}

function handleAddNode(nodeType: string) {
  const newNode = {
    id: uuid.v4(),
    type: nodeType,
    position: { x: Math.random() * 400, y: 200 + Math.random() * 300 },
    data: { nodeType },
  }
  addNodes([newNode])
}

function handleAddTrigger(eventType: string) {
  const triggerId = `trigger-${eventType}`

  // Check if this trigger already exists
  const existingTrigger = nodes.value.find((n) => n.id === triggerId)
  if (existingTrigger) {
    alert('This trigger already exists in the workflow')
    return
  }

  const newTrigger = {
    id: triggerId,
    type: 'trigger',
    position: { x: nodes.value.filter((n) => n.type === 'trigger').length * 350, y: 0 },
    data: { event: eventType as DiscordEventType },
  }
  addNodes([newTrigger])
}

function goBack() {
  if (workflowStore.hasUnsavedChanges) {
    const proceed = confirm('You have unsaved changes. Are you sure you want to leave?')
    if (!proceed) return
  }
  router.push({ name: 'WorkflowList', params: { guildId: props.guildId } })
}
</script>

<template>
  <div class="editor-container">
    <div class="editor-header">
      <button class="back-btn" @click="goBack">← Back to Workflows</button>
      <h1 class="title">Workflow Editor</h1>
    </div>

    <WorkflowToolbar
      @save="handleSave"
      @validate="handleValidate"
      @add-node="handleAddNode"
      @add-trigger="handleAddTrigger"
    />

    <div class="editor-canvas">
      <VueFlow
        :nodes="nodes"
        :edges="edges"
        :node-types="nodeTypes"
        @nodes-change="applyNodeChanges"
        @edges-change="applyEdgeChanges"
        :default-edge-options="{ type: 'default', animated: true }"
      >
        <Background pattern-color="#374151" :gap="16" />
        <Controls />

        <Panel position="top-right" class="save-indicator">
          <span v-if="workflowStore.hasUnsavedChanges" class="unsaved">Unsaved changes</span>
          <span v-else class="saved">Saved</span>
        </Panel>
      </VueFlow>
    </div>
  </div>
</template>

<style>
.vue-flow {
  background: #020617;
}

.vue-flow__node {
  cursor: move;
}

.vue-flow__edge {
  stroke: #3b82f6;
  stroke-width: 2px;
}

.vue-flow__edge.selected {
  stroke: #60a5fa;
}

.vue-flow__edge-path {
  stroke: #3b82f6;
}

.vue-flow__connectionline {
  stroke: #3b82f6;
}

.vue-flow__handle {
  width: 12px;
  height: 12px;
  background: #3b82f6;
  border: 2px solid #020617;
}

.vue-flow__handle:hover {
  background: #60a5fa;
}
</style>

<style scoped>
.editor-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: #020617;
}

.editor-header {
  padding: 1rem 1.5rem;
  background: rgba(15, 23, 42, 0.96);
  border-bottom: 1px solid rgba(31, 41, 55, 0.9);
}

.back-btn {
  background: transparent;
  border: none;
  color: #9ca3af;
  cursor: pointer;
  font-size: 0.9rem;
  padding: 0;
  margin-bottom: 0.5rem;
  transition: color 0.15s ease;
}

.back-btn:hover {
  color: #e5e7eb;
}

.title {
  font-size: 1.5rem;
  font-weight: 600;
  color: #e5e7eb;
  margin: 0;
}

.editor-canvas {
  flex: 1;
  position: relative;
}

.save-indicator {
  background: rgba(15, 23, 42, 0.95);
  padding: 0.5rem 1rem;
  border-radius: 0.5rem;
  font-size: 0.85rem;
  font-weight: 500;
  border: 1px solid rgba(148, 163, 184, 0.3);
}

.unsaved {
  color: #fbbf24;
}

.saved {
  color: #22c55e;
}
</style>
