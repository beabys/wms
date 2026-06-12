import { computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { authClient } from '@/api/authClient'
import { configureApi, request } from '@/composables/useApi'
import type { LoginRequest, LoginResponse } from '@/api/types'

export function useAuth() {
  const store = useAuthStore()

  // Configure API layer with auth store integration
  configureApi({
    getAccessToken: () => store.accessToken,
    getRefreshToken: () => store.refreshToken,
    onRefresh: async () => {
      try {
        const res = await authClient.refresh(store.refreshToken!)
        store.setTokens(res.access_token, res.refresh_token)
        return res.access_token
      } catch {
        store.$reset()
        return null
      }
    },
    onLogout: () => {
      store.$reset()
    },
  })

  const isAuthenticated = computed(() => store.isAuthenticated)
  const user = computed(() => store.user)
  const loading = computed(() => store.loading)

  async function login(email: string, password: string): Promise<void> {
    store.loading = true
    try {
      const data: LoginResponse = await authClient.login({ email, password } as LoginRequest)
      store.setTokens(data.access_token, data.refresh_token)
      // Fetch user profile after login
      const me = await authClient.getMe()
      store.setUser(me)
    } catch (err) {
      store.$reset()
      throw err
    } finally {
      store.loading = false
    }
  }

  async function logout(): Promise<void> {
    try {
      await authClient.logout()
    } catch {
      // Ignore logout API errors — still clear local state
    }
    store.$reset()
  }

  async function checkAuth(): Promise<boolean> {
    if (!store.accessToken) return false
    store.loading = true
    try {
      const me = await authClient.getMe()
      store.setUser(me)
      return true
    } catch {
      store.$reset()
      return false
    } finally {
      store.loading = false
    }
  }

  return {
    isAuthenticated,
    user,
    loading,
    login,
    logout,
    checkAuth,
    request,
  }
}
