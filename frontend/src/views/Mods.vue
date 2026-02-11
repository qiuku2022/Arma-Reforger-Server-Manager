<template>
  <div class="mods-page">
    <div class="card add-mod-section">
      <h3 class="section-title">➕ 添加模组</h3>
      <div class="add-mod-form">
        <input v-model="newMod.id" placeholder="模组 ID (Workshop ID)" />
        <input v-model="newMod.name" placeholder="模组名称" />
        <input v-model="newMod.version" placeholder="版本号 (可选，如 1.0.1)" />
        <button class="btn-primary" @click="addMod">添加</button>
      </div>
    </div>

    <div class="mods-lists">
      <div class="card mods-list">
        <div class="card-header">
          <h3 class="section-title">📦 未启用模组</h3>
          <span class="count">{{ disabledMods.length }}</span>
        </div>
        <div class="mods-container">
          <div v-for="mod in disabledMods" :key="mod.id" class="mod-item">
            <span :class="['status-dot', mod.downloaded ? 'status-online' : 'status-offline']"></span>
            <div class="mod-info">
              <span class="mod-name">{{ mod.name || mod.id }}</span>
              <span class="mod-id">{{ mod.id }}</span>
              <span v-if="mod.version" class="mod-version">v{{ mod.version }}</span>
            </div>
            <div class="mod-actions">
              <button class="btn-success btn-sm" @click="enableMod(mod.id)">启用 →</button>
              <button class="btn-danger btn-sm" @click="deleteMod(mod.id)">🗑</button>
            </div>
          </div>
          <div v-if="disabledMods.length === 0" class="empty">暂无未启用模组</div>
        </div>
      </div>

      <div class="card mods-list">
        <div class="card-header">
          <h3 class="section-title">✅ 已启用模组</h3>
          <span class="count">{{ enabledMods.length }}</span>
        </div>
        <div class="mods-container">
          <div v-for="mod in enabledMods" :key="mod.id" class="mod-item">
            <span :class="['status-dot', mod.downloaded ? 'status-online' : 'status-offline']"></span>
            <div class="mod-info">
              <span class="mod-name">{{ mod.name || mod.id }}</span>
              <span class="mod-id">{{ mod.id }}</span>
              <span v-if="mod.version" class="mod-version">v{{ mod.version }}</span>
            </div>
            <div class="mod-actions">
              <button class="btn-secondary btn-sm" @click="disableMod(mod.id)">← 禁用</button>
              <button class="btn-danger btn-sm" @click="deleteMod(mod.id)">🗑</button>
            </div>
          </div>
          <div v-if="enabledMods.length === 0" class="empty">暂无已启用模组</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import * as api from '@/api'

interface Mod {
  id: string
  name: string
  version: string
  enabled: boolean
  downloaded: boolean
}

const mods = ref<Mod[]>([])
const newMod = ref({ id: '', name: '', version: '' })

const enabledMods = computed(() => mods.value.filter(m => m.enabled))
const disabledMods = computed(() => mods.value.filter(m => !m.enabled))

const loadMods = async () => {
  try {
    mods.value = await api.getMods() as Mod[]
  } catch (e) {
    console.error(e)
  }
}

const addMod = async () => {
  if (!newMod.value.id) return
  try {
    await api.addMod({ 
      id: newMod.value.id, 
      name: newMod.value.name || newMod.value.id, 
      version: newMod.value.version || '' 
    })
    newMod.value = { id: '', name: '', version: '' }
    await loadMods()
  } catch (e: any) {
    alert(e.message)
  }
}

const enableMod = async (id: string) => {
  try {
    await api.enableMod(id)
    await loadMods()
  } catch (e: any) {
    alert(e.message)
  }
}

const disableMod = async (id: string) => {
  try {
    await api.disableMod(id)
    await loadMods()
  } catch (e: any) {
    alert(e.message)
  }
}

const deleteMod = async (id: string) => {
  if (!confirm('确定要删除此模组吗？')) return
  try {
    await api.deleteMod(id)
    await loadMods()
  } catch (e: any) {
    alert(e.message)
  }
}

onMounted(loadMods)
</script>

<style scoped>
.mods-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.add-mod-section {
  padding: 20px;
}

.section-title {
  margin: 0;
}

.add-mod-form {
  display: flex;
  gap: 12px;
  margin-top: 16px;
}

.add-mod-form input {
  flex: 1;
}

.mods-lists {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

.mods-list {
  display: flex;
  flex-direction: column;
}

.count {
  background: var(--bg-primary);
  padding: 2px 10px;
  border-radius: 10px;
  font-size: 12px;
}

.mods-container {
  flex: 1;
  min-height: 300px;
  max-height: 500px;
  overflow-y: auto;
}

.mod-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border-bottom: 1px solid var(--border-color);
}

.mod-item:last-child {
  border-bottom: none;
}

.mod-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.mod-name {
  font-weight: 500;
}

.mod-id {
  font-size: 12px;
  color: var(--text-secondary);
  font-family: monospace;
}

.mod-version {
  font-size: 12px;
  color: var(--success-color);
}

.mod-actions {
  display: flex;
  gap: 6px;
}

.btn-sm {
  padding: 4px 10px;
  font-size: 12px;
}

.empty {
  padding: 40px;
  text-align: center;
  color: var(--text-secondary);
}
</style>
