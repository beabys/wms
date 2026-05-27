<script setup lang="ts">
import { computed } from 'vue'
defineOptions({ name: 'InboundDetail' })

export interface InboundDetailShipment {
  id: string
  customerName: string
  skuCount: number
  totalUnits: number
  status: string
  expectedDate: string
}

import StatusBadge from '../common/StatusBadge.vue'

const props = withDefaults(defineProps<{
  shipment: InboundDetailShipment | null
  currentStep: number
  loading?: boolean
}>(), {
  loading: false
})

const emit = defineEmits<{
  stepChange: [step: number]
}>()

const steps = ['Submitted', 'Inspection', 'Approve / Flag / Hold', 'Release']

const canAdvance = computed(() => props.currentStep < steps.length - 1 && !props.loading)
const canGoBack = computed(() => props.currentStep > 0 && !props.loading)
</script>

<template>
  <div class="bg-white border rounded-lg p-6 space-y-6 dark:bg-gray-800 dark:border-gray-700">
    <div v-if="loading" class="text-center text-gray-400 dark:text-gray-400 py-8">Loading...</div>

    <template v-else-if="shipment">
      <div class="flex items-center justify-between">
        <div>
          <h3 class="text-lg font-semibold text-gray-800 dark:text-gray-100">Inbound #{{ shipment.id }}</h3>
          <p class="text-sm text-gray-500 dark:text-gray-400">{{ shipment.customerName }} &middot; {{ shipment.skuCount }} SKUs &middot; {{ shipment.totalUnits }} units</p>
        </div>
        <StatusBadge :status="shipment.status as any" />
      </div>

      <div class="flex items-center gap-1">
        <div
          v-for="(step, idx) in steps"
          :key="idx"
          class="flex items-center"
        >
          <div
            class="flex items-center gap-2 px-3 py-1.5 rounded text-sm font-medium"
            :class="{
              'bg-blue-600 text-white': idx === currentStep,
              'bg-green-100 text-green-700 dark:bg-green-900/50 dark:text-green-300': idx < currentStep,
              'bg-gray-100 text-gray-400 dark:bg-gray-700 dark:text-gray-500': idx > currentStep
            }"
          >
            <span v-if="idx < currentStep">&#10003;</span>
            <span>{{ step }}</span>
          </div>
          <div v-if="idx < steps.length - 1" class="w-6 h-0.5 mx-1 bg-gray-300 dark:bg-gray-600" />
        </div>
      </div>

      <div class="border rounded-lg p-4 bg-gray-50 min-h-[120px] dark:bg-gray-700 dark:border-gray-600">
        <p class="text-sm text-gray-600 dark:text-gray-300">
          <template v-if="currentStep === 0">Awaiting submission review. Verify documentation and schedule inspection.</template>
          <template v-else-if="currentStep === 1">Physical inspection in progress. Check quantity, condition, and packaging.</template>
          <template v-else-if="currentStep === 2">Decision required: approve inventory, flag discrepancy, or hold for review.</template>
          <template v-else-if="currentStep === 3">All items approved. Inventory released to available stock.</template>
        </p>
      </div>

      <div class="flex justify-between">
        <button
          :disabled="!canGoBack"
          class="px-4 py-2 border rounded text-sm hover:bg-gray-50 disabled:opacity-40 dark:bg-gray-700 dark:text-gray-300 dark:border-gray-600 dark:hover:bg-gray-600"
          @click="emit('stepChange', currentStep - 1)"
        >
          Previous
        </button>
        <button
          :disabled="!canAdvance"
          class="px-4 py-2 bg-blue-600 text-white rounded text-sm hover:bg-blue-700 disabled:opacity-60"
          @click="emit('stepChange', currentStep + 1)"
        >
          {{ currentStep === steps.length - 1 ? 'Completed' : 'Next Step' }}
        </button>
      </div>
    </template>

    <div v-else class="text-center text-gray-400 dark:text-gray-400 py-8">Select a shipment to view details</div>
  </div>
</template>
