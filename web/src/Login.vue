<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from './stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const loading = ref(false)
const error = ref<string | null>(null)

async function handleLogin() {
  loading.value = true
  error.value = null

  try {
    const oauthUrl = await authStore.getDiscordOAuthUrl()

    if (!oauthUrl) {
      error.value = 'Failed to get Discord OAuth URL. Please try again.'
      loading.value = false
      return
    }

    // Redirect to Discord OAuth
    window.location.href = oauthUrl
  } catch (err) {
    console.error('Error initiating login:', err)
    error.value = err instanceof Error ? err.message : 'An unexpected error occurred'
    loading.value = false
  }
}
</script>

<template>
  <main class="page">
    <section class="card">
      <div class="header">
        <h1>Welcome to codis</h1>
        <p class="subtitle">Sign in with your Discord account to continue</p>
      </div>

      <div v-if="error" class="error-container">
        <p class="error-message">{{ error }}</p>
      </div>

      <button
        type="button"
        class="btn btn-discord"
        :disabled="loading"
        @click="handleLogin"
      >
        <span v-if="loading">Connecting...</span>
        <span v-else>Login with Discord</span>
      </button>

      <p class="hint">
        By signing in, you agree to authorize codis to access your Discord account information.
      </p>
    </section>
  </main>
</template>

<style scoped>
.page {
  min-height: 100vh;
  display: flex;
  align-items: center;
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

.card {
  width: 100%;
  max-width: 420px;
  background: rgba(15, 23, 42, 0.96);
  border-radius: 1rem;
  padding: 3rem 2.5rem;
  box-shadow:
    0 24px 60px rgba(0, 0, 0, 0.6),
    0 0 0 1px rgba(148, 163, 184, 0.2);
  text-align: center;
}

.header {
  margin-bottom: 2.5rem;
}

h1 {
  font-size: 2rem;
  font-weight: 600;
  margin-bottom: 0.5rem;
  background: linear-gradient(135deg, #e5e7eb, #9ca3af);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.subtitle {
  color: #9ca3af;
  font-size: 0.95rem;
  line-height: 1.5;
}

.error-container {
  margin-bottom: 1.5rem;
  padding: 0.75rem 1rem;
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.3);
  border-radius: 0.75rem;
}

.error-message {
  color: #f97373;
  font-size: 0.9rem;
  margin: 0;
}

.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  width: 100%;
  padding: 0.75rem 1.5rem;
  border-radius: 999px;
  border: none;
  cursor: pointer;
  font-size: 1rem;
  font-weight: 500;
  transition:
    transform 0.08s ease-out,
    box-shadow 0.08s ease-out,
    filter 0.12s ease-out;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-discord {
  background: linear-gradient(135deg, #5865f2, #4752c4);
  color: white;
  box-shadow:
    0 10px 25px rgba(88, 101, 242, 0.45),
    0 0 0 1px rgba(88, 101, 242, 0.6);
}

.btn-discord:hover:not(:disabled) {
  transform: translateY(-1px);
  filter: brightness(1.05);
  box-shadow:
    0 14px 35px rgba(88, 101, 242, 0.6),
    0 0 0 1px rgba(129, 140, 248, 0.9);
}

.btn-discord:active:not(:disabled) {
  transform: translateY(0);
  box-shadow:
    0 6px 18px rgba(71, 82, 196, 0.6),
    0 0 0 1px rgba(88, 101, 242, 0.6);
}

.hint {
  margin-top: 1.5rem;
  font-size: 0.8rem;
  color: #6b7280;
  line-height: 1.5;
}

@media (max-width: 480px) {
  .card {
    padding: 2rem 1.5rem;
  }

  h1 {
    font-size: 1.5rem;
  }
}
</style>

