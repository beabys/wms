<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { AuthClient } from '@wms/api-client'

const router = useRouter()
const authClient = new AuthClient()
const email = ref('')
const password = ref('')
const loading = ref(false)
const error = ref<string | null>(null)
const success = ref(false)

async function handleRegister(): Promise<void> {
  loading.value = true
  error.value = null
  try {
    await authClient.register(email.value, password.value)
    success.value = true
    setTimeout(() => router.push('/login'), 1500)
  } catch (e: any) {
    error.value = e.message || 'Registration failed'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-100 dark:bg-gray-950">
    <div class="w-full max-w-sm bg-white dark:bg-gray-900 rounded-lg shadow p-8">
      <h1 class="text-xl font-semibold mb-6">Create Account</h1>

      <div v-if="success" class="text-green-600 mb-4">
        Account created! Redirecting to login...
      </div>

      <form v-else @submit.prevent="handleRegister">
        <div class="mb-4">
          <label class="block text-sm font-medium text-gray-700 mb-1">Email</label>
          <input v-model="email" type="email" required
            class="w-full border rounded px-3 py-2" placeholder="admin@wms.com" />
        </div>
        <div class="mb-4">
          <label class="block text-sm font-medium text-gray-700 mb-1">Password</label>
          <input v-model="password" type="password" required
            class="w-full border rounded px-3 py-2" placeholder="password123" />
        </div>
        <p v-if="error" class="text-red-600 text-sm mb-4">{{ error }}</p>
        <button type="submit" :disabled="loading"
          class="w-full bg-blue-600 text-white rounded py-2 hover:bg-blue-700 disabled:opacity-50">
          {{ loading ? 'Registering...' : 'Register' }}
        </button>
      </form>

      <p class="text-sm text-gray-500 mt-4 text-center">
        Already have an account?
        <router-link to="/login" class="text-blue-600">Login</router-link>
      </p>
    </div>
  </div>
</template>
