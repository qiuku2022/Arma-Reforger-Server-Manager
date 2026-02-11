<template>
  <div class="rcon-page">
    <div class="players-section card">
      <div class="card-header">
        <h3 class="section-title">👥 玩家列表</h3>
        <button class="btn-secondary" @click="loadPlayers">🔄 刷新</button>
      </div>
      <div class="players-container">
        <div v-for="player in players" :key="player.id" class="player-item">
          <span 
            class="status-dot" 
            :class="player.online ? 'status-online' : 'status-offline'"
          ></span>
          <div class="player-info">
            <span class="player-name">{{ player.name || 'Unknown' }}</span>
            <span class="player-id">ID: {{ player.id }}</span>
          </div>
          <span class="player-time">{{ formatTime(player.online_time) }}</span>
          <div class="player-actions">
            <button 
              class="btn-danger btn-sm" 
              @click="kick(player.id)" 
              :disabled="!player.online"
            >踢出</button>
            <button class="btn-danger btn-sm" @click="ban(player.id)">拉黑</button>
          </div>
        </div>
        <div v-if="!players || players.length === 0" class="empty">
          暂无玩家数据
        </div>
      </div>
    </div>
    
    <div class="console-section card">
      <h3 class="section-title">💻 RCON 控制台</h3>
      <div class="console-input">
        <input 
          v-model="command" 
          placeholder="输入 RCON 命令..." 
          @keyup.enter="sendCommand"
        />
        <button class="btn-primary" @click="sendCommand">发送</button>
      </div>
      <div class="console-output" ref="consoleOutputEl">
        <div 
          v-for="(line, index) in consoleOutput" 
          :key="index" 
          class="console-line"
        >{{ line }}</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick } from 'vue'
import * as api from '@/api'

interface Player {
  id: string
  name: string
  online: boolean
  online_time: number
}

const players = ref<Player[]>([])
const command = ref('')
const consoleOutput = ref<string[]>([])
const consoleOutputEl = ref<HTMLElement | null>(null)

const formatTime = (seconds: number): string => {
  if (!seconds || seconds <= 0) return '--'
  const hours = Math.floor(seconds / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)
  return `${hours}h ${minutes}m`
}

const scrollToBottom = () => {
  nextTick(() => {
    if (consoleOutputEl.value) {
      consoleOutputEl.value.scrollTop = consoleOutputEl.value.scrollHeight
    }
  })
}

const addConsoleLine = (line: string) => {
  consoleOutput.value.push(line)
  // 限制数量
  if (consoleOutput.value.length > 100) {
    consoleOutput.value = consoleOutput.value.slice(-100)
  }
  scrollToBottom()
}

const loadPlayers = async () => {
  try {
    const data = await api.getPlayers()
    // 确保数据是数组
    if (Array.isArray(data)) {
      players.value = data
    } else {
      players.value = []
      console.warn('Invalid players data:', data)
    }
  } catch (e) {
    console.error('Load players error:', e)
    players.value = []
    addConsoleLine(`[Error] 无法获取玩家列表: ${e}`)
  }
}

const kick = async (id: string) => {
  if (!id || !confirm('确定要踢出该玩家吗？')) return
  try {
    await api.kickPlayer(id)
    addConsoleLine(`[KICK] Player ${id} kicked`)
    await loadPlayers()
  } catch (e: any) {
    const msg = e?.message || String(e)
    addConsoleLine(`[Error] 踢出失败: ${msg}`)
  }
}

const ban = async (id: string) => {
  if (!id || !confirm('确定要封禁该玩家吗？')) return
  try {
    await api.banPlayer(id)
    addConsoleLine(`[BAN] Player ${id} banned`)
    await loadPlayers()
  } catch (e: any) {
    const msg = e?.message || String(e)
    addConsoleLine(`[Error] 封禁失败: ${msg}`)
  }
}

const sendCommand = async () => {
  const cmd = command.value.trim()
  if (!cmd) return
  
  try {
    addConsoleLine(`> ${cmd}`)
    const result: any = await api.sendRCONCommand(cmd)
    const response = result?.response || 'OK'
    addConsoleLine(response)
    command.value = ''
  } catch (e: any) {
    const msg = e?.message || String(e)
    addConsoleLine(`[Error] ${msg}`)
  }
}

onMounted(() => {
  loadPlayers()
})
</script>

<style scoped>
.rcon-page {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
  height: calc(100vh - 100px);
}

.section-title {
  margin: 0;
}

.players-section {
  display: flex;
  flex-direction: column;
}

.players-container {
  flex: 1;
  overflow-y: auto;
}

.player-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border-bottom: 1px solid var(--border-color);
}

.player-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.player-name {
  font-weight: 500;
}

.player-id {
  font-size: 12px;
  color: var(--text-secondary);
  font-family: monospace;
}

.player-time {
  font-size: 12px;
  color: var(--text-secondary);
  min-width: 50px;
  text-align: right;
}

.player-actions {
  display: flex;
  gap: 6px;
}

.btn-sm {
  padding: 4px 10px;
  font-size: 12px;
}

.console-section {
  display: flex;
  flex-direction: column;
}

.console-input {
  display: flex;
  gap: 10px;
  margin-bottom: 16px;
}

.console-input input {
  flex: 1;
}

.console-output {
  flex: 1;
  min-height: 200px;
  background: #1e1e1e;
  border-radius: 8px;
  padding: 12px;
  overflow-y: auto;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 13px;
}

.console-line {
  color: #d4d4d4;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-all;
}

.empty {
  padding: 40px;
  text-align: center;
  color: var(--text-secondary);
}
</style>
