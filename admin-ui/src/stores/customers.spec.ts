import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useCustomerStore } from './customers'

const mockCustomer = {
  id: '1',
  company_name: 'Acme Corp',
  email: 'acme@test.com',
  phone: '123456789',
  vat_number: 'VAT123',
  status: 'pending',
  created_at: '2024-01-01',
}

describe('useCustomerStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.restoreAllMocks()
  })

  it('starts with default state', () => {
    const store = useCustomerStore()
    expect(store.customers).toEqual([])
    expect(store.pagination).toEqual({ page: 1, page_size: 10, total_items: 0 })
    expect(store.loading).toBe(false)
    expect(store.error).toBe('')
  })

  it('fetchCustomers populates customers and pagination', async () => {
    const mockResponse = {
      customers: [mockCustomer],
      pagination: { page: 1, page_size: 10, total_items: 1 },
    }
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true, data: mockResponse }),
    }))

    const store = useCustomerStore()
    await store.fetchCustomers()
    expect(store.customers).toEqual([mockCustomer])
    expect(store.pagination).toEqual({ page: 1, page_size: 10, total_items: 1 })
    expect(store.loading).toBe(false)
    expect(store.error).toBe('')
  })

  it('fetchCustomers sets error on failure', async () => {
    vi.stubGlobal('fetch', vi.fn().mockRejectedValue(new Error('Network error')))

    const store = useCustomerStore()
    await store.fetchCustomers()
    expect(store.customers).toEqual([])
    expect(store.error).toBe('Network error')
    expect(store.loading).toBe(false)
  })

  it('approveCustomer updates the customer status', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true, data: { ...mockCustomer, status: 'active' } }),
    }))

    const store = useCustomerStore()
    store.customers = [mockCustomer]
    await store.approveCustomer('1')
    expect(store.customers[0].status).toBe('active')
  })

  it('rejectCustomer updates the customer status', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true, data: { ...mockCustomer, status: 'rejected' } }),
    }))

    const store = useCustomerStore()
    store.customers = [mockCustomer]
    await store.rejectCustomer('1')
    expect(store.customers[0].status).toBe('rejected')
  })

  it('suspendCustomer sets status to suspended', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true, data: null }),
    }))

    const store = useCustomerStore()
    store.customers = [mockCustomer]
    await store.suspendCustomer('1')
    expect(store.customers[0].status).toBe('suspended')
  })

  it('restoreCustomer updates the customer status', async () => {
    const suspended = { ...mockCustomer, status: 'suspended' }
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true, data: { ...suspended, status: 'active' } }),
    }))

    const store = useCustomerStore()
    store.customers = [suspended]
    await store.restoreCustomer('1')
    expect(store.customers[0].status).toBe('active')
  })

  it('$reset restores default state', () => {
    const store = useCustomerStore()
    store.customers = [mockCustomer]
    store.pagination = { page: 2, page_size: 20, total_items: 5 }
    store.loading = true
    store.error = 'some error'
    store.$reset()
    expect(store.customers).toEqual([])
    expect(store.pagination).toEqual({ page: 1, page_size: 10, total_items: 0 })
    expect(store.loading).toBe(false)
    expect(store.error).toBe('')
  })
})
