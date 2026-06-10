<script setup lang="ts">
import WmsCard from '@/components/common/WmsCard.vue'
import WmsButton from '@/components/common/WmsButton.vue'

withDefaults(defineProps<{
  visible: boolean
  title?: string
  message?: string
  confirmText?: string
  cancelText?: string
  variant?: 'primary' | 'danger'
}>(), {
  visible: false,
  title: 'Confirm',
  message: 'Are you sure?',
  confirmText: 'Confirm',
  cancelText: 'Cancel',
  variant: 'primary',
})

const emit = defineEmits<{
  confirm: []
  cancel: []
}>()
</script>

<template>
  <Teleport to="body">
    <div v-if="visible" class="wms-confirm-overlay" @click.self="emit('cancel')">
      <WmsCard padding="lg" class="wms-confirm-card">
        <template #header>
          <h2 class="wms-confirm-card__title">{{ title }}</h2>
        </template>

        <p v-if="message && !$slots.default" class="wms-confirm-card__message">
          {{ message }}
        </p>
        <div v-else-if="$slots.default" class="wms-confirm-card__message">
          <slot />
        </div>

        <template #footer>
          <div class="wms-confirm-card__actions">
            <WmsButton variant="secondary" @click="emit('cancel')">
              {{ cancelText }}
            </WmsButton>
            <WmsButton :variant="variant" @click="emit('confirm')">
              {{ confirmText }}
            </WmsButton>
          </div>
        </template>
      </WmsCard>
    </div>
  </Teleport>
</template>

<style scoped>
.wms-confirm-overlay {
  position: fixed;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 1000;
}

.wms-confirm-card {
  width: 100%;
  max-width: 28rem;
}

.wms-confirm-card__title {
  margin: 0;
  font-size: var(--text-lg);
  font-weight: var(--font-bold);
  color: var(--color-text-primary);
}

.wms-confirm-card__message {
  margin: 0;
  color: var(--color-text-secondary);
  line-height: var(--leading-relaxed);
}

.wms-confirm-card__actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--spacing-3);
}
</style>
