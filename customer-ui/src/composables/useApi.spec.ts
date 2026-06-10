import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest'
import { configureApi, request, ApiError } from './useApi'

describe('useApi', () => {
  afterEach(() => { vi.restoreAllMocks() })

  beforeEach(() => {
    configureApi({
      getAccessToken: () => null,
      getRefreshToken: () => null,
      onRefresh: async () => null,
      onLogout: () => {},
    })
  })

  it('makes a successful request and returns data', async () => {
    const mockData = { id: '1', name: 'Test' }
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true, data: mockData }),
    }))
    const result = await request<{ id: string; name: string }>('/api/test')
    expect(result).toEqual(mockData)
  })

  it('throws ApiError on unsuccessful response', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 400,
      json: async () => ({ success: false, error: 'Bad request' }),
    }))
    await expect(request('/api/test')).rejects.toThrow(ApiError)
  })

  it('attaches Authorization header when token exists', async () => {
    configureApi({
      getAccessToken: () => 'test-token',
      getRefreshToken: () => null,
      onRefresh: async () => null,
      onLogout: () => {},
    })
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true, data: {} }),
    }))
    await request('/api/test')
    expect(fetch).toHaveBeenCalledWith('/api/test', expect.objectContaining({
      headers: expect.objectContaining({ Authorization: 'Bearer test-token' }),
    }))
  })

  it('calls onRefresh on 401 and retries', async () => {
    const refreshFn = vi.fn().mockResolvedValue('new-token')
    const logoutFn = vi.fn()
    configureApi({
      getAccessToken: () => 'expired-token',
      getRefreshToken: () => 'refresh-token',
      onRefresh: refreshFn,
      onLogout: logoutFn,
    })

    let callCount = 0
    vi.stubGlobal('fetch', vi.fn().mockImplementation(() => {
      callCount++
      return callCount === 1
        ? Promise.resolve({ status: 401, json: async () => ({ success: false, error: 'Unauthorized' }) })
        : Promise.resolve({ status: 200, json: async () => ({ success: true, data: { recovered: true } }) })
    }))

    const result = await request<{ recovered: boolean }>('/api/test')
    expect(result).toEqual({ recovered: true })
    expect(refreshFn).toHaveBeenCalledOnce()
    expect(logoutFn).not.toHaveBeenCalled()
  })

  it('calls onLogout when refresh fails on 401', async () => {
    const refreshFn = vi.fn().mockResolvedValue(null)
    const logoutFn = vi.fn()
    configureApi({
      getAccessToken: () => 'expired-token',
      getRefreshToken: () => 'refresh-token',
      onRefresh: refreshFn,
      onLogout: logoutFn,
    })
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 401,
      json: async () => ({ success: false, error: 'Unauthorized' }),
    }))
    await expect(request('/api/test')).rejects.toThrow('Session expired')
    expect(refreshFn).toHaveBeenCalledOnce()
    expect(logoutFn).toHaveBeenCalledOnce()
  })
})
