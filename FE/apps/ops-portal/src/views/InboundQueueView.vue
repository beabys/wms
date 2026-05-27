<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { InboundQueue, InboundDetail } from '@wms/ui'
import { InboundClient } from '@wms/api-client'
import { useAuthStore } from '../stores/auth'
import type { InboundItem } from '@wms/api-client/types/responses'

const auth = useAuthStore()
const inboundClient = new InboundClient(import.meta.env.VITE_API_BASE_URL || '')

const selectedId = ref<string | null>(null)
const loading = ref(false)
const shipments = ref<InboundItem[]>([])
const selectedShipment = ref<InboundItem | null>(null)
const currentStep = ref(0)

onMounted(async () => {
  await fetchQueue()
})

async function fetchQueue(): Promise<void> {
  if (!auth.token) return
  loading.value = true
  try {
    inboundClient.setToken(auth.token)
    const result = await inboundClient.getQueue()
    shipments.value = result.inbounds || []
  } catch {
    shipments.value = []
  } finally {
    loading.value = false
  }
}

async function handleSelect(id: string): Promise<void> {
  selectedId.value = id
  currentStep.value = 0
  if (!auth.token) return
  try {
    inboundClient.setToken(auth.token)
    const detail = await inboundClient.getById(id)
    selectedShipment.value = detail
  } catch {
    selectedShipment.value = null
  }
}

function handleStepChange(step: number): void {
  currentStep.value = step
}
</script>

<template>
  <div class="py-8 px-4 space-y-6">
    <InboundQueue
      :shipments="(shipments || []).map(s => ({
        id: s.id,
        customerName: s.customer_name,
        productCount: (s.items || []).length,
        status: s.status,
        expectedDate: s.expected_date,
        receivedDate: undefined
      }))"
      :loading="loading"
      @select="handleSelect"
    />
    <InboundDetail
      v-if="selectedShipment"
      :shipment="{
        id: selectedShipment.id,
        customerName: selectedShipment.customer_name,
        skuCount: (selectedShipment.items || []).length,
        totalUnits: (selectedShipment.items || []).reduce((sum, i) => sum + i.quantity_declared, 0),
        status: selectedShipment.status,
        expectedDate: selectedShipment.expected_date
      }"
      :current-step="currentStep"
      @step-change="handleStepChange"
    />
  </div>
</template>
