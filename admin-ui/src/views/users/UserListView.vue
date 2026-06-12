<script setup lang="ts">
import { onMounted, ref, watch, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/users'
import { useToastStore } from '@/stores/toast'
import UserTable from '@/components/users/UserTable.vue'
import UserFilters from '@/components/users/UserFilters.vue'
import WmsCard from '@/components/common/WmsCard.vue'
import WmsButton from '@/components/common/WmsButton.vue'

const router = useRouter()
const store = useUserStore()
const toastStore = useToastStore()

const filters = ref<{ role?: string }>({})

onMounted(() => {
  store.fetchUsers()
})

watch(filters, (newVal) => {
  store.fetchUsers({ role: newVal.role, page: 1, page_size: 10 })
}, { deep: true })

watch(() => store.error, (err) => {
  if (err) toastStore.addToast(err, 'error')
})

function goToCreate() {
  router.push('/users/create')
}

function nextPage() {
  store.fetchUsers({ role: filters.value.role, page: store.pagination.page + 1, page_size: store.pagination.page_size })
}

function prevPage() {
  if (store.pagination.page > 1) {
    store.fetchUsers({ role: filters.value.role, page: store.pagination.page - 1, page_size: store.pagination.page_size })
  }
}

const totalPages = computed(() =>
  Math.max(1, Math.ceil(store.pagination.total_items / store.pagination.page_size))
)
</script>

<template>
  <div class="user-list-view">
    <header class="user-list-view__header">
      <h1 class="user-list-view__title">Users</h1>
      <WmsButton variant="primary" size="md" @click="goToCreate">
        Create User
      </WmsButton>
    </header>

    <WmsCard padding="md">
      <template #header>
        <UserFilters v-model="filters" :loading="store.loading" />
      </template>
      <UserTable :users="store.users" :loading="store.loading" />
      <template #footer>
        <div class="user-list-view__pagination">
          <span class="user-list-view__page-info">
            Page {{ store.pagination.page }} of {{ totalPages }}
            ({{ store.pagination.total_items }} total)
          </span>
          <div class="user-list-view__page-controls">
            <WmsButton
              variant="secondary"
              size="sm"
              :disabled="store.pagination.page <= 1 || store.loading"
              @click="prevPage"
            >
              Previous
            </WmsButton>
            <WmsButton
              variant="secondary"
              size="sm"
              :disabled="store.pagination.page >= totalPages || store.loading"
              @click="nextPage"
            >
              Next
            </WmsButton>
          </div>
        </div>
      </template>
    </WmsCard>

  </div>
</template>

<style scoped>
.user-list-view {
  padding: var(--spacing-6);
  max-width: 80rem;
  margin: 0 auto;
}

.user-list-view__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--spacing-6);
}

.user-list-view__title {
  font-size: var(--text-2xl);
  font-weight: var(--font-bold);
  color: var(--color-text-primary);
}

.user-list-view__pagination {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.user-list-view__page-info {
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
}

.user-list-view__page-controls {
  display: flex;
  gap: var(--spacing-2);
}
</style>
