<script setup lang="ts">
import { ref } from 'vue'
import WmsInput from '@/components/common/WmsInput.vue'

withDefaults(defineProps<{
  label?: string
  placeholder?: string
  error?: string
  modelValue?: string
  disabled?: boolean
  showStrength?: boolean
  strength?: number
  strengthLabel?: string
}>(), {
  label: 'Password',
  placeholder: 'Enter password',
  error: '',
  modelValue: '',
  disabled: false,
  showStrength: false,
  strength: 0,
  strengthLabel: '',
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
  <div class="wms-password-field">
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
    <div v-if="showStrength && modelValue" class="wms-password-strength">
      <div class="wms-password-strength__bar">
        <div
          class="wms-password-strength__fill"
          :class="`wms-password-strength__fill--${strength}`"
          :style="{ width: (strength / 5) * 100 + '%' }"
        />
      </div>
      <span class="wms-password-strength__label">{{ strengthLabel }}</span>
    </div>
  </div>
</template>

<style scoped>
.wms-password-field {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-1);
}

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

.wms-password-strength {
  display: flex;
  align-items: center;
  gap: var(--spacing-2);
}

.wms-password-strength__bar {
  flex: 1;
  height: 4px;
  background-color: var(--color-bg-tertiary);
  border-radius: var(--radius-full);
  overflow: hidden;
}

.wms-password-strength__fill {
  height: 100%;
  border-radius: var(--radius-full);
  transition: width 0.3s ease, background-color 0.3s ease;
}

.wms-password-strength__fill--0 { width: 0; }
.wms-password-strength__fill--1 { background-color: var(--color-error); }
.wms-password-strength__fill--2 { background-color: var(--color-warning); }
.wms-password-strength__fill--3 { background-color: #eab308; }
.wms-password-strength__fill--4 { background-color: var(--color-success); }
.wms-password-strength__fill--5 { background-color: var(--color-primary); }

.wms-password-strength__label {
  font-size: var(--text-xs);
  color: var(--color-text-muted);
  white-space: nowrap;
  min-width: 5rem;
  text-align: right;
}
</style>
