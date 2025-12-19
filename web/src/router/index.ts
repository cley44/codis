import { createRouter, createWebHistory } from 'vue-router'
import DiscordOAuthDebug from '../DiscordOAuthDebug.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/discord/callback',
      name: 'DiscordOAuthDebug',
      component: DiscordOAuthDebug,
    },
  ],
})

export default router
