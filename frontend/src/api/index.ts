const BASE_URL = '/api'

async function request<T>(url: string, options?: RequestInit): Promise<T> {
  const response = await fetch(BASE_URL + url, {
    headers: {
      'Content-Type': 'application/json',
    },
    ...options,
  })
  const data = await response.json()
  if (data.code !== 0) {
    throw new Error(data.message)
  }
  return data.data
}

// 系统信息
export const getSystemInfo = () => request('/system/info')

// SteamCMD
export const getSteamCMDStatus = () => request('/steamcmd/status')
export const installSteamCMD = () => request('/steamcmd/install', { method: 'POST' })
export const updateSteamCMD = () => request('/steamcmd/update', { method: 'POST' })
export const deleteSteamCMD = () => request('/steamcmd', { method: 'DELETE' })

// 服务端
export const getServerStatus = () => request('/server/status')
export const installServer = () => request('/server/install', { method: 'POST' })
export const updateServer = () => request('/server/update', { method: 'POST' })
export const deleteServer = () => request('/server', { method: 'DELETE' })
export const startServer = () => request('/server/start', { method: 'POST' })
export const stopServer = () => request('/server/stop', { method: 'POST' })
export const restartServer = () => request('/server/restart', { method: 'POST' })

// 配置
export const getConfig = () => request('/config')
export const saveConfig = (config: any) => request('/config', { method: 'POST', body: JSON.stringify(config) })
export const getPresets = () => request<string[]>('/config/presets')
export const getPreset = (name: string) => request<any>(`/config/presets/${name}`)
export const savePreset = (name: string, config: any) => request('/config/presets', { method: 'POST', body: JSON.stringify({ name, config }) })
export const deletePreset = (name: string) => request(`/config/presets/${name}`, { method: 'DELETE' })
export const getScenarios = () => request('/config/scenarios')
export const exportConfig = () => window.open(BASE_URL + '/config/export')

// 模组
export const getMods = () => request('/mods')
export const addMod = (mod: any) => request('/mods', { method: 'POST', body: JSON.stringify(mod) })
export const deleteMod = (id: string) => request(`/mods/${id}`, { method: 'DELETE' })
export const enableMod = (id: string) => request(`/mods/${id}/enable`, { method: 'POST' })
export const disableMod = (id: string) => request(`/mods/${id}/disable`, { method: 'POST' })
export const checkModFiles = (id: string) => request(`/mods/${id}/check`)

// RCON
export const getPlayers = () => request('/rcon/players')
export const getRCONStatus = () => request('/rcon/status')
export const getRCONLogs = () => request('/rcon/logs')
export const kickPlayer = (id: string) => request(`/rcon/kick/${id}`, { method: 'POST' })
export const banPlayer = (id: string) => request(`/rcon/ban/${id}`, { method: 'POST' })
export const sendRCONCommand = (command: string) => request('/rcon/command', { method: 'POST', body: JSON.stringify({ command }) })

// 设置
export const getSettings = () => request('/settings')
export const saveSettings = (settings: any) => request('/settings', { method: 'POST', body: JSON.stringify(settings) })
