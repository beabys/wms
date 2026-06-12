<script setup lang="ts">
withDefaults(defineProps<{
  padding?: 'none' | 'sm' | 'md' | 'lg'
}>(), {
  padding: 'md',
})
</script>

<template>
  <div class="wms-card" :class="`wms-card--${padding}`">
    <div v-if="$slots.header" class="wms-card__header">
      <slot name="header" />
    </div>
    <div class="wms-card__body">
      <slot />
    </div>
    <div v-if="$slots.footer" class="wms-card__footer">
      <slot name="footer" />
    </div>
  </div>
</template>

<style scoped>
.wms-card {
  background-color: var(--color-card-bg);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  overflow: hidden;
}

.wms-card--none {
  padding: 0;
}

.wms-card--sm {
  padding: var(--spacing-3);
}

.wms-card--md {
  padding: var(--spacing-5);
}

.wms-card--lg {
  padding: var(--spacing-8);
}

.wms-card__header {
  padding-bottom: var(--spacing-4);
  margin-bottom: var(--spacing-4);
  border-bottom: 1px solid var(--color-border);
}

.wms-card__footer {
  padding-top: var(--spacing-4);
  margin-top: var(--spacing-4);
  border-top: 1px solid var(--color-border);
}

.wms-card__body {
  min-height: 1rem;
}
</style>
