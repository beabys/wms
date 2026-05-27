<script setup lang="ts">
import { computed } from 'vue'
import DataTable from '../common/DataTable.vue'
import type { ColumnDef } from '../common/DataTable.vue'
defineOptions({ name: 'StockTable' })

export interface StockTableEntry {
  id: string
  product_sku: string
  product_name: string
  bin_location: string
  quantity: number
  reserved_quantity: number
  status: string
  lot_number: string
  expiry_date: string | null
  low_stock_threshold: number
}

const props = defineProps<{
  stockEntries: StockTableEntry[]
  loading?: boolean
}>()

const emit = defineEmits<{
  adjust: [entry: StockTableEntry]
  reserve: [entry: StockTableEntry]
  release: [entry: StockTableEntry]
}>()

const columns: ColumnDef[] = [
  { key: 'product_sku', label: 'Product SKU', sortable: true },
  { key: 'product_name', label: 'Product Name' },
  { key: 'bin_location', label: 'Bin Location' },
  { key: 'quantity', label: 'Quantity', sortable: true },
  { key: 'reserved_quantity', label: 'Reserved' },
  { key: 'status', label: 'Status' },
  { key: 'lot_number', label: 'Lot#' },
  { key: 'expiry_date', label: 'Expiry' },
]

function rowClass(entry: StockTableEntry): Record<string, boolean> {
  const qty = entry.quantity
  const threshold = entry.low_stock_threshold
  return {
    'bg-red-50 dark:bg-red-900/20': qty <= threshold,
    'bg-yellow-50 dark:bg-yellow-900/20': qty > threshold && qty < threshold * 2,
  }
}

const tableData = computed(() =>
  props.stockEntries.map(e => ({
    ...e,
    _rowClass: rowClass(e),
  }))
)
</script>

<template>
  <div class="bg-white border rounded-lg overflow-hidden dark:bg-gray-900 dark:border-gray-700">
    <div class="px-4 py-3 border-b bg-gray-50 dark:bg-gray-800 dark:border-gray-700 flex items-center justify-between">
      <h3 class="font-semibold text-gray-800 dark:text-gray-100">Stock Levels</h3>
    </div>

    <table class="min-w-full divide-y divide-gray-200 dark:text-gray-100">
      <thead class="bg-gray-50 dark:bg-gray-800">
        <tr>
          <th v-for="col in columns" :key="col.key" class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase dark:text-gray-300">
            {{ col.label }}
          </th>
          <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase dark:text-gray-300">Actions</th>
        </tr>
      </thead>
      <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
        <tr v-if="loading">
          <td :colspan="columns.length + 1" class="px-4 py-8 text-center text-gray-400">Loading...</td>
        </tr>
        <tr v-else-if="stockEntries.length === 0">
          <td :colspan="columns.length + 1" class="px-4 py-8 text-center text-gray-400">No stock entries</td>
        </tr>
        <tr
          v-for="entry in stockEntries"
          :key="entry.id"
          class="hover:bg-gray-50 dark:border-gray-700 dark:hover:bg-gray-800"
          :class="rowClass(entry)"
        >
          <td class="px-4 py-3 text-sm font-mono text-blue-600 dark:text-blue-400">{{ entry.product_sku }}</td>
          <td class="px-4 py-3 text-sm text-gray-700 dark:text-gray-300">{{ entry.product_name }}</td>
          <td class="px-4 py-3 text-sm text-gray-500 dark:text-gray-400">{{ entry.bin_location }}</td>
          <td class="px-4 py-3 text-sm font-semibold" :class="entry.quantity <= entry.low_stock_threshold ? 'text-red-600' : 'text-gray-700 dark:text-gray-300'">
            {{ entry.quantity }}
          </td>
          <td class="px-4 py-3 text-sm text-gray-500 dark:text-gray-400">{{ entry.reserved_quantity }}</td>
          <td class="px-4 py-3 text-sm">
            <span class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium border"
              :class="{
                'bg-green-100 text-green-800 border-green-300 dark:bg-green-900/50 dark:text-green-300': entry.status === 'available',
                'bg-yellow-100 text-yellow-800 border-yellow-300 dark:bg-yellow-900/50 dark:text-yellow-300': entry.status === 'reserved',
                'bg-gray-100 text-gray-800 border-gray-300 dark:bg-gray-800 dark:text-gray-300': entry.status === 'inactive',
              }"
            >
              {{ entry.status }}
            </span>
          </td>
          <td class="px-4 py-3 text-sm text-gray-500 dark:text-gray-400 font-mono">{{ entry.lot_number || '-' }}</td>
          <td class="px-4 py-3 text-sm text-gray-500 dark:text-gray-400">{{ entry.expiry_date || '-' }}</td>
          <td class="px-4 py-3 text-sm whitespace-nowrap">
            <button
              class="text-xs text-blue-600 hover:text-blue-800 mr-2 dark:text-blue-400 dark:hover:text-blue-300"
              @click="emit('adjust', entry)"
            >Adjust</button>
            <button
              class="text-xs text-green-600 hover:text-green-800 mr-2 dark:text-green-400 dark:hover:text-green-300"
              @click="emit('reserve', entry)"
            >Reserve</button>
            <button
              class="text-xs text-orange-600 hover:text-orange-800 dark:text-orange-400 dark:hover:text-orange-300"
              @click="emit('release', entry)"
            >Release</button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
