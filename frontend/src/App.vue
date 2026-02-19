<template>
  <div class="app-container">
    <!-- 未登录时只显示登录页面内容 -->
    <template v-if="!authStore.isAuthenticated && authStore.authEnabled">
      <router-view />
    </template>

    <!-- 已登录或认证未启用时显示完整界面 -->
    <template v-else>
      <aside class="sidebar">
        <div class="logo">
          <span class="logo-icon">🎮</span>
          <span class="logo-text">ARSM</span>
        </div>
        <nav class="nav">
          <router-link to="/" class="nav-item">
            <span class="nav-icon">📊</span>
            <span>仪表盘</span>
          </router-link>
          <router-link to="/config" class="nav-item">
            <span class="nav-icon">⚙️</span>
            <span>服务端配置</span>
          </router-link>
          <router-link to="/mods" class="nav-item">
            <span class="nav-icon">📦</span>
            <span>模组管理</span>
          </router-link>
          <router-link to="/rcon" class="nav-item">
            <span class="nav-icon">👥</span>
            <span>RCON</span>
          </router-link>
          <router-link to="/settings" class="nav-item">
            <span class="nav-icon">🔧</span>
            <span>设置</span>
          </router-link>
        </nav>

        <!-- 用户信息区域 -->
        <div class="user-section" v-if="authStore.authEnabled">
          <div class="user-info">
            <span class="user-avatar">👤</span>
            <div class="user-details">
              <span class="username">{{ authStore.username }}</span>
              <span class="role">{{ authStore.isAdmin ? '管理员' : '用户' }}</span>
            </div>
          </div>
          <button class="logout-btn" @click="handleLogout" title="退出登录">
            <span>🚪</span>
          </button>
        </div>
      </aside>
      <main class="main-content">
        <router-view />
      </main>
    </template>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useAuthStore } from './stores/auth'

const router = useRouter()
const authStore = useAuthStore()

async function handleLogout() {
  if (confirm('确定要退出登录吗？')) {
    await authStore.logout()
    router.push('/login')
  }
}
</script>

<style scoped>
.app-container {
  display: flex;
  min-height: 100vh;
}

.sidebar {
  width: 220px;
  background: var(--bg-secondary);
  border-right: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
}

.logo {
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 10px;
  border-bottom: 1px solid var(--border-color);
}

.logo-icon {
  font-size: 24px;
}

.logo-text {
  font-size: 20px;
  font-weight: bold;
  color: var(--primary-color);
}

.nav {
  padding: 10px;
  display: flex;
  flex-direction: column;
  gap: 4px;
  flex: 1;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 16px;
  border-radius: 8px;
  color: var(--text-secondary);
  text-decoration: none;
  transition: all 0.2s;
}

.nav-item:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}

.nav-item.router-link-active {
  background: var(--primary-color);
  color: white;
}

.nav-icon {
  font-size: 18px;
}

.main-content {
  flex: 1;
  padding: 24px;
  overflow-y: auto;
  background: var(--bg-primary);
}

/* 用户信息区域 */
.user-section {
  padding: 12px;
  border-top: 1px solid var(--border-color);
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 10px;
  flex: 1;
  min-width: 0;
}

.user-avatar {
  font-size: 24px;
  flex-shrink: 0;
}

.user-details {
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.username {
  font-weight: 500;
  color: var(--text-primary);
  font-size: 14px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.role {
  font-size: 12px;
  color: var(--text-secondary);
}

.logout-btn {
  background: transparent;
  border: none;
  color: var(--text-secondary);
  cursor: pointer;
  padding: 8px;
  border-radius: 6px;
  transition: all 0.2s;
  flex-shrink: 0;
}

.logout-btn:hover {
  background: rgba(255, 77, 79, 0.1);
  color: #ff4d4f;
}
</style>
