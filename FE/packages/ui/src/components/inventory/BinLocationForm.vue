<script setup lang="ts">
import { ref } from 'vue'
defineOptions({ name: 'BinLocationForm' })

export interface BinLocationFormData {
  warehouse_zone: string
  aisle: string
  rack: string
  shelf: string
}

const props = defineProps<{
  loading?: boolean
  errors?: Record<string, string>
}>()

const emit = defineEmits<{
  submit: [data: BinLocationFormData]
  cancel: []
}>()

const form = ref<BinLocationFormData>({
  warehouse_zone: '',
  aisle: '',
  rack: '',
  shelf: '',
})

const zones = ['A', 'B', 'C', 'D']

function handleSubmit(): void {
  if (!form.value.warehouse_zone || !form.value.aisle.trim() || !form.value.rack.trim() || !form.value.shelf.trim()) return
  emit('submit', { ...form.value })
}
</script>

<template>
  <div class="bg-white border rounded-lg p-6 space-y-6 dark:bg-gray-800 dark:border-gray-700">
    <h2 class="text-xl font-semibold text-gray-800 dark:text-gray-100">Create Bin Location</h2>

    <div v-if="errors && errors._general" class="text-sm text-red-600 bg-red-50 border border-red-200 rounded px-3 py-2 dark:text-red-400 dark:bg-red-900/30 dark:border-red-800">
      {{ errors._general }}
    </div>

    <div class="grid grid-cols-2 gap-4">
      <div>
        <label class="block text-xs text-gray-500 mb-1 dark:text-gray-400">Warehouse Zone *</label>
        <select
          v-model="form.warehouse_zone"
          class="w-full px-2 py-1.5 border rounded text-sm dark:bg-gray-700 dark:text-gray-100 dark:border-gray-600"
        >
          <option value="">-- Select --</option>
          <option v-for="z in zones" :key="z" :value="z">Zone {{ z }}</option>
        </select>
      </div>
      <div>
        <label class="block text-xs text-gray-500 mb-1 dark:text-gray-400">Aisle *</label>
        <input
          v-model="form.aisle"
          placeholder="e.g. A1"
          class="w-full px-2 py-1.5 border rounded text-sm dark:bg-gray-700 dark:text-gray-100 dark:border-gray-600"
        />
      </div>
      <div>
        <label class="block text-xs text-gray-500 mb-1 dark:text-gray-400">Rack *</label>
        <input
          v-model="form.rack"
          placeholder="e.g. R01"
          class="w-full px-2 py-1.5 border rounded text-sm dark:bg-gray-700 dark:text-gray-100 dark:border-gray-600"
        />
      </div>
      <div>
        <label class="block text-xs text-gray-500 mb-1 dark:text-gray-400">Shelf *</label>
        <input
          v-model="form.shelf"
          placeholder="e.g. S1"
          class="w-full px-2 py-1.5 border rounded text-sm dark:bg-gray-700 dark:text-gray-100 dark:border-gray-600"
        />
      </div>
    </div>

    <div class="flex justify-end gap-3 pt-2">
      <button
        class="px-4 py-2 border border-gray-300 rounded text-sm text-gray-600 hover:bg-gray-50 dark:border-gray-600 dark:text-gray-300 dark:hover:bg-gray-700"
        @click="emit('cancel')"
      >
        Cancel
      </button>
      <button
        :disabled="loading || !form.warehouse_zone || !form.aisle.trim() || !form.rack.trim() || !form.shelf.trim()"
        class="px-6 py-2 bg-blue-600 text-white rounded text-sm font-medium hover:bg-blue-700 disabled:opacity-60"
        @click="handleSubmit"
      >
        {{ loading ? 'Creating...' : 'Create Location' }}
      </button>
    </div>
  </div>
</template>
