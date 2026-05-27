<script setup lang="ts">
defineOptions({ name: 'LowStockAlert' })

export interface LowStockAlertItem {
  product_id: string
  product_sku: string
  product_name: string
  current_quantity: number
  threshold: number
  status: 'critical' | 'warning'
}

const props = defineProps<{
  alerts: LowStockAlertItem[]
}>()

const emit = defineEmits<{
  view: [productId: string]
}>()

function alertColor(status: 'critical' | 'warning'): string {
  return status === 'critical'
    ? 'bg-red-50 border-red-200 text-red-800 dark:bg-red-900/30 dark:border-red-800 dark:text-red-300'
    : 'bg-yellow-50 border-yellow-200 text-yellow-800 dark:bg-yellow-900/30 dark:border-yellow-800 dark:text-yellow-300'
}

function alertDot(status: 'critical' | 'warning'): string {
  return status === 'critical' ? 'bg-red-500' : 'bg-yellow-500'
}
</script>

<template>
  <div class="space-y-2">
    <div v-if="alerts.length === 0" class="text-sm text-gray-500 dark:text-gray-400 py-2">
      All stock levels are healthy.
    </div>

    <div
      v-for="alert in alerts"
      :key="alert.product_id"
      class="border rounded-lg px-4 py-3 flex items-center justify-between"
      :class="alertColor(alert.status)"
    >
      <div class="flex items-center gap-3">
        <span class="w-2 h-2 rounded-full shrink-0" :class="alertDot(alert.status)"></span>
        <div>
          <p class="text-sm font-medium">{{ alert.product_name }}</p>
          <p class="text-xs opacity-75">
            SKU: {{ alert.product_sku }} &middot; Current: {{ alert.current_quantity }} &middot; Threshold: {{ alert.threshold }}
          </p>
        </div>
      </div>
      <button
        class="text-xs font-medium underline hover:no-underline whitespace-nowrap"
        @click="emit('view', alert.product_id)"
      >
        View Product
      </button>
    </div>
  </div>
</template>
