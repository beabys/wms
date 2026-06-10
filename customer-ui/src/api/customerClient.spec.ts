import { describe, it, expect, beforeEach, vi } from 'vitest'
import { customerClient } from './customerClient'

describe('customerClient', () => {
  beforeEach(() => { vi.restoreAllMocks() })

  it('register calls POST /api/v1/customers/register', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true, data: { id: '1', company_name: 'Acme', email: 'a@b.com', status: 'pending' } }),
    }))
    const result = await customerClient.register({ token: 't', company_name: 'Acme', email: 'a@b.com', password: 'pass1234' })
    expect(result.id).toBe('1')
    expect(fetch).toHaveBeenCalledWith('/api/v1/customers/register', expect.objectContaining({ method: 'POST' }))
    vi.restoreAllMocks()
  })

  it('getCustomer calls GET /api/v1/customers/:id', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true, data: { id: '1', company_name: 'Acme', email: 'a@b.com', phone: '', vat_number: '', status: 'active', created_at: '' } }),
    }))
    const result = await customerClient.getCustomer('1')
    expect(result.id).toBe('1')
    vi.restoreAllMocks()
  })

  it('getMyCustomer calls GET /api/v1/customers/me', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true, data: { id: '2', company_name: 'MyCo', email: 'me@co.com', phone: '', vat_number: '', status: 'active', created_at: '' } }),
    }))
    const result = await customerClient.getMyCustomer()
    expect(result.id).toBe('2')
    expect(result.company_name).toBe('MyCo')
    expect(fetch).toHaveBeenCalledWith('/api/v1/customers/me', expect.any(Object))
    vi.restoreAllMocks()
  })
})
