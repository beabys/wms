import { defineStore } from 'pinia'
import { ref } from 'vue'
import { authClient } from '@/api/authClient'
import type { InviteEntry, Pagination, InviteUserResponse } from '@/api/types'

export const useInviteStore = defineStore('invites', () => {
  const invites = ref<InviteEntry[]>([])
  const pagination = ref<Pagination>({ page: 1, page_size: 10, total_items: 0 })
  const loading = ref(false)
  const error = ref('')

  async function fetchInvites(params?: {
    page?: number
    page_size?: number
    status?: string
    expired?: boolean
    created_after?: number
    created_before?: number
  }) {
    loading.value = true
    error.value = ''
    try {
      const res = await authClient.listInvites(params)
      invites.value = res.invites
      pagination.value = res.pagination
    } catch (err: any) {
      error.value = err?.message || 'Failed to fetch invites'
      invites.value = []
    } finally {
      loading.value = false
    }
  }

  async function cancelInvite(token: string) {
    try {
      await authClient.cancelInvite(token)
      await fetchInvites()
      return null
    } catch (err: any) {
      return err?.message || 'Failed to cancel invitation'
    }
  }

  async function inviteUser(email: string): Promise<InviteUserResponse | string> {
    try {
      const result = await authClient.inviteUser(email)
      return result
    } catch (err: any) {
      const msg = err?.message || 'Failed to send invitation'
      // Backend returns "email already has pending invite" when duplicate
      if (msg.toLowerCase().includes('pending invite') || msg.toLowerCase().includes('already has')) {
        return 'This email already has a pending invitation'
      }
      return msg
    }
  }

  function $reset() {
    invites.value = []
    pagination.value = { page: 1, page_size: 10, total_items: 0 }
    loading.value = false
    error.value = ''
  }

  return {
    invites,
    pagination,
    loading,
    error,
    fetchInvites,
    cancelInvite,
    inviteUser,
    $reset,
  }
})
