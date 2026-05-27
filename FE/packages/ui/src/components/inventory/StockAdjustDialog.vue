<script setup lang="ts">
import { ref } from 'vue'
defineOptions({ name: 'StockAdjustDialog' })

export interface StockAdjustData {
  delta: number
  reason: string
  notes: string
}

const props = defineProps<{
  productSku: string
  productName: string
  currentQuantity: number
  loading?: boolean
}>()

const emit = defineEmits<{
  confirm: [data: StockAdjustData]
  cancel: []
}>()

const delta = ref<number | null>(null)
const reason = ref('manual')
const notes = ref('')

const reasons = [
  { value: 'inbound_approval', label: 'Inbound Approval' },
  { value: 'manual', label: 'Manual Adjustment' },
  { value: 'return', label: 'Return' },
  { value: 'disposal', label: 'Disposal' },
]

function handleConfirm(): void {
  if (delta.value === null || delta.value === 0) return
  emit('confirm', {
    delta: delta.value,
    reason: reason.value,
    notes: notes.value,
  })
}
</script>

<template>
  <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/40">
    <div class="bg-white rounded-lg shadow-xl w-full max-w-md mx-4 dark:bg-gray-800">
      <div class="px-6 py-4 border-b dark:border-gray-700">
        <h3 class="text-lg font-semibold text-gray-800 dark:text-gray-100">Adjust Stock</h3>
      </div>

      <div class="px-6 py-4 space-y-4">
        <div class="text-sm text-gray-600 dark:text-gray-300">
          <p><strong>Product:</strong> {{ productSku }} - {{ productName }}</p>
          <p><strong>Current Quantity:</strong> {{ currentQuantity }}</p>
        </div>

        <div>
          <label class="block text-xs text-gray-500 mb-1 dark:text-gray-400">Delta (positive to add, negative to remove)</label>
          <input
            v-model.number="delta"
            type="number"
            placeholder="e.g. 10 or -5"
            class="w-full px-2 py-1.5 border rounded text-sm dark:bg-gray-700 dark:text-gray-100 dark:border-gray-600"
          />
        </div>

        <div>
          <label class="block text-xs text-gray-500 mb-1 dark:text-gray-400">Reason</label>
          <select
            v-model="reason"
            class="w-full px-2 py-1.5 border rounded text-sm dark:bg-gray-700 dark:text-gray-100 dark:border-gray-600"
          >
            <option v-for="r in reasons" :key="r.value" :value="r.value">{{ r.label }}</option>
          </select>
        </div>

        <div>
          <label class="block text-xs text-gray-500 mb-1 dark:text-gray-400">Notes</label>
          <textarea
            v-model="notes"
            rows="3"
            placeholder="Optional notes"
            class="w-full px-2 py-1.5 border rounded text-sm dark:bg-gray-700 dark:text-gray-100 dark:border-gray-600"
          ></textarea>
        </div>
      </div>

      <div class="px-6 py-4 border-t dark:border-gray-700 flex justify-end gap-3">
        <button
          class="px-4 py-2 border border-gray-300 rounded text-sm text-gray-600 hover:bg-gray-50 dark:border-gray-600 dark:text-gray-300 dark:hover:bg-gray-700"
          @click="emit('cancel')"
        >
          Cancel
        </button>
        <button
          :disabled="loading || delta === null || delta === 0"
          class="px-6 py-2 bg-blue-600 text-white rounded text-sm font-medium hover:bg-blue-700 disabled:opacity-60"
          @click="handleConfirm"
        >
          {{ loading ? 'Adjusting...' : 'Confirm Adjustment' }}
        </button>
      </div>
    </div>
  </div>
</template>
