<script setup lang="ts">
import { computed, watch } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useDashboardStore } from '../stores/dashboard'

const auth = useAuthStore()
const dashboard = useDashboardStore()

const isEmpty = computed(
  () =>
    !dashboard.loading &&
    dashboard.activeInbounds === 0 &&
    dashboard.totalInbounds === 0 &&
    dashboard.stockEntries === 0
)

watch(
  () => auth.token,
  (newToken) => {
    dashboard.setToken(newToken)
    if (auth.user) {
      dashboard.fetchStats(auth.user.id)
    }
  },
  { immediate: true }
)
</script>

<template>
  <div class="p-6">
    <div class="flex items-center justify-between mb-8">
      <div>
        <h1 class="text-2xl font-semibold text-gray-800 dark:text-gray-100">Welcome, {{ auth.user?.name }}</h1>
        <p class="text-gray-500 dark:text-gray-400">Customer Dashboard</p>
      </div>
      <div v-if="dashboard.loading" class="text-sm text-gray-500 dark:text-gray-400">Loading...</div>
    </div>
    <div v-if="dashboard.error" class="mb-4 text-red-500 text-sm">Error: {{ dashboard.error }}</div>
    <div class="grid grid-cols-3 gap-6">
      <div class="bg-white rounded-lg shadow p-5 dark:bg-gray-800 dark:text-gray-100">
        <h3 class="text-sm font-medium text-gray-500 dark:text-gray-400">Active Inbounds</h3>
        <p v-if="isEmpty" class="text-sm text-gray-400 mt-1 italic">No data yet</p>
        <p v-else class="text-2xl font-bold mt-1">{{ dashboard.activeInbounds }}</p>
      </div>
      <div class="bg-white rounded-lg shadow p-5 dark:bg-gray-800 dark:text-gray-100">
        <h3 class="text-sm font-medium text-gray-500 dark:text-gray-400">Total Inbounds</h3>
        <p v-if="isEmpty" class="text-sm text-gray-400 mt-1 italic">No data yet</p>
        <p v-else class="text-2xl font-bold mt-1">{{ dashboard.totalInbounds }}</p>
      </div>
      <div class="bg-white rounded-lg shadow p-5 dark:bg-gray-800 dark:text-gray-100">
        <h3 class="text-sm font-medium text-gray-500 dark:text-gray-400">Stock Entries</h3>
        <p v-if="isEmpty" class="text-sm text-gray-400 mt-1 italic">No data yet</p>
        <p v-else class="text-2xl font-bold mt-1">{{ dashboard.stockEntries }}</p>
      </div>
    </div>
  </div>
</template>
