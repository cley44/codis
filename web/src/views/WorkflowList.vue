<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useWorkflowStore } from '@/stores/workflow'

const props = defineProps<{
  guildId: string
}>()

const router = useRouter()
const workflowStore = useWorkflowStore()

onMounted(async () => {
  try {
    await workflowStore.fetchWorkflows(props.guildId)
  } catch (err) {
    console.error('Failed to fetch workflows:', err)
  }
})

async function handleCreateWorkflow() {
  try {
    const workflow = await workflowStore.createWorkflow(props.guildId)
    router.push({
      name: 'WorkflowEditor',
      params: { guildId: props.guildId, workflowId: workflow.id },
    })
  } catch (err) {
    console.error('Failed to create workflow:', err)
    alert('Failed to create workflow. Please try again.')
  }
}

function handleEditWorkflow(workflowId: string) {
  router.push({
    name: 'WorkflowEditor',
    params: { guildId: props.guildId, workflowId },
  })
}

async function handleDeleteWorkflow(workflowId: string) {
  if (confirm('Are you sure you want to delete this workflow?')) {
    try {
      await workflowStore.deleteWorkflow(workflowId)
    } catch (err) {
      console.error('Failed to delete workflow:', err)
      alert('Failed to delete workflow. Please try again.')
    }
  }
}

function goBack() {
  router.push({ name: 'GuildsDashboard' })
}
</script>

<template>
  <main class="page">
    <section class="container">
      <div class="header">
        <div>
          <button class="back-btn" @click="goBack">← Back to Guilds</button>
          <h1>Workflows</h1>
          <p class="subtitle">Manage automation workflows for this guild</p>
        </div>
        <button class="btn btn-create" @click="handleCreateWorkflow">+ Create Workflow</button>
      </div>

      <div v-if="workflowStore.loading" class="loading">
        <p>Loading workflows...</p>
      </div>

      <div v-else-if="workflowStore.error" class="error-container">
        <p class="error-message">{{ workflowStore.error }}</p>
      </div>

      <div v-else-if="workflowStore.workflows.length === 0" class="empty">
        <p>No workflows yet. Create your first workflow to get started!</p>
      </div>

      <div v-else class="workflows-grid">
        <div v-for="workflow in workflowStore.workflows" :key="workflow.id" class="workflow-card">
          <div class="card-header">
            <h2 class="workflow-name">Workflow {{ workflow.id.slice(0, 8) }}</h2>
            <div class="badges">
              <span class="badge">{{ workflow.nodes?.length || 0 }} nodes</span>
              <span class="badge">{{ workflow.starting_discord_events.length }} triggers</span>
            </div>
          </div>

          <div class="card-footer">
            <button class="btn btn-edit" @click="handleEditWorkflow(workflow.id)">Edit</button>
            <button class="btn btn-delete" @click="handleDeleteWorkflow(workflow.id)">
              Delete
            </button>
          </div>
        </div>
      </div>
    </section>
  </main>
</template>

<style scoped>
.page {
  min-height: 100vh;
  padding: 3rem 1.5rem;
  background: radial-gradient(circle at top left, #1f2933, #020617);
  color: #e5e7eb;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 2rem;
}

.back-btn {
  background: transparent;
  border: none;
  color: #9ca3af;
  cursor: pointer;
  font-size: 0.9rem;
  padding: 0;
  margin-bottom: 1rem;
  transition: color 0.15s ease;
}

.back-btn:hover {
  color: #e5e7eb;
}

h1 {
  font-size: 2rem;
  font-weight: 600;
  margin-bottom: 0.5rem;
}

.subtitle {
  color: #9ca3af;
  font-size: 1rem;
}

.btn-create {
  background: linear-gradient(135deg, #5865f2, #3b82f6);
  color: white;
  padding: 0.65rem 1.25rem;
  border-radius: 0.5rem;
  border: none;
  cursor: pointer;
  font-weight: 500;
  transition: filter 0.15s ease;
}

.btn-create:hover {
  filter: brightness(1.1);
}

.loading,
.error-container,
.empty {
  text-align: center;
  padding: 3rem 1.5rem;
  color: #9ca3af;
}

.error-message {
  color: #ef4444;
  margin-bottom: 1rem;
  font-size: 0.95rem;
}

.workflows-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 1.5rem;
}

.workflow-card {
  background: rgba(15, 23, 42, 0.96);
  border-radius: 1rem;
  padding: 1.5rem;
  border: 1px solid rgba(148, 163, 184, 0.2);
  box-shadow: 0 24px 60px rgba(0, 0, 0, 0.6);
  transition:
    transform 0.15s ease-out,
    box-shadow 0.15s ease-out;
}

.workflow-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 28px 70px rgba(0, 0, 0, 0.7);
}

.card-header {
  margin-bottom: 1rem;
}

.workflow-name {
  font-size: 1.1rem;
  font-weight: 600;
  margin-bottom: 0.5rem;
}

.badges {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.badge {
  display: inline-block;
  padding: 0.25rem 0.6rem;
  border-radius: 999px;
  font-size: 0.75rem;
  font-weight: 500;
  background: rgba(55, 65, 81, 0.6);
  color: #9ca3af;
  border: 1px solid rgba(148, 163, 184, 0.3);
}

.card-footer {
  display: flex;
  gap: 0.75rem;
  margin-top: 1rem;
  padding-top: 1rem;
  border-top: 1px solid rgba(55, 65, 81, 0.5);
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

.btn-edit {
  flex: 1;
  background: linear-gradient(135deg, #5865f2, #3b82f6);
  color: white;
}

.btn-edit:hover {
  filter: brightness(1.1);
}

.btn-delete {
  background: rgba(239, 68, 68, 0.2);
  color: #ef4444;
  border: 1px solid rgba(239, 68, 68, 0.4);
}

.btn-delete:hover {
  background: rgba(239, 68, 68, 0.3);
}

@media (max-width: 768px) {
  .workflows-grid {
    grid-template-columns: 1fr;
  }

  h1 {
    font-size: 1.5rem;
  }

  .header {
    flex-direction: column;
    gap: 1rem;
  }

  .btn-create {
    width: 100%;
  }
}
</style>
