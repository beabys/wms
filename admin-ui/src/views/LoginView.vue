<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import WmsLoginForm from '@/components/auth/WmsLoginForm.vue'
import WmsThemeToggle from '@/components/layout/WmsThemeToggle.vue'
import { useAuth } from '@/composables/useAuth'

const router = useRouter()
const { login } = useAuth()
const loading = ref(false)
const error = ref('')

async function handleSubmit(email: string, password: string) {
  loading.value = true
  error.value = ''
  try {
    await login(email, password)
    router.push('/dashboard')
  } catch (err: any) {
    error.value = err?.message || 'Login failed. Please try again.'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="login-view">
    <div class="login-view__toggle">
      <WmsThemeToggle />
    </div>
    <div class="login-view__card">
      <WmsLoginForm
        :loading="loading"
        :error="error"
        @submit="handleSubmit"
      />
    </div>
  </div>
</template>

<style scoped>
.login-view {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  padding: var(--spacing-4);
  background-color: var(--color-bg-secondary);
}

.login-view__toggle {
  position: fixed;
  top: var(--spacing-4);
  right: var(--spacing-4);
}

.login-view__card {
  width: 100%;
  max-width: 28rem;
}
</style>
