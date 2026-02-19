<template>
  <div class="settings-page">
    <div class="card">
      <div class="card-header">
        <h3 class="section-title">📂 安装路径</h3>
        <button class="btn-secondary" @click="checkStatus">🔄 检测状态</button>
      </div>
      <div class="form-group">
        <label class="form-label">SteamCMD 路径</label>
        <div class="path-input">
          <input v-model="settings.steamcmd_path" />
          <span :class="['status-indicator', steamcmdInstalled ? 'installed' : 'not-installed']">
            {{ steamcmdInstalled ? '✓ 已安装' : '✗ 未安装' }}
          </span>
        </div>
        <p class="path-hint">{{ isWindows ? '示例: C:\\steamcmd' : '示例: /home/user/steamcmd' }}</p>
      </div>
      <div class="form-group">
        <label class="form-label">游戏服务端路径</label>
        <div class="path-input">
          <input v-model="settings.server_path" />
          <span :class="['status-indicator', serverInstalled ? 'installed' : 'not-installed']">
            {{ serverInstalled ? '✓ 已安装' : '✗ 未安装' }}
          </span>
        </div>
        <p class="path-hint">{{ isWindows ? '示例: C:\\ArmaReforgerServer' : '示例: /home/user/arma-reforger-server' }}</p>
      </div>
      <div class="form-actions">
        <button class="btn-primary" @click="save">保存设置</button>
      </div>
    </div>

    <div class="card">
      <h3 class="section-title">🔐 账号安全</h3>
      <form @submit.prevent="updatePassword" class="password-form">
        <div class="form-group">
          <label class="form-label">修改用户名 (可选)</label>
          <input v-model="profileForm.new_username" placeholder="留空则保持不变" />
        </div>
        <div class="form-group">
          <label class="form-label">当前密码</label>
          <input v-model="profileForm.old_password" type="password" required />
        </div>
        <div class="form-group">
          <label class="form-label">新密码</label>
          <input v-model="profileForm.new_password" type="password" required minlength="6" />
        </div>
        <div class="form-actions">
          <button type="submit" class="btn-primary" :disabled="updating">
            {{ updating ? '更新中...' : '更新账号信息' }}
          </button>
        </div>
      </form>
    </div>

    <div class="card">
      <h3 class="section-title">ℹ️ 关于</h3>
      <div class="about-info">
        <p><strong>ARSM</strong> - Arma Reforger Server Manager</p>
        <p>版本: 1.0.0</p>
        <p>一款轻量级的 Arma Reforger 服务器管理工具</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, reactive } from 'vue'
import * as api from '@/api'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'

const authStore = useAuthStore()
const router = useRouter()

const settings = ref({
  steamcmd_path: '',
  server_path: '',
  default_preset: ''
})

const profileForm = reactive({
  new_username: '',
  old_password: '',
  new_password: ''
})

const updating = ref(false)
const steamcmdInstalled = ref(false)
const serverInstalled = ref(false)
const systemOS = ref('')

const isWindows = computed(() => systemOS.value === 'windows')

const loadSettings = async () => {
  try {
    const data = await api.getSettings() as any
    settings.value = data
  } catch (e) {
    console.error(e)
  }
}

const checkStatus = async () => {
  try {
    const [steamcmd, server, system] = await Promise.all([
      api.getSteamCMDStatus(),
      api.getServerStatus(),
      api.getSystemInfo()
    ]) as any[]
    steamcmdInstalled.value = steamcmd.installed
    serverInstalled.value = server.installed
    systemOS.value = system.os
  } catch (e) {
    console.error(e)
  }
}

const save = async () => {
  try {
    await api.saveSettings(settings.value)
    alert('设置已保存')
  } catch (e: any) {
    alert(e.message)
  }
}

const updatePassword = async () => {
  updating.value = true
  try {
    const success = await authStore.changePasswordWithUsername(
      profileForm.old_password,
      profileForm.new_password,
      profileForm.new_username
    )
    if (success) {
      alert('账号信息已成功更新，请重新登录')
      authStore.clearAuth()
      router.push('/login')
    }
  } catch (e: any) {
    alert(e.message || '更新失败')
  } finally {
    updating.value = false
  }
}

onMounted(() => {
  loadSettings()
  checkStatus()
})
</script>

<style scoped>
.settings-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
  max-width: 800px;
}

.section-title {
  margin: 0;
}

.path-input {
  display: flex;
  gap: 12px;
  align-items: center;
}

.path-input input {
  flex: 1;
}

.status-indicator {
  font-size: 12px;
  padding: 4px 10px;
  border-radius: 4px;
  white-space: nowrap;
}

.status-indicator.installed {
  background: rgba(34, 197, 94, 0.1);
  color: var(--success-color);
}

.status-indicator.not-installed {
  background: rgba(100, 116, 139, 0.1);
  color: var(--text-secondary);
}

.path-hint {
  margin-top: 4px;
  font-size: 12px;
  color: var(--text-secondary);
}

.form-actions {
  margin-top: 16px;
}

.about-info {
  color: var(--text-secondary);
  line-height: 1.8;
}

.about-info strong {
  color: var(--text-primary);
}
</style>
