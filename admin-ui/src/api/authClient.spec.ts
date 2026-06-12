import { describe, it, expect, beforeEach, vi } from 'vitest'
import { authClient } from './authClient'

describe('authClient', () => {
  beforeEach(() => {
    vi.restoreAllMocks()
  })

  it('login calls POST /api/v1/auth/login with credentials', async () => {
    const mockResponse = { access_token: 'token', refresh_token: 'refresh', expires_in: 3600 }
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true, data: mockResponse }),
    }))

    const result = await authClient.login({ email: 'test@test.com', password: 'password123' })
    expect(result).toEqual(mockResponse)
    expect(fetch).toHaveBeenCalledWith('/api/v1/auth/login', expect.objectContaining({
      method: 'POST',
      body: JSON.stringify({ email: 'test@test.com', password: 'password123' }),
    }))
  })

  it('refresh calls POST /api/v1/auth/refresh', async () => {
    const mockResponse = { access_token: 'new-token', refresh_token: 'new-refresh', expires_in: 3600 }
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true, data: mockResponse }),
    }))

    const result = await authClient.refresh('refresh-token')
    expect(result).toEqual(mockResponse)
    expect(fetch).toHaveBeenCalledWith('/api/v1/auth/refresh', expect.objectContaining({
      method: 'POST',
      body: JSON.stringify({ refresh_token: 'refresh-token' }),
    }))
  })

  it('getMe calls GET /api/v1/auth/me', async () => {
    const mockUser = { id: '1', email: 'a@b.com', name: 'Admin', role: 'admin', created_at: '', updated_at: '' }
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true, data: mockUser }),
    }))

    const result = await authClient.getMe()
    expect(result).toEqual(mockUser)
    expect(fetch).toHaveBeenCalledWith('/api/v1/auth/me', expect.anything())
  })

  it('listUsers calls GET /api/v1/users with params', async () => {
    const mockResponse = { users: [], pagination: { page: 1, page_size: 10, total_items: 0 } }
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true, data: mockResponse }),
    }))

    const result = await authClient.listUsers({ role: 'admin', page: 1, page_size: 10 })
    expect(result).toEqual(mockResponse)
    expect(fetch).toHaveBeenCalledWith('/api/v1/users?role=admin&page=1&page_size=10', expect.anything())
  })

  it('listUsers calls GET /api/v1/users without params', async () => {
    const mockResponse = { users: [], pagination: { page: 1, page_size: 10, total_items: 0 } }
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true, data: mockResponse }),
    }))

    const result = await authClient.listUsers()
    expect(result).toEqual(mockResponse)
    expect(fetch).toHaveBeenCalledWith('/api/v1/users', expect.anything())
  })

  it('inviteUser calls POST /api/v1/auth/invite', async () => {
    const mockResponse = { token: 'invite-token', invite_link: 'http://example.com/invite/token', expires_at: 9999999999 }
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true, data: mockResponse }),
    }))

    const result = await authClient.inviteUser('invite@test.com')
    expect(result).toEqual(mockResponse)
    expect(fetch).toHaveBeenCalledWith('/api/v1/auth/invite', expect.objectContaining({
      method: 'POST',
      body: JSON.stringify({ email: 'invite@test.com' }),
    }))
  })

  it('createUser calls POST /api/v1/users', async () => {
    const mockUser = { id: '2', email: 'new@b.com', name: 'New', role: 'staff', created_at: '', updated_at: '' }
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true, data: mockUser }),
    }))

    const result = await authClient.createUser({ email: 'new@b.com', password: 'pass123', name: 'New', role: 'staff' })
    expect(result).toEqual(mockUser)
    expect(fetch).toHaveBeenCalledWith('/api/v1/users', expect.objectContaining({
      method: 'POST',
      body: expect.stringContaining('"email":"new@b.com"'),
    }))
  })

  it('logout calls POST /api/v1/auth/logout', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true, data: {} }),
    }))

    await authClient.logout()
    expect(fetch).toHaveBeenCalledWith('/api/v1/auth/logout', expect.objectContaining({
      method: 'POST',
    }))
  })

  it('inviteUser calls POST /api/v1/auth/invite with email', async () => {
    const mockResponse = { token: 'invite-token', invite_link: 'http://example.com/invite?token=invite-token', expires_at: 9999999999 }
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true, data: mockResponse }),
    }))

    const result = await authClient.inviteUser('new@company.com')
    expect(result).toEqual(mockResponse)
    expect(fetch).toHaveBeenCalledWith('/api/v1/auth/invite', expect.objectContaining({
      method: 'POST',
      body: JSON.stringify({ email: 'new@company.com' }),
    }))
  })
})
