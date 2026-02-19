<template>
  <div class="login-container">
    <div class="login-box">
      <div class="login-header">
        <span class="logo-icon">🎮</span>
        <h1>ARSM</h1>
        <p>Arma Reforger Server Manager</p>
      </div>

      <!-- 加载状态 -->
      <div v-if="checking" class="loading-state">
        <div class="spinner"></div>
        <p>正在检查认证状态...</p>
      </div>

      <!-- 认证未启用提示 -->
      <div v-else-if="!authEnabled" class="info-state">
        <div class="info-icon">✓</div>
        <h2>认证已禁用</h2>
        <p>当前系统未启用登录认证，直接进入...</p>
      </div>

      <!-- 登录表单 -->
      <form v-else @submit.prevent="handleLogin" class="login-form">
        <h2>账号登录</h2>

        <!-- 默认密码警告 -->
        <div v-if="defaultPassword" class="warning-box">
          <strong>⚠️ 安全警告</strong>
          <p>当前正在使用默认密码 (admin/admin)，请登录后立即修改密码！</p>
        </div>

        <div class="form-group">
          <label for="username">用户名</label>
          <input
            id="username"
            v-model="form.username"
            type="text"
            placeholder="请输入用户名"
            required
            autofocus
          />
        </div>

        <div class="form-group">
          <label for="password">密码</label>
          <input
            id="password"
            v-model="form.password"
            type="password"
            placeholder="请输入密码"
            required
          />
        </div>

        <div v-if="error" class="error-message">
          {{ error }}
        </div>

        <button
          type="submit"
          class="login-btn"
          :disabled="loading"
        >
          <span v-if="loading" class="spinner-small"></span>
          <span v-else>登录</span>
        </button>
      </form>

      <div class="login-footer">
        <p>ARSM v1.0.0</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const form = reactive({
  username: '',
  password: '',
})

const loading = ref(false)
const error = ref('')
const checking = ref(true)
const authEnabled = ref(true)
const defaultPassword = ref(false)

onMounted(async () => {
  // 检查认证状态
  const status = await authStore.init()
  
  checking.value = false
  authEnabled.value = status.enabled
  defaultPassword.value = status.default_password || false

  // 如果认证未启用或已登录，跳转到首页
  if (!status.enabled || status.authenticated) {
    router.push('/')
  }
})

async function handleLogin() {
  error.value = ''
  loading.value = true

  const success = await authStore.login(form.username, form.password)

  if (success) {
    router.push('/')
  } else {
    error.value = authStore.error || '登录失败'
  }

  loading.value = false
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, var(--bg-primary) 0%, var(--bg-secondary) 100%);
  padding: 20px;
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
}

.login-box {
  width: 100%;
  max-width: 400px;
  background: var(--bg-secondary);
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
  padding: 40px;
}

.login-header {
  text-align: center;
  margin-bottom: 30px;
}

.logo-icon {
  font-size: 48px;
  display: block;
  margin-bottom: 10px;
}

.login-header h1 {
  font-size: 28px;
  color: var(--primary-color);
  margin: 0;
}

.login-header p {
  color: var(--text-secondary);
  margin: 8px 0 0;
}

.login-form h2 {
  text-align: center;
  margin-bottom: 24px;
  color: var(--text-primary);
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  color: var(--text-secondary);
  font-size: 14px;
}

.form-group input {
  width: 100%;
  padding: 12px 16px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  background: var(--bg-primary);
  color: var(--text-primary);
  font-size: 16px;
  transition: border-color 0.2s;
}

.form-group input:focus {
  outline: none;
  border-color: var(--primary-color);
}

.form-group input::placeholder {
  color: var(--text-secondary);
  opacity: 0.5;
}

.error-message {
  color: #ff4d4f;
  background: rgba(255, 77, 79, 0.1);
  padding: 12px;
  border-radius: 8px;
  margin-bottom: 20px;
  font-size: 14px;
}

.warning-box {
  background: rgba(255, 170, 0, 0.1);
  border: 1px solid rgba(255, 170, 0, 0.3);
  border-radius: 8px;
  padding: 16px;
  margin-bottom: 20px;
}

.warning-box strong {
  color: #ffaa00;
  display: block;
  margin-bottom: 8px;
}

.warning-box p {
  color: var(--text-secondary);
  margin: 0;
  font-size: 14px;
}

.login-btn {
  width: 100%;
  padding: 14px;
  background: var(--primary-color);
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.login-btn:hover:not(:disabled) {
  background: var(--primary-hover);
}

.login-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.loading-state,
.info-state {
  text-align: center;
  padding: 40px 20px;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid var(--border-color);
  border-top-color: var(--primary-color);
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 20px;
}

.spinner-small {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.info-icon {
  width: 60px;
  height: 60px;
  background: var(--success-color, #52c41a);
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 32px;
  margin: 0 auto 20px;
}

.info-state h2 {
  color: var(--text-primary);
  margin-bottom: 10px;
}

.info-state p {
  color: var(--text-secondary);
}

.login-footer {
  text-align: center;
  margin-top: 30px;
  padding-top: 20px;
  border-top: 1px solid var(--border-color);
}

.login-footer p {
  color: var(--text-secondary);
  font-size: 12px;
  margin: 0;
}
</style>
