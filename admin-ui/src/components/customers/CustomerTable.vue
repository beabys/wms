<script setup lang="ts">
import type { CustomerResponse } from '@/api/types'
import WmsButton from '@/components/common/WmsButton.vue'
import CustomerStatusBadge from '@/components/customers/CustomerStatusBadge.vue'

defineProps<{
  customers: CustomerResponse[]
  loading: boolean
}>()

const emit = defineEmits<{
  approve: [customer: CustomerResponse]
  reject: [customer: CustomerResponse]
  suspend: [customer: CustomerResponse]
  restore: [customer: CustomerResponse]
  rowClick: [customer: CustomerResponse]
}>()

function formatDate(dateStr: string): string {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString()
}
</script>

<template>
  <div class="customer-table-wrapper">
    <div v-if="loading" class="customer-table__loading">
      Loading customers...
    </div>

    <table v-else class="customer-table">
      <thead>
        <tr>
          <th>Company</th>
          <th>Email</th>
          <th>Status</th>
          <th>Created</th>
          <th>Actions</th>
        </tr>
      </thead>
      <tbody>
        <tr v-if="customers.length === 0">
          <td colspan="5" class="customer-table__empty">
            No customers found.
          </td>
        </tr>
        <tr v-for="c in customers" :key="c.id" class="customer-table__row" @click="emit('rowClick', c)">
          <td class="customer-table__cell-company">{{ c.company_name }}</td>
          <td>{{ c.email }}</td>
          <td>
            <CustomerStatusBadge :status="c.status" />
          </td>
          <td>{{ formatDate(c.created_at) }}</td>
          <td class="customer-table__actions">
            <WmsButton
              v-if="c.status === 'pending'"
              size="sm"
              variant="primary"
              @click="emit('approve', c)"
            >
              Approve
            </WmsButton>
            <WmsButton
              v-if="c.status === 'pending'"
              size="sm"
              variant="danger"
              @click="emit('reject', c)"
            >
              Reject
            </WmsButton>
            <WmsButton
              v-if="c.status === 'active'"
              size="sm"
              variant="danger"
              @click="emit('suspend', c)"
            >
              Suspend
            </WmsButton>
            <WmsButton
              v-if="c.status === 'suspended'"
              size="sm"
              variant="primary"
              @click="emit('restore', c)"
            >
              Restore
            </WmsButton>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<style scoped>
.customer-table-wrapper {
  overflow-x: auto;
}

.customer-table__loading {
  padding: var(--spacing-8);
  text-align: center;
  color: var(--color-text-secondary);
}

.customer-table {
  width: 100%;
  border-collapse: collapse;
}

.customer-table thead th {
  padding: var(--spacing-3) var(--spacing-4);
  text-align: left;
  font-size: var(--text-xs);
  font-weight: var(--font-semibold);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: var(--color-text-muted);
  border-bottom: 1px solid var(--color-border);
  background-color: var(--color-bg-secondary);
}

.customer-table tbody td {
  padding: var(--spacing-3) var(--spacing-4);
  font-size: var(--text-sm);
  color: var(--color-text-primary);
  border-bottom: 1px solid var(--color-border);
}

.customer-table__cell-company {
  font-weight: var(--font-medium);
}

.customer-table__empty {
  text-align: center;
  color: var(--color-text-muted);
  padding: var(--spacing-8) !important;
}

.customer-table__row {
  cursor: pointer;
  transition: background-color 0.1s ease;
}

.customer-table__row:hover {
  background-color: var(--color-bg-secondary);
}

.customer-table__actions {
  display: flex;
  gap: var(--spacing-2);
  white-space: nowrap;
}
</style>
