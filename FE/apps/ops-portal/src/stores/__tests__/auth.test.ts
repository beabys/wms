import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'

// ── localStorage polyfill (node environment) ────────
const { localStorageMock } = vi.hoisted(() => {
  let store: Record<string, string> = {}
  return {
    localStorageMock: {
      getItem: vi.fn((key: string) => store[key] ?? null),
      setItem: vi.fn((key: string, value: string) => {
        store[key] = value
      }),
      removeItem: vi.fn((key: string) => {
        delete store[key]
      }),
      clear: vi.fn(() => {
        store = {}
      }),
      get length() {
        return Object.keys(store).length
      },
      key: vi.fn((i: number) => Object.keys(store)[i] ?? null),
    },
  }
})

Object.defineProperty(globalThis, 'localStorage', {
  value: localStorageMock,
  writable: true,
  configurable: true,
})

// ── Mock AuthClient ─────────────────────────────────
const { mockLogin, mockSetToken } = vi.hoisted(() => ({
  mockLogin: vi.fn(),
  mockSetToken: vi.fn(),
}))

vi.mock('@wms/api-client', () => ({
  AuthClient: vi.fn().mockImplementation(() => ({
    login: mockLogin,
    setToken: mockSetToken,
  })),
}))

import { useAuthStore } from '../auth'

describe('useAuthStore', () => {
  let store: ReturnType<typeof useAuthStore>

  beforeEach(() => {
    localStorageMock.clear()
    vi.clearAllMocks()
    setActivePinia(createPinia())
    store = useAuthStore()
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('has correct initial state — no token, not authenticated', () => {
    expect(store.token).toBeNull()
    expect(store.user).toBeNull()
    expect(store.isAuthenticated).toBe(false)
  })

  it('login sets token, user, and isAuthenticated', async () => {
    const loginResult = {
      access_token: 'ops-token-456',
      refresh_token: 'ops-refresh',
      user: { id: 'u2', email: 'ops@example.com', role: 'ops' as const },
    }
    mockLogin.mockResolvedValueOnce(loginResult)

    await store.login('ops@example.com', 'ops-secret')

    expect(mockLogin).toHaveBeenCalledWith('ops@example.com', 'ops-secret')
    expect(store.token).toBe('ops-token-456')
    expect(store.user).toEqual({
      id: 'u2',
      name: 'ops@example.com',
      email: 'ops@example.com',
      role: 'ops',
    })
    expect(store.isAuthenticated).toBe(true)
    expect(mockSetToken).toHaveBeenCalledWith('ops-token-456')
    expect(localStorageMock.setItem).toHaveBeenCalledWith(
      'wms_ops_token',
      'ops-token-456',
    )
    expect(localStorageMock.setItem).toHaveBeenCalledWith(
      'wms_ops_user',
      JSON.stringify({
        id: 'u2',
        name: 'ops@example.com',
        email: 'ops@example.com',
        role: 'ops',
      }),
    )
  })

  it('logout clears token, user, and isAuthenticated', async () => {
    // Arrange: login first
    mockLogin.mockResolvedValueOnce({
      access_token: 'ops-token-456',
      refresh_token: 'ops-refresh',
      user: { id: 'u2', email: 'ops@example.com', role: 'ops' as const },
    })
    await store.login('ops@example.com', 'ops-secret')

    // Act
    store.logout()

    // Assert
    expect(store.token).toBeNull()
    expect(store.user).toBeNull()
    expect(store.isAuthenticated).toBe(false)
    expect(mockSetToken).toHaveBeenCalledWith(null)
    expect(localStorageMock.removeItem).toHaveBeenCalledWith('wms_ops_token')
    expect(localStorageMock.removeItem).toHaveBeenCalledWith('wms_ops_user')
  })
})
