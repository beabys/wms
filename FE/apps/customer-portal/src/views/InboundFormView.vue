<script setup lang="ts">
import { ref } from 'vue'
import { useNotification } from '@wms/composables'
import { InboundForm } from '@wms/ui'
import { InboundClient } from '@wms/api-client'
import { useAuthStore } from '../stores/auth'

const inboundClient = new InboundClient(import.meta.env.VITE_API_BASE_URL || '')
const auth = useAuthStore()

const { success: notifySuccess, error: notifyError } = useNotification()
const loading = ref(false)
const products = ref<any[]>([])

function handleAddProduct(): void {
  products.value.push({
    id: 'p' + Date.now(),
    sku: '',
    name: '',
    quantity: 1,
    packaging: 'box',
  })
}

function handleRemoveProduct(id: string): void {
  products.value = products.value.filter(p => p.id !== id)
}

function handleUpdateProduct(id: string, field: string, value: unknown): void {
  const p = products.value.find(p => p.id === id)
  if (p) p[field] = value
}

async function handleSubmit(): Promise<void> {
  if (!auth.token) return
  loading.value = true
  try {
    inboundClient.setToken(auth.token)
    await inboundClient.create({
      customer_id: auth.user?.id || '',
      expected_date: new Date().toISOString().slice(0, 10),
      items: products.value.map(p => ({
        sku: p.sku,
        quantity_declared: p.quantity,
      })),
    })
    notifySuccess('Shipment submitted successfully!')
    products.value = []
  } catch (e: any) {
    notifyError(e.message || 'Failed to submit shipment')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="max-w-3xl mx-auto py-8 px-4">
    <InboundForm
      :products="products"
      :loading="loading"
      @add-product="handleAddProduct"
      @remove-product="handleRemoveProduct"
      @update-product="handleUpdateProduct"
      @submit="handleSubmit"
    />
  </div>
</template>
