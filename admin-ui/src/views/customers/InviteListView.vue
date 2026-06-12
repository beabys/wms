<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import WmsCard from '@/components/common/WmsCard.vue'
import WmsBadge from '@/components/common/WmsBadge.vue'
import WmsButton from '@/components/common/WmsButton.vue'
import WmsConfirmDialog from '@/components/common/WmsConfirmDialog.vue'
import InviteListFilters from '@/components/customers/InviteListFilters.vue'
import type { InviteFilterValues } from '@/components/customers/InviteListFilters.vue'
import InviteCustomerModal from '@/views/customers/InviteCustomerModal.vue'
import { useInviteStore } from '@/stores/invites'
import { useToastStore } from '@/stores/toast'
import type { InviteEntry } from '@/api/types'

const store = useInviteStore()
const toastStore = useToastStore()

const showInviteModal = ref(false)
const showCancelConfirm = ref(false)
const cancelTarget = ref<InviteEntry | null>(null)
const filters = ref<InviteFilterValues>({})

onMounted(() => {
  store.fetchInvites()
})

function onFiltersChange(val: InviteFilterValues) {
  filters.value = val
  store.fetchInvites({ ...val, page: 1 })
}

function goToPage(page: number) {
  store.fetchInvites({ ...filters.value, page })
}

async function copyLink(invite: InviteEntry) {
  const link = `${window.location.origin}/signup?token=${invite.token}`
  try {
    await navigator.clipboard.writeText(link)
    toastStore.addToast('Invite link copied!', 'success')
  } catch {
    toastStore.addToast('Failed to copy link', 'error')
  }
}

function statusBadgeVariant(invite: InviteEntry): 'success' | 'error' | 'neutral' | 'warning' {
  if (invite.status === 'used') return 'neutral'
  if (invite.status === 'cancelled') return 'error'
  if (invite.status === 'pending' && invite.expires_at < Math.floor(Date.now() / 1000)) return 'warning'
  return 'success'
}

function statusLabel(invite: InviteEntry): string {
  if (invite.status === 'used') return 'Used'
  if (invite.status === 'cancelled') return 'Cancelled'
  if (invite.expires_at < Math.floor(Date.now() / 1000)) return 'Expired'
  return 'Pending'
}

function handleCancel(invite: InviteEntry) {
  cancelTarget.value = invite
  showCancelConfirm.value = true
}

async function confirmCancel() {
  if (!cancelTarget.value) return
  const invite = cancelTarget.value
  showCancelConfirm.value = false
  cancelTarget.value = null
  const err = await store.cancelInvite(invite.token)
  if (err) {
    toastStore.addToast(err, 'error')
  } else {
    toastStore.addToast('Invitation cancelled', 'success')
  }
}

const totalPages = computed(() => {
  const size = store.pagination.page_size || 1
  return Math.max(1, Math.ceil(store.pagination.total_items / size))
})
</script>

<template>
  <div class="invite-list-view">
    <header class="invite-list-view__header">
      <h1 class="invite-list-view__title">Sent Invitations</h1>
      <WmsButton variant="primary" @click="showInviteModal = true">
        Send New Invite
      </WmsButton>
    </header>

    <WmsCard padding="none" class="invite-list-view__card">
      <InviteListFilters
        :model-value="filters"
        @update:model-value="onFiltersChange"
      />

      <!-- Loading state -->
      <div v-if="store.loading" class="invite-list-view__state">
        Loading invitations...
      </div>

      <!-- Error state -->
      <div v-else-if="store.error" class="invite-list-view__state invite-list-view__state--error">
        {{ store.error }}
        <WmsButton size="sm" variant="secondary" @click="store.fetchInvites()">
          Retry
        </WmsButton>
      </div>

      <!-- Empty state -->
      <div v-else-if="store.invites.length === 0" class="invite-list-view__state">
        No invitations found.
      </div>

      <!-- Table -->
      <table v-else class="invite-list-view__table">
        <thead>
          <tr>
            <th>Email</th>
            <th>Status</th>
            <th>Invited By</th>
            <th>Created</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="invite in store.invites" :key="invite.id">
            <td>{{ invite.email }}</td>
            <td>
              <WmsBadge :variant="statusBadgeVariant(invite)">
                {{ statusLabel(invite) }}
              </WmsBadge>
            </td>
            <td>{{ invite.invited_by }}</td>
            <td>{{ new Date(invite.created_at * 1000).toLocaleDateString() }}</td>
            <td class="invite-list-view__actions-cell">
              <WmsButton size="sm" variant="secondary" @click="copyLink(invite)">
                Copy Link
              </WmsButton>
              <WmsButton
                v-if="invite.status === 'pending'"
                size="sm"
                variant="danger"
                @click="handleCancel(invite)"
              >
                Cancel
              </WmsButton>
            </td>
          </tr>
        </tbody>
      </table>

      <!-- Pagination -->
      <div v-if="store.pagination.total_items > 0" class="invite-list-view__pagination">
        <WmsButton
          size="sm"
          variant="secondary"
          :disabled="store.pagination.page <= 1"
          @click="goToPage(store.pagination.page - 1)"
        >
          Previous
        </WmsButton>
        <span class="invite-list-view__page-info">
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

    <InviteCustomerModal
      :visible="showInviteModal"
      @close="showInviteModal = false"
    />

    <WmsConfirmDialog
      :visible="showCancelConfirm"
      title="Cancel Invitation"
      :message="`Cancel invitation for ${cancelTarget?.email}?`"
      confirm-text="Cancel Invitation"
      cancel-text="Go Back"
      variant="danger"
      @confirm="confirmCancel"
      @cancel="showCancelConfirm = false"
    />
  </div>
</template>

<style scoped>
.invite-list-view {
  padding: var(--spacing-6);
  max-width: 80rem;
  margin: 0 auto;
}

.invite-list-view__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--spacing-6);
}

.invite-list-view__title {
  font-size: var(--text-2xl);
  font-weight: var(--font-bold);
  color: var(--color-text-primary);
  margin: 0;
}

.invite-list-view__card {
  overflow: hidden;
}

.invite-list-view__state {
  padding: var(--spacing-8);
  text-align: center;
  color: var(--color-text-secondary);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--spacing-3);
}

.invite-list-view__state--error {
  color: var(--color-error-text);
}

.invite-list-view__table {
  width: 100%;
  border-collapse: collapse;
}

.invite-list-view__table th {
  text-align: left;
  padding: var(--spacing-3) var(--spacing-4);
  font-size: var(--text-xs);
  font-weight: var(--font-semibold);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: var(--color-text-secondary);
  background-color: var(--color-bg-tertiary);
  border-bottom: 1px solid var(--color-border);
}

.invite-list-view__table td {
  padding: var(--spacing-3) var(--spacing-4);
  font-size: var(--text-sm);
  color: var(--color-text-primary);
  border-bottom: 1px solid var(--color-border);
}

.invite-list-view__actions-cell {
  display: flex;
  gap: var(--spacing-2);
  align-items: center;
}

.invite-list-view__pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--spacing-4);
  padding: var(--spacing-4);
  border-top: 1px solid var(--color-border);
}

.invite-list-view__page-info {
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
}
</style>
