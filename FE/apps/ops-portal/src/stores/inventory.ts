import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { InventoryClient } from '@wms/api-client'
import type {
  Product,
  StockEntry,
  BinLocation,
  LowStockAlert,
} from '@wms/api-client'

const inventoryClient = new InventoryClient(import.meta.env.VITE_API_BASE_URL || '')

export const useInventoryStore = defineStore('inventory', () => {
  const token = ref<string | null>(localStorage.getItem('wms_ops_token'))

  // ── Products ───────────────────────────────────────
  const products = ref<Product[]>([])
  const productsLoading = ref(false)
  const productsError = ref<string | null>(null)

  // ── Stock ──────────────────────────────────────────
  const stockEntries = ref<StockEntry[]>([])
  const stockLoading = ref(false)
  const stockError = ref<string | null>(null)

  // ── Bin Locations ──────────────────────────────────
  const binLocations = ref<BinLocation[]>([])
  const binLocationsLoading = ref(false)
  const binLocationsError = ref<string | null>(null)

  // ── Low Stock Alerts ───────────────────────────────
  const lowStockAlerts = ref<LowStockAlert[]>([])
  const alertsLoading = ref(false)
  const alertsError = ref<string | null>(null)

  const loading = computed(() =>
    productsLoading.value || stockLoading.value || binLocationsLoading.value || alertsLoading.value
  )

  const error = computed(() =>
    productsError.value || stockError.value || binLocationsError.value || alertsError.value
  )

  function setToken(t: string | null): void {
    token.value = t
    inventoryClient.setToken(t)
  }

  // ── Product actions ────────────────────────────────
  async function fetchProducts(params?: { page?: number; page_size?: number; search?: string; category?: string }): Promise<void> {
    if (!token.value) return
    productsLoading.value = true
    productsError.value = null
    try {
      inventoryClient.setToken(token.value)
      const result = await inventoryClient.listProducts(params)
      products.value = result.products
    } catch (e) {
      productsError.value = (e as Error).message
      products.value = []
    } finally {
      productsLoading.value = false
    }
  }

  async function createProduct(data: import('@wms/api-client').CreateProductRequest): Promise<Product | null> {
    if (!token.value) return null
    try {
      inventoryClient.setToken(token.value)
      const product = await inventoryClient.createProduct(data)
      await fetchProducts()
      return product
    } catch (e) {
      productsError.value = (e as Error).message
      return null
    }
  }

  async function updateProduct(id: string, data: import('@wms/api-client').UpdateProductRequest): Promise<Product | null> {
    if (!token.value) return null
    try {
      inventoryClient.setToken(token.value)
      const product = await inventoryClient.updateProduct(id, data)
      await fetchProducts()
      return product
    } catch (e) {
      productsError.value = (e as Error).message
      return null
    }
  }

  async function archiveProduct(id: string): Promise<boolean> {
    if (!token.value) return false
    try {
      inventoryClient.setToken(token.value)
      await inventoryClient.archiveProduct(id)
      await fetchProducts()
      return true
    } catch (e) {
      productsError.value = (e as Error).message
      return false
    }
  }

  // ── Stock actions ──────────────────────────────────
  async function fetchStock(params?: { page_size?: number; page_token?: string; product_id?: string; bin_location_id?: string; status?: string }): Promise<void> {
    if (!token.value) return
    stockLoading.value = true
    stockError.value = null
    try {
      inventoryClient.setToken(token.value)
      const result = await inventoryClient.listStock(params)
      stockEntries.value = result.entries
    } catch (e) {
      stockError.value = (e as Error).message
      stockEntries.value = []
    } finally {
      stockLoading.value = false
    }
  }

  async function addStock(data: import('@wms/api-client').AddStockRequest): Promise<StockEntry | null> {
    if (!token.value) return null
    try {
      inventoryClient.setToken(token.value)
      const entry = await inventoryClient.addStock(data)
      await fetchStock()
      return entry
    } catch (e) {
      stockError.value = (e as Error).message
      return null
    }
  }

  async function adjustStock(id: string, data: import('@wms/api-client').AdjustStockRequest): Promise<StockEntry | null> {
    if (!token.value) return null
    try {
      inventoryClient.setToken(token.value)
      const entry = await inventoryClient.adjustStock(id, data)
      await fetchStock()
      return entry
    } catch (e) {
      stockError.value = (e as Error).message
      return null
    }
  }

  async function reserveStock(id: string, data: import('@wms/api-client').ReserveStockRequest): Promise<StockEntry | null> {
    if (!token.value) return null
    try {
      inventoryClient.setToken(token.value)
      const entry = await inventoryClient.reserveStock(id, data)
      await fetchStock()
      return entry
    } catch (e) {
      stockError.value = (e as Error).message
      return null
    }
  }

  async function releaseStock(id: string, quantity?: number): Promise<StockEntry | null> {
    if (!token.value) return null
    try {
      inventoryClient.setToken(token.value)
      const entry = await inventoryClient.releaseStock(id, quantity)
      await fetchStock()
      return entry
    } catch (e) {
      stockError.value = (e as Error).message
      return null
    }
  }

  // ── Bin Location actions ───────────────────────────
  async function fetchBinLocations(params?: { page_size?: number; page_token?: string; warehouse_zone?: string }): Promise<void> {
    if (!token.value) return
    binLocationsLoading.value = true
    binLocationsError.value = null
    try {
      inventoryClient.setToken(token.value)
      const result = await inventoryClient.listBinLocations(params)
      binLocations.value = result.locations
    } catch (e) {
      binLocationsError.value = (e as Error).message
      binLocations.value = []
    } finally {
      binLocationsLoading.value = false
    }
  }

  async function createBinLocation(data: import('@wms/api-client').CreateBinLocationRequest): Promise<BinLocation | null> {
    if (!token.value) return null
    try {
      inventoryClient.setToken(token.value)
      const location = await inventoryClient.createBinLocation(data)
      await fetchBinLocations()
      return location
    } catch (e) {
      binLocationsError.value = (e as Error).message
      return null
    }
  }

  // ── Low Stock Alerts ───────────────────────────────
  async function fetchLowStockAlerts(): Promise<void> {
    if (!token.value) return
    alertsLoading.value = true
    alertsError.value = null
    try {
      inventoryClient.setToken(token.value)
      const result = await inventoryClient.getLowStockAlerts()
      lowStockAlerts.value = result.alerts
    } catch (e) {
      alertsError.value = (e as Error).message
      lowStockAlerts.value = []
    } finally {
      alertsLoading.value = false
    }
  }

  return {
    // state
    products, productsLoading, productsError,
    stockEntries, stockLoading, stockError,
    binLocations, binLocationsLoading, binLocationsError,
    lowStockAlerts, alertsLoading, alertsError,
    loading, error,
    // actions
    setToken,
    fetchProducts, createProduct, updateProduct, archiveProduct,
    fetchStock, addStock, adjustStock, reserveStock, releaseStock,
    fetchBinLocations, createBinLocation,
    fetchLowStockAlerts,
  }
})
