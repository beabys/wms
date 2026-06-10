<script setup lang="ts">
import type { UserResponse } from '@/api/types'
import WmsBadge from '@/components/common/WmsBadge.vue'

defineProps<{
  users: UserResponse[]
  loading: boolean
}>()

function roleVariant(role: string): 'success' | 'warning' | 'error' | 'info' | 'neutral' {
  switch (role) {
    case 'admin': return 'info'
    case 'warehouse_staff': return 'success'
    case 'billing_manager': return 'warning'
    default: return 'neutral'
  }
}

function formatDate(dateStr: string): string {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return d.toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' })
}
</script>

<template>
  <div class="user-table">
    <div v-if="loading" class="user-table__loading">Loading users...</div>
    <table v-else-if="users.length > 0" class="user-table__table">
      <thead>
        <tr>
          <th>Email</th>
          <th>Name</th>
          <th>Role</th>
          <th>Created</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="user in users" :key="user.id">
          <td>{{ user.email }}</td>
          <td>{{ user.name }}</td>
          <td>
            <WmsBadge :variant="roleVariant(user.role)" size="sm">
              {{ user.role.replace(/_/g, ' ') }}
            </WmsBadge>
          </td>
          <td>{{ formatDate(user.created_at) }}</td>
        </tr>
      </tbody>
    </table>
    <div v-else class="user-table__empty">No users found.</div>
  </div>
</template>

<style scoped>
.user-table__loading,
.user-table__empty {
  padding: var(--spacing-8);
  text-align: center;
  color: var(--color-text-muted);
}

.user-table__table {
  width: 100%;
  border-collapse: collapse;
}

.user-table__table th,
.user-table__table td {
  padding: var(--spacing-3) var(--spacing-4);
  text-align: left;
  border-bottom: 1px solid var(--color-border);
}

.user-table__table th {
  font-weight: var(--font-medium);
  font-size: var(--text-xs);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: var(--color-text-secondary);
}

.user-table__table td {
  font-size: var(--text-sm);
  color: var(--color-text-primary);
}

.user-table__table tbody tr:hover {
  background-color: var(--color-bg-secondary);
}
</style>
