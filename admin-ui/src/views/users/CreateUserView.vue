<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/users'
import { useToastStore } from '@/stores/toast'
import type { CreateUserRequest } from '@/api/types'
import UserForm from '@/components/users/UserForm.vue'
import WmsCard from '@/components/common/WmsCard.vue'
import WmsButton from '@/components/common/WmsButton.vue'

const router = useRouter()
const store = useUserStore()
const toastStore = useToastStore()

const formData = ref<CreateUserRequest>({
  email: '',
  password: '',
  name: '',
  role: '',
})

const submitting = ref(false)

async function onSubmit() {
  submitting.value = true
  try {
    await store.createUser(formData.value)
    toastStore.addToast('User created successfully', 'success')
    router.push('/users')
  } catch (err: any) {
    toastStore.addToast(err?.message || 'Failed to create user', 'error')
  } finally {
    submitting.value = false
  }
}

function goBack() {
  router.push('/users')
}
</script>

<template>
  <div class="create-user-view">
    <header class="create-user-view__header">
      <h1 class="create-user-view__title">Create User</h1>
    </header>

    <WmsCard padding="md">
      <UserForm
        v-model="formData"
        :loading="submitting"
        @submit="onSubmit"
      />
      <template #footer>
        <div class="create-user-view__actions">
          <WmsButton variant="ghost" @click="goBack">
            Cancel
          </WmsButton>
        </div>
      </template>
    </WmsCard>
  </div>
</template>

<style scoped>
.create-user-view {
  padding: var(--spacing-6);
  max-width: 40rem;
  margin: 0 auto;
}

.create-user-view__header {
  margin-bottom: var(--spacing-6);
}

.create-user-view__title {
  font-size: var(--text-2xl);
  font-weight: var(--font-bold);
  color: var(--color-text-primary);
}

.create-user-view__actions {
  display: flex;
  justify-content: flex-end;
}
</style>
