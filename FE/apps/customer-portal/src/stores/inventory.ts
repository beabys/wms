import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { InventoryClient } from '@wms/api-client'
import type { StockEntry } from '@wms/api-client'

const inventoryClient = new InventoryClient(import.meta.env.VITE_API_BASE_URL || '')

export const useInventoryStore = defineStore('inventory', () => {
  const token = ref<string | null>(localStorage.getItem('wms_token'))
  const customerId = ref<string | null>(localStorage.getItem('wms_customer_id'))
  const stockEntries = ref<StockEntry[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  function setToken(t: string | null): void {
    token.value = t
    inventoryClient.setToken(t)
  }

  function setCustomerId(id: string): void {
    customerId.value = id
    localStorage.setItem('wms_customer_id', id)
  }

  async function fetchMyStock(params?: { page_size?: number; page_token?: string }): Promise<void> {
    if (!token.value || !customerId.value) return
    loading.value = true
    error.value = null
    try {
      inventoryClient.setToken(token.value)
      const result = await inventoryClient.getCustomerStock(customerId.value, params)
      stockEntries.value = result.entries
    } catch (e) {
      error.value = (e as Error).message
      stockEntries.value = []
    } finally {
      loading.value = false
    }
  }

  return {
    stockEntries,
    loading,
    error,
    setToken,
    setCustomerId,
    fetchMyStock,
  }
})
