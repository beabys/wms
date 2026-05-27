<script setup lang="ts">
import { useAuthStore } from '../stores/auth'
const auth = useAuthStore()
</script>

<template>
  <div class="p-6">
    <h1 class="text-2xl font-semibold text-gray-800 mb-2 dark:text-gray-100">Operations Dashboard</h1>
    <p class="text-gray-500 mb-8 dark:text-gray-400">Welcome, {{ auth.user?.name }}</p>
    <div class="grid grid-cols-3 gap-6">
      <div class="bg-white rounded-lg shadow p-5 dark:bg-gray-800 dark:text-gray-100">
        <h3 class="text-sm font-medium text-gray-500 dark:text-gray-400">Pending Approvals</h3>
        <p class="text-2xl font-bold mt-1">-</p>
      </div>
      <div class="bg-white rounded-lg shadow p-5 dark:bg-gray-800 dark:text-gray-100">
        <h3 class="text-sm font-medium text-gray-500 dark:text-gray-400">Inbound Queue</h3>
        <p class="text-2xl font-bold mt-1">-</p>
      </div>
      <div class="bg-white rounded-lg shadow p-5 dark:bg-gray-800 dark:text-gray-100">
        <h3 class="text-sm font-medium text-gray-500 dark:text-gray-400">Active Orders</h3>
        <p class="text-2xl font-bold mt-1">-</p>
      </div>
    </div>
  </div>
</template>
