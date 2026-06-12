<script setup lang="ts">
import { computed } from 'vue'
import WmsCard from '@/components/common/WmsCard.vue'
import WmsButton from '@/components/common/WmsButton.vue'
import type { CustomerResponse } from '@/api/types'

const props = withDefaults(defineProps<{
  visible: boolean
  customer: CustomerResponse | null
  action: 'approve' | 'reject' | 'suspend'
}>(), {
  visible: false,
  action: 'approve',
})

const emit = defineEmits<{
  confirm: []
  cancel: []
}>()

const title = computed(() => {
  switch (props.action) {
    case 'approve': return 'Approve Customer'
    case 'reject': return 'Reject Customer'
    case 'suspend': return 'Suspend Customer'
  }
})

const confirmVariant = computed<'primary' | 'danger'>(() => {
  return props.action === 'approve' ? 'primary' : 'danger'
})

const confirmText = computed(() => {
  switch (props.action) {
    case 'approve': return 'Approve'
    case 'reject': return 'Reject'
    case 'suspend': return 'Suspend'
  }
})

const companyName = computed(() => props.customer?.company_name || 'this customer')
</script>

<template>
  <Teleport to="body">
    <div v-if="visible" class="dialog-overlay" @click.self="emit('cancel')">
      <WmsCard padding="lg" class="dialog-card">
        <template #header>
          <h2 class="dialog-card__title">{{ title }}</h2>
        </template>

        <p class="dialog-card__message">
          <slot>
            Are you sure you want to <strong>{{ action }}</strong> "{{ companyName }}"?
          </slot>
        </p>

        <template #footer>
          <div class="dialog-card__actions">
            <WmsButton variant="secondary" @click="emit('cancel')">
              Cancel
            </WmsButton>
            <WmsButton :variant="confirmVariant" @click="emit('confirm')">
              {{ confirmText }}
            </WmsButton>
          </div>
        </template>
      </WmsCard>
    </div>
  </Teleport>
</template>

<style scoped>
.dialog-overlay {
  position: fixed;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 1000;
}

.dialog-card {
  width: 100%;
  max-width: 28rem;
}

.dialog-card__title {
  margin: 0;
  font-size: var(--text-lg);
  font-weight: var(--font-bold);
  color: var(--color-text-primary);
}

.dialog-card__message {
  margin: 0;
  color: var(--color-text-secondary);
  line-height: var(--leading-relaxed);
}

.dialog-card__actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--spacing-3);
}
</style>
