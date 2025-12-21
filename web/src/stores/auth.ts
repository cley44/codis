import { ref, computed } from 'vue'
import { defineStore } from 'pinia'

const BACKEND_BASE_URL = 'http://localhost:80'

interface User {
  ID: string
  Username: string
  GlobalName?: string
  Email?: string
  Avatar?: string
}

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const isAuthenticated = computed(() => user.value !== null)
  const isLoading = ref(false)

  async function checkAuth(): Promise<boolean> {
    isLoading.value = true
    try {
      const res = await fetch(`${BACKEND_BASE_URL}/profil`, {
        credentials: 'include',
      })

      if (res.status === 401 || !res.ok) {
        user.value = null
        return false
      }

      const userData = await res.json()
      user.value = userData as User
      return true
    } catch (err) {
      console.error('Error checking auth:', err)
      user.value = null
      return false
    } finally {
      isLoading.value = false
    }
  }

  async function login(code: string, state: string): Promise<boolean> {
    isLoading.value = true
    try {
      const res = await fetch(`${BACKEND_BASE_URL}/discord/callback`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({ code, state }),
      })

      if (!res.ok) {
        const text = await res.text()
        console.error('Login failed:', res.status, text)
        return false
      }

      const userData = await res.json()
      user.value = userData as User
      return true
    } catch (err) {
      console.error('Error during login:', err)
      return false
    } finally {
      isLoading.value = false
    }
  }

  async function logout(): Promise<void> {
    try {
      // Call backend to clear session
      await fetch(`${BACKEND_BASE_URL}/logout`, {
        method: 'POST',
        credentials: 'include',
      })
    } catch (err) {
      console.error('Error during logout:', err)
      // Continue with local logout even if backend call fails
    } finally {
      // Clear local state
      user.value = null
    }
  }

  async function getDiscordOAuthUrl(): Promise<string | null> {
    try {
      const res = await fetch(`${BACKEND_BASE_URL}/discord/invite_link`, {
        credentials: 'include',
      })

      if (!res.ok) {
        console.error('Failed to get OAuth URL:', res.status)
        return null
      }

      const data = await res.json()
      return (data as { discord_invite_link?: string }).discord_invite_link || null
    } catch (err) {
      console.error('Error getting OAuth URL:', err)
      return null
    }
  }

  return {
    user,
    isAuthenticated,
    isLoading,
    checkAuth,
    login,
    logout,
    getDiscordOAuthUrl,
  }
})

