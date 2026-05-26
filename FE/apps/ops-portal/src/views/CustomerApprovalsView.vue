<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { CustomerApprovalQueue } from '@wms/ui'
import { CustomerClient } from '@wms/api-client'
import { useAuthStore } from '../stores/auth'
import type { CustomerItem } from '@wms/api-client/types/responses'

const auth = useAuthStore()
const customerClient = new CustomerClient(import.meta.env.VITE_API_BASE_URL || '')

const customers = ref<CustomerItem[]>([])
const loading = ref(false)

onMounted(async () => {
  await fetchQueue()
})

async function fetchQueue(): Promise<void> {
  if (!auth.token) return
  loading.value = true
  try {
    customerClient.setToken(auth.token)
    const result = await customerClient.getQueue()
    customers.value = result.customers
  } catch {
    customers.value = []
  } finally {
    loading.value = false
  }
}

async function handleApprove(id: string): Promise<void> {
  if (!auth.token) return
  try {
    customerClient.setToken(auth.token)
    await customerClient.approve(id)
    customers.value = customers.value.filter(c => c.id !== id)
  } catch {
    // error handled silently
  }
}

async function handleReject(id: string): Promise<void> {
  if (!auth.token) return
  try {
    customerClient.setToken(auth.token)
    await customerClient.suspend(id, 'Rejected by admin')
    customers.value = customers.value.filter(c => c.id !== id)
  } catch {
    // error handled silently
  }
}
</script>

<template>
  <div class="py-8 px-4 max-w-5xl mx-auto">
    <h1 class="text-2xl font-semibold text-gray-800 mb-6">Customer Approvals</h1>
    <CustomerApprovalQueue
      :customers="customers.map(c => ({
        id: c.id,
        name: c.company_name,
        email: c.email,
        company: c.company_name,
        registeredAt: c.created_at,
        status: c.status as 'pending' | 'approved' | 'rejected'
      }))"
      :loading="loading"
      @approve="handleApprove"
      @reject="handleReject"
    />
  </div>
</template>
