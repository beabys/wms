<script setup lang="ts">
defineOptions({ name: 'InboundQueue' })

export interface InboundShipment {
  id: string
  customerName: string
  productCount: number
  status: string
  expectedDate: string
  receivedDate?: string
}

import StatusBadge from '../common/StatusBadge.vue'

const props = defineProps<{
  shipments: InboundShipment[]
  loading?: boolean
}>()

const emit = defineEmits<{
  select: [id: string]
}>()

const columns = [
  { key: 'id', label: 'ID', width: '120px' },
  { key: 'customerName', label: 'Customer' },
  { key: 'productCount', label: 'Products', width: '100px' },
  { key: 'status', label: 'Status' },
  { key: 'expectedDate', label: 'Expected' },
  { key: 'receivedDate', label: 'Received' }
]
</script>

<template>
  <div class="bg-white border rounded-lg overflow-hidden dark:bg-gray-900 dark:border-gray-700">
    <div class="px-4 py-3 border-b bg-gray-50 dark:bg-gray-800 dark:border-gray-700">
      <h3 class="font-semibold text-gray-800 dark:text-gray-100">Inbound Queue</h3>
    </div>

    <table class="min-w-full divide-y divide-gray-200 dark:text-gray-100">
      <thead class="bg-gray-50 dark:bg-gray-800">
        <tr>
          <th v-for="col in columns" :key="col.key" class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase dark:text-gray-300">
            {{ col.label }}
          </th>
        </tr>
      </thead>
      <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
        <tr v-if="loading">
          <td :colspan="columns.length" class="px-4 py-8 text-center text-gray-400 dark:text-gray-400">Loading...</td>
        </tr>
        <tr v-else-if="shipments.length === 0">
          <td :colspan="columns.length" class="px-4 py-8 text-center text-gray-400 dark:text-gray-400">No inbound shipments</td>
        </tr>
        <tr
          v-for="s in shipments"
          :key="s.id"
          class="hover:bg-blue-50 cursor-pointer dark:border-gray-700 dark:hover:bg-gray-800"
          @click="emit('select', s.id)"
        >
          <td class="px-4 py-3 text-sm font-mono text-blue-600 dark:text-blue-400">{{ s.id }}</td>
          <td class="px-4 py-3 text-sm text-gray-700 dark:text-gray-300">{{ s.customerName }}</td>
          <td class="px-4 py-3 text-sm text-gray-500 dark:text-gray-400">{{ s.productCount }}</td>
          <td class="px-4 py-3 text-sm"><StatusBadge :status="s.status as any" /></td>
          <td class="px-4 py-3 text-sm text-gray-500 dark:text-gray-400">{{ s.expectedDate }}</td>
          <td class="px-4 py-3 text-sm text-gray-500 dark:text-gray-400">{{ s.receivedDate || '-' }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
