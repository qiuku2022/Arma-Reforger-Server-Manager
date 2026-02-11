import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'dashboard',
      component: () => import('@/views/Dashboard.vue')
    },
    {
      path: '/config',
      name: 'config',
      component: () => import('@/views/Config.vue')
    },
    {
      path: '/mods',
      name: 'mods',
      component: () => import('@/views/Mods.vue')
    },
    {
      path: '/rcon',
      name: 'rcon',
      component: () => import('@/views/Rcon.vue')
    },
    {
      path: '/settings',
      name: 'settings',
      component: () => import('@/views/Settings.vue')
    }
  ]
})

export default router
