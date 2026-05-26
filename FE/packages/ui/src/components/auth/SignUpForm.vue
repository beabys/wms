<script setup lang="ts">
import { ref, watch } from 'vue'
defineOptions({ name: 'SignUpForm' })

const props = defineProps<{
  validating?: boolean
  success?: boolean
  error?: string | null
}>()

const emit = defineEmits<{
  submit: [name: string, email: string, password: string, confirmPassword: string]
}>()

const name = ref('')
const email = ref('')
const password = ref('')
const confirmPassword = ref('')
const fieldErrors = ref<Record<string, string>>({})

watch([name, email, password, confirmPassword], () => {
  const errors: Record<string, string> = {}
  if (password.value && password.value.length < 8) {
    errors.password = 'Min 8 characters'
  }
  if (confirmPassword.value && password.value !== confirmPassword.value) {
    errors.confirmPassword = 'Passwords do not match'
  }
  if (email.value && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email.value)) {
    errors.email = 'Invalid email format'
  }
  if (!name.value.trim()) {
    errors.name = 'Name is required'
  }
  fieldErrors.value = errors
})

function onSubmit(): void {
  if (Object.keys(fieldErrors.value).length > 0) return
  emit('submit', name.value, email.value, password.value, confirmPassword.value)
}
</script>

<template>
  <form class="space-y-4" @submit.prevent="onSubmit">
    <h2 class="text-xl font-semibold text-gray-800 dark:text-gray-100">Create Account</h2>

    <div v-if="success" class="text-sm text-green-700 bg-green-50 border border-green-200 rounded px-3 py-2 dark:text-green-300 dark:bg-green-900/30 dark:border-green-800">
      Account created successfully! Check your email to confirm.
    </div>

    <div>
      <label class="block text-sm font-medium text-gray-700 mb-1 dark:text-gray-300">Full Name</label>
      <input
        v-model="name"
        type="text"
        required
        placeholder="Jane Doe"
        class="w-full px-3 py-2 border rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 border-gray-300 dark:border-gray-600"
        :class="{ 'border-red-400': fieldErrors.name }"
        :disabled="validating || success"
      />
      <p v-if="fieldErrors.name" class="text-xs text-red-500 dark:text-red-400 mt-1">{{ fieldErrors.name }}</p>
    </div>

    <div>
      <label class="block text-sm font-medium text-gray-700 mb-1 dark:text-gray-300">Email</label>
      <input
        v-model="email"
        type="email"
        required
        placeholder="you@example.com"
        class="w-full px-3 py-2 border rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 border-gray-300 dark:border-gray-600"
        :class="{ 'border-red-400': fieldErrors.email }"
        :disabled="validating || success"
      />
      <p v-if="fieldErrors.email" class="text-xs text-red-500 dark:text-red-400 mt-1">{{ fieldErrors.email }}</p>
    </div>

    <div>
      <label class="block text-sm font-medium text-gray-700 mb-1 dark:text-gray-300">Password</label>
      <input
        v-model="password"
        type="password"
        required
        placeholder="Min 8 characters"
        class="w-full px-3 py-2 border rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 border-gray-300 dark:border-gray-600"
        :class="{ 'border-red-400': fieldErrors.password }"
        :disabled="validating || success"
      />
      <p v-if="fieldErrors.password" class="text-xs text-red-500 dark:text-red-400 mt-1">{{ fieldErrors.password }}</p>
    </div>

    <div>
      <label class="block text-sm font-medium text-gray-700 mb-1 dark:text-gray-300">Confirm Password</label>
      <input
        v-model="confirmPassword"
        type="password"
        required
        placeholder="Repeat password"
        class="w-full px-3 py-2 border rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 border-gray-300 dark:border-gray-600"
        :class="{ 'border-red-400': fieldErrors.confirmPassword }"
        :disabled="validating || success"
      />
      <p v-if="fieldErrors.confirmPassword" class="text-xs text-red-500 dark:text-red-400 mt-1">{{ fieldErrors.confirmPassword }}</p>
    </div>

    <button
      type="submit"
      :disabled="validating || success"
      class="w-full py-2 px-4 bg-blue-600 text-white rounded-lg font-medium text-sm hover:bg-blue-700 disabled:opacity-60 disabled:cursor-not-allowed flex items-center justify-center gap-2"
    >
      <span v-if="validating" class="inline-block w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
      {{ validating ? 'Validating...' : 'Create Account' }}
    </button>

    <p v-if="error" class="text-sm text-red-600 dark:text-red-400 bg-red-50 dark:bg-red-950 border border-red-200 dark:border-red-800 rounded px-3 py-2">
      {{ error }}
    </p>
  </form>
</template>
