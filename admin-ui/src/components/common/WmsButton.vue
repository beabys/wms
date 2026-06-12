<script setup lang="ts">
withDefaults(defineProps<{
  variant?: 'primary' | 'secondary' | 'ghost' | 'danger'
  size?: 'sm' | 'md' | 'lg'
  loading?: boolean
  disabled?: boolean
  type?: 'button' | 'submit'
}>(), {
  variant: 'primary',
  size: 'md',
  loading: false,
  disabled: false,
  type: 'button',
})

const emit = defineEmits<{
  click: [e: MouseEvent]
}>()
</script>

<template>
  <button
    :class="['wms-btn', `wms-btn--${variant}`, `wms-btn--${size}`]"
    :disabled="disabled || loading"
    :type="type"
    @click="emit('click', $event)"
  >
    <span v-if="loading" class="wms-btn__spinner" />
    <slot />
  </button>
</template>

<style scoped>
.wms-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--spacing-2);
  border: 1px solid transparent;
  border-radius: var(--radius-md);
  font-weight: var(--font-medium);
  cursor: pointer;
  transition: background-color 0.15s ease, color 0.15s ease, border-color 0.15s ease, opacity 0.15s ease;
  white-space: nowrap;
  line-height: var(--leading-tight);
}

.wms-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* Sizes */
.wms-btn--sm {
  padding: var(--spacing-1-5) var(--spacing-3);
  font-size: var(--text-sm);
  height: 2rem;
}

.wms-btn--md {
  padding: var(--spacing-2) var(--spacing-4);
  font-size: var(--text-sm);
  height: 2.5rem;
}

.wms-btn--lg {
  padding: var(--spacing-3) var(--spacing-6);
  font-size: var(--text-base);
  height: 3rem;
}

/* Variants */
.wms-btn--primary {
  background-color: var(--color-primary);
  color: #ffffff;
  border-color: var(--color-primary);
}

.wms-btn--primary:hover:not(:disabled) {
  background-color: var(--color-primary-hover);
  border-color: var(--color-primary-hover);
}

.wms-btn--secondary {
  background-color: var(--color-bg-secondary);
  color: var(--color-text-primary);
  border-color: var(--color-border);
}

.wms-btn--secondary:hover:not(:disabled) {
  background-color: var(--color-bg-tertiary);
  border-color: var(--color-text-muted);
}

.wms-btn--ghost {
  background-color: transparent;
  color: var(--color-text-primary);
  border-color: transparent;
}

.wms-btn--ghost:hover:not(:disabled) {
  background-color: var(--color-bg-secondary);
}

.wms-btn--danger {
  background-color: var(--color-error);
  color: #ffffff;
  border-color: var(--color-error);
}

.wms-btn--danger:hover:not(:disabled) {
  opacity: 0.9;
}

/* Spinner */
.wms-btn__spinner {
  width: 1rem;
  height: 1rem;
  border: 2px solid currentColor;
  border-right-color: transparent;
  border-radius: var(--radius-full);
  animation: spin 0.6s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
