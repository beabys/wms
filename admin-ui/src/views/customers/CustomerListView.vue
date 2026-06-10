<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import WmsCard from '@/components/common/WmsCard.vue'
import WmsButton from '@/components/common/WmsButton.vue'
import CustomerTable from '@/components/customers/CustomerTable.vue'
import ApproveRejectDialog from '@/components/customers/ApproveRejectDialog.vue'
import WmsConfirmDialog from '@/components/common/WmsConfirmDialog.vue'
import InviteCustomerModal from '@/views/customers/InviteCustomerModal.vue'
import { useCustomerStore } from '@/stores/customers'
import { useToastStore } from '@/stores/toast'
import type { CustomerResponse } from '@/api/types'

const router = useRouter()
const store = useCustomerStore()
const toastStore = useToastStore()

const statusFilter = ref('')
const showInviteModal = ref(false)
const selectedCustomer = ref<CustomerResponse | null>(null)
const dialogAction = ref<'approve' | 'reject' | 'suspend' | 'restore'>('approve')
const showDialog = ref(false)
const showRestoreConfirm = ref(false)

onMounted(() => {
  store.fetchCustomers()
})

function onFilterChange() {
  store.fetchCustomers({ status: statusFilter.value || undefined })
}

function goToPage(page: number) {
  store.fetchCustomers({ status: statusFilter.value || undefined, page })
}

function openApprove(customer: CustomerResponse) {
  selectedCustomer.value = customer
  dialogAction.value = 'approve'
  showDialog.value = true
}

function openReject(customer: CustomerResponse) {
  selectedCustomer.value = customer
  dialogAction.value = 'reject'
  showDialog.value = true
}

function openSuspend(customer: CustomerResponse) {
  selectedCustomer.value = customer
  dialogAction.value = 'suspend'
  showDialog.value = true
}

function openRestore(customer: CustomerResponse) {
  selectedCustomer.value = customer
  dialogAction.value = 'restore'
  showRestoreConfirm.value = true
}

function onRowClick(customer: CustomerResponse) {
  router.push(`/customers/${customer.id}`)
}

async function handleDialogConfirm() {
  if (!selectedCustomer.value) return
  try {
    switch (dialogAction.value) {
      case 'approve':
        await store.approveCustomer(selectedCustomer.value.id)
        toastStore.addToast('Customer approved', 'success')
        break
      case 'reject':
        await store.rejectCustomer(selectedCustomer.value.id)
        toastStore.addToast('Customer rejected', 'success')
        break
      case 'suspend':
        await store.suspendCustomer(selectedCustomer.value.id)
        toastStore.addToast('Customer suspended', 'success')
        break
      case 'restore':
        await store.restoreCustomer(selectedCustomer.value.id)
        toastStore.addToast('Customer restored', 'success')
        break
    }
  } catch {
    toastStore.addToast(`Failed to ${dialogAction.value} customer`, 'error')
  } finally {
    showDialog.value = false
    showRestoreConfirm.value = false
    selectedCustomer.value = null
  }
}

function handleDialogCancel() {
  showDialog.value = false
  showRestoreConfirm.value = false
  selectedCustomer.value = null
}

const totalPages = computed(() => {
  return Math.max(1, Math.ceil(store.pagination.total_items / store.pagination.page_size))
})
</script>

<template>
  <div class="customer-list-view">
    <header class="customer-list-view__header">
      <h1 class="customer-list-view__title">Customers</h1>
      <WmsButton variant="primary" @click="showInviteModal = true">
        Invite Customer
      </WmsButton>
    </header>

    <WmsCard padding="none" class="customer-list-view__card">
      <div class="customer-list-view__toolbar">
        <label class="customer-list-view__filter">
          <span>Status:</span>
          <select v-model="statusFilter" class="customer-list-view__select" @change="onFilterChange">
            <option value="">All</option>
            <option value="pending">Pending</option>
            <option value="active">Active</option>
            <option value="rejected">Rejected</option>
            <option value="suspended">Suspended</option>
          </select>
        </label>
      </div>

      <CustomerTable
        :customers="store.customers"
        :loading="store.loading"
        @approve="openApprove"
        @reject="openReject"
        @suspend="openSuspend"
        @restore="openRestore"
        @row-click="onRowClick"
      />

      <div v-if="store.pagination.total_items > 0" class="customer-list-view__pagination">
        <WmsButton
          size="sm"
          variant="secondary"
          :disabled="store.pagination.page <= 1"
          @click="goToPage(store.pagination.page - 1)"
        >
          Previous
        </WmsButton>
        <span class="customer-list-view__page-info">
          Page {{ store.pagination.page }} of {{ totalPages }}
        </span>
        <WmsButton
          size="sm"
          variant="secondary"
          :disabled="store.pagination.page >= totalPages"
          @click="goToPage(store.pagination.page + 1)"
        >
          Next
        </WmsButton>
      </div>
    </WmsCard>

    <ApproveRejectDialog
      :visible="showDialog"
      :customer="selectedCustomer"
      :action="dialogAction === 'restore' ? 'approve' : dialogAction"
      @confirm="handleDialogConfirm"
      @cancel="handleDialogCancel"
    />

    <WmsConfirmDialog
      :visible="showRestoreConfirm"
      title="Restore Customer"
      message="Are you sure you want to restore this customer?"
      variant="primary"
      confirm-text="Restore"
      @confirm="handleDialogConfirm"
      @cancel="handleDialogCancel"
    />

    <InviteCustomerModal
      :visible="showInviteModal"
      @close="showInviteModal = false"
    />
  </div>
</template>

<style scoped>
.customer-list-view {
  padding: var(--spacing-6);
  max-width: 80rem;
  margin: 0 auto;
}

.customer-list-view__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--spacing-6);
}

.customer-list-view__title {
  font-size: var(--text-2xl);
  font-weight: var(--font-bold);
  color: var(--color-text-primary);
  margin: 0;
}

.customer-list-view__card {
  overflow: hidden;
}

.customer-list-view__toolbar {
  display: flex;
  align-items: center;
  gap: var(--spacing-4);
  padding: var(--spacing-4);
  border-bottom: 1px solid var(--color-border);
  background-color: var(--color-card-bg);
}

.customer-list-view__filter {
  display: flex;
  align-items: center;
  gap: var(--spacing-2);
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
}

.customer-list-view__select {
  padding: var(--spacing-1-5) var(--spacing-3);
  border: 1px solid var(--color-input-border);
  border-radius: var(--radius-md);
  background-color: var(--color-input-bg);
  color: var(--color-text-primary);
  font-size: var(--text-sm);
  min-height: 2.25rem;
}

.customer-list-view__pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--spacing-4);
  padding: var(--spacing-4);
  border-top: 1px solid var(--color-border);
}

.customer-list-view__page-info {
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
}

</style>
