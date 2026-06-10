import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useCustomerStore } from './customer'

const mockCustomer = {
  id: '1',
  company_name: 'Acme',
  email: 'a@b.com',
  phone: '123',
  vat_number: 'VAT123',
  status: 'active',
  created_at: '2024-01-01',
}

describe('useCustomerStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.restoreAllMocks()
  })

  it('initial state', () => {
    const store = useCustomerStore()
    expect(store.customer).toBeNull()
    expect(store.loading).toBe(false)
    expect(store.error).toBe('')
    expect(store.isApproved).toBe(false)
    expect(store.approvalStatus).toBe('unknown')
  })

  it('isApproved true when status active', () => {
    const store = useCustomerStore()
    store.customer = { ...mockCustomer, status: 'active' }
    expect(store.isApproved).toBe(true)
    expect(store.approvalStatus).toBe('active')
  })

  it('isApproved false when status pending', () => {
    const store = useCustomerStore()
    store.customer = { ...mockCustomer, status: 'pending' }
    expect(store.isApproved).toBe(false)
    expect(store.approvalStatus).toBe('pending')
  })

  it('isApproved false when status rejected', () => {
    const store = useCustomerStore()
    store.customer = { ...mockCustomer, status: 'rejected' }
    expect(store.isApproved).toBe(false)
    expect(store.approvalStatus).toBe('rejected')
  })

  it('isApproved false when status suspended', () => {
    const store = useCustomerStore()
    store.customer = { ...mockCustomer, status: 'suspended' }
    expect(store.isApproved).toBe(false)
    expect(store.approvalStatus).toBe('suspended')
  })

  it('approvalStatus returns unknown when customer null', () => {
    const store = useCustomerStore()
    expect(store.approvalStatus).toBe('unknown')
  })

  it('fetchMyCustomer fetches and sets customer', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true, data: mockCustomer }),
    }))
    const store = useCustomerStore()
    await store.fetchMyCustomer()
    expect(store.customer).toEqual(mockCustomer)
    expect(store.loading).toBe(false)
    expect(store.error).toBe('')
    vi.restoreAllMocks()
  })

  it('fetchMyCustomer sets error on failure', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 404,
      json: async () => ({ success: false, error: 'Not found' }),
    }))
    const store = useCustomerStore()
    await store.fetchMyCustomer()
    expect(store.customer).toBeNull()
    expect(store.error).toBe('Not found')
    expect(store.loading).toBe(false)
    vi.restoreAllMocks()
  })

  it('fetchMyCustomer sets generic error when no message', async () => {
    vi.stubGlobal('fetch', vi.fn().mockRejectedValue({}))
    const store = useCustomerStore()
    await store.fetchMyCustomer()
    expect(store.error).toBe('Failed to fetch customer')
    expect(store.loading).toBe(false)
    vi.restoreAllMocks()
  })

  it('$reset clears all state', () => {
    const store = useCustomerStore()
    store.customer = mockCustomer as any
    store.loading = true
    store.error = 'something'
    store.$reset()
    expect(store.customer).toBeNull()
    expect(store.loading).toBe(false)
    expect(store.error).toBe('')
  })
})
