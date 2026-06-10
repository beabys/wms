import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuth } from './useAuth'
import { useAuthStore } from '@/stores/auth'
import { useCustomerStore } from '@/stores/customer'

describe('useAuth', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
    vi.restoreAllMocks()
  })

  it('returns not authenticated initially', () => {
    const auth = useAuth()
    expect(auth.isAuthenticated.value).toBe(false)
    expect(auth.user.value).toBeNull()
  })

  it('login sets tokens and user on success', async () => {
    const mockLogin = { access_token: 'access-123', refresh_token: 'refresh-123', expires_in: 3600 }
    const mockUser = { id: '1', email: 'cust@test.com', name: 'Customer', role: 'customer', created_at: '', updated_at: '' }
    const mockCustomer = { id: '1', company_name: 'Acme', email: 'cust@test.com', phone: '', vat_number: '', status: 'active', created_at: '' }

    vi.stubGlobal('fetch', vi.fn()
      .mockResolvedValueOnce({ status: 200, json: async () => ({ success: true, data: mockLogin }) })
      .mockResolvedValueOnce({ status: 200, json: async () => ({ success: true, data: mockUser }) })
      .mockResolvedValueOnce({ status: 200, json: async () => ({ success: true, data: mockCustomer }) }))

    const auth = useAuth()
    await auth.login('cust@test.com', 'password123')

    const store = useAuthStore()
    expect(store.accessToken).toBe('access-123')
    expect(store.refreshToken).toBe('refresh-123')
    expect(store.user).toEqual(mockUser)

    const customerStore = useCustomerStore()
    expect(customerStore.customer).toEqual(mockCustomer)
    vi.restoreAllMocks()
  })

  it('login throws on failure', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 401,
      json: async () => ({ success: false, error: 'Invalid credentials' }),
    }))
    const auth = useAuth()
    await expect(auth.login('cust@test.com', 'wrong')).rejects.toThrow('Invalid credentials')
    vi.restoreAllMocks()
  })

  it('register submits to customer endpoint', async () => {
    const registerData = { token: 'invite-123', company_name: 'Acme', email: 'a@b.com', password: 'pass1234', phone: '', vat_number: '' }
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true, data: { id: '1', company_name: 'Acme', email: 'a@b.com', status: 'pending' } }),
    }))
    const auth = useAuth()
    await auth.register(registerData)
    expect(fetch).toHaveBeenCalledWith('/api/v1/customers/register', expect.objectContaining({ method: 'POST' }))
    vi.restoreAllMocks()
  })

  it('logout calls API and clears state', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true }),
    }))
    const store = useAuthStore()
    store.setTokens('access', 'refresh')
    const auth = useAuth()
    await auth.logout()
    expect(fetch).toHaveBeenCalledWith('/api/v1/auth/logout', expect.objectContaining({ method: 'POST' }))
    expect(store.accessToken).toBeNull()
    vi.restoreAllMocks()
  })

  it('logout clears state even when API fails', async () => {
    vi.stubGlobal('fetch', vi.fn().mockRejectedValue(new Error('Network error')))
    const store = useAuthStore()
    store.setTokens('access', 'refresh')
    const auth = useAuth()
    await auth.logout()
    expect(store.accessToken).toBeNull()
    vi.restoreAllMocks()
  })

  it('checkAuth returns false when no token', async () => {
    const auth = useAuth()
    const result = await auth.checkAuth()
    expect(result).toBe(false)
  })

  it('checkAuth returns true with valid token', async () => {
    const store = useAuthStore()
    store.setTokens('valid', 'refresh')
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true, data: { id: '1', email: 'a@b.com', name: 'C', role: 'customer', created_at: '', updated_at: '' } }),
    }))
    const auth = useAuth()
    const result = await auth.checkAuth()
    expect(result).toBe(true)
    vi.restoreAllMocks()
  })
})
