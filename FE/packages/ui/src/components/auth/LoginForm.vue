<script setup lang="ts">
import { ref } from 'vue'
defineOptions({ name: 'LoginForm' })

const props = defineProps<{
  loading?: boolean
  error?: string | null
}>()

const emit = defineEmits<{
  submit: [email: string, password: string]
}>()

const email = ref('')
const password = ref('')

function onSubmit(): void {
  if (!email.value || !password.value) return
  emit('submit', email.value, password.value)
}
</script>

<template>
  <form class="space-y-4" @submit.prevent="onSubmit">
    <h2 class="text-xl font-semibold text-gray-800 dark:text-gray-100">Sign In</h2>

    <div>
      <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Email</label>
      <input
        v-model="email"
        type="email"
        required
        placeholder="you@example.com"
        class="w-full px-3 py-2 border rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 border-gray-300 dark:border-gray-600"
        :disabled="loading"
      />
    </div>

    <div>
      <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Password</label>
      <input
        v-model="password"
        type="password"
        required
        placeholder="&#8226;&#8226;&#8226;&#8226;&#8226;&#8226;&#8226;&#8226;"
        class="w-full px-3 py-2 border rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 border-gray-300 dark:border-gray-600"
        :disabled="loading"
      />
    </div>

    <button
      type="submit"
      :disabled="loading"
      class="w-full py-2 px-4 bg-blue-600 text-white rounded-lg font-medium text-sm hover:bg-blue-700 disabled:opacity-60 disabled:cursor-not-allowed flex items-center justify-center gap-2"
    >
      <span v-if="loading" class="inline-block w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
      {{ loading ? 'Signing in...' : 'Sign In' }}
    </button>

    <p v-if="error" class="text-sm text-red-600 dark:text-red-400 bg-red-50 dark:bg-red-950 border border-red-200 dark:border-red-800 rounded px-3 py-2">
      {{ error }}
    </p>
  </form>
</template>
