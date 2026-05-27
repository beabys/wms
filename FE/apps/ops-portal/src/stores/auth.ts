import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { AuthClient } from '@wms/api-client'

const authClient = new AuthClient(import.meta.env.VITE_API_BASE_URL || '')

export interface User {
  id: string
  name: string
  email: string
  role: 'customer' | 'ops'
}

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem('wms_ops_token'))
  const user = ref<User | null>(
    localStorage.getItem('wms_ops_user')
      ? JSON.parse(localStorage.getItem('wms_ops_user')!)
      : null
  )

  const isAuthenticated = computed(() => !!token.value && !!user.value)

  async function login(email: string, password: string): Promise<void> {
    const result = await authClient.login(email, password)
    token.value = result.access_token
    authClient.setToken(result.access_token)
    user.value = {
      id: result.user.id,
      name: result.user.email,
      email: result.user.email,
      role: result.user.role as 'customer' | 'ops',
    }
    localStorage.setItem('wms_ops_token', token.value)
    localStorage.setItem('wms_ops_user', JSON.stringify(user.value))
  }

  function logout(): void {
    token.value = null
    user.value = null
    authClient.setToken(null)
    localStorage.removeItem('wms_ops_token')
    localStorage.removeItem('wms_ops_user')
  }

  return { token, user, isAuthenticated, login, logout }
})
