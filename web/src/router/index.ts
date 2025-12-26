import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import DiscordOAuthDebug from '../DiscordOAuthDebug.vue'
import GuildsDashboard from '../GuildsDashboard.vue'
import Login from '../Login.vue'
import OAuthCallback from '../OAuthCallback.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: Login,
      meta: { requiresGuest: true },
    },
    {
      path: '/discord/callback',
      name: 'OAuthCallback',
      component: OAuthCallback,
    },
    {
      path: '/',
      name: 'GuildsDashboard',
      component: GuildsDashboard,
      meta: { requiresAuth: true },
    },
    {
      path: '/dashboard',
      name: 'GuildsDashboardAlt',
      component: GuildsDashboard,
      meta: { requiresAuth: true },
    },
    {
      path: '/debug',
      name: 'DiscordOAuthDebug',
      component: DiscordOAuthDebug,
    },
    {
      path: '/guilds/:guildId/workflows',
      name: 'WorkflowList',
      component: () => import('../views/WorkflowList.vue'),
      meta: { requiresAuth: true },
      props: true,
    },
    {
      path: '/guilds/:guildId/workflows/:workflowId',
      name: 'WorkflowEditor',
      component: () => import('../views/WorkflowEditor.vue'),
      meta: { requiresAuth: true },
      props: true,
    },
  ],
})

// Route guards
router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()

  // Check authentication status if not already checked
  if (!authStore.isAuthenticated && authStore.user === null && !authStore.isLoading) {
    await authStore.checkAuth()
  }

  // Protected routes - require authentication
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next({ name: 'Login', query: { redirect: to.fullPath } })
    return
  }

  // Guest-only routes - redirect if already authenticated
  if (to.meta.requiresGuest && authStore.isAuthenticated) {
    next({ name: 'GuildsDashboard' })
    return
  }

  next()
})

export default router
