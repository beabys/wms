import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from './auth'

describe('useAuthStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
  })

  it('starts with no tokens', () => {
    const store = useAuthStore()
    expect(store.accessToken).toBeNull()
    expect(store.isAuthenticated).toBe(false)
  })

  it('setTokens stores in localStorage with customer key', () => {
    const store = useAuthStore()
    store.setTokens('access-cust', 'refresh-cust')
    expect(localStorage.getItem('wms_customer_access_token')).toBe('access-cust')
    expect(localStorage.getItem('wms_customer_refresh_token')).toBe('refresh-cust')
  })

  it('clearTokens removes everything', () => {
    const store = useAuthStore()
    store.setTokens('a', 'r')
    store.clearTokens()
    expect(store.accessToken).toBeNull()
    expect(localStorage.getItem('wms_customer_access_token')).toBeNull()
  })

  it('$reset clears all state', () => {
    const store = useAuthStore()
    store.setTokens('a', 'r')
    store.$reset()
    expect(store.isAuthenticated).toBe(false)
  })
})
