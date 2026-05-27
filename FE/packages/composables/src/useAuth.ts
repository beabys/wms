import { ref, computed } from 'vue'

export interface AuthUser {
  id: string
  name: string
  email: string
  role: 'customer' | 'ops'
}

const token = ref<string | null>(null)
const user = ref<AuthUser | null>(null)

export function useAuth() {
  const isAuthenticated = computed(() => !!token.value && !!user.value)

  async function login(email: string, password: string): Promise<void> {
    // Stub — replace with real API call
    await new Promise(r => setTimeout(r, 500))
    token.value = 'jwt-' + Date.now()
    user.value = { id: '1', name: 'User', email, role: 'customer' }
  }

  function logout(): void {
    token.value = null
    user.value = null
  }

  return { token, user, isAuthenticated, login, logout }
}
