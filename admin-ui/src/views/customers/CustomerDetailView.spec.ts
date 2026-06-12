import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import CustomerDetailView from './CustomerDetailView.vue'

function flushPromises() {
  return new Promise(resolve => setTimeout(resolve, 0))
}

const teleportStub = { template: '<div><slot /></div>' }

const stubs = {
  WmsConfirmDialog: true,
  Teleport: teleportStub,
}

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
  status: 'active',
  created_at: '2024-01-01',
}

const mockAuditLog = {
  entries: [
    { id: 'a1', customer_id: '1', action: 'approved', performed_by: 'admin@test.com', details: 'Approved by admin', created_at: 1700000000 },
  ],
  pagination: { page: 1, page_size: 10, total_items: 1 },
}

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/customers/:id', name: 'CustomerDetail', component: CustomerDetailView },
    { path: '/customers', name: 'Customers', component: { template: '<div>Customers</div>' } },
  ],
})

describe('CustomerDetailView', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.restoreAllMocks()
  })

  async function createWrapper() {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValueOnce({
      status: 200,
      json: async () => ({ success: true, data: mockCustomer }),
    }).mockResolvedValueOnce({
      status: 200,
      json: async () => ({ success: true, data: mockAuditLog }),
    }))

    router.push('/customers/1')
    await router.isReady()

    const wrapper = mount(CustomerDetailView, {
      global: {
        plugins: [router],
        stubs,
      },
    })
    await flushPromises()
    await wrapper.vm.$nextTick()
    return wrapper
  }

  it('renders company name in header', async () => {
    const wrapper = await createWrapper()
    expect(wrapper.text()).toContain('Acme Corp')
  })

  it('renders status badge', async () => {
    const wrapper = await createWrapper()
    expect(wrapper.text()).toContain('Active')
  })

  it('renders Edit button', async () => {
    const wrapper = await createWrapper()
    const editBtn = wrapper.findAll('button').find(b => b.text() === 'Edit')
    expect(editBtn).toBeTruthy()
  })

  it('renders company info fields', async () => {
    const wrapper = await createWrapper()
    expect(wrapper.text()).toContain('Acme Corp')
    expect(wrapper.text()).toContain('acme@test.com')
    expect(wrapper.text()).toContain('123456789')
  })

  it('renders actions card', async () => {
    const wrapper = await createWrapper()
    expect(wrapper.text()).toContain('Actions')
  })

  it('renders audit log section', async () => {
    const wrapper = await createWrapper()
    expect(wrapper.text()).toContain('Audit Log')
  })

  it('renders audit entries when present', async () => {
    const wrapper = await createWrapper()
    expect(wrapper.text()).toContain('Approved')
    expect(wrapper.text()).toContain('admin@test.com')
  })

  it('shows Suspend button for active customer', async () => {
    const wrapper = await createWrapper()
    const buttons = wrapper.findAll('button').map(b => b.text())
    expect(buttons).toContain('Suspend')
  })

  it('shows approve/reject for pending customer', async () => {
    setActivePinia(createPinia())
    const pendingCustomer = { ...mockCustomer, status: 'pending' }
    vi.stubGlobal('fetch', vi.fn().mockResolvedValueOnce({
      status: 200,
      json: async () => ({ success: true, data: pendingCustomer }),
    }).mockResolvedValueOnce({
      status: 200,
      json: async () => ({ success: true, data: mockAuditLog }),
    }))

    router.push('/customers/1')
    await router.isReady()

    const wrapper = mount(CustomerDetailView, {
      global: {
        plugins: [router],
        stubs,
      },
    })
    await flushPromises()
    await wrapper.vm.$nextTick()

    const buttons = wrapper.findAll('button').map(b => b.text())
    expect(buttons).toContain('Approve')
    expect(buttons).toContain('Reject')
  })

  it('shows Restore button for suspended customer', async () => {
    setActivePinia(createPinia())
    const suspendedCustomer = { ...mockCustomer, status: 'suspended' }
    vi.stubGlobal('fetch', vi.fn().mockResolvedValueOnce({
      status: 200,
      json: async () => ({ success: true, data: suspendedCustomer }),
    }).mockResolvedValueOnce({
      status: 200,
      json: async () => ({ success: true, data: mockAuditLog }),
    }))

    router.push('/customers/1')
    await router.isReady()

    const wrapper = mount(CustomerDetailView, {
      global: {
        plugins: [router],
        stubs,
      },
    })
    await flushPromises()
    await wrapper.vm.$nextTick()

    const buttons = wrapper.findAll('button').map(b => b.text())
    expect(buttons).toContain('Restore')
  })

  it('shows loading state', async () => {
    setActivePinia(createPinia())
    // Don't resolve fetch to keep loading
    vi.stubGlobal('fetch', vi.fn().mockReturnValue(new Promise(() => {})))

    router.push('/customers/1')
    await router.isReady()

    const wrapper = mount(CustomerDetailView, {
      global: {
        plugins: [router],
        stubs,
      },
    })
    await wrapper.vm.$nextTick()
    expect(wrapper.text()).toContain('Loading customer details...')
  })

  it('shows error state when fetch fails', async () => {
    setActivePinia(createPinia())
    vi.stubGlobal('fetch', vi.fn().mockRejectedValue(new Error('Failed to load')))

    router.push('/customers/1')
    await router.isReady()

    const wrapper = mount(CustomerDetailView, {
      global: {
        plugins: [router],
        stubs,
      },
    })
    await flushPromises()
    await wrapper.vm.$nextTick()
    expect(wrapper.text()).toContain('Back to Customers')
  })

  it('switches to edit mode on Edit click', async () => {
    const wrapper = await createWrapper()

    const editBtn = wrapper.findAll('button').find(b => b.text() === 'Edit')
    await editBtn?.trigger('click')
    await wrapper.vm.$nextTick()

    expect(wrapper.text()).toContain('Save')
    expect(wrapper.text()).toContain('Cancel')
  })
})
