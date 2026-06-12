import { describe, it, expect, beforeEach, vi } from 'vitest'
import { customerClient } from './customerClient'

describe('customerClient', () => {
  beforeEach(() => {
    vi.restoreAllMocks()
  })

  const mockCustomer = { id: '1', company_name: 'Acme', email: 'acme@test.com', phone: '', vat_number: '', status: 'active', created_at: '2024-01-01' }

  it('listCustomers calls GET /api/v1/customers with no params', async () => {
    const mockResponse = { customers: [mockCustomer], pagination: { page: 1, page_size: 10, total_items: 1 } }
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true, data: mockResponse }),
    }))

    const result = await customerClient.listCustomers()
    expect(result).toEqual(mockResponse)
    expect(fetch).toHaveBeenCalledWith('/api/v1/customers', expect.anything())
  })

  it('listCustomers passes query params', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true, data: { customers: [], pagination: { page: 1, page_size: 20, total_items: 0 } } }),
    }))

    await customerClient.listCustomers({ status: 'pending', page: 2, page_size: 20 })
    expect(fetch).toHaveBeenCalledWith('/api/v1/customers?status=pending&page=2&page_size=20', expect.anything())
  })

  it('getCustomer calls GET /api/v1/customers/:id', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true, data: mockCustomer }),
    }))

    const result = await customerClient.getCustomer('1')
    expect(result).toEqual(mockCustomer)
    expect(fetch).toHaveBeenCalledWith('/api/v1/customers/1', expect.anything())
  })

  it('approveCustomer calls POST /api/v1/customers/:id/approve', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true, data: { id: '1', status: 'active' } }),
    }))

    const result = await customerClient.approveCustomer('1')
    expect(result).toEqual({ id: '1', status: 'active' })
    expect(fetch).toHaveBeenCalledWith('/api/v1/customers/1/approve', expect.objectContaining({ method: 'POST' }))
  })

  it('rejectCustomer calls POST /api/v1/customers/:id/reject', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true, data: { id: '1', status: 'rejected' } }),
    }))

    const result = await customerClient.rejectCustomer('1')
    expect(result).toEqual({ id: '1', status: 'rejected' })
    expect(fetch).toHaveBeenCalledWith('/api/v1/customers/1/reject', expect.objectContaining({ method: 'POST' }))
  })

  it('suspendCustomer calls POST /api/v1/customers/:id/suspend', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true, data: null }),
    }))

    await customerClient.suspendCustomer('1')
    expect(fetch).toHaveBeenCalledWith('/api/v1/customers/1/suspend', expect.objectContaining({ method: 'POST' }))
  })

  it('restoreCustomer calls POST /api/v1/customers/:id/restore', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true, data: { id: '1', status: 'active' } }),
    }))

    const result = await customerClient.restoreCustomer('1')
    expect(result).toEqual({ id: '1', status: 'active' })
    expect(fetch).toHaveBeenCalledWith('/api/v1/customers/1/restore', expect.objectContaining({ method: 'POST' }))
  })

  it('listAuditLogs calls GET /api/v1/customers/:id/audit', async () => {
    const mockResponse = { entries: [], pagination: { page: 1, page_size: 10, total_items: 0 } }
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true, data: mockResponse }),
    }))

    const result = await customerClient.listAuditLogs('1')
    expect(result).toEqual(mockResponse)
    expect(fetch).toHaveBeenCalledWith('/api/v1/customers/1/audit', expect.anything())
  })

  it('listAuditLogs passes query params', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true, data: { entries: [], pagination: { page: 2, page_size: 20, total_items: 0 } } }),
    }))

    await customerClient.listAuditLogs('1', { page: 2, page_size: 20 })
    expect(fetch).toHaveBeenCalledWith('/api/v1/customers/1/audit?page=2&page_size=20', expect.anything())
  })

  it('updateCustomer calls PUT /api/v1/customers/:id', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true, data: { id: '1', company_name: 'Updated' } }),
    }))

    const result = await customerClient.updateCustomer('1', { company_name: 'Updated' })
    expect(result).toEqual({ id: '1', company_name: 'Updated' })
    expect(fetch).toHaveBeenCalledWith('/api/v1/customers/1', expect.objectContaining({
      method: 'PUT',
      body: JSON.stringify({ company_name: 'Updated' }),
    }))
  })
})
