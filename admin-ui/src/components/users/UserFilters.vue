<script setup lang="ts">
const props = withDefaults(defineProps<{
  modelValue: { role?: string }
  loading?: boolean
}>(), {
  loading: false,
})

const emit = defineEmits<{
  'update:modelValue': [value: { role?: string }]
}>()

const roles = [
  { value: '', label: 'All Roles' },
  { value: 'admin', label: 'Admin' },
  { value: 'warehouse_staff', label: 'Warehouse Staff' },
  { value: 'billing_manager', label: 'Billing Manager' },
]

function onRoleChange(e: Event) {
  const role = (e.target as HTMLSelectElement).value
  emit('update:modelValue', { ...props.modelValue, role: role || undefined })
}
</script>

<template>
  <div class="user-filters">
    <div class="user-filters__group">
      <label class="user-filters__label" for="filter-role">Role</label>
      <select
        id="filter-role"
        class="user-filters__select"
        :value="modelValue.role || ''"
        :disabled="loading"
        @change="onRoleChange"
      >
        <option v-for="r in roles" :key="r.value" :value="r.value">
          {{ r.label }}
        </option>
      </select>
    </div>
  </div>
</template>

<style scoped>
.user-filters {
  display: flex;
  gap: var(--spacing-4);
  align-items: flex-end;
}

.user-filters__group {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-1);
}

.user-filters__label {
  font-size: var(--text-sm);
  font-weight: var(--font-medium);
  color: var(--color-text-secondary);
}

.user-filters__select {
  padding: var(--spacing-2-5) var(--spacing-3);
  border: 1px solid var(--color-input-border);
  border-radius: var(--radius-md);
  background-color: var(--color-input-bg);
  color: var(--color-text-primary);
  font-size: var(--text-sm);
  min-width: 10rem;
  min-height: 2.5rem;
  outline: none;
  transition: border-color 0.15s ease, box-shadow 0.15s ease;
}

.user-filters__select:focus {
  border-color: var(--color-input-focus);
  box-shadow: 0 0 0 3px var(--color-primary-light);
}

.user-filters__select:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
