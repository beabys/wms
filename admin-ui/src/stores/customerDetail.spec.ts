import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useCustomerDetailStore } from './customerDetail'

const mockCustomer = {
  id: '1',
  company_name: 'Acme Corp',
  email: 'acme@test.com',
  phone: '123456789',
  vat_number: 'VAT123',
  address: '123 Main St',
  city: 'Springfield',
  postal_code: '12345',
  country: 'US',
  status: 'pending',
  created_at: '2024-01-01',
}

describe('useCustomerDetailStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.restoreAllMocks()
  })

  it('starts with default state', () => {
    const store = useCustomerDetailStore()
    expect(store.customer).toBeNull()
    expect(store.auditLogs).toEqual([])
    expect(store.loading).toBe(false)
    expect(store.error).toBe('')
  })

  it('fetchCustomer populates customer', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true, data: mockCustomer }),
    }))

    const store = useCustomerDetailStore()
    await store.fetchCustomer('1')
    expect(store.customer).toEqual(mockCustomer)
    expect(store.loading).toBe(false)
    expect(store.error).toBe('')
  })

  it('fetchCustomer sets error on failure', async () => {
    vi.stubGlobal('fetch', vi.fn().mockRejectedValue(new Error('Network error')))

    const store = useCustomerDetailStore()
    await store.fetchCustomer('1')
    expect(store.customer).toBeNull()
    expect(store.error).toBe('Network error')
    expect(store.loading).toBe(false)
  })

  it('updateCustomer updates customer', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true, data: { ...mockCustomer, company_name: 'Updated Corp' } }),
    }))

    const store = useCustomerDetailStore()
    store.customer = { ...mockCustomer }
    await store.updateCustomer('1', { company_name: 'Updated Corp' })
    expect(store.customer?.company_name).toBe('Updated Corp')
  })

  it('approveCustomer updates status', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true, data: { ...mockCustomer, status: 'active' } }),
    }))

    const store = useCustomerDetailStore()
    store.customer = { ...mockCustomer }
    await store.approveCustomer('1')
    expect(store.customer?.status).toBe('active')
  })

  it('rejectCustomer updates status', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true, data: { ...mockCustomer, status: 'rejected' } }),
    }))

    const store = useCustomerDetailStore()
    store.customer = { ...mockCustomer }
    await store.rejectCustomer('1')
    expect(store.customer?.status).toBe('rejected')
  })

  it('suspendCustomer sets status to suspended', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true, data: null }),
    }))

    const store = useCustomerDetailStore()
    store.customer = { ...mockCustomer, status: 'active' }
    await store.suspendCustomer('1')
    expect(store.customer?.status).toBe('suspended')
  })

  it('restoreCustomer updates status', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true, data: { ...mockCustomer, status: 'active' } }),
    }))

    const store = useCustomerDetailStore()
    store.customer = { ...mockCustomer, status: 'suspended' }
    await store.restoreCustomer('1')
    expect(store.customer?.status).toBe('active')
  })

  it('fetchAuditLogs populates audit logs', async () => {
    const mockResponse = {
      entries: [
        { id: '1', customer_id: '1', action: 'approved', performed_by: 'admin@test.com', details: '', created_at: 1700000000 },
      ],
      pagination: { page: 1, page_size: 10, total_items: 1 },
    }
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true, data: mockResponse }),
    }))

    const store = useCustomerDetailStore()
    await store.fetchAuditLogs('1')
    expect(store.auditLogs).toEqual(mockResponse.entries)
    expect(store.auditPagination).toEqual(mockResponse.pagination)
  })

  it('fetchAuditLogs sets empty on failure', async () => {
    vi.stubGlobal('fetch', vi.fn().mockRejectedValue(new Error('Network error')))

    const store = useCustomerDetailStore()
    await store.fetchAuditLogs('1')
    expect(store.auditLogs).toEqual([])
  })

  it('$reset restores default state', () => {
    const store = useCustomerDetailStore()
    store.customer = mockCustomer as any
    store.auditLogs = [{ id: '1', customer_id: '1', action: 'test', performed_by: 'a', details: '', created_at: 0 }]
    store.loading = true
    store.error = 'some error'
    store.$reset()
    expect(store.customer).toBeNull()
    expect(store.auditLogs).toEqual([])
    expect(store.loading).toBe(false)
    expect(store.error).toBe('')
  })
})
