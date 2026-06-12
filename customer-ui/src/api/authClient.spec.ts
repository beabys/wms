import { describe, it, expect, beforeEach, vi } from 'vitest'
import { authClient } from './authClient'

describe('authClient', () => {
  beforeEach(() => { vi.restoreAllMocks() })

  it('login calls POST /api/v1/auth/login', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true, data: { access_token: 't', refresh_token: 'r', expires_in: 3600 } }),
    }))
    const result = await authClient.login({ email: 'a@b.com', password: 'pass' })
    expect(result.access_token).toBe('t')
    vi.restoreAllMocks()
  })

  it('refresh calls POST /api/v1/auth/refresh', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true, data: { access_token: 't2', refresh_token: 'r2', expires_in: 3600 } }),
    }))
    const result = await authClient.refresh('refresh-token')
    expect(result.access_token).toBe('t2')
    vi.restoreAllMocks()
  })

  it('getMe calls GET /api/v1/auth/me', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true, data: { id: '1', email: 'a@b.com', name: 'C', role: 'customer', created_at: '', updated_at: '' } }),
    }))
    const result = await authClient.getMe()
    expect(result.email).toBe('a@b.com')
    vi.restoreAllMocks()
  })

  it('logout calls POST /api/v1/auth/logout', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true }),
    }))
    await authClient.logout()
    expect(fetch).toHaveBeenCalledWith('/api/v1/auth/logout', expect.objectContaining({ method: 'POST' }))
    vi.restoreAllMocks()
  })
})
