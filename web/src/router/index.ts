import { createRouter, createWebHistory } from 'vue-router'
import DiscordOAuthDebug from '../DiscordOAuthDebug.vue'
import GuildsDashboard from '../GuildsDashboard.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'GuildsDashboard',
      component: GuildsDashboard,
    },
    {
      path: '/dashboard',
      name: 'GuildsDashboardAlt',
      component: GuildsDashboard,
    },
    {
      path: '/discord/callback',
      name: 'DiscordOAuthDebug',
      component: DiscordOAuthDebug,
    },
  ],
})

export default router
