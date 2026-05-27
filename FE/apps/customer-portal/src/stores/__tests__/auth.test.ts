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
      access_token: 'test-token-123',
      refresh_token: 'refresh-xyz',
      user: { id: 'u1', email: 'alice@example.com', role: 'customer' as const },
    }
    mockLogin.mockResolvedValueOnce(loginResult)

    await store.login('alice@example.com', 'secret')

    expect(mockLogin).toHaveBeenCalledWith('alice@example.com', 'secret')
    expect(store.token).toBe('test-token-123')
    expect(store.user).toEqual({
      id: 'u1',
      name: 'alice@example.com',
      email: 'alice@example.com',
      role: 'customer',
    })
    expect(store.isAuthenticated).toBe(true)
    expect(mockSetToken).toHaveBeenCalledWith('test-token-123')
    expect(localStorageMock.setItem).toHaveBeenCalledWith(
      'wms_token',
      'test-token-123',
    )
    expect(localStorageMock.setItem).toHaveBeenCalledWith(
      'wms_user',
      JSON.stringify({
        id: 'u1',
        name: 'alice@example.com',
        email: 'alice@example.com',
        role: 'customer',
      }),
    )
  })

  it('logout clears token, user, and isAuthenticated', async () => {
    // Arrange: login first
    mockLogin.mockResolvedValueOnce({
      access_token: 'test-token-123',
      refresh_token: 'refresh-xyz',
      user: { id: 'u1', email: 'alice@example.com', role: 'customer' as const },
    })
    await store.login('alice@example.com', 'secret')

    // Act
    store.logout()

    // Assert
    expect(store.token).toBeNull()
    expect(store.user).toBeNull()
    expect(store.isAuthenticated).toBe(false)
    expect(mockSetToken).toHaveBeenCalledWith(null)
    expect(localStorageMock.removeItem).toHaveBeenCalledWith('wms_token')
    expect(localStorageMock.removeItem).toHaveBeenCalledWith('wms_user')
  })
})
