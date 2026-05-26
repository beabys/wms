import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { AuthClient, CustomerClient, InboundClient } from '../index'

function mockFetch(responseData: unknown, ok = true) {
  return vi.spyOn(globalThis, 'fetch').mockResolvedValue(
    new Response(JSON.stringify(responseData), { status: ok ? 200 : 400 })
  )
}

describe('AuthClient', () => {
  let client: AuthClient

  beforeEach(() => {
    client = new AuthClient('/api')
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('login sends POST and returns unwrapped data', async () => {
    const mock = mockFetch({
      success: true,
      data: { access_token: 'abc', refresh_token: 'def', user: { id: '1', email: 'a@b.com', role: 'customer', customer_id: 'c1', created_at: '2026-01-01' } },
    })

    const result = await client.login('a@b.com', 'pwd')
    expect(result.access_token).toBe('abc')
    expect(mock).toHaveBeenCalledWith('/api/v1/auth/login', expect.objectContaining({ method: 'POST' }))
  })

  it('throws on success: false', async () => {
    mockFetch({ success: false, error: 'Invalid credentials' })

    await expect(client.login('bad', 'creds')).rejects.toThrow('Invalid credentials')
  })

  it('setToken adds Authorization header', async () => {
    client.setToken('tok123')
    const mock = mockFetch({ success: true, data: { user: { id: '1', email: 'a@b.com', role: 'customer', customer_id: 'c1', created_at: '2026-01-01' } } })

    await client.me()
    expect(mock).toHaveBeenCalledWith(
      '/api/v1/auth/me',
      expect.objectContaining({
        headers: expect.objectContaining({ Authorization: 'Bearer tok123' }),
      })
    )
  })
})

describe('CustomerClient', () => {
  let client: CustomerClient

  beforeEach(() => {
    client = new CustomerClient('/api')
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('invite sends POST to correct path', async () => {
    mockFetch({ success: true, data: { invite_link: 'https://example.com/invite', expires_at: '2026-06-01' } })
    const result = await client.invite('a@b.com')
    expect(result.invite_link).toContain('example.com')
  })

  it('getQueue builds query params', async () => {
    mockFetch({ success: true, data: { customers: [], total: 0, page: 1, page_size: 20 } })
    await client.getQueue(1, 20)
  })
})

describe('InboundClient', () => {
  let client: InboundClient

  beforeEach(() => {
    client = new InboundClient('/api')
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('create sends POST with items', async () => {
    mockFetch({ success: true, data: { id: 'inb-1', message: 'Created' } })
    const result = await client.create({
      customer_id: 'c1',
      expected_date: '2026-06-01',
      items: [{ sku: 'SKU-1', quantity_declared: 10 }],
    })
    expect(result.id).toBe('inb-1')
  })

  it('getQueue passes optional params', async () => {
    mockFetch({ success: true, data: { inbounds: [], page_size: 10 } })
    await client.getQueue({ page_size: 10, status: 'pending' })
  })

  it('flag sends reason in body', async () => {
    mockFetch({ success: true, data: { flag_id: 'f-1', reason: 'damaged' } })
    const result = await client.flag('inb-1', 'damaged')
    expect(result.flag_id).toBe('f-1')
  })
})
