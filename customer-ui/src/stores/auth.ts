import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { UserResponse } from '@/api/types'

const TOKEN_KEY = 'wms_customer_access_token'
const REFRESH_KEY = 'wms_customer_refresh_token'

export const useAuthStore = defineStore('auth', () => {
  const accessToken = ref<string | null>(localStorage.getItem(TOKEN_KEY))
  const refreshToken = ref<string | null>(localStorage.getItem(REFRESH_KEY))
  const user = ref<UserResponse | null>(null)
  const loading = ref(false)

  const isAuthenticated = computed(() => !!accessToken.value)

  function setTokens(access: string, refresh: string) {
    accessToken.value = access
    refreshToken.value = refresh
    localStorage.setItem(TOKEN_KEY, access)
    localStorage.setItem(REFRESH_KEY, refresh)
  }

  function clearTokens() {
    accessToken.value = null
    refreshToken.value = null
    user.value = null
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(REFRESH_KEY)
  }

  function setUser(u: UserResponse | null) {
    user.value = u
  }

  function $reset() {
    clearTokens()
  }

  return {
    accessToken,
    refreshToken,
    user,
    loading,
    isAuthenticated,
    setTokens,
    clearTokens,
    setUser,
    $reset,
  }
})
