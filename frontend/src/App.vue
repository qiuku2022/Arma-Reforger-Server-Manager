<template>
  <div class="app-container">
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
    </aside>
    <main class="main-content">
      <router-view />
    </main>
  </div>
</template>

<script setup lang="ts">
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
</style>
