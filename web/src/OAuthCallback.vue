<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from './stores/auth'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const loading = ref(true)
const error = ref<string | null>(null)

onMounted(async () => {
  // Extract code and state from URL query parameters
  const code = route.query.code as string | undefined
  const state = route.query.state as string | undefined

  if (!code || !state) {
    error.value = 'Missing authorization code or state. Please try logging in again.'
    loading.value = false
    return
  }

  try {
    // Exchange code for session
    const success = await authStore.login(code, state)

    if (success) {
      // Redirect to dashboard on success
      router.push('/dashboard')
    } else {
      error.value = 'Failed to complete login. Please try again.'
      loading.value = false
    }
  } catch (err) {
    console.error('Error during OAuth callback:', err)
    error.value = err instanceof Error ? err.message : 'An unexpected error occurred'
    loading.value = false
  }
})

function handleRetry() {
  router.push('/login')
}
</script>

<template>
  <main class="page">
    <section class="card">
      <div v-if="loading" class="loading-container">
        <div class="spinner"></div>
        <h2>Completing login...</h2>
        <p class="subtitle">Please wait while we authenticate your account</p>
      </div>

      <div v-else-if="error" class="error-container">
        <h2>Login Failed</h2>
        <p class="error-message">{{ error }}</p>
        <button type="button" class="btn" @click="handleRetry">Try Again</button>
      </div>
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

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}

.spinner {
  width: 48px;
  height: 48px;
  border: 3px solid rgba(88, 101, 242, 0.2);
  border-top-color: #5865f2;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

h2 {
  font-size: 1.5rem;
  font-weight: 600;
  margin: 0;
}

.subtitle {
  color: #9ca3af;
  font-size: 0.9rem;
  margin: 0;
}

.error-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}

.error-message {
  color: #f97373;
  font-size: 0.95rem;
  margin: 0;
  padding: 0.75rem 1rem;
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.3);
  border-radius: 0.75rem;
  width: 100%;
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
</style>
