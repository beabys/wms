import { defineStore } from 'pinia'
import { ref } from 'vue'
import { InboundClient, InventoryClient } from '@wms/api-client'

const inboundClient = new InboundClient(import.meta.env.VITE_API_BASE_URL || '')
const inventoryClient = new InventoryClient(import.meta.env.VITE_API_BASE_URL || '')

export const useDashboardStore = defineStore('dashboard', () => {
  const loading = ref(false)
  const error = ref<string | null>(null)
  const activeInbounds = ref(0)
  const totalInbounds = ref(0)
  const stockEntries = ref(0)

  function setToken(t: string | null): void {
    inboundClient.setToken(t)
    inventoryClient.setToken(t)
  }

  async function fetchStats(customerId: string): Promise<void> {
    loading.value = true
    error.value = null
    try {
      // Single call — derive total + active from same response
      const inboundRes = await inboundClient.getCustomerInbounds(customerId)
      totalInbounds.value = inboundRes.inbounds.length
      activeInbounds.value = inboundRes.inbounds.filter(
        (item) => item.status === 'submitted'
      ).length

      // Stock entries count (approximate if paginated)
      const stockRes = await inventoryClient.listStock()
      stockEntries.value = stockRes.entries.length
    } catch (e) {
      error.value = e instanceof Error ? e.message : String(e)
    } finally {
      loading.value = false
    }
  }

  return {
    loading,
    error,
    activeInbounds,
    totalInbounds,
    stockEntries,
    setToken,
    fetchStats,
  }
})
