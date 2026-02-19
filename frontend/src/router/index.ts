import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/Login.vue'),
      meta: { public: true }
    },
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

// 路由守卫 - 认证检查
router.beforeEach(async (to, _from, next) => {
  const authStore = useAuthStore()
  
  // 公开路由直接放行
  if (to.meta.public) {
    next()
    return
  }

  // 检查认证状态
  const status = await authStore.checkAuthStatus()
  
  // 如果认证未启用，直接放行
  if (!status.enabled) {
    next()
    return
  }

  // 需要认证但未登录，跳转到登录页
  if (!authStore.isAuthenticated) {
    next('/login')
    return
  }

  next()
})

export default router
