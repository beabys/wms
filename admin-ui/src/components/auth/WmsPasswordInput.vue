<script setup lang="ts">
import { ref } from 'vue'
import WmsInput from '@/components/common/WmsInput.vue'

withDefaults(defineProps<{
  label?: string
  placeholder?: string
  error?: string
  modelValue?: string
  disabled?: boolean
}>(), {
  label: 'Password',
  placeholder: 'Enter password',
  error: '',
  modelValue: '',
  disabled: false,
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
  blur: [e: FocusEvent]
}>()

const showPassword = ref(false)

function onBlur(e: FocusEvent) {
  emit('blur', e)
}
</script>

<template>
  <WmsInput
    :label="label"
    :placeholder="placeholder"
    :error="error"
    :type="showPassword ? 'text' : 'password'"
    :model-value="modelValue"
    :disabled="disabled"
    @update:model-value="emit('update:modelValue', $event)"
    @blur="onBlur"
  >
    <template #suffix>
      <button
        type="button"
        class="wms-password-toggle"
        :title="showPassword ? 'Hide password' : 'Show password'"
        @click="showPassword = !showPassword"
      >
        <!-- Eye icon when hidden -->
        <svg
          v-if="!showPassword"
          class="wms-password-toggle__icon"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" />
          <circle cx="12" cy="12" r="3" />
        </svg>
        <!-- Eye-off icon when shown -->
        <svg
          v-else
          class="wms-password-toggle__icon"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94" />
          <path d="M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19" />
          <line x1="1" y1="1" x2="23" y2="23" />
        </svg>
      </button>
    </template>
  </WmsInput>
</template>

<style scoped>
.wms-password-toggle {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  background: none;
  border: none;
  color: var(--color-text-muted);
  cursor: pointer;
  padding: var(--spacing-1);
  border-radius: var(--radius-sm);
  transition: color 0.15s ease;
}

.wms-password-toggle:hover {
  color: var(--color-text-primary);
}

.wms-password-toggle__icon {
  width: 1.125rem;
  height: 1.125rem;
}
</style>
