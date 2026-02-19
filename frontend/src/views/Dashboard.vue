<template>
  <div class="dashboard">
    <!-- 系统信息 + 服务控制 -->
    <div class="top-section">
      <div class="card system-info">
        <div class="card-header">
          <h3 class="card-title">📊 系统信息</h3>
        </div>
        <div class="info-grid">
          <div class="info-item">
            <span class="info-label">操作系统</span>
            <span class="info-value">{{ systemInfo.os }} {{ systemInfo.arch }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">主机名</span>
            <span class="info-value">{{ systemInfo.hostname }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">CPU 使用率</span>
            <span class="info-value">{{ systemInfo.cpu_usage?.toFixed(1) || '--' }}%</span>
          </div>
          <div class="info-item">
            <span class="info-label">内存使用</span>
            <span class="info-value">{{ formatBytes(systemInfo.memory_used) }} / {{ formatBytes(systemInfo.memory_total) }}</span>
          </div>
        </div>
      </div>

      <div class="card server-control">
        <div class="card-header">
          <h3 class="card-title">🎮 服务端控制</h3>
          <span :class="['status-badge', serverStatus.running ? 'running' : 'stopped']">
            {{ serverStatus.running ? '运行中' : '已停止' }}
          </span>
        </div>
        <div class="control-buttons">
          <button class="btn-success" @click="handleStart" :disabled="serverStatus.running || loading">
            ▶ 启动
          </button>
          <button class="btn-danger" @click="handleStop" :disabled="!serverStatus.running || loading">
            ⏹ 停止
          </button>
          <button class="btn-primary" @click="handleRestart" :disabled="!serverStatus.running || loading">
            🔄 重启
          </button>
        </div>
      </div>
    </div>

    <!-- SteamCMD 和 服务端管理 -->
    <div class="middle-section">
      <div class="card">
        <div class="card-header">
          <h3 class="card-title">🔧 SteamCMD</h3>
          <span :class="['status-dot', steamCMDStatus.installed ? 'status-online' : 'status-offline']"></span>
        </div>
        <div class="action-buttons">
          <button class="btn-secondary" @click="checkSteamCMD">🔍 检测</button>
          <button class="btn-primary" @click="handleInstallSteamCMD" :disabled="steamCMDStatus.installed || loading">
            📥 下载
          </button>
          <button class="btn-secondary" @click="handleUpdateSteamCMD" :disabled="!steamCMDStatus.installed || loading">
            🔄 更新
          </button>
          <button class="btn-danger" @click="handleDeleteSteamCMD" :disabled="!steamCMDStatus.installed || loading">
            🗑 删除
          </button>
        </div>
      </div>

      <div class="card">
        <div class="card-header">
          <h3 class="card-title">🎯 游戏服务端</h3>
          <span :class="['status-dot', serverStatus.installed ? 'status-online' : 'status-offline']"></span>
        </div>
        <div class="action-buttons">
          <button class="btn-secondary" @click="checkServer">🔍 检测</button>
          <button class="btn-primary" @click="handleInstallServer" :disabled="serverStatus.installed || loading">
            📥 下载
          </button>
          <button class="btn-secondary" @click="handleUpdateServer" :disabled="!serverStatus.installed || loading">
            🔄 更新
          </button>
          <button class="btn-danger" @click="handleDeleteServer" :disabled="!serverStatus.installed || loading">
            🗑 删除
          </button>
        </div>
      </div>
    </div>

    <!-- 日志 -->
    <div class="card log-section">
      <div class="card-header">
        <h3 class="card-title">📋 实时日志</h3>
        <div class="log-controls">
          <button :class="['btn-icon', autoScroll ? 'btn-active' : 'btn-secondary']" @click="toggleAutoScroll" :title="autoScroll ? '自动滚动开启' : '自动滚动关闭'">
            {{ autoScroll ? '📜⬇️' : '📜' }}
          </button>
          <button class="btn-secondary" @click="clearLogs" title="清除日志">🗑️</button>
          <button class="btn-secondary" @click="scrollToBottom" title="滚动到最新">⬇️到底</button>
        </div>
      </div>
      <div class="log-container" ref="logContainer">
        <div v-for="(log, index) in logs" :key="index" class="log-line" v-html="formatLog(log)"></div>
        <div v-if="logs.length === 0" class="log-empty">
          暂无日志
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import * as api from '@/api'
import { AnsiUp } from 'ansi_up'

const ansiUp = new AnsiUp()
const formatLog = (text: string) => {
  return ansiUp.ansi_to_html(text)
}

const systemInfo = ref<any>({})
const serverStatus = ref<any>({ installed: false, running: false })
const steamCMDStatus = ref<any>({ installed: false })
const LOGS_STORAGE_KEY = 'arsm_logs'
const LOGS_MAX_SIZE = 1000
const logs = ref<string[]>([])
const autoScroll = ref(true)

// 从 localStorage 加载历史日志
const loadLogsFromStorage = () => {
  try {
    const stored = localStorage.getItem(LOGS_STORAGE_KEY)
    if (stored) {
      logs.value = JSON.parse(stored)
    }
  } catch (e) {
    console.error('加载日志失败:', e)
  }
}

// 保存日志到 localStorage
const saveLogsToStorage = () => {
  try {
    localStorage.setItem(LOGS_STORAGE_KEY, JSON.stringify(logs.value))
  } catch (e) {
    // 存储空间不足时，清理一半日志
    if (e instanceof Error && e.name === 'QuotaExceededError') {
      logs.value = logs.value.slice(-Math.floor(LOGS_MAX_SIZE / 2))
      localStorage.setItem(LOGS_STORAGE_KEY, JSON.stringify(logs.value))
    }
  }
}

// 添加单条日志
const addLog = (log: string) => {
  logs.value.push(log)
  // 限制数量
  if (logs.value.length > LOGS_MAX_SIZE) {
    logs.value = logs.value.slice(-LOGS_MAX_SIZE)
  }
  // 异步保存，避免频繁 I/O 阻塞
  requestAnimationFrame(saveLogsToStorage)
}
const loading = ref(false)
const logContainer = ref<HTMLElement>()

let ws: WebSocket | null = null

const formatBytes = (bytes: number) => {
  if (!bytes) return '--'
  const gb = bytes / (1024 * 1024 * 1024)
  return gb.toFixed(1) + ' GB'
}

const fetchSystemInfo = async () => {
  try {
    systemInfo.value = await api.getSystemInfo()
  } catch (e) {
    console.error(e)
  }
}

const checkSteamCMD = async () => {
  try {
    steamCMDStatus.value = await api.getSteamCMDStatus()
  } catch (e) {
    console.error(e)
  }
}

const checkServer = async () => {
  try {
    serverStatus.value = await api.getServerStatus()
  } catch (e) {
    console.error(e)
  }
}

const handleStart = async () => {
  loading.value = true
  try {
    await api.startServer()
    await checkServer()
  } catch (e: any) {
    alert(e.message)
  } finally {
    loading.value = false
  }
}

const handleStop = async () => {
  loading.value = true
  try {
    await api.stopServer()
    await checkServer()
  } catch (e: any) {
    alert(e.message)
  } finally {
    loading.value = false
  }
}

const handleRestart = async () => {
  loading.value = true
  try {
    await api.restartServer()
    await checkServer()
  } catch (e: any) {
    alert(e.message)
  } finally {
    loading.value = false
  }
}

const handleInstallSteamCMD = async () => {
  loading.value = true
  try {
    await api.installSteamCMD()
    await checkSteamCMD()
  } catch (e: any) {
    alert(e.message)
  } finally {
    loading.value = false
  }
}

const handleUpdateSteamCMD = async () => {
  loading.value = true
  try {
    await api.updateSteamCMD()
  } catch (e: any) {
    alert(e.message)
  } finally {
    loading.value = false
  }
}

const handleDeleteSteamCMD = async () => {
  if (!confirm('确定要删除 SteamCMD 吗？')) return
  loading.value = true
  try {
    await api.deleteSteamCMD()
    await checkSteamCMD()
  } catch (e: any) {
    alert(e.message)
  } finally {
    loading.value = false
  }
}

const handleInstallServer = async () => {
  loading.value = true
  try {
    await api.installServer()
    await checkServer()
  } catch (e: any) {
    alert(e.message)
  } finally {
    loading.value = false
  }
}

const handleUpdateServer = async () => {
  loading.value = true
  try {
    await api.updateServer()
  } catch (e: any) {
    alert(e.message)
  } finally {
    loading.value = false
  }
}

const handleDeleteServer = async () => {
  if (!confirm('确定要删除游戏服务端吗？')) return
  loading.value = true
  try {
    await api.deleteServer()
    await checkServer()
  } catch (e: any) {
    alert(e.message)
  } finally {
    loading.value = false
  }
}

const clearLogs = () => {
  if (!confirm('确定要清除所有日志吗？')) return
  logs.value = []
  localStorage.removeItem(LOGS_STORAGE_KEY)
}

const toggleAutoScroll = () => {
  autoScroll.value = !autoScroll.value
}

const scrollToBottom = () => {
  nextTick(() => {
    if (logContainer.value) {
      logContainer.value.scrollTop = logContainer.value.scrollHeight
    }
  })
}

const connectWebSocket = () => {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  ws = new WebSocket(`${protocol}//${window.location.host}/ws/logs`)
  
  ws.onmessage = (event) => {
    addLog(event.data)
    if (autoScroll.value) {
      nextTick(() => {
        if (logContainer.value) {
          logContainer.value.scrollTop = logContainer.value.scrollHeight
        }
      })
    }
  }

  ws.onclose = () => {
    setTimeout(connectWebSocket, 3000)
  }
}

onMounted(() => {
  // 先加载历史日志
  loadLogsFromStorage()
  fetchSystemInfo()
  checkSteamCMD()
  checkServer()
  connectWebSocket()
  
  // 每秒刷新系统信息
  const interval = setInterval(fetchSystemInfo, 1000)
  
  onUnmounted(() => {
    clearInterval(interval)
    if (ws) ws.close()
  })
})
</script>

<style scoped>
.dashboard {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.top-section {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

.middle-section {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

.info-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.info-label {
  font-size: 12px;
  color: var(--text-secondary);
}

.info-value {
  font-size: 14px;
  font-weight: 500;
}

.status-badge {
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.status-badge.running {
  background: rgba(34, 197, 94, 0.1);
  color: var(--success-color);
}

.status-badge.stopped {
  background: rgba(100, 116, 139, 0.1);
  color: var(--text-secondary);
}

.control-buttons {
  display: flex;
  gap: 10px;
}

.action-buttons {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.log-section {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.log-container {
  flex: 1;
  min-height: 300px;
  max-height: 400px;
  overflow-y: auto;
  background: #1e1e1e;
  border-radius: 8px;
  padding: 12px;
  font-family: 'Consolas', 'Monaco', monospace;
  font-size: 13px;
  line-height: 1.5;
}

.log-line {
  color: #d4d4d4;
  white-space: pre-wrap;
  word-break: break-all;
}

.log-empty {
  color: #666;
  text-align: center;
  padding: 40px;
}

.log-controls {
  display: flex;
  gap: 8px;
  align-items: center;
}

.btn-icon {
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  padding: 6px 10px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s;
}

.btn-icon:hover {
  background: var(--bg-hover);
}

.btn-active {
  background: var(--primary-color) !important;
  color: white !important;
  border-color: var(--primary-color) !important;
}

button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
