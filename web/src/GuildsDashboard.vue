<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const BACKEND_BASE_URL = 'http://localhost:80'

interface DiscordGuild {
  ID: string
  Name: string
  IconURL: string | null
  BannerURL: string | null
  Owner: boolean
  BotInviteLink: string
  BotPresent: boolean
}

const guilds = ref<DiscordGuild[]>([])
const loading = ref(true)
const error = ref<string | null>(null)

function getGuildIconUrl(guild: DiscordGuild): string {
  if (guild.IconURL) {
    // If IconURL is already a full URL, return it
    if (guild.IconURL.startsWith('http')) {
      return guild.IconURL
    }
    // Otherwise construct Discord CDN URL
    return `https://cdn.discordapp.com/icons/${guild.ID}/${guild.IconURL}.png`
  }
  // Return placeholder or default icon
  return ''
}

async function fetchGuilds() {
  loading.value = true
  error.value = null
  try {
    const res = await fetch(`${BACKEND_BASE_URL}/discord/guilds`, {
      credentials: 'include',
    })

    if (res.status === 401) {
      error.value = 'Unauthorized. Please log in with Discord.'
      loading.value = false
      return
    }

    if (!res.ok) {
      const text = await res.text()
      error.value = `Failed to fetch guilds: ${res.status} ${text}`
      loading.value = false
      return
    }

    const data = await res.json()
    guilds.value = data as DiscordGuild[]
  } catch (err: unknown) {
    if (err instanceof Error) {
      error.value = err.message
    } else {
      error.value = 'An unexpected error occurred'
    }
  } finally {
    loading.value = false
  }
}

function handleInviteBot(guild: DiscordGuild) {
  if (guild.BotInviteLink) {
    window.open(guild.BotInviteLink, '_blank')
  }
}

onMounted(() => {
  fetchGuilds()
})
</script>

<template>
  <main class="page">
    <section class="container">
      <div class="header">
        <h1>My Guilds</h1>
        <p class="subtitle">Manage your Discord servers</p>
      </div>

      <div v-if="loading" class="loading">
        <p>Loading guilds...</p>
      </div>

      <div v-else-if="error" class="error-container">
        <p class="error-message">{{ error }}</p>
        <button type="button" class="btn" @click="fetchGuilds">Retry</button>
      </div>

      <div v-else-if="guilds.length === 0" class="empty">
        <p>No guilds found. You need to be an administrator of at least one Discord server.</p>
      </div>

      <div v-else class="guilds-grid">
        <div v-for="guild in guilds" :key="guild.ID" class="guild-card">
          <div class="card-header">
            <div class="guild-icon-container">
              <img
                v-if="getGuildIconUrl(guild)"
                :src="getGuildIconUrl(guild)"
                :alt="`${guild.Name} icon`"
                class="guild-icon"
              />
              <div v-else class="guild-icon-placeholder">
                {{ guild.Name.charAt(0).toUpperCase() }}
              </div>
            </div>
            <div class="guild-info">
              <h2 class="guild-name">{{ guild.Name }}</h2>
              <div class="badges">
                <span v-if="guild.Owner" class="badge badge-owner">Owner</span>
                <span
                  :class="['badge', guild.BotPresent ? 'badge-bot-present' : 'badge-bot-absent']"
                >
                  {{ guild.BotPresent ? 'Bot Active' : 'Bot Not Installed' }}
                </span>
              </div>
            </div>
          </div>

          <div class="card-footer">
            <button
              v-if="guild.BotPresent"
              type="button"
              class="btn btn-workflows"
              @click="router.push({ name: 'WorkflowList', params: { guildId: guild.ID } })"
            >
              Manage Workflows
            </button>
            <button
              v-else
              type="button"
              class="btn btn-invite"
              @click="handleInviteBot(guild)"
            >
              Invite Bot
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
  display: flex;
  align-items: flex-start;
  justify-content: center;
  padding: 3rem 1.5rem;
  background: radial-gradient(circle at top left, #1f2933, #020617);
  color: #e5e7eb;
  font-family:
    system-ui,
    -apple-system,
    BlinkMacSystemFont,
    'SF Pro Text',
    sans-serif;
}

.container {
  width: 100%;
  max-width: 1200px;
}

.header {
  margin-bottom: 2rem;
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

.loading,
.error-container,
.empty {
  text-align: center;
  padding: 3rem 1.5rem;
  color: #9ca3af;
}

.error-message {
  color: #f97373;
  margin-bottom: 1rem;
  font-size: 0.95rem;
}

.guilds-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 1.5rem;
}

.guild-card {
  background: rgba(15, 23, 42, 0.96);
  border-radius: 1rem;
  padding: 1.5rem;
  box-shadow:
    0 24px 60px rgba(0, 0, 0, 0.6),
    0 0 0 1px rgba(148, 163, 184, 0.2);
  transition:
    transform 0.15s ease-out,
    box-shadow 0.15s ease-out;
  display: flex;
  flex-direction: column;
}

.guild-card:hover {
  transform: translateY(-2px);
  box-shadow:
    0 28px 70px rgba(0, 0, 0, 0.7),
    0 0 0 1px rgba(148, 163, 184, 0.3);
}

.card-header {
  display: flex;
  gap: 1rem;
  margin-bottom: 1rem;
}

.guild-icon-container {
  flex-shrink: 0;
}

.guild-icon {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  object-fit: cover;
}

.guild-icon-placeholder {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  background: linear-gradient(135deg, #5865f2, #3b82f6);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
  font-weight: 600;
  color: white;
}

.guild-info {
  flex: 1;
  min-width: 0;
}

.guild-name {
  font-size: 1.1rem;
  font-weight: 600;
  margin-bottom: 0.5rem;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
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
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.badge-owner {
  background: rgba(251, 191, 36, 0.2);
  color: #fbbf24;
  border: 1px solid rgba(251, 191, 36, 0.3);
}

.badge-bot-present {
  background: rgba(34, 197, 94, 0.2);
  color: #22c55e;
  border: 1px solid rgba(34, 197, 94, 0.3);
}

.badge-bot-absent {
  background: rgba(239, 68, 68, 0.2);
  color: #ef4444;
  border: 1px solid rgba(239, 68, 68, 0.3);
}

.card-footer {
  margin-top: auto;
  padding-top: 1rem;
  border-top: 1px solid rgba(55, 65, 81, 0.9);
  display: flex;
  align-items: center;
  justify-content: flex-end;
}

.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.35rem;
  padding: 0.55rem 1.1rem;
  border-radius: 999px;
  border: none;
  cursor: pointer;
  background: linear-gradient(135deg, #5865f2, #3b82f6);
  color: white;
  font-size: 0.9rem;
  font-weight: 500;
  box-shadow:
    0 10px 25px rgba(59, 130, 246, 0.45),
    0 0 0 1px rgba(59, 130, 246, 0.6);
  transition:
    transform 0.08s ease-out,
    box-shadow 0.08s ease-out,
    filter 0.12s ease-out;
}

.btn:hover {
  transform: translateY(-1px);
  filter: brightness(1.05);
  box-shadow:
    0 14px 35px rgba(59, 130, 246, 0.6),
    0 0 0 1px rgba(96, 165, 250, 0.9);
}

.btn:active {
  transform: translateY(0);
  box-shadow:
    0 6px 18px rgba(37, 99, 235, 0.6),
    0 0 0 1px rgba(59, 130, 246, 0.6);
}

.btn-invite {
  width: 100%;
}

.bot-status-active {
  color: #22c55e;
  font-size: 0.85rem;
  font-weight: 500;
}

@media (max-width: 768px) {
  .guilds-grid {
    grid-template-columns: 1fr;
  }

  h1 {
    font-size: 1.5rem;
  }
}
</style>

