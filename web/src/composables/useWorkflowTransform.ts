import type {
  Workflow,
  Node,
  VueFlowNode,
  VueFlowEdge,
  ValidationResult,
} from '@/types/workflow'

const NODE_WIDTH = 300
const NODE_HEIGHT = 150
const VERTICAL_SPACING = 200
const HORIZONTAL_SPACING = 350

export function useWorkflowTransform() {
  // Convert backend workflow to vue-flow format
  function backendToVueFlow(workflow: Workflow): {
    nodes: VueFlowNode[]
    edges: VueFlowEdge[]
  } {
    const vueFlowNodes: VueFlowNode[] = []
    const vueFlowEdges: VueFlowEdge[] = []

    // Ensure arrays are defined
    const startingEvents = workflow.starting_discord_events || []
    const startingNodeIds = (workflow.starting_nodes_ids || []) as string[]
    const workflowNodes = workflow.nodes || []

    // Create trigger nodes from starting_discord_events
    startingEvents.forEach((event, index) => {
      const triggerId = `trigger-${event}`
      vueFlowNodes.push({
        id: triggerId,
        type: 'trigger',
        position: { x: index * HORIZONTAL_SPACING, y: 0 },
        data: { event },
      })

      // Connect trigger to starting nodes
      startingNodeIds.forEach((nodeId) => {
        vueFlowEdges.push({
          id: `${triggerId}-${nodeId}`,
          source: triggerId,
          target: nodeId,
        })
      })
    })

    // Create action nodes from workflow.nodes
    const nodeMap = new Map(workflowNodes.map((n) => [n.id, n]))
    const processedNodes = new Set<string>()

    // Calculate positions using topological layout
    startingNodeIds.forEach((startNodeId, chainIndex) => {
      let currentNodeId: string | null = startNodeId
      let depth = 1

      while (currentNodeId && nodeMap.has(currentNodeId)) {
        if (processedNodes.has(currentNodeId)) break

        const node: Node = nodeMap.get(currentNodeId)!
        vueFlowNodes.push({
          id: node.id,
          type: node.type,
          position: {
            x: chainIndex * HORIZONTAL_SPACING,
            y: depth * VERTICAL_SPACING,
          },
          data: {
            nodeType: node.type,
            nodeId: node.id,
          },
        })

        // Create edge to next node
        if (node.next_node_id) {
          vueFlowEdges.push({
            id: `${node.id}-${node.next_node_id}`,
            source: node.id,
            target: node.next_node_id,
          })
        }

        processedNodes.add(currentNodeId)
        currentNodeId = node.next_node_id
        depth++
      }
    })

    return { nodes: vueFlowNodes, edges: vueFlowEdges }
  }

  // Convert vue-flow format back to backend format
  function vueFlowToBackend(
    workflowId: string,
    guildId: string,
    vueFlowNodes: VueFlowNode[],
    vueFlowEdges: VueFlowEdge[],
  ): { workflow: Partial<Workflow>; nodes: Node[] } {
    // Extract trigger nodes
    const triggerNodes = vueFlowNodes.filter((n) => n.type === 'trigger')
    const starting_discord_events = triggerNodes.map((n) => n.data.event!)

    // Build edge map for quick lookups
    const edgeMap = new Map<string, string[]>()
    vueFlowEdges.forEach((edge) => {
      if (!edgeMap.has(edge.source)) {
        edgeMap.set(edge.source, [])
      }
      edgeMap.get(edge.source)!.push(edge.target)
    })

    // Find starting nodes (nodes connected from triggers)
    const starting_nodes_ids = new Set<string>()
    triggerNodes.forEach((trigger) => {
      const targets = edgeMap.get(trigger.id) || []
      targets.forEach((t) => starting_nodes_ids.add(t))
    })

    // Convert action nodes to backend format
    const actionNodes = vueFlowNodes.filter((n) => n.type !== 'trigger')
    const nodes: Node[] = actionNodes.map((vfNode) => {
      const outgoingEdges = edgeMap.get(vfNode.id) || []
      const next_node_id: string | null =
        outgoingEdges.length > 0 ? (outgoingEdges[0] ?? null) : null

      return {
        id: vfNode.data.nodeId || vfNode.id,
        workflow_id: workflowId,
        type: vfNode.data.nodeType!,
        next_node_id,
        created_at: '', // Backend will set
        updated_at: '',
      }
    })

    return {
      workflow: {
        id: workflowId,
        guild_id: guildId,
        starting_nodes_ids: Array.from(starting_nodes_ids),
        starting_discord_events,
      },
      nodes,
    }
  }

  // Validate workflow before saving
  function validateWorkflow(nodes: VueFlowNode[], edges: VueFlowEdge[]): ValidationResult {
    const errors: string[] = []
    const warnings: string[] = []

    // Check for trigger nodes
    const triggerNodes = nodes.filter((n) => n.type === 'trigger')
    if (triggerNodes.length === 0) {
      errors.push('Workflow must have at least one trigger event')
    }

    // Check for cycles
    const visited = new Set<string>()
    const recursionStack = new Set<string>()
    const edgeMap = new Map<string, string[]>()
    edges.forEach((edge) => {
      if (!edgeMap.has(edge.source)) {
        edgeMap.set(edge.source, [])
      }
      edgeMap.get(edge.source)!.push(edge.target)
    })

    function hasCycle(nodeId: string): boolean {
      visited.add(nodeId)
      recursionStack.add(nodeId)

      const neighbors = edgeMap.get(nodeId) || []
      for (const neighbor of neighbors) {
        if (!visited.has(neighbor)) {
          if (hasCycle(neighbor)) return true
        } else if (recursionStack.has(neighbor)) {
          return true
        }
      }

      recursionStack.delete(nodeId)
      return false
    }

    for (const node of nodes) {
      if (!visited.has(node.id)) {
        if (hasCycle(node.id)) {
          errors.push('Workflow contains cycles - loops are not allowed')
          break
        }
      }
    }

    // Check for disconnected nodes (excluding triggers)
    const actionNodes = nodes.filter((n) => n.type !== 'trigger')
    const connectedNodes = new Set<string>()
    edges.forEach((edge) => {
      connectedNodes.add(edge.source)
      connectedNodes.add(edge.target)
    })

    actionNodes.forEach((node) => {
      if (!connectedNodes.has(node.id)) {
        warnings.push(`Node "${node.id}" is not connected to any other nodes`)
      }
    })

    // Check for multiple outgoing edges (backend supports only next_node_id)
    edgeMap.forEach((targets, source) => {
      if (targets.length > 1) {
        const sourceNode = nodes.find((n) => n.id === source)
        if (sourceNode?.type !== 'trigger') {
          errors.push(`Node "${source}" has multiple outgoing connections. Only one is allowed.`)
        }
      }
    })

    return {
      valid: errors.length === 0,
      errors,
      warnings,
    }
  }

  return {
    backendToVueFlow,
    vueFlowToBackend,
    validateWorkflow,
  }
}
