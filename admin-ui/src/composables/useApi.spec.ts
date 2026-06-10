import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest'
import { configureApi, request, ApiError } from './useApi'

describe('useApi', () => {
  afterEach(() => {
    vi.restoreAllMocks()
  })

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
    await expect(request('/api/test')).rejects.toThrow('Bad request')
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
    const mockFetch = vi.fn().mockImplementation(() => {
      callCount++
      if (callCount === 1) {
        return Promise.resolve({
          status: 401,
          json: async () => ({ success: false, error: 'Unauthorized' }),
        })
      }
      return Promise.resolve({
        status: 200,
        json: async () => ({ success: true, data: { recovered: true } }),
      })
    })
    vi.stubGlobal('fetch', mockFetch)

    const result = await request<{ recovered: boolean }>('/api/test')
    expect(result).toEqual({ recovered: true })
    expect(refreshFn).toHaveBeenCalledOnce()
    expect(logoutFn).not.toHaveBeenCalled()
    expect(callCount).toBe(2)
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

  it('throws ApiError with status on network error envelope', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 500,
      json: async () => ({ success: false, error: 'Server error' }),
    }))

    try {
      await request('/api/test')
    } catch (err) {
      expect(err).toBeInstanceOf(ApiError)
      expect((err as ApiError).status).toBe(500)
    }
  })
})
