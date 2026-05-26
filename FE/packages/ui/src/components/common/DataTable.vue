<script setup lang="ts">
defineOptions({ name: 'WmsDataTable' })

export interface ColumnDef {
  key: string
  label: string
  sortable?: boolean
  width?: string
}

const props = defineProps<{
  columns: ColumnDef[]
  data: Record<string, unknown>[]
  sortKey?: string
  sortDir?: 'asc' | 'desc'
  loading?: boolean
  emptyText?: string
}>()

const emit = defineEmits<{
  sort: [key: string]
}>()

function handleSort(key: string): void {
  const col = props.columns.find(c => c.key === key)
  if (col?.sortable) {
    emit('sort', key)
  }
}
</script>

<template>
  <div class="w-full overflow-x-auto dark:bg-gray-900">
    <table class="min-w-full divide-y divide-gray-200 dark:text-gray-100">
      <thead class="bg-gray-50 dark:bg-gray-800">
        <tr>
          <th
            v-for="col in columns"
            :key="col.key"
            :style="col.width ? { width: col.width } : undefined"
            class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer select-none dark:text-gray-200"
            :class="{ 'hover:bg-gray-100 dark:hover:bg-gray-700': col.sortable }"
            @click="col.sortable && handleSort(col.key)"
          >
            <span class="inline-flex items-center gap-1">
              {{ col.label }}
              <span v-if="col.sortable && sortKey === col.key" class="text-gray-700 dark:text-gray-300">
                {{ sortDir === 'asc' ? '▲' : '▼' }}
              </span>
            </span>
          </th>
        </tr>
      </thead>
      <tbody class="bg-white dark:bg-gray-900 divide-y divide-gray-200 dark:divide-gray-700">
        <tr v-if="loading">
          <td :colspan="columns.length" class="px-4 py-8 text-center text-gray-400 dark:text-gray-400">
            Loading...
          </td>
        </tr>
        <tr v-else-if="data.length === 0">
          <td :colspan="columns.length" class="px-4 py-8 text-center text-gray-400 dark:text-gray-400">
            {{ emptyText || 'No data' }}
          </td>
        </tr>
        <tr v-for="(row, idx) in data" :key="idx" class="hover:bg-gray-50 dark:border-gray-700 dark:hover:bg-gray-800">
          <td v-for="col in columns" :key="col.key" class="px-4 py-3 text-sm text-gray-700 dark:text-gray-300">
            {{ row[col.key] }}
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
