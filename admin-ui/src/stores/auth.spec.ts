import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from './auth'

describe('useAuthStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
  })

  it('starts with no tokens and not authenticated', () => {
    const store = useAuthStore()
    expect(store.accessToken).toBeNull()
    expect(store.refreshToken).toBeNull()
    expect(store.isAuthenticated).toBe(false)
    expect(store.user).toBeNull()
  })

  it('setTokens stores tokens in state and localStorage', () => {
    const store = useAuthStore()
    store.setTokens('access-123', 'refresh-456')
    expect(store.accessToken).toBe('access-123')
    expect(store.refreshToken).toBe('refresh-456')
    expect(store.isAuthenticated).toBe(true)
    expect(localStorage.getItem('wms_admin_access_token')).toBe('access-123')
    expect(localStorage.getItem('wms_admin_refresh_token')).toBe('refresh-456')
  })

  it('clearTokens removes tokens from state and localStorage', () => {
    const store = useAuthStore()
    store.setTokens('access', 'refresh')
    store.setUser({ id: '1', email: 'a@b.com', name: 'A', role: 'admin', created_at: '', updated_at: '' })
    store.clearTokens()
    expect(store.accessToken).toBeNull()
    expect(store.refreshToken).toBeNull()
    expect(store.user).toBeNull()
    expect(localStorage.getItem('wms_admin_access_token')).toBeNull()
  })

  it('setUser updates user state', () => {
    const store = useAuthStore()
    const user = { id: '1', email: 'admin@test.com', name: 'Admin', role: 'admin', created_at: '2024-01-01', updated_at: '2024-01-01' }
    store.setUser(user)
    expect(store.user).toEqual(user)
  })

  it('$reset clears all state', () => {
    const store = useAuthStore()
    store.setTokens('access', 'refresh')
    store.setUser({ id: '1', email: 'a@b.com', name: 'A', role: 'admin', created_at: '', updated_at: '' })
    store.$reset()
    expect(store.accessToken).toBeNull()
    expect(store.refreshToken).toBeNull()
    expect(store.user).toBeNull()
    expect(store.isAuthenticated).toBe(false)
  })
})
