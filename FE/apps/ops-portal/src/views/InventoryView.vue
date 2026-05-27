<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { StockTable, LowStockAlert, StockAdjustDialog } from '@wms/ui'
import { useAuthStore } from '../stores/auth'
import { useInventoryStore } from '../stores/inventory'

interface StockTableEntry {
  id: string; product_sku: string; product_name: string; bin_location: string
  quantity: number; reserved_quantity: number; status: string
  lot_number: string; expiry_date: string | null; low_stock_threshold: number
}

interface StockAdjustData {
  delta: number; reason: string; notes: string
}

const auth = useAuthStore()
const inv = useInventoryStore()

const search = ref('')
const statusFilter = ref('')
const showAdjustDialog = ref(false)
const selectedEntry = ref<StockTableEntry | null>(null)
const adjustLoading = ref(false)

onMounted(async () => {
  if (auth.token) {
    inv.setToken(auth.token)
    await Promise.all([
      inv.fetchStock(),
      inv.fetchLowStockAlerts(),
    ])
  }
})

watch(() => auth.token, async (tok) => {
  if (tok) {
    inv.setToken(tok)
    await Promise.all([
      inv.fetchStock(),
      inv.fetchLowStockAlerts(),
    ])
  }
})

const stockTableEntries = ref<StockTableEntry[]>([])

watch(() => inv.stockEntries, (entries) => {
  stockTableEntries.value = entries.map(e => ({
    id: e.id,
    product_sku: '', // would need product lookup
    product_name: '',
    bin_location: e.bin_location_id,
    quantity: e.quantity,
    reserved_quantity: e.reserved_quantity,
    status: e.status,
    lot_number: e.lot_number,
    expiry_date: e.expiry_date,
    low_stock_threshold: 10,
  }))
}, { immediate: true })

const alertItems = ref<Array<{ product_id: string; product_sku: string; product_name: string; current_quantity: number; threshold: number; status: 'critical' | 'warning' }>>([])

watch(() => inv.lowStockAlerts, (alerts) => {
  alertItems.value = alerts.map(a => ({
    product_id: a.product.id,
    product_sku: a.product.sku,
    product_name: a.product.name,
    current_quantity: a.current_quantity,
    threshold: a.threshold,
    status: a.status,
  }))
}, { immediate: true })

function handleAdjust(entry: StockTableEntry): void {
  selectedEntry.value = entry
  showAdjustDialog.value = true
}

function handleReserve(entry: StockTableEntry): void {
  inv.reserveStock(entry.id, { quantity: entry.quantity })
}

function handleRelease(entry: StockTableEntry): void {
  inv.releaseStock(entry.id)
}

async function handleAdjustConfirm(data: StockAdjustData): Promise<void> {
  if (!selectedEntry.value) return
  adjustLoading.value = true
  await inv.adjustStock(selectedEntry.value.id, {
    delta: data.delta,
    reason: data.reason,
    notes: data.notes,
  })
  adjustLoading.value = false
  showAdjustDialog.value = false
  selectedEntry.value = null
}

function handleAdjustCancel(): void {
  showAdjustDialog.value = false
  selectedEntry.value = null
}

function handleViewProduct(productId: string): void {
  // Navigate to products view — handled by router
}
</script>

<template>
  <div class="py-8 px-4 space-y-6">
    <h1 class="text-2xl font-semibold text-gray-800 dark:text-gray-100">Inventory</h1>

    <!-- Filters -->
    <div class="flex gap-4 items-center">
      <input
        v-model="search"
        placeholder="Search by SKU or name..."
        class="px-3 py-1.5 border rounded text-sm w-64 dark:bg-gray-800 dark:text-gray-100 dark:border-gray-600"
      />
      <select
        v-model="statusFilter"
        class="px-3 py-1.5 border rounded text-sm dark:bg-gray-800 dark:text-gray-100 dark:border-gray-600"
      >
        <option value="">All Status</option>
        <option value="available">Available</option>
        <option value="reserved">Reserved</option>
        <option value="inactive">Inactive</option>
      </select>
    </div>

    <!-- Low Stock Alerts -->
    <div class="bg-white border rounded-lg p-4 dark:bg-gray-800 dark:border-gray-700">
      <h3 class="text-sm font-semibold text-gray-700 mb-3 dark:text-gray-300">Low Stock Alerts</h3>
      <LowStockAlert
        :alerts="alertItems"
        @view="handleViewProduct"
      />
    </div>

    <!-- Stock Table -->
    <StockTable
      :stock-entries="stockTableEntries"
      :loading="inv.stockLoading"
      @adjust="handleAdjust"
      @reserve="handleReserve"
      @release="handleRelease"
    />

    <!-- Adjust Dialog -->
    <StockAdjustDialog
      v-if="showAdjustDialog && selectedEntry"
      :product-sku="selectedEntry.product_sku"
      :product-name="selectedEntry.product_name"
      :current-quantity="selectedEntry.quantity"
      :loading="adjustLoading"
      @confirm="handleAdjustConfirm"
      @cancel="handleAdjustCancel"
    />
  </div>
</template>
