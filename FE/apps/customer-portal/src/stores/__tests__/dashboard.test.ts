import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useDashboardStore } from '../dashboard'

function response(data: unknown, ok = true): Response {
  return new Response(JSON.stringify(data), {
    status: ok ? 200 : 400,
    headers: { 'Content-Type': 'application/json' },
  })
}

describe('useDashboardStore', () => {
  let store: ReturnType<typeof useDashboardStore>

  beforeEach(() => {
    setActivePinia(createPinia())
    store = useDashboardStore()
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('has correct initial state', () => {
    expect(store.loading).toBe(false)
    expect(store.error).toBeNull()
    expect(store.activeInbounds).toBe(0)
    expect(store.totalInbounds).toBe(0)
    expect(store.stockEntries).toBe(0)
  })

  it('setToken propagates to both API clients (Authorization header sent on fetch)', async () => {
    store.setToken('my-token')

    vi.spyOn(globalThis, 'fetch')
      .mockResolvedValueOnce(response({ success: true, data: { inbounds: [] } }))
      .mockResolvedValueOnce(response({ success: true, data: { entries: [] } }))

    await store.fetchStats('customer-1')

    expect(vi.mocked(fetch)).toHaveBeenCalledWith(
      '/v1/inbounds/customer/customer-1',
      expect.objectContaining({
        headers: expect.objectContaining({ Authorization: 'Bearer my-token' }),
      })
    )
  })

  it('fetchStats populates counts on success', async () => {
    const inboundData = {
      success: true,
      data: {
        inbounds: [
          { id: '1', status: 'submitted' },
          { id: '2', status: 'submitted' },
          { id: '3', status: 'received' },
        ],
      },
    }
    const stockData = {
      success: true,
      data: {
        entries: [
          { sku: 'SKU-1', quantity: 10 },
          { sku: 'SKU-2', quantity: 5 },
        ],
      },
    }

    vi.spyOn(globalThis, 'fetch')
      .mockResolvedValueOnce(response(inboundData))
      .mockResolvedValueOnce(response(stockData))

    await store.fetchStats('customer-1')

    expect(store.loading).toBe(false)
    expect(store.error).toBeNull()
    expect(store.activeInbounds).toBe(2)   // 2 with status 'submitted'
    expect(store.totalInbounds).toBe(3)    // 3 total
    expect(store.stockEntries).toBe(2)     // 2 stock entries
  })

  it('fetchStats sets error on API failure (success: false)', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValueOnce(
      response({ success: false, error: 'Customer not found' })
    )

    await store.fetchStats('unknown')

    expect(store.loading).toBe(false)
    expect(store.error).toBe('Customer not found')
    expect(store.activeInbounds).toBe(0)
    expect(store.totalInbounds).toBe(0)
    expect(store.stockEntries).toBe(0)
  })

  it('fetchStats sets error on network failure', async () => {
    vi.spyOn(globalThis, 'fetch').mockRejectedValueOnce(new Error('Network error'))

    await store.fetchStats('customer-1')

    expect(store.loading).toBe(false)
    expect(store.error).toBe('Network error')
  })

  it('fetchStats handles non-Error throws gracefully', async () => {
    vi.spyOn(globalThis, 'fetch').mockRejectedValueOnce('string error')

    await store.fetchStats('customer-1')

    expect(store.loading).toBe(false)
    expect(store.error).toBe('string error')
  })
})
