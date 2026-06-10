import { defineStore } from 'pinia'
import { ref } from 'vue'
import { authClient } from '@/api/authClient'
import type { UserResponse, Pagination, CreateUserRequest } from '@/api/types'

export const useUserStore = defineStore('users', () => {
  const users = ref<UserResponse[]>([])
  const pagination = ref<Pagination>({ page: 1, page_size: 10, total_items: 0 })
  const loading = ref(false)
  const error = ref('')

  async function fetchUsers(params?: { role?: string; page?: number; page_size?: number }) {
    loading.value = true
    error.value = ''
    try {
      const res = await authClient.listUsers(params)
      users.value = res.users
      pagination.value = res.pagination
    } catch (err: any) {
      error.value = err?.message || 'Failed to fetch users'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function createUser(data: CreateUserRequest): Promise<UserResponse> {
    loading.value = true
    error.value = ''
    try {
      const user = await authClient.createUser(data)
      return user
    } catch (err: any) {
      error.value = err?.message || 'Failed to create user'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function inviteUser(email: string) {
    loading.value = true
    error.value = ''
    try {
      return await authClient.inviteUser(email)
    } catch (err: any) {
      error.value = err?.message || 'Failed to invite user'
      throw err
    } finally {
      loading.value = false
    }
  }

  function $reset() {
    users.value = []
    pagination.value = { page: 1, page_size: 10, total_items: 0 }
    loading.value = false
    error.value = ''
  }

  return {
    users,
    pagination,
    loading,
    error,
    fetchUsers,
    createUser,
    inviteUser,
    $reset,
  }
})
