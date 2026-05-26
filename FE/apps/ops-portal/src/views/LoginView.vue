<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { LoginForm } from '@wms/ui'

const router = useRouter()
const loading = ref(false)
const error = ref<string | null>(null)

async function handleLogin(email: string, password: string): Promise<void> {
  loading.value = true
  error.value = null
  try {
    const { useAuthStore } = await import('../stores/auth')
    const auth = useAuthStore()
    await auth.login(email, password)
    router.push('/dashboard')
  } catch (e: any) {
    error.value = e.message || 'Login failed'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-100 dark:bg-gray-950">
    <div class="w-full max-w-sm">
      <LoginForm :loading="loading" :error="error" @submit="handleLogin" />
    </div>
  </div>
</template>
