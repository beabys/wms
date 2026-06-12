<script setup lang="ts">
import { computed } from 'vue'
import type { CreateUserRequest } from '@/api/types'
import WmsInput from '@/components/common/WmsInput.vue'
import WmsButton from '@/components/common/WmsButton.vue'

const props = withDefaults(defineProps<{
  modelValue: CreateUserRequest
  loading?: boolean
  error?: string
}>(), {
  loading: false,
  error: '',
})

const emit = defineEmits<{
  'update:modelValue': [value: CreateUserRequest]
  'submit': []
}>()

const roles = [
  { value: 'admin', label: 'Admin' },
  { value: 'warehouse_staff', label: 'Warehouse Staff' },
  { value: 'billing_manager', label: 'Billing Manager' },
]

const emailError = computed(() => {
  if (!props.modelValue.email) return 'Email is required'
  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(props.modelValue.email)) return 'Invalid email format'
  return ''
})

const passwordError = computed(() => {
  if (!props.modelValue.password) return 'Password is required'
  if (props.modelValue.password.length < 8) return 'Password must be at least 8 characters'
  return ''
})

const nameError = computed(() => {
  if (!props.modelValue.name) return 'Name is required'
  return ''
})

const isValid = computed(() => !emailError.value && !passwordError.value && !nameError.value)

function update(field: keyof CreateUserRequest, value: string) {
  emit('update:modelValue', { ...props.modelValue, [field]: value })
}

function onSubmit() {
  if (isValid.value) {
    emit('submit')
  }
}
</script>

<template>
  <form class="user-form" @submit.prevent="onSubmit">
    <div v-if="error" class="user-form__error">{{ error }}</div>

    <WmsInput
      label="Email"
      type="email"
      placeholder="user@example.com"
      :modelValue="modelValue.email"
      :error="emailError"
      @update:modelValue="update('email', $event)"
    />

    <WmsInput
      label="Password"
      type="password"
      placeholder="Min. 8 characters"
      :modelValue="modelValue.password"
      :error="passwordError"
      @update:modelValue="update('password', $event)"
    />

    <WmsInput
      label="Name"
      placeholder="Full name"
      :modelValue="modelValue.name"
      :error="nameError"
      @update:modelValue="update('name', $event)"
    />

    <div class="user-form__field">
      <label class="user-form__label" for="role-select">Role</label>
      <select
        id="role-select"
        class="user-form__select"
        :value="modelValue.role"
        @change="update('role', ($event.target as HTMLSelectElement).value)"
      >
        <option value="" disabled>Select a role</option>
        <option v-for="r in roles" :key="r.value" :value="r.value">
          {{ r.label }}
        </option>
      </select>
    </div>

    <div class="user-form__actions">
      <WmsButton type="submit" :loading="loading" :disabled="!isValid">
        Create User
      </WmsButton>
    </div>
  </form>
</template>

<style scoped>
.user-form {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-4);
}

.user-form__error {
  padding: var(--spacing-3);
  background-color: var(--color-error-bg);
  color: var(--color-error-text);
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
}

.user-form__field {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-1);
}

.user-form__label {
  font-size: var(--text-sm);
  font-weight: var(--font-medium);
  color: var(--color-text-primary);
}

.user-form__select {
  padding: var(--spacing-2-5) var(--spacing-3);
  border: 1px solid var(--color-input-border);
  border-radius: var(--radius-md);
  background-color: var(--color-input-bg);
  color: var(--color-text-primary);
  font-size: var(--text-sm);
  min-height: 2.5rem;
  outline: none;
  transition: border-color 0.15s ease, box-shadow 0.15s ease;
}

.user-form__select:focus {
  border-color: var(--color-input-focus);
  box-shadow: 0 0 0 3px var(--color-primary-light);
}

.user-form__actions {
  display: flex;
  justify-content: flex-end;
  padding-top: var(--spacing-2);
}
</style>
