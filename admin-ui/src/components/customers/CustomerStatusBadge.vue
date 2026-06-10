<script setup lang="ts">
import { computed } from 'vue'
import WmsBadge from '@/components/common/WmsBadge.vue'

const props = withDefaults(defineProps<{
  status: string
}>(), {
  status: 'pending',
})

const variantMap: Record<string, 'warning' | 'success' | 'error' | 'neutral'> = {
  pending: 'warning',
  active: 'success',
  rejected: 'error',
  suspended: 'neutral',
}

const variant = computed(() => variantMap[props.status] || 'neutral')

const label = computed(() => props.status.charAt(0).toUpperCase() + props.status.slice(1))
</script>

<template>
  <WmsBadge :variant="variant" size="sm">
    {{ label }}
  </WmsBadge>
</template>
