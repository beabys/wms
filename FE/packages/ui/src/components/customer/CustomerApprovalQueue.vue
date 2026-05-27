<script setup lang="ts">
defineOptions({ name: 'CustomerApprovalQueue' })

export interface CustomerApproval {
  id: string
  name: string
  email: string
  company: string
  registeredAt: string
  status: 'pending' | 'approved' | 'rejected'
}

const props = defineProps<{
  customers: CustomerApproval[]
  loading?: boolean
}>()

const emit = defineEmits<{
  approve: [id: string]
  reject: [id: string]
}>()
</script>

<template>
  <div class="bg-white border rounded-lg overflow-hidden dark:bg-gray-900 dark:border-gray-700">
    <div class="px-4 py-3 border-b bg-gray-50 dark:bg-gray-800 dark:border-gray-700">
      <h3 class="font-semibold text-gray-800 dark:text-gray-100">Pending Approvals</h3>
    </div>

    <table class="min-w-full divide-y divide-gray-200 dark:text-gray-100">
      <thead class="bg-gray-50 dark:bg-gray-800">
        <tr>
          <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase dark:text-gray-300">Name</th>
          <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase dark:text-gray-300">Email</th>
          <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase dark:text-gray-300">Company</th>
          <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase dark:text-gray-300">Registered</th>
          <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase dark:text-gray-300">Actions</th>
        </tr>
      </thead>
      <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
        <tr v-if="loading">
          <td colspan="5" class="px-4 py-8 text-center text-gray-400 dark:text-gray-400">Loading...</td>
        </tr>
        <tr v-else-if="customers.length === 0">
          <td colspan="5" class="px-4 py-8 text-center text-gray-400 dark:text-gray-400">No pending approvals</td>
        </tr>
        <tr v-for="c in customers" :key="c.id" class="hover:bg-gray-50 dark:border-gray-700 dark:hover:bg-gray-800">
          <td class="px-4 py-3 text-sm text-gray-700 dark:text-gray-300">{{ c.name }}</td>
          <td class="px-4 py-3 text-sm text-gray-500 dark:text-gray-400">{{ c.email }}</td>
          <td class="px-4 py-3 text-sm text-gray-700 dark:text-gray-300">{{ c.company }}</td>
          <td class="px-4 py-3 text-sm text-gray-500 dark:text-gray-400">{{ c.registeredAt }}</td>
          <td class="px-4 py-3 text-sm space-x-2">
            <button
              class="px-3 py-1 bg-green-600 text-white rounded text-xs font-medium hover:bg-green-700"
              @click="emit('approve', c.id)"
            >
              Approve
            </button>
            <button
              class="px-3 py-1 bg-red-600 text-white rounded text-xs font-medium hover:bg-red-700"
              @click="emit('reject', c.id)"
            >
              Reject
            </button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
