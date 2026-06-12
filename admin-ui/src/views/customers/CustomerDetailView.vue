<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import WmsCard from '@/components/common/WmsCard.vue'
import WmsButton from '@/components/common/WmsButton.vue'
import WmsInput from '@/components/common/WmsInput.vue'
import WmsConfirmDialog from '@/components/common/WmsConfirmDialog.vue'
import CustomerStatusBadge from '@/components/customers/CustomerStatusBadge.vue'
import { useCustomerDetailStore } from '@/stores/customerDetail'
import { useToastStore } from '@/stores/toast'

const route = useRoute()
const router = useRouter()
const store = useCustomerDetailStore()
const toastStore = useToastStore()

const editing = ref(false)
const confirmAction = ref<'approve' | 'reject' | 'suspend' | 'restore' | null>(null)
const showConfirm = ref(false)
const editForm = ref({
  company_name: '',
  email: '',
  phone: '',
  vat_number: '',
  address: '',
  city: '',
  postal_code: '',
  country: '',
})

const customerId = computed(() => route.params.id as string)

const auditTotalPages = computed(() =>
  Math.max(1, Math.ceil(store.auditPagination.total_items / store.auditPagination.page_size)),
)

onMounted(async () => {
  await store.fetchCustomer(customerId.value)
  await store.fetchAuditLogs(customerId.value)
  if (store.customer) {
    syncForm()
  }
})

function syncForm() {
  if (!store.customer) return
  editForm.value = {
    company_name: store.customer.company_name || '',
    email: store.customer.email || '',
    phone: store.customer.phone || '',
    vat_number: store.customer.vat_number || '',
    address: store.customer.address || '',
    city: store.customer.city || '',
    postal_code: store.customer.postal_code || '',
    country: store.customer.country || '',
  }
}

function startEditing() {
  syncForm()
  editing.value = true
}

function cancelEditing() {
  editing.value = false
  syncForm()
}

async function saveEditing() {
  try {
    await store.updateCustomer(customerId.value, editForm.value)
    toastStore.addToast('Customer updated', 'success')
    editing.value = false
  } catch {
    toastStore.addToast('Failed to update customer', 'error')
  }
}

function openConfirm(action: 'approve' | 'reject' | 'suspend' | 'restore') {
  confirmAction.value = action
  showConfirm.value = true
}

async function handleConfirm() {
  if (!confirmAction.value) return
  try {
    switch (confirmAction.value) {
      case 'approve':
        await store.approveCustomer(customerId.value)
        toastStore.addToast('Customer approved', 'success')
        break
      case 'reject':
        await store.rejectCustomer(customerId.value)
        toastStore.addToast('Customer rejected', 'success')
        break
      case 'suspend':
        await store.suspendCustomer(customerId.value)
        toastStore.addToast('Customer suspended', 'success')
        break
      case 'restore':
        await store.restoreCustomer(customerId.value)
        toastStore.addToast('Customer restored', 'success')
        break
    }
    await store.fetchAuditLogs(customerId.value)
  } catch {
    toastStore.addToast(`Failed to ${confirmAction.value} customer`, 'error')
  } finally {
    showConfirm.value = false
    confirmAction.value = null
  }
}

function handleCancel() {
  showConfirm.value = false
  confirmAction.value = null
}

function goToAuditPage(page: number) {
  store.fetchAuditLogs(customerId.value, { page })
}

function formatTimestamp(ts: number): string {
  return new Date(ts * 1000).toLocaleString()
}

function goBack() {
  router.push('/customers')
}

const confirmTitle = computed(() => {
  switch (confirmAction.value) {
    case 'approve': return 'Approve Customer'
    case 'reject': return 'Reject Customer'
    case 'suspend': return 'Suspend Customer'
    case 'restore': return 'Restore Customer'
    default: return 'Confirm'
  }
})

const confirmMessage = computed(() => {
  switch (confirmAction.value) {
    case 'approve': return 'Are you sure you want to approve this customer?'
    case 'reject': return 'Are you sure you want to reject this customer?'
    case 'suspend': return 'Are you sure you want to suspend this customer?'
    case 'restore': return 'Are you sure you want to restore this customer?'
    default: return 'Are you sure?'
  }
})

const confirmVariant = computed<'primary' | 'danger'>(() => {
  return confirmAction.value === 'reject' || confirmAction.value === 'suspend' ? 'danger' : 'primary'
})
</script>

<template>
  <div class="customer-detail">
    <div v-if="store.loading" class="customer-detail__loading">
      Loading customer details...
    </div>

    <div v-else-if="store.error && !store.customer" class="customer-detail__error">
      <p>{{ store.error }}</p>
      <WmsButton variant="secondary" @click="goBack">Back to Customers</WmsButton>
    </div>

    <template v-else-if="store.customer">
      <!-- Header -->
      <div class="customer-detail__header">
        <div class="customer-detail__header-left">
          <button class="customer-detail__back" @click="goBack">&larr; Back</button>
          <h1 class="customer-detail__title">{{ store.customer.company_name }}</h1>
          <CustomerStatusBadge :status="store.customer.status" />
        </div>
        <div class="customer-detail__header-right">
          <WmsButton
            v-if="!editing"
            variant="secondary"
            @click="startEditing"
          >
            Edit
          </WmsButton>
          <template v-else>
            <WmsButton variant="secondary" @click="cancelEditing">Cancel</WmsButton>
            <WmsButton variant="primary" @click="saveEditing">Save</WmsButton>
          </template>
        </div>
      </div>

      <div class="customer-detail__grid">
        <!-- Company Info Card -->
        <WmsCard class="customer-detail__card">
          <template #header>
            <h2 class="customer-detail__card-title">Company Information</h2>
          </template>
          <div class="customer-detail__fields">
            <WmsInput
              label="Company Name"
              v-if="editing"
              v-model="editForm.company_name"
            />
            <div v-else class="customer-detail__field">
              <span class="customer-detail__field-label">Company Name</span>
              <span class="customer-detail__field-value">{{ store.customer.company_name }}</span>
            </div>

            <WmsInput
              label="Email"
              v-if="editing"
              v-model="editForm.email"
              type="email"
            />
            <div v-else class="customer-detail__field">
              <span class="customer-detail__field-label">Email</span>
              <span class="customer-detail__field-value">{{ store.customer.email }}</span>
            </div>

            <WmsInput
              label="Phone"
              v-if="editing"
              v-model="editForm.phone"
            />
            <div v-else class="customer-detail__field">
              <span class="customer-detail__field-label">Phone</span>
              <span class="customer-detail__field-value">{{ store.customer.phone || '-' }}</span>
            </div>

            <WmsInput
              label="VAT Number"
              v-if="editing"
              v-model="editForm.vat_number"
            />
            <div v-else class="customer-detail__field">
              <span class="customer-detail__field-label">VAT Number</span>
              <span class="customer-detail__field-value">{{ store.customer.vat_number || '-' }}</span>
            </div>

            <WmsInput
              label="Address"
              v-if="editing"
              v-model="editForm.address"
            />
            <div v-else class="customer-detail__field">
              <span class="customer-detail__field-label">Address</span>
              <span class="customer-detail__field-value">{{ store.customer.address || '-' }}</span>
            </div>

            <WmsInput
              label="City"
              v-if="editing"
              v-model="editForm.city"
            />
            <div v-else class="customer-detail__field">
              <span class="customer-detail__field-label">City</span>
              <span class="customer-detail__field-value">{{ store.customer.city || '-' }}</span>
            </div>

            <WmsInput
              label="Postal Code"
              v-if="editing"
              v-model="editForm.postal_code"
            />
            <div v-else class="customer-detail__field">
              <span class="customer-detail__field-label">Postal Code</span>
              <span class="customer-detail__field-value">{{ store.customer.postal_code || '-' }}</span>
            </div>

            <WmsInput
              label="Country"
              v-if="editing"
              v-model="editForm.country"
            />
            <div v-else class="customer-detail__field">
              <span class="customer-detail__field-label">Country</span>
              <span class="customer-detail__field-value">{{ store.customer.country || '-' }}</span>
            </div>
          </div>
        </WmsCard>

        <!-- Actions Card -->
        <WmsCard class="customer-detail__card">
          <template #header>
            <h2 class="customer-detail__card-title">Actions</h2>
          </template>
          <div class="customer-detail__actions">
            <WmsButton
              v-if="store.customer.status === 'pending'"
              variant="primary"
              @click="openConfirm('approve')"
            >
              Approve
            </WmsButton>
            <WmsButton
              v-if="store.customer.status === 'pending'"
              variant="danger"
              @click="openConfirm('reject')"
            >
              Reject
            </WmsButton>
            <WmsButton
              v-if="store.customer.status === 'active'"
              variant="danger"
              @click="openConfirm('suspend')"
            >
              Suspend
            </WmsButton>
            <WmsButton
              v-if="store.customer.status === 'suspended'"
              variant="primary"
              @click="openConfirm('restore')"
            >
              Restore
            </WmsButton>
          </div>
        </WmsCard>
      </div>

      <!-- Audit Log Card -->
      <WmsCard class="customer-detail__card">
        <template #header>
          <h2 class="customer-detail__card-title">Audit Log</h2>
        </template>
        <div class="customer-detail__audit">
          <table v-if="store.auditLogs.length > 0" class="customer-detail__audit-table">
            <thead>
              <tr>
                <th>Action</th>
                <th>Performed By</th>
                <th>Details</th>
                <th>Timestamp</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="entry in store.auditLogs" :key="entry.id">
                <td><CustomerStatusBadge :status="entry.action" /></td>
                <td>{{ entry.performed_by }}</td>
                <td>{{ entry.details || '-' }}</td>
                <td>{{ formatTimestamp(entry.created_at) }}</td>
              </tr>
            </tbody>
          </table>
          <p v-else class="customer-detail__audit-empty">No audit entries found.</p>

          <div
            v-if="store.auditPagination.total_items > 0"
            class="customer-detail__audit-pagination"
          >
            <WmsButton
              size="sm"
              variant="secondary"
              :disabled="store.auditPagination.page <= 1"
              @click="goToAuditPage(store.auditPagination.page - 1)"
            >
              Previous
            </WmsButton>
            <span class="customer-detail__audit-page-info">
              Page {{ store.auditPagination.page }} of {{ auditTotalPages }}
            </span>
            <WmsButton
              size="sm"
              variant="secondary"
              :disabled="store.auditPagination.page >= auditTotalPages"
              @click="goToAuditPage(store.auditPagination.page + 1)"
            >
              Next
            </WmsButton>
          </div>
        </div>
      </WmsCard>
    </template>

    <WmsConfirmDialog
      :visible="showConfirm"
      :title="confirmTitle"
      :message="confirmMessage"
      :variant="confirmVariant"
      @confirm="handleConfirm"
      @cancel="handleCancel"
    />
  </div>
</template>

<style scoped>
.customer-detail {
  padding: var(--spacing-6);
  max-width: 80rem;
  margin: 0 auto;
}

.customer-detail__loading {
  padding: var(--spacing-8);
  text-align: center;
  color: var(--color-text-secondary);
}

.customer-detail__error {
  padding: var(--spacing-8);
  text-align: center;
  color: var(--color-error);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--spacing-4);
}

.customer-detail__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--spacing-6);
  flex-wrap: wrap;
  gap: var(--spacing-4);
}

.customer-detail__header-left {
  display: flex;
  align-items: center;
  gap: var(--spacing-3);
}

.customer-detail__header-right {
  display: flex;
  gap: var(--spacing-2);
}

.customer-detail__back {
  background: none;
  border: none;
  color: var(--color-primary);
  cursor: pointer;
  font-size: var(--text-sm);
  padding: 0;
}

.customer-detail__back:hover {
  text-decoration: underline;
}

.customer-detail__title {
  font-size: var(--text-2xl);
  font-weight: var(--font-bold);
  color: var(--color-text-primary);
  margin: 0;
}

.customer-detail__grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: var(--spacing-6);
  margin-bottom: var(--spacing-6);
}

@media (min-width: 768px) {
  .customer-detail__grid {
    grid-template-columns: 2fr 1fr;
  }
}

.customer-detail__card {
  margin-bottom: var(--spacing-6);
}

.customer-detail__card-title {
  margin: 0;
  font-size: var(--text-lg);
  font-weight: var(--font-semibold);
  color: var(--color-text-primary);
}

.customer-detail__fields {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-4);
}

.customer-detail__field {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-1);
}

.customer-detail__field-label {
  font-size: var(--text-xs);
  font-weight: var(--font-semibold);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: var(--color-text-muted);
}

.customer-detail__field-value {
  font-size: var(--text-sm);
  color: var(--color-text-primary);
}

.customer-detail__actions {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-3);
}

.customer-detail__audit-table {
  width: 100%;
  border-collapse: collapse;
}

.customer-detail__audit-table th {
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

.customer-detail__audit-table td {
  padding: var(--spacing-3) var(--spacing-4);
  font-size: var(--text-sm);
  color: var(--color-text-primary);
  border-bottom: 1px solid var(--color-border);
}

.customer-detail__audit-empty {
  padding: var(--spacing-8);
  text-align: center;
  color: var(--color-text-muted);
}

.customer-detail__audit-pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--spacing-4);
  padding: var(--spacing-4);
  border-top: 1px solid var(--color-border);
}

.customer-detail__audit-page-info {
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
}
</style>
