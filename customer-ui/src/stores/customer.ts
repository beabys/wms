import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { customerClient } from '@/api/customerClient'
import type { CustomerResponse } from '@/api/types'

export const useCustomerStore = defineStore('customer', () => {
  const customer = ref<CustomerResponse | null>(null)
  const loading = ref(false)
  const error = ref('')

  const isApproved = computed(() => customer.value?.status === 'active')
  const approvalStatus = computed(() => customer.value?.status || 'unknown')

  async function fetchMyCustomer() {
    loading.value = true
    error.value = ''
    try {
      customer.value = await customerClient.getMyCustomer()
    } catch (err: any) {
      error.value = err?.message || 'Failed to fetch customer'
    } finally {
      loading.value = false
    }
  }

  function $reset() {
    customer.value = null
    loading.value = false
    error.value = ''
  }

  return { customer, loading, error, isApproved, approvalStatus, fetchMyCustomer, $reset }
})
