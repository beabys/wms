import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import DashboardView from './DashboardView.vue'

function mockFetch(response: any) {
  vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
    status: 200,
    json: async () => ({ success: true, data: response }),
  }))
}

function flushPromises() {
  return new Promise(resolve => setTimeout(resolve, 0))
}

describe('DashboardView', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.restoreAllMocks()
  })

  it('renders dashboard heading', async () => {
    mockFetch({ id: '1', company_name: 'Acme', email: 'a@b.com', phone: '', vat_number: '', status: 'active', created_at: '' })
    const wrapper = mount(DashboardView)
    await flushPromises()
    expect(wrapper.text()).toContain('Dashboard')
  })

  it('shows welcome message when approved', async () => {
    mockFetch({ id: '1', company_name: 'Acme', email: 'a@b.com', phone: '', vat_number: '', status: 'active', created_at: '' })
    const wrapper = mount(DashboardView)
    await flushPromises()
    expect(wrapper.text()).toContain('Welcome to the WMS Customer Portal')
  })

  it('shows pending approval message when status pending', async () => {
    mockFetch({ id: '1', company_name: 'Acme', email: 'a@b.com', phone: '', vat_number: '', status: 'pending', created_at: '' })
    const wrapper = mount(DashboardView)
    await flushPromises()
    expect(wrapper.text()).toContain('Account Pending Approval')
    expect(wrapper.text()).toContain('awaiting approval')
  })

  it('shows rejected message when status rejected', async () => {
    mockFetch({ id: '1', company_name: 'Acme', email: 'a@b.com', phone: '', vat_number: '', status: 'rejected', created_at: '' })
    const wrapper = mount(DashboardView)
    await flushPromises()
    expect(wrapper.text()).toContain('Account Rejected')
  })

  it('shows suspended message when status suspended', async () => {
    mockFetch({ id: '1', company_name: 'Acme', email: 'a@b.com', phone: '', vat_number: '', status: 'suspended', created_at: '' })
    const wrapper = mount(DashboardView)
    await flushPromises()
    expect(wrapper.text()).toContain('Account Suspended')
  })

  it('does not render theme toggle', async () => {
    mockFetch({ id: '1', company_name: 'Acme', email: 'a@b.com', phone: '', vat_number: '', status: 'active', created_at: '' })
    const wrapper = mount(DashboardView)
    await flushPromises()
    expect(wrapper.find('.wms-theme-toggle').exists()).toBe(false)
  })
})
