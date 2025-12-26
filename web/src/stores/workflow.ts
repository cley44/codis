import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import type { Workflow, VueFlowNode, VueFlowEdge } from '@/types/workflow'

const BACKEND_BASE_URL = 'http://localhost:80'

export const useWorkflowStore = defineStore('workflow', () => {
  // State
  const workflows = ref<Workflow[]>([])
  const currentWorkflow = ref<Workflow | null>(null)
  const currentNodes = ref<VueFlowNode[]>([])
  const currentEdges = ref<VueFlowEdge[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)
  const hasUnsavedChanges = ref(false)

  // Computed
  const currentGuildId = computed(() => currentWorkflow.value?.guild_id || null)

  // Actions
  async function fetchWorkflows(guildId: string): Promise<void> {
    loading.value = true
    error.value = null
    try {
      const res = await fetch(`${BACKEND_BASE_URL}/workflows?guild_id=${guildId}`, {
        credentials: 'include',
      })

      if (!res.ok) {
        throw new Error(`Failed to fetch workflows: ${res.status}`)
      }

      workflows.value = await res.json()
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Unknown error'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function fetchWorkflow(workflowId: string): Promise<void> {
    loading.value = true
    error.value = null
    try {
      const res = await fetch(`${BACKEND_BASE_URL}/workflows/${workflowId}`, {
        credentials: 'include',
      })

      if (!res.ok) {
        throw new Error(`Failed to fetch workflow: ${res.status}`)
      }

      currentWorkflow.value = await res.json()
      hasUnsavedChanges.value = false
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Unknown error'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function createWorkflow(guildId: string): Promise<Workflow> {
    loading.value = true
    error.value = null
    try {
      const res = await fetch(`${BACKEND_BASE_URL}/workflows`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({
          guild_id: guildId,
          starting_nodes_ids: [],
          starting_discord_events: [],
        }),
      })

      if (!res.ok) {
        throw new Error(`Failed to create workflow: ${res.status}`)
      }

      const workflow = await res.json()
      workflows.value.push(workflow)
      return workflow
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Unknown error'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function saveWorkflow(
    startingNodesIds: string[],
    startingDiscordEvents: string[],
    nodes: any[],
  ): Promise<void> {
    if (!currentWorkflow.value) {
      throw new Error('No workflow to save')
    }

    loading.value = true
    error.value = null
    try {
      // Update workflow metadata and nodes
      const workflowRes = await fetch(`${BACKEND_BASE_URL}/workflows/${currentWorkflow.value.id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({
          starting_nodes_ids: startingNodesIds,
          starting_discord_events: startingDiscordEvents,
          nodes: nodes,
        }),
      })

      if (!workflowRes.ok) {
        throw new Error(`Failed to update workflow: ${workflowRes.status}`)
      }

      // Update local state
      currentWorkflow.value.starting_nodes_ids = startingNodesIds
      currentWorkflow.value.starting_discord_events = startingDiscordEvents as any
      currentWorkflow.value.nodes = nodes as any

      hasUnsavedChanges.value = false
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Unknown error'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function deleteWorkflow(workflowId: string): Promise<void> {
    loading.value = true
    error.value = null
    try {
      const res = await fetch(`${BACKEND_BASE_URL}/workflows/${workflowId}`, {
        method: 'DELETE',
        credentials: 'include',
      })

      if (!res.ok) {
        throw new Error(`Failed to delete workflow: ${res.status}`)
      }

      workflows.value = workflows.value.filter((w) => w.id !== workflowId)
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Unknown error'
      throw err
    } finally {
      loading.value = false
    }
  }

  function setNodes(nodes: VueFlowNode[]): void {
    currentNodes.value = nodes
    hasUnsavedChanges.value = true
  }

  function setEdges(edges: VueFlowEdge[]): void {
    currentEdges.value = edges
    hasUnsavedChanges.value = true
  }

  function resetCurrentWorkflow(): void {
    currentWorkflow.value = null
    currentNodes.value = []
    currentEdges.value = []
    hasUnsavedChanges.value = false
  }

  return {
    // State
    workflows,
    currentWorkflow,
    currentNodes,
    currentEdges,
    loading,
    error,
    hasUnsavedChanges,
    // Computed
    currentGuildId,
    // Actions
    fetchWorkflows,
    fetchWorkflow,
    createWorkflow,
    saveWorkflow,
    deleteWorkflow,
    setNodes,
    setEdges,
    resetCurrentWorkflow,
  }
})
