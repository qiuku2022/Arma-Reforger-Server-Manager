import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import axios from 'axios'

// API 基础 URL
const API_BASE = ''

// 创建 axios 实例
const api = axios.create({
  baseURL: API_BASE,
  timeout: 30000,
})

// 请求拦截器 - 自动添加 Token
api.interceptors.request.use(
  (config) => {
    const token = sessionStorage.getItem('arsm_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器 - 处理 401 错误
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Token 过期或无效，清除登录状态
      sessionStorage.removeItem('arsm_token')
      sessionStorage.removeItem('arsm_user')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export { api }

interface User {
  username: string
  role: string
}

interface AuthStatus {
  enabled: boolean
  authenticated: boolean
  username?: string
  role?: string
  default_password?: boolean
}

export const useAuthStore = defineStore('auth', () => {
  // State
  const token = ref<string | null>(sessionStorage.getItem('arsm_token'))
  const user = ref<User | null>(JSON.parse(sessionStorage.getItem('arsm_user') || 'null'))
  const authEnabled = ref<boolean>(true)
  const loading = ref(false)
  const error = ref<string | null>(null)

  // Getters
  const isAuthenticated = computed(() => !!token.value && !!user.value)
  const isAdmin = computed(() => user.value?.role === 'admin')
  const username = computed(() => user.value?.username || '')

  // Actions
  async function checkAuthStatus(): Promise<AuthStatus> {
    try {
      const response = await api.get('/api/auth/status')
      const data = response.data.data
      authEnabled.value = data.enabled
      
      // 如果认证未启用，自动设置已登录状态
      if (!data.enabled) {
        return data
      }
      
      // 如果服务端显示已登录但本地没有 token，清除状态
      if (data.authenticated && !token.value) {
        clearAuth()
      }
      
      return data
    } catch (err) {
      console.error('检查认证状态失败:', err)
      return { enabled: true, authenticated: false }
    }
  }

  async function login(username: string, password: string): Promise<boolean> {
    loading.value = true
    error.value = null

    try {
      const response = await api.post('/api/auth/login', {
        username,
        password,
      })

      const data = response.data.data

      // 认证未启用
      if (!data.enabled) {
        authEnabled.value = false
        return true
      }

      // 登录成功
      if (data.token) {
        token.value = data.token
        user.value = {
          username: data.username,
          role: data.role,
        }

        // 保存到 sessionStorage
        sessionStorage.setItem('arsm_token', data.token)
        sessionStorage.setItem('arsm_user', JSON.stringify(user.value))

        return true
      }

      error.value = '登录失败'
      return false
    } catch (err: any) {
      error.value = err.response?.data?.message || '登录失败'
      return false
    } finally {
      loading.value = false
    }
  }

  async function logout(): Promise<void> {
    try {
      await api.post('/api/auth/logout')
    } catch (err) {
      console.error('登出失败:', err)
    } finally {
      clearAuth()
    }
  }

  function clearAuth() {
    token.value = null
    user.value = null
    sessionStorage.removeItem('arsm_token')
    sessionStorage.removeItem('arsm_user')
  }

  async function changePassword(oldPassword: string, newPassword: string): Promise<boolean> {
    try {
      await api.post('/api/auth/password', {
        old_password: oldPassword,
        new_password: newPassword,
      })
      return true
    } catch (err: any) {
      error.value = err.response?.data?.message || '修改密码失败'
      return false
    }
  }

  async function changePasswordWithUsername(oldPassword: string, newPassword: string, newUsername?: string): Promise<boolean> {
    try {
      await api.post('/api/auth/password', {
        old_password: oldPassword,
        new_password: newPassword,
        new_username: newUsername || undefined
      })
      return true
    } catch (err: any) {
      error.value = err.response?.data?.message || '更新账号信息失败'
      return false
    }
  }

  // 初始化时检查认证状态
  async function init() {
    const status = await checkAuthStatus()
    
    // 如果认证已启用且本地有 token，验证 token 是否有效
    if (status.enabled && token.value) {
      try {
        await api.get('/api/auth/profile')
      } catch (err) {
        // Token 无效，清除
        clearAuth()
      }
    }
    
    return status
  }

  return {
    // State
    token,
    user,
    authEnabled,
    loading,
    error,
    // Getters
    isAuthenticated,
    isAdmin,
    username,
    // Actions
    login,
    logout,
    checkAuthStatus,
    changePassword,
    changePasswordWithUsername,
    clearAuth,
    init,
  }
})
