import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import { useCustomerStore } from '@/stores/customers'
import CustomerListView from './CustomerListView.vue'

// happy-dom doesn't support <Teleport> natively, stub it to render inline
const teleportStub = { template: '<div><slot /></div>' }

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/customers', name: 'Customers', component: CustomerListView },
    { path: '/customers/:id', name: 'CustomerDetail', component: { template: '<div>Detail</div>' } },
  ],
})

const stubs = {
  ApproveRejectDialog: true,
  WmsConfirmDialog: true,
  InviteCustomerModal: true,
  Teleport: teleportStub,
}

vi.mock('@/api/customerClient', () => ({
  customerClient: {
    listCustomers: vi.fn().mockResolvedValue({ customers: [], pagination: { page: 1, page_size: 10, total_items: 0 } }),
    approveCustomer: vi.fn(),
    rejectCustomer: vi.fn(),
    suspendCustomer: vi.fn(),
    restoreCustomer: vi.fn(),
  },
}))

describe('CustomerListView', () => {
  async function createWrapper() {
    setActivePinia(createPinia())
    router.push('/customers')
    await router.isReady()
    return mount(CustomerListView, {
      global: { plugins: [router], stubs },
    })
  }

  it('renders title', async () => {
    const wrapper = await createWrapper()
    expect(wrapper.text()).toContain('Customers')
  })

  it('renders Invite Customer button', async () => {
    const wrapper = await createWrapper()
    const buttons = wrapper.findAll('button')
    const inviteBtn = buttons.find(b => b.text() === 'Invite Customer')
    expect(inviteBtn).toBeTruthy()
  })

  it('renders status filter select', async () => {
    const wrapper = await createWrapper()
    expect(wrapper.find('select').exists()).toBe(true)
  })

  it('renders customer table', async () => {
    const wrapper = await createWrapper()
    const table = wrapper.findComponent({ name: 'CustomerTable' })
    expect(table.exists()).toBe(true)
  })

  it('shows pagination when items exist', async () => {
    const wrapper = await createWrapper()
    const store = useCustomerStore()
    store.pagination = { page: 1, page_size: 10, total_items: 25 }
    store.customers = [
      { id: '1', company_name: 'A', email: 'a@a.com', phone: '', vat_number: '', status: 'active', created_at: '' },
    ]
    await wrapper.vm.$nextTick()
    expect(wrapper.text()).toContain('Page')
    expect(wrapper.text()).toContain('3')
  })

  it('shows invite modal when button clicked', async () => {
    const wrapper = await createWrapper()
    const inviteBtn = wrapper.findAll('button').find(b => b.text() === 'Invite Customer')
    await inviteBtn?.trigger('click')
    await wrapper.vm.$nextTick()
    const modal = wrapper.findComponent({ name: 'InviteCustomerModal' })
    expect(modal.props('visible')).toBe(true)
  })

  it('shows restore confirmation dialog', async () => {
    const wrapper = await createWrapper()
    const store = useCustomerStore()
    store.customers = [
      { id: '3', company_name: 'Suspended Co', email: 's@test.com', phone: '', vat_number: '', status: 'suspended', created_at: '' },
    ]
    await wrapper.vm.$nextTick()
    // Trigger restore via table
    const table = wrapper.findComponent({ name: 'CustomerTable' })
    table.vm.$emit('restore', store.customers[0])
    await wrapper.vm.$nextTick()
    const confirmDialog = wrapper.findComponent({ name: 'WmsConfirmDialog' })
    expect(confirmDialog.props('visible')).toBe(true)
  })
})
