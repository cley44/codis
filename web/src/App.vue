<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from './stores/auth'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

// Hide header on login and callback pages
const showHeader = computed(() => {
  return route.name !== 'Login' && route.name !== 'OAuthCallback'
})

async function handleLogout() {
  await authStore.logout()
  router.push('/login')
}
</script>

<template>
  <div class="app">
    <header v-if="showHeader" class="header">
      <RouterLink to="/" class="title-link">
        <h1 class="title">codis</h1>
      </RouterLink>
      <nav class="nav">
        <RouterLink v-if="authStore.isAuthenticated" to="/" class="link">Dashboard</RouterLink>
        <RouterLink to="/debug" class="link">Discord OAuth Debug</RouterLink>
        <div v-if="authStore.isAuthenticated" class="user-section">
          <span v-if="authStore.user" class="username">{{
            authStore.user.GlobalName || authStore.user.Username
          }}</span>
          <button type="button" class="btn-logout" @click="handleLogout">Logout</button>
        </div>
      </nav>
    </header>

    <main class="main">
      <RouterView />
    </main>
  </div>
</template>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html,
body {
  margin: 0;
  padding: 0;
  width: 100%;
  height: 100%;
}

#app {
  width: 100%;
  height: 100%;
}
</style>

<style scoped>
.app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: #020617;
  color: #e5e7eb;
  font-family:
    system-ui,
    -apple-system,
    BlinkMacSystemFont,
    'SF Pro Text',
    sans-serif;
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.75rem 1.5rem;
  border-bottom: 1px solid rgba(31, 41, 55, 0.9);
  background: rgba(15, 23, 42, 0.96);
}

.title-link {
  text-decoration: none;
  color: inherit;
}

.title {
  font-size: 1.1rem;
  font-weight: 600;
}

.nav {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.link {
  font-size: 0.9rem;
  color: #9ca3af;
  text-decoration: none;
  transition: color 0.15s ease-out;
}

.link:hover {
  color: #e5e7eb;
}

.user-section {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-left: 0.75rem;
  padding-left: 0.75rem;
  border-left: 1px solid rgba(55, 65, 81, 0.9);
}

.username {
  font-size: 0.9rem;
  color: #9ca3af;
}

.btn-logout {
  font-size: 0.85rem;
  color: #9ca3af;
  background: transparent;
  border: 1px solid rgba(55, 65, 81, 0.9);
  border-radius: 0.5rem;
  padding: 0.35rem 0.75rem;
  cursor: pointer;
  transition:
    color 0.15s ease-out,
    border-color 0.15s ease-out,
    background 0.15s ease-out;
}

.btn-logout:hover {
  color: #e5e7eb;
  border-color: rgba(148, 163, 184, 0.3);
  background: rgba(55, 65, 81, 0.3);
}

.main {
  flex: 1;
}
</style>
