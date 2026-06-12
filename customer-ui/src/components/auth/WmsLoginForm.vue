<script setup lang="ts">
import { ref, computed } from 'vue'
import WmsInput from '@/components/common/WmsInput.vue'
import WmsButton from '@/components/common/WmsButton.vue'
import WmsPasswordInput from '@/components/auth/WmsPasswordInput.vue'

withDefaults(defineProps<{
  loading?: boolean
  error?: string
}>(), {
  loading: false,
  error: '',
})

const emit = defineEmits<{
  submit: [email: string, password: string]
}>()

const email = ref('')
const password = ref('')
const touchedEmail = ref(false)
const touchedPassword = ref(false)

const emailError = computed(() => {
  if (!touchedEmail.value && !email.value) return ''
  if (!email.value) return 'Email is required'
  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email.value)) return 'Invalid email format'
  return ''
})

const passwordError = computed(() => {
  if (!touchedPassword.value && !password.value) return ''
  if (!password.value) return 'Password is required'
  if (password.value.length < 8) return 'Password must be at least 8 characters'
  return ''
})

const isValid = computed(() => !emailError.value && !passwordError.value && email.value && password.value)

function onSubmit() {
  touchedEmail.value = true
  touchedPassword.value = true
  if (!isValid.value) return
  emit('submit', email.value, password.value)
}
</script>

<template>
  <form class="wms-login-form" @submit.prevent="onSubmit">
    <h2 class="wms-login-form__title">Welcome Back</h2>
    <p class="wms-login-form__subtitle">Sign in to manage your inventory.</p>

    <div class="wms-login-form__fields">
      <WmsInput
        v-model="email"
        label="Email"
        type="email"
        placeholder="you@example.com"
        :error="emailError"
        @blur="touchedEmail = true"
      />
      <WmsPasswordInput
        v-model="password"
        label="Password"
        placeholder="Enter your password"
        :error="passwordError"
        @blur="touchedPassword = true"
      />
    </div>

    <WmsButton
      type="submit"
      variant="primary"
      size="lg"
      :loading="loading"
      :disabled="loading"
    >
      Sign In
    </WmsButton>

    <p v-if="error" class="wms-login-form__error">{{ error }}</p>
  </form>
</template>

<style scoped>
.wms-login-form {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-5);
  width: 100%;
  max-width: 24rem;
  padding: var(--spacing-8);
  background-color: var(--color-card-bg);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-xl);
}

.wms-login-form__title {
  font-size: var(--text-2xl);
  font-weight: var(--font-bold);
  color: var(--color-text-primary);
  margin: 0;
}

.wms-login-form__subtitle {
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
  margin: 0;
}

.wms-login-form__fields {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-4);
}

.wms-login-form__error {
  font-size: var(--text-sm);
  color: var(--color-error);
  text-align: center;
}
</style>
