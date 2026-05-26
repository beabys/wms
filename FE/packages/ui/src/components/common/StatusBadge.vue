<script setup lang="ts">
defineOptions({ name: 'StatusBadge' })

const props = defineProps<{
  status: 'pending' | 'approved' | 'flagged' | 'held' | 'shipped' | 'completed'
}>()

const statusColors: Record<string, string> = {
  pending: 'bg-yellow-100 text-yellow-800 border-yellow-300 dark:bg-yellow-900/50 dark:text-yellow-300',
  approved: 'bg-green-100 text-green-800 border-green-300 dark:bg-green-900/50 dark:text-green-300',
  flagged: 'bg-red-100 text-red-800 border-red-300 dark:bg-red-900/50 dark:text-red-300',
  held: 'bg-orange-100 text-orange-800 border-orange-300 dark:bg-red-900/50 dark:text-red-300',
  shipped: 'bg-blue-100 text-blue-800 border-blue-300 dark:bg-blue-900/50 dark:text-blue-300',
  completed: 'bg-gray-100 text-gray-800 border-gray-300 dark:bg-green-900/50 dark:text-green-300'
}
</script>

<template>
  <span
    class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium border"
    :class="statusColors[status] || statusColors.pending"
  >
    {{ status.charAt(0).toUpperCase() + status.slice(1) }}
  </span>
</template>
