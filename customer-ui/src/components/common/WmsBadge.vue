<script setup lang="ts">
withDefaults(defineProps<{
  variant?: 'success' | 'warning' | 'error' | 'info' | 'neutral'
  size?: 'sm' | 'md'
}>(), {
  variant: 'neutral',
  size: 'sm',
})
</script>

<template>
  <span class="wms-badge" :class="[`wms-badge--${variant}`, `wms-badge--${size}`]">
    <slot />
  </span>
</template>

<style scoped>
.wms-badge {
  display: inline-flex;
  align-items: center;
  font-weight: var(--font-medium);
  border-radius: var(--radius-full);
  white-space: nowrap;
}

.wms-badge--sm {
  padding: var(--spacing-0-5) var(--spacing-2);
  font-size: var(--text-xs);
}

.wms-badge--md {
  padding: var(--spacing-1) var(--spacing-3);
  font-size: var(--text-sm);
}

.wms-badge--success {
  background-color: var(--color-success-bg);
  color: var(--color-success-text);
}

.wms-badge--warning {
  background-color: var(--color-warning-bg);
  color: var(--color-warning-text);
}

.wms-badge--error {
  background-color: var(--color-error-bg);
  color: var(--color-error-text);
}

.wms-badge--info {
  background-color: var(--color-info-bg);
  color: var(--color-info-text);
}

.wms-badge--neutral {
  background-color: var(--color-bg-tertiary);
  color: var(--color-text-secondary);
}
</style>
