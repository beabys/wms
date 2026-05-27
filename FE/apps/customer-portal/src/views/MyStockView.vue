<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useInventoryStore } from '../stores/inventory'

const auth = useAuthStore()
const inv = useInventoryStore()

onMounted(async () => {
  if (auth.token && auth.user?.id) {
    inv.setToken(auth.token)
    inv.setCustomerId(auth.user.id)
    await inv.fetchMyStock()
  }
})

watch(() => auth.token, async (tok) => {
  if (tok && auth.user?.id) {
    inv.setToken(tok)
    inv.setCustomerId(auth.user.id)
    await inv.fetchMyStock()
  }
})

const columns = [
  { key: 'lot_number', label: 'Lot#' },
  { key: 'quantity', label: 'Quantity' },
  { key: 'reserved_quantity', label: 'Reserved' },
  { key: 'status', label: 'Status' },
  { key: 'expiry_date', label: 'Expiry' },
  { key: 'last_updated', label: 'Last Updated' },
]
</script>

<template>
  <div class="py-8 px-4 space-y-6">
    <h1 class="text-2xl font-semibold text-gray-800 dark:text-gray-100">My Stock</h1>

    <div class="bg-white border rounded-lg overflow-hidden dark:bg-gray-900 dark:border-gray-700">
      <table class="min-w-full divide-y divide-gray-200 dark:text-gray-100">
        <thead class="bg-gray-50 dark:bg-gray-800">
          <tr>
            <th v-for="col in columns" :key="col.key" class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase dark:text-gray-300">
              {{ col.label }}
            </th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
          <tr v-if="inv.loading">
            <td :colspan="columns.length" class="px-4 py-8 text-center text-gray-400">Loading...</td>
          </tr>
          <tr v-else-if="inv.stockEntries.length === 0">
            <td :colspan="columns.length" class="px-4 py-8 text-center text-gray-400">No stock entries</td>
          </tr>
          <tr
            v-for="entry in inv.stockEntries"
            :key="entry.id"
            class="hover:bg-gray-50 dark:border-gray-700 dark:hover:bg-gray-800"
          >
            <td class="px-4 py-3 text-sm font-mono text-gray-500 dark:text-gray-400">{{ entry.lot_number || '-' }}</td>
            <td class="px-4 py-3 text-sm font-semibold text-gray-700 dark:text-gray-300">{{ entry.quantity }}</td>
            <td class="px-4 py-3 text-sm text-gray-500 dark:text-gray-400">{{ entry.reserved_quantity }}</td>
            <td class="px-4 py-3 text-sm">
              <span
                class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium border"
                :class="{
                  'bg-green-100 text-green-800 border-green-300 dark:bg-green-900/50 dark:text-green-300': entry.status === 'available',
                  'bg-yellow-100 text-yellow-800 border-yellow-300 dark:bg-yellow-900/50 dark:text-yellow-300': entry.status === 'reserved',
                  'bg-gray-100 text-gray-800 border-gray-300 dark:bg-gray-800 dark:text-gray-300': entry.status === 'inactive',
                }"
              >
                {{ entry.status }}
              </span>
            </td>
            <td class="px-4 py-3 text-sm text-gray-500 dark:text-gray-400">{{ entry.expiry_date || '-' }}</td>
            <td class="px-4 py-3 text-sm text-gray-500 dark:text-gray-400">{{ entry.last_updated }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
