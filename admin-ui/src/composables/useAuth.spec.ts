import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuth } from './useAuth'
import { useAuthStore } from '@/stores/auth'

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
    const mockLoginResponse = {
      access_token: 'access-123',
      refresh_token: 'refresh-123',
      expires_in: 3600,
    }
    const mockUser = { id: '1', email: 'admin@test.com', name: 'Admin', role: 'admin', created_at: '', updated_at: '' }

    const mockFetch = vi.fn()
      .mockResolvedValueOnce({
        status: 200,
        json: async () => ({ success: true, data: mockLoginResponse }),
      })
      .mockResolvedValueOnce({
        status: 200,
        json: async () => ({ success: true, data: mockUser }),
      })
    vi.stubGlobal('fetch', mockFetch)

    const auth = useAuth()
    await auth.login('admin@test.com', 'password123')

    const store = useAuthStore()
    expect(store.accessToken).toBe('access-123')
    expect(store.refreshToken).toBe('refresh-123')
    expect(store.user).toEqual(mockUser)
  })

  it('login throws error on failure and clears store', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 401,
      json: async () => ({ success: false, error: 'Invalid credentials' }),
    }))

    const auth = useAuth()
    await expect(auth.login('admin@test.com', 'wrong')).rejects.toThrow('Invalid credentials')

    const store = useAuthStore()
    expect(store.accessToken).toBeNull()
  })

  it('logout calls API and clears all stored tokens', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true, data: {} }),
    }))

    const store = useAuthStore()
    store.setTokens('access', 'refresh')
    store.setUser({ id: '1', email: 'a@b.com', name: 'A', role: 'admin', created_at: '', updated_at: '' })

    const auth = useAuth()
    await auth.logout()

    expect(fetch).toHaveBeenCalledWith('/api/v1/auth/logout', expect.objectContaining({ method: 'POST' }))
    expect(store.accessToken).toBeNull()
    expect(store.refreshToken).toBeNull()
    expect(store.user).toBeNull()
  })

  it('checkAuth returns false when no token', async () => {
    const auth = useAuth()
    const result = await auth.checkAuth()
    expect(result).toBe(false)
  })

  it('checkAuth returns true when token is valid', async () => {
    const store = useAuthStore()
    store.setTokens('valid-token', 'refresh')
    const mockUser = { id: '1', email: 'admin@test.com', name: 'Admin', role: 'admin', created_at: '', updated_at: '' }

    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true, data: mockUser }),
    }))

    const auth = useAuth()
    const result = await auth.checkAuth()
    expect(result).toBe(true)
    expect(store.user).toEqual(mockUser)
  })

  it('checkAuth returns false when token is invalid', async () => {
    const store = useAuthStore()
    store.setTokens('invalid-token', 'refresh')

    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 401,
      json: async () => ({ success: false, error: 'Unauthorized' }),
    }))

    const auth = useAuth()
    const result = await auth.checkAuth()
    expect(result).toBe(false)
    expect(store.accessToken).toBeNull()
  })
})
