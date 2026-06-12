import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { CustomerResponse, Pagination } from '@/api/types'
import { customerClient } from '@/api/customerClient'

export const useCustomerStore = defineStore('customers', () => {
  const customers = ref<CustomerResponse[]>([])
  const pagination = ref<Pagination>({ page: 1, page_size: 10, total_items: 0 })
  const loading = ref(false)
  const error = ref('')

  async function fetchCustomers(params?: { status?: string; page?: number; page_size?: number }) {
    loading.value = true
    error.value = ''
    try {
      const res = await customerClient.listCustomers(params)
      customers.value = res.customers
      pagination.value = res.pagination
    } catch (err: any) {
      error.value = err?.message || 'Failed to fetch customers'
      customers.value = []
    } finally {
      loading.value = false
    }
  }

  async function approveCustomer(id: string) {
    error.value = ''
    try {
      const updated = await customerClient.approveCustomer(id)
      const idx = customers.value.findIndex((c) => c.id === id)
      if (idx !== -1) {
        customers.value[idx] = { ...customers.value[idx], ...updated }
      }
    } catch (err: any) {
      error.value = err?.message || 'Failed to approve customer'
      throw err
    }
  }

  async function rejectCustomer(id: string, reason?: string) {
    error.value = ''
    try {
      const updated = await customerClient.rejectCustomer(id, reason)
      const idx = customers.value.findIndex((c) => c.id === id)
      if (idx !== -1) {
        customers.value[idx] = { ...customers.value[idx], ...updated }
      }
    } catch (err: any) {
      error.value = err?.message || 'Failed to reject customer'
      throw err
    }
  }

  async function suspendCustomer(id: string) {
    error.value = ''
    try {
      await customerClient.suspendCustomer(id)
      const idx = customers.value.findIndex((c) => c.id === id)
      if (idx !== -1) {
        customers.value[idx] = { ...customers.value[idx], status: 'suspended' }
      }
    } catch (err: any) {
      error.value = err?.message || 'Failed to suspend customer'
      throw err
    }
  }

  async function restoreCustomer(id: string) {
    error.value = ''
    try {
      const updated = await customerClient.restoreCustomer(id)
      const idx = customers.value.findIndex((c) => c.id === id)
      if (idx !== -1) {
        customers.value[idx] = { ...customers.value[idx], ...updated }
      }
    } catch (err: any) {
      error.value = err?.message || 'Failed to restore customer'
      throw err
    }
  }

  function $reset() {
    customers.value = []
    pagination.value = { page: 1, page_size: 10, total_items: 0 }
    loading.value = false
    error.value = ''
  }

  return {
    customers,
    pagination,
    loading,
    error,
    fetchCustomers,
    approveCustomer,
    rejectCustomer,
    suspendCustomer,
    restoreCustomer,
    $reset,
  }
})
