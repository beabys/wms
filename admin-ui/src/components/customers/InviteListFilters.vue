<script setup lang="ts">
import { ref, watch } from 'vue'
import WmsButton from '@/components/common/WmsButton.vue'

export interface InviteFilterValues {
  status?: string
  expired?: boolean
  created_after?: number
  created_before?: number
}

const props = defineProps<{
  modelValue: InviteFilterValues
}>()

const emit = defineEmits<{
  'update:modelValue': [value: InviteFilterValues]
}>()

const statusFilter = ref<string>('')
const expiredFilter = ref<string>('')
const dateFrom = ref<string>('')
const dateTo = ref<string>('')

function toUnixTimestamp(dateStr: string): number | undefined {
  if (!dateStr) return undefined
  return Math.floor(new Date(dateStr + 'T00:00:00').getTime() / 1000)
}

function fromUnixTimestamp(ts: number | undefined): string {
  if (!ts) return ''
  const d = new Date(ts * 1000)
  return d.toISOString().slice(0, 10)
}

watch(
  () => props.modelValue,
  (val) => {
    statusFilter.value = val.status || ''
    expiredFilter.value = val.expired === undefined ? '' : String(val.expired)
    dateFrom.value = fromUnixTimestamp(val.created_after)
    dateTo.value = fromUnixTimestamp(val.created_before)
  },
  { immediate: true },
)

function applyFilters() {
  const filters: InviteFilterValues = {}
  if (statusFilter.value) filters.status = statusFilter.value
  if (expiredFilter.value !== '') filters.expired = expiredFilter.value === 'true'
  const after = toUnixTimestamp(dateFrom.value)
  const before = toUnixTimestamp(dateTo.value)
  if (after !== undefined) filters.created_after = after
  if (before !== undefined) filters.created_before = before
  emit('update:modelValue', filters)
}

function resetFilters() {
  statusFilter.value = ''
  expiredFilter.value = ''
  dateFrom.value = ''
  dateTo.value = ''
  emit('update:modelValue', {})
}
</script>

<template>
  <div class="invite-list-filters">
    <label class="invite-list-filters__field">
      <span>Status:</span>
      <select v-model="statusFilter" class="invite-list-filters__select">
        <option value="">All</option>
        <option value="pending">Pending</option>
        <option value="used">Used</option>
        <option value="cancelled">Cancelled</option>
      </select>
    </label>

    <label class="invite-list-filters__field">
      <span>Expired:</span>
      <select v-model="expiredFilter" class="invite-list-filters__select">
        <option value="">All</option>
        <option value="true">Expired</option>
        <option value="false">Unexpired</option>
      </select>
    </label>

    <label class="invite-list-filters__field">
      <span>From:</span>
      <input
        v-model="dateFrom"
        type="date"
        class="invite-list-filters__date"
      />
    </label>

    <label class="invite-list-filters__field">
      <span>To:</span>
      <input
        v-model="dateTo"
        type="date"
        class="invite-list-filters__date"
      />
    </label>

    <div class="invite-list-filters__actions">
      <WmsButton size="sm" variant="primary" @click="applyFilters">
        Apply
      </WmsButton>
      <WmsButton size="sm" variant="secondary" @click="resetFilters">
        Reset
      </WmsButton>
    </div>
  </div>
</template>

<style scoped>
.invite-list-filters {
  display: flex;
  flex-wrap: wrap;
  align-items: flex-end;
  gap: var(--spacing-4);
  padding: var(--spacing-4);
  border-bottom: 1px solid var(--color-border);
  background-color: var(--color-card-bg);
}

.invite-list-filters__field {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-1);
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
}

.invite-list-filters__select,
.invite-list-filters__date {
  padding: var(--spacing-1-5) var(--spacing-3);
  border: 1px solid var(--color-input-border);
  border-radius: var(--radius-md);
  background-color: var(--color-input-bg);
  color: var(--color-text-primary);
  font-size: var(--text-sm);
  min-height: 2.25rem;
}

.invite-list-filters__actions {
  display: flex;
  gap: var(--spacing-2);
  align-items: flex-end;
  padding-bottom: var(--spacing-0-5);
}
</style>
