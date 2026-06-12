import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { CustomerResponse, AuditEntry, Pagination } from '@/api/types'
import { customerClient } from '@/api/customerClient'

export const useCustomerDetailStore = defineStore('customerDetail', () => {
  const customer = ref<CustomerResponse | null>(null)
  const auditLogs = ref<AuditEntry[]>([])
  const auditPagination = ref<Pagination>({ page: 1, page_size: 10, total_items: 0 })
  const loading = ref(false)
  const error = ref('')

  async function fetchCustomer(id: string) {
    loading.value = true
    error.value = ''
    try {
      customer.value = await customerClient.getCustomer(id)
    } catch (err: any) {
      error.value = err?.message || 'Failed to fetch customer'
      customer.value = null
    } finally {
      loading.value = false
    }
  }

  async function updateCustomer(id: string, data: Partial<CustomerResponse>) {
    error.value = ''
    try {
      const updated = await customerClient.updateCustomer(id, data)
      customer.value = updated
    } catch (err: any) {
      error.value = err?.message || 'Failed to update customer'
      throw err
    }
  }

  async function approveCustomer(id: string) {
    error.value = ''
    try {
      const updated = await customerClient.approveCustomer(id)
      customer.value = updated
    } catch (err: any) {
      error.value = err?.message || 'Failed to approve customer'
      throw err
    }
  }

  async function rejectCustomer(id: string, reason?: string) {
    error.value = ''
    try {
      const updated = await customerClient.rejectCustomer(id, reason)
      customer.value = updated
    } catch (err: any) {
      error.value = err?.message || 'Failed to reject customer'
      throw err
    }
  }

  async function suspendCustomer(id: string) {
    error.value = ''
    try {
      await customerClient.suspendCustomer(id)
      if (customer.value) {
        customer.value = { ...customer.value, status: 'suspended' }
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
      customer.value = updated
    } catch (err: any) {
      error.value = err?.message || 'Failed to restore customer'
      throw err
    }
  }

  async function fetchAuditLogs(id: string, params?: { page?: number; page_size?: number }) {
    try {
      const res = await customerClient.listAuditLogs(id, params)
      auditLogs.value = res.entries
      auditPagination.value = res.pagination
    } catch (err: any) {
      // audit log fetch failure is non-fatal
      auditLogs.value = []
    }
  }

  function $reset() {
    customer.value = null
    auditLogs.value = []
    auditPagination.value = { page: 1, page_size: 10, total_items: 0 }
    loading.value = false
    error.value = ''
  }

  return {
    customer,
    auditLogs,
    auditPagination,
    loading,
    error,
    fetchCustomer,
    updateCustomer,
    approveCustomer,
    rejectCustomer,
    suspendCustomer,
    restoreCustomer,
    fetchAuditLogs,
    $reset,
  }
})
