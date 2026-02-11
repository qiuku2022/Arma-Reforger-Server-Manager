<template>
  <div class="config-page">
    <div class="card-header page-header">
      <h2>⚙️ 服务端配置</h2>
      <div class="header-actions">
        <div class="preset-controls">
          <select v-model="selectedPreset" @change="loadPreset" class="preset-select">
            <option value="">选择预设...</option>
            <option v-for="p in presets" :key="p" :value="p">
              {{ p }} {{ p === defaultPreset ? '(默认)' : '' }}
            </option>
          </select>
          <button class="btn-icon" @click="setAsDefault" title="设为默认" :disabled="!selectedPreset">
            <span :class="selectedPreset && selectedPreset === defaultPreset ? 'star-filled' : 'star-empty'">⭐</span>
          </button>
          <button class="btn-icon btn-danger-icon" @click="deletePreset" title="删除预设" :disabled="!selectedPreset">
            🗑️
          </button>
        </div>
        <div class="separator"></div>
        <button class="btn-secondary" @click="showSavePreset = true">💾 保存预设</button>
        <button class="btn-secondary" @click="importConfig">📥 导入</button>
        <button class="btn-secondary" @click="exportConfig">📤 导出</button>
        <button class="btn-primary" @click="save">保存配置</button>
      </div>
    </div>

    <div class="config-sections">
      <!-- 基础网络配置 -->
      <div class="card">
        <h3 class="section-title">🌐 网络配置</h3>
        <div class="form-row">
          <div class="form-group">
            <label class="form-label">绑定地址</label>
            <input v-model="config.bindAddress" placeholder="留空自动绑定" />
          </div>
          <div class="form-group">
            <label class="form-label">绑定端口</label>
            <input v-model.number="config.bindPort" type="number" />
          </div>
        </div>
        <div class="form-row">
          <div class="form-group">
            <label class="form-label">公网地址</label>
            <input v-model="config.publicAddress" placeholder="服务器公网IP" />
          </div>
          <div class="form-group">
            <label class="form-label">公网端口</label>
            <input v-model.number="config.publicPort" type="number" />
          </div>
        </div>
      </div>

      <!-- A2S 配置 -->
      <div class="card">
        <h3 class="section-title">📡 A2S 查询</h3>
        <div class="form-row">
          <div class="form-group">
            <label class="form-label">A2S 地址</label>
            <input v-model="config.a2s.address" placeholder="留空使用默认" />
          </div>
          <div class="form-group">
            <label class="form-label">A2S 端口</label>
            <input v-model.number="config.a2s.port" type="number" />
          </div>
        </div>
      </div>

      <!-- RCON 配置 -->
      <div class="card">
        <div class="card-header">
          <h3 class="section-title">🔐 RCON</h3>
          <div class="toggle-switch" :class="{ active: rconEnabled }" @click="rconEnabled = !rconEnabled"></div>
        </div>
        <template v-if="rconEnabled">
          <div class="form-row">
            <div class="form-group">
              <label class="form-label">RCON 地址</label>
              <input v-model="config.rcon.address" placeholder="0.0.0.0" />
            </div>
            <div class="form-group">
              <label class="form-label">RCON 端口</label>
              <input v-model.number="config.rcon.port" type="number" />
            </div>
          </div>
          <div class="form-row">
            <div class="form-group">
              <label class="form-label">RCON 密码 (至少3字符，无空格)</label>
              <input v-model="config.rcon.password" type="password" />
            </div>
            <div class="form-group">
              <label class="form-label">权限级别</label>
              <select v-model="config.rcon.permission">
                <option value="admin">admin</option>
                <option value="monitor">monitor</option>
              </select>
            </div>
          </div>
          <div class="form-row">
            <div class="form-group">
              <label class="form-label">RCON 白名单 (逗号分隔)</label>
              <input v-model="whitelistInput" @change="updateWhitelist" placeholder="IP地址..." />
            </div>
            <div class="form-group">
              <label class="form-label">RCON 黑名单 (逗号分隔)</label>
              <input v-model="blacklistInput" @change="updateBlacklist" placeholder="IP地址..." />
            </div>
          </div>
        </template>
      </div>

      <!-- 游戏配置 -->
      <div class="card">
        <h3 class="section-title">🎮 游戏设置</h3>
        <div class="form-group">
          <label class="form-label">服务器名称</label>
          <input v-model="config.game.name" />
        </div>
        <div class="form-row">
          <div class="form-group">
            <label class="form-label">服务器密码</label>
            <input v-model="config.game.password" type="password" />
          </div>
          <div class="form-group">
            <label class="form-label">管理员密码</label>
            <input v-model="config.game.passwordAdmin" type="password" />
          </div>
        </div>
        <div class="form-group">
          <label class="form-label">管理员 ID 列表 (逗号分隔)</label>
          <input v-model="adminsInput" @change="updateAdmins" placeholder="76561198..." />
        </div>
        <div class="form-group">
          <label class="form-label">最大玩家数</label>
          <input v-model.number="config.game.maxPlayers" type="number" />
        </div>
        
        <!-- 场景选择 -->
        <div class="form-group">
          <label class="form-label">任务场景</label>
          <div class="scenario-selector">
            <label class="radio-label">
              <input type="radio" v-model="scenarioMode" value="official" />
              <span>官方场景</span>
            </label>
            <label class="radio-label">
              <input type="radio" v-model="scenarioMode" value="custom" />
              <span>自定义 ID</span>
            </label>
          </div>
          <select v-if="scenarioMode === 'official'" v-model="config.game.scenarioId" class="scenario-select">
            <option v-for="s in scenarios" :key="s.id" :value="s.id">
              {{ s.name }} ({{ s.map }})
            </option>
          </select>
          <input v-else v-model="config.game.scenarioId" placeholder="{GUID}Path/To/Config.conf" />
        </div>

        <!-- 布尔开关 -->
        <div class="toggle-grid">
          <div class="toggle-item">
            <span>服务器可见</span>
            <div class="toggle-switch" :class="{ active: config.game.visible }" @click="config.game.visible = !config.game.visible"></div>
          </div>
          <div class="toggle-item">
            <span>跨平台</span>
            <div class="toggle-switch" :class="{ active: config.game.crossPlatform }" @click="config.game.crossPlatform = !config.game.crossPlatform"></div>
          </div>
        </div>
      </div>

      <!-- 游戏属性 -->
      <div class="card">
        <h3 class="section-title">🎯 游戏属性</h3>
        <div class="form-row">
          <div class="form-group">
            <label class="form-label">最大视距</label>
            <input v-model.number="config.game.gameProperties.serverMaxViewDistance" type="number" />
          </div>
          <div class="form-group">
            <label class="form-label">最小草地距离</label>
            <input v-model.number="config.game.gameProperties.serverMinGrassDistance" type="number" />
          </div>
        </div>
        <div class="form-group">
          <label class="form-label">网络视距</label>
          <input v-model.number="config.game.gameProperties.networkViewDistance" type="number" />
        </div>
        <div class="toggle-grid">
          <div class="toggle-item">
            <span>禁用第三人称</span>
            <div class="toggle-switch" :class="{ active: config.game.gameProperties.disableThirdPerson }" @click="config.game.gameProperties.disableThirdPerson = !config.game.gameProperties.disableThirdPerson"></div>
          </div>
          <div class="toggle-item">
            <span>快速验证</span>
            <div class="toggle-switch" :class="{ active: config.game.gameProperties.fastValidation }" @click="config.game.gameProperties.fastValidation = !config.game.gameProperties.fastValidation"></div>
          </div>
          <div class="toggle-item">
            <span>BattlEye 反作弊</span>
            <div class="toggle-switch" :class="{ active: config.game.gameProperties.battlEye }" @click="config.game.gameProperties.battlEye = !config.game.gameProperties.battlEye"></div>
          </div>
          <div class="toggle-item">
            <span>禁用无线电通话UI</span>
            <div class="toggle-switch" :class="{ active: config.game.gameProperties.VONDisableUI }" @click="config.game.gameProperties.VONDisableUI = !config.game.gameProperties.VONDisableUI"></div>
          </div>
          <div class="toggle-item">
            <span>禁用直接对话UI</span>
            <div class="toggle-switch" :class="{ active: config.game.gameProperties.VONDisableDirectSpeechUI }" @click="config.game.gameProperties.VONDisableDirectSpeechUI = !config.game.gameProperties.VONDisableDirectSpeechUI"></div>
          </div>
          <div class="toggle-item">
            <span>跨阵营语音</span>
            <div class="toggle-switch" :class="{ active: config.game.gameProperties.VONCanTransmitCrossFaction }" @click="config.game.gameProperties.VONCanTransmitCrossFaction = !config.game.gameProperties.VONCanTransmitCrossFaction"></div>
          </div>
          <div class="toggle-item">
            <span>大厅玩家同步</span>
            <div class="toggle-switch" :class="{ active: config.operating.lobbyPlayerSynchronise }" @click="config.operating.lobbyPlayerSynchronise = !config.operating.lobbyPlayerSynchronise"></div>
          </div>
        </div>
        <div class="form-group" style="margin-top: 16px;">
          <label class="form-label">排队人数限制</label>
          <input v-model.number="config.operating.joinQueue.maxSize" type="number" class="short-input" />
        </div>
      </div>
    </div>

    <!-- 保存预设对话框 -->
    <div v-if="showSavePreset" class="modal-overlay" @click.self="showSavePreset = false">
      <div class="modal">
        <h3>保存预设</h3>
        <input v-model="presetName" placeholder="预设名称" />
        <div class="modal-actions">
          <button class="btn-secondary" @click="showSavePreset = false">取消</button>
          <button class="btn-primary" @click="savePreset">保存</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import * as api from '@/api'

const config = ref<any>({
  bindAddress: '',
  bindPort: 2001,
  publicAddress: '',
  publicPort: 2001,
  a2s: { address: '', port: 17777 },
  rcon: { address: '', port: 19999, password: '', permission: 'admin', blacklist: [], whitelist: [] },
  game: {
    name: 'Arma Reforger Server',
    password: '',
    passwordAdmin: '',
    admins: [],
    scenarioId: '',
    maxPlayers: 64,
    visible: true,
    crossPlatform: false,
    supportedPlatforms: ['PLATFORM_PC'],
    gameProperties: {
      serverMaxViewDistance: 2500,
      serverMinGrassDistance: 50,
      networkViewDistance: 1500,
      disableThirdPerson: false,
      fastValidation: true,
      battlEye: true,
      VONDisableUI: false,
      VONDisableDirectSpeechUI: false,
      VONCanTransmitCrossFaction: false
    },
    mods: []
  },
  operating: {
    lobbyPlayerSynchronise: true,
    joinQueue: { maxSize: 64 },
    disableNavmeshStreaming: []
  }
})

const rconEnabled = ref(false)
const scenarios = ref<any[]>([])
const scenarioMode = ref('official')
const presets = ref<string[]>([])
const selectedPreset = ref('')
const defaultPreset = ref('')
const showSavePreset = ref(false)
const presetName = ref('')
const adminsInput = ref('')
const whitelistInput = ref('')
const blacklistInput = ref('')

// 获取全局设置以读取默认预设
const loadGlobalSettings = async () => {
  try {
    const settings = await api.getSettings() as any
    defaultPreset.value = settings.default_preset || ''
    
    // 如果有默认预设，且当前不是手动操作（这里假设页面加载时），自动加载
    // 但为了避免覆盖 live config，只有在用户明确希望或初次加载时才覆盖？
    // 用户需求是 "每次打开网页就默认加载"，所以直接覆盖。
    if (defaultPreset.value) {
      selectedPreset.value = defaultPreset.value
      await loadPreset()
    }
  } catch (e) {
    console.error(e)
  }
}

const loadConfig = async () => {
  try {
    // 先加载实际配置
    config.value = await api.getConfig()
    
    // 然后加载全局设置，如果有默认预设，会覆盖上面的 config
    await loadGlobalSettings()
    
    // 刷新 UI
    updateUIFromConfig()
  } catch (e) {
    console.error(e)
  }
}

const updateUIFromConfig = () => {
  // 判断 RCON 是否启用：rcon 对象存在且 password 不为空
  rconEnabled.value = !!(config.value.rcon?.password)
  adminsInput.value = (config.value.game.admins || []).join(',')
  whitelistInput.value = (config.value.rcon?.whitelist || []).join(',')
  blacklistInput.value = (config.value.rcon?.blacklist || []).join(',')
  // 如果 rcon 不存在，初始化一个默认结构（仅内存中使用，不保存）
  if (!config.value.rcon) {
    config.value.rcon = {
      address: '',
      port: 19999,
      password: '',
      permission: 'admin',
      blacklist: [],
      whitelist: []
    }
  }
  const isOfficial = scenarios.value.some(s => s.id === config.value.game.scenarioId)
  scenarioMode.value = isOfficial ? 'official' : 'custom'
}

const updateAdmins = () => {
  config.value.game.admins = adminsInput.value.split(',').map(s => s.trim()).filter(s => s)
}

const updateWhitelist = () => {
  config.value.rcon.whitelist = whitelistInput.value.split(',').map(s => s.trim()).filter(s => s)
}

const updateBlacklist = () => {
  config.value.rcon.blacklist = blacklistInput.value.split(',').map(s => s.trim()).filter(s => s)
}

const loadScenarios = async () => {
  try {
    scenarios.value = await api.getScenarios() as any[]
  } catch (e) {
    console.error(e)
  }
}

const loadPresets = async () => {
  try {
    presets.value = await api.getPresets()
  } catch (e) {
    console.error(e)
  }
}

const loadPreset = async () => {
  if (!selectedPreset.value) return
  try {
    const presetConfig = await api.getPreset(selectedPreset.value)
    config.value = presetConfig
    updateUIFromConfig()
    // alert(`预设 "${selectedPreset.value}" 加载成功`) // 自动加载时不弹窗？或者仅手动时弹窗？
    // 为了体验，如果是自动加载（在 loadGlobalSettings 里调用的），最好不弹窗。
    // 这里简单处理：如果是用户点击触发 change 事件，会有弹窗需求。
    // 但 loadGlobalSettings 直接调用了 loadPreset。
    // 暂时不改 loadPreset 签名，通过判断 event 是否存在？不行。
    // 简单起见，移除 alert，或者只在非默认加载时 alert。
  } catch (e: any) {
    alert(e.message)
  }
}

const setAsDefault = async () => {
  if (!selectedPreset.value) return
  try {
    const settings = await api.getSettings() as any
    // 如果已经是默认，则取消默认
    if (defaultPreset.value === selectedPreset.value) {
      settings.default_preset = ""
      defaultPreset.value = ""
    } else {
      settings.default_preset = selectedPreset.value
      defaultPreset.value = selectedPreset.value
    }
    await api.saveSettings(settings)
  } catch (e: any) {
    alert(e.message)
  }
}

const deletePreset = async () => {
  if (!selectedPreset.value) return
  if (!confirm(`确定要删除预设 "${selectedPreset.value}" 吗？`)) return
  try {
    await api.deletePreset(selectedPreset.value)
    if (defaultPreset.value === selectedPreset.value) {
        // 如果删除了默认预设，清除设置
        const settings = await api.getSettings() as any
        settings.default_preset = ""
        await api.saveSettings(settings)
        defaultPreset.value = ""
    }
    selectedPreset.value = ""
    await loadPresets()
  } catch (e: any) {
    alert(e.message)
  }
}

const save = async () => {
  try {
    // 根据 rconEnabled 决定是否保留 rcon 段
    const configToSave = { ...config.value }
    if (!rconEnabled.value) {
      // RCON 关闭时，删除 rcon 字段
      configToSave.rcon = null
    } else {
      // RCON 打开时，确保 rcon 对象存在且字段完整
      if (!configToSave.rcon) {
        configToSave.rcon = {
          address: '',
          port: 19999,
          password: '',
          permission: 'admin',
          blacklist: [],
          whitelist: []
        }
      }
    }
    await api.saveConfig(configToSave)
    alert('配置已保存')
  } catch (e: any) {
    alert(e.message)
  }
}

const savePreset = async () => {
  if (!presetName.value) return
  try {
    await api.savePreset(presetName.value, config.value)
    showSavePreset.value = false
    presetName.value = ''
    await loadPresets()
  } catch (e: any) {
    alert(e.message)
  }
}

const importConfig = () => {
  const input = document.createElement('input')
  input.type = 'file'
  input.accept = '.json'
  input.onchange = async (e: any) => {
    const file = e.target.files[0]
    if (!file) return
    const text = await file.text()
    config.value = JSON.parse(text)
    adminsInput.value = (config.value.game.admins || []).join(',')
    whitelistInput.value = (config.value.rcon.whitelist || []).join(',')
    blacklistInput.value = (config.value.rcon.blacklist || []).join(',')
  }
  input.click()
}

const exportConfig = () => {
  api.exportConfig()
}

onMounted(() => {
  loadConfig()
  loadScenarios()
  loadPresets()
})
</script>

<style scoped>
.config-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0;
  margin-bottom: 10px;
}

.header-actions {
  display: flex;
  gap: 10px;
  align-items: center;
}

.preset-controls {
  display: flex;
  gap: 8px;
  align-items: center;
}

.btn-icon {
  background: none;
  border: 1px solid var(--border-color);
  padding: 8px;
  border-radius: 6px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-icon:hover {
  background: var(--bg-hover);
}

.btn-icon:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-danger-icon:hover {
  background: #fee2e2;
  border-color: #fca5a5;
}

.star-filled {
  filter: grayscale(0%);
}

.star-empty {
  filter: grayscale(100%);
  opacity: 0.5;
}

.separator {
  width: 1px;
  height: 24px;
  background: var(--border-color);
  margin: 0 8px;
}

.config-sections {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.section-title {
  margin-bottom: 16px;
  font-size: 16px;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.toggle-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
  margin-top: 16px;
}

.toggle-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  background: var(--bg-primary);
  border-radius: 8px;
}

.scenario-selector {
  display: flex;
  gap: 24px;
  margin-bottom: 12px;
}

.scenario-selector label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  font-size: 14px;
  white-space: nowrap;
}

.scenario-select {
  width: 100%;
}

.preset-select {
  width: 250px;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal {
  background: var(--bg-secondary);
  padding: 24px;
  border-radius: 12px;
  width: 400px;
}

.modal h3 {
  margin-bottom: 16px;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 16px;
}

.short-input {
  width: 120px;
}
</style>
