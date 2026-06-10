<script setup lang="ts">
withDefaults(defineProps<{
  label?: string
  placeholder?: string
  error?: string
  modelValue?: string
  type?: 'text' | 'email' | 'password'
  disabled?: boolean
  readonly?: boolean
}>(), {
  label: '',
  placeholder: '',
  error: '',
  modelValue: '',
  type: 'text',
  disabled: false,
  readonly: false,
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
  'blur': [e: FocusEvent]
}>()

function onInput(e: Event) {
  const target = e.target as HTMLInputElement
  emit('update:modelValue', target.value)
}

function onBlur(e: FocusEvent) {
  emit('blur', e)
}
</script>

<template>
  <div class="wms-input">
    <label v-if="label" class="wms-input__label">{{ label }}</label>
    <div class="wms-input__wrapper" :class="{ 'wms-input__wrapper--error': error }">
      <slot name="prefix" />
      <input
        :class="['wms-input__field', { 'wms-input__field--with-prefix': !!$slots.prefix }]"
        :type="type"
        :placeholder="placeholder"
        :value="modelValue"
        :disabled="disabled"
        :readonly="readonly"
        @input="onInput"
        @blur="onBlur"
      />
      <slot name="suffix" />
    </div>
    <p v-if="error" class="wms-input__error">{{ error }}</p>
  </div>
</template>

<style scoped>
.wms-input {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-1);
  width: 100%;
}

.wms-input__label {
  font-size: var(--text-sm);
  font-weight: var(--font-medium);
  color: var(--color-text-primary);
}

.wms-input__wrapper {
  display: flex;
  align-items: center;
  gap: var(--spacing-2);
  padding: 0 var(--spacing-3);
  border: 1px solid var(--color-input-border);
  border-radius: var(--radius-md);
  background-color: var(--color-input-bg);
  transition: border-color 0.15s ease, box-shadow 0.15s ease;
}

.wms-input__wrapper:focus-within {
  border-color: var(--color-input-focus);
  box-shadow: 0 0 0 3px var(--color-primary-light);
}

.wms-input__wrapper--error {
  border-color: var(--color-error);
}

.wms-input__wrapper--error:focus-within {
  box-shadow: 0 0 0 3px var(--color-error-bg);
}

.wms-input__field {
  flex: 1;
  padding: var(--spacing-2-5) 0;
  border: none;
  outline: none;
  background: transparent;
  color: var(--color-text-primary);
  font-size: var(--text-sm);
  min-height: 2.5rem;
}

.wms-input__field::placeholder {
  color: var(--color-text-muted);
}

.wms-input__field:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.wms-input__field--with-prefix {
  padding-left: 0;
}

.wms-input__error {
  font-size: var(--text-xs);
  color: var(--color-error);
}
</style>
