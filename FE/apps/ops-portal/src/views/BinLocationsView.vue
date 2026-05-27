<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { BinLocationForm } from '@wms/ui'
import { useAuthStore } from '../stores/auth'
import { useInventoryStore } from '../stores/inventory'
import type { CreateBinLocationRequest } from '@wms/api-client'

interface BinLocationFormData {
  warehouse_zone: string; aisle: string; rack: string; shelf: string
}

const auth = useAuthStore()
const inv = useInventoryStore()

const showForm = ref(false)
const formLoading = ref(false)
const formErrors = ref<Record<string, string>>({})

onMounted(async () => {
  if (auth.token) {
    inv.setToken(auth.token)
    await inv.fetchBinLocations()
  }
})

watch(() => auth.token, async (tok) => {
  if (tok) {
    inv.setToken(tok)
    await inv.fetchBinLocations()
  }
})

async function handleSubmit(data: BinLocationFormData): Promise<void> {
  formLoading.value = true
  formErrors.value = {}

  const createData: CreateBinLocationRequest = {
    warehouse_zone: data.warehouse_zone,
    aisle: data.aisle,
    rack: data.rack,
    shelf: data.shelf,
  }

  const result = await inv.createBinLocation(createData)
  if (result) {
    showForm.value = false
  } else {
    formErrors.value = { _general: 'Failed to create bin location' }
  }

  formLoading.value = false
}

function handleCancel(): void {
  showForm.value = false
}
</script>

<template>
  <div class="py-8 px-4 space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="text-2xl font-semibold text-gray-800 dark:text-gray-100">Bin Locations</h1>
      <button
        v-if="!showForm"
        class="px-4 py-2 bg-blue-600 text-white rounded text-sm font-medium hover:bg-blue-700"
        @click="showForm = true"
      >
        + New Location
      </button>
    </div>

    <!-- Create Form -->
    <BinLocationForm
      v-if="showForm"
      :loading="formLoading"
      :errors="formErrors"
      @submit="handleSubmit"
      @cancel="handleCancel"
    />

    <!-- Location List -->
    <div v-else class="bg-white border rounded-lg overflow-hidden dark:bg-gray-900 dark:border-gray-700">
      <table class="min-w-full divide-y divide-gray-200 dark:text-gray-100">
        <thead class="bg-gray-50 dark:bg-gray-800">
          <tr>
            <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase dark:text-gray-300">Zone</th>
            <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase dark:text-gray-300">Aisle</th>
            <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase dark:text-gray-300">Rack</th>
            <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase dark:text-gray-300">Shelf</th>
            <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase dark:text-gray-300">Status</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
          <tr v-if="inv.binLocationsLoading">
            <td colspan="5" class="px-4 py-8 text-center text-gray-400">Loading...</td>
          </tr>
          <tr v-else-if="inv.binLocations.length === 0">
            <td colspan="5" class="px-4 py-8 text-center text-gray-400">No bin locations</td>
          </tr>
          <tr
            v-for="loc in inv.binLocations"
            :key="loc.id"
            class="hover:bg-gray-50 dark:border-gray-700 dark:hover:bg-gray-800"
          >
            <td class="px-4 py-3 text-sm font-semibold text-gray-700 dark:text-gray-300">Zone {{ loc.warehouse_zone }}</td>
            <td class="px-4 py-3 text-sm text-gray-500 dark:text-gray-400">{{ loc.aisle }}</td>
            <td class="px-4 py-3 text-sm text-gray-500 dark:text-gray-400">{{ loc.rack }}</td>
            <td class="px-4 py-3 text-sm text-gray-500 dark:text-gray-400">{{ loc.shelf }}</td>
            <td class="px-4 py-3 text-sm">
              <span
                class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium border"
                :class="loc.is_active
                  ? 'bg-green-100 text-green-800 border-green-300 dark:bg-green-900/50 dark:text-green-300'
                  : 'bg-gray-100 text-gray-800 border-gray-300 dark:bg-gray-800 dark:text-gray-300'"
              >
                {{ loc.is_active ? 'Active' : 'Inactive' }}
              </span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
