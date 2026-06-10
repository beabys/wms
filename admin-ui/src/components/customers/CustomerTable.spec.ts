import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import CustomerTable from './CustomerTable.vue'
import type { CustomerResponse } from '@/api/types'

const pendingCustomer: CustomerResponse = {
  id: '1', company_name: 'Acme Corp', email: 'acme@test.com', phone: '', vat_number: 'VAT001', status: 'pending', created_at: '2024-01-15',
}
const activeCustomer: CustomerResponse = {
  id: '2', company_name: 'Globex Inc', email: 'globex@test.com', phone: '', vat_number: 'VAT002', status: 'active', created_at: '2024-02-01',
}
const suspendedCustomer: CustomerResponse = {
  id: '3', company_name: 'Suspended Co', email: 'suspended@test.com', phone: '', vat_number: 'VAT003', status: 'suspended', created_at: '2024-03-01',
}

describe('CustomerTable', () => {
  it('renders loading text when loading', () => {
    const wrapper = mount(CustomerTable, {
      props: { customers: [], loading: true },
    })
    expect(wrapper.text()).toContain('Loading customers...')
  })

  it('renders empty state when no customers', () => {
    const wrapper = mount(CustomerTable, {
      props: { customers: [], loading: false },
    })
    expect(wrapper.text()).toContain('No customers found.')
  })

  it('renders customer rows', () => {
    const wrapper = mount(CustomerTable, {
      props: { customers: [pendingCustomer, activeCustomer], loading: false },
    })
    expect(wrapper.text()).toContain('Acme Corp')
    expect(wrapper.text()).toContain('Globex Inc')
  })

  it('emits approve with customer object', async () => {
    const wrapper = mount(CustomerTable, {
      props: { customers: [pendingCustomer], loading: false },
    })
    const approveBtn = wrapper.findAll('button').filter(b => b.text() === 'Approve')
    await approveBtn[0].trigger('click')
    expect(wrapper.emitted('approve')).toBeTruthy()
    expect(wrapper.emitted('approve')![0]).toEqual([pendingCustomer])
  })

  it('emits reject with customer object', async () => {
    const wrapper = mount(CustomerTable, {
      props: { customers: [pendingCustomer], loading: false },
    })
    const rejectBtn = wrapper.findAll('button').filter(b => b.text() === 'Reject')
    await rejectBtn[0].trigger('click')
    expect(wrapper.emitted('reject')).toBeTruthy()
    expect(wrapper.emitted('reject')![0]).toEqual([pendingCustomer])
  })

  it('emits suspend with customer object for active customers', async () => {
    const wrapper = mount(CustomerTable, {
      props: { customers: [activeCustomer], loading: false },
    })
    const suspendBtn = wrapper.findAll('button').filter(b => b.text() === 'Suspend')
    await suspendBtn[0].trigger('click')
    expect(wrapper.emitted('suspend')).toBeTruthy()
    expect(wrapper.emitted('suspend')![0]).toEqual([activeCustomer])
  })

  it('emits restore with customer object for suspended customers', async () => {
    const wrapper = mount(CustomerTable, {
      props: { customers: [suspendedCustomer], loading: false },
    })
    const restoreBtn = wrapper.findAll('button').filter(b => b.text() === 'Restore')
    await restoreBtn[0].trigger('click')
    expect(wrapper.emitted('restore')).toBeTruthy()
    expect(wrapper.emitted('restore')![0]).toEqual([suspendedCustomer])
  })

  it('shows restore button for suspended customer', () => {
    const wrapper = mount(CustomerTable, {
      props: { customers: [suspendedCustomer, activeCustomer, pendingCustomer], loading: false },
    })
    const buttons = wrapper.findAll('button').map(b => b.text())
    expect(buttons).toContain('Restore')
    expect(buttons).toContain('Approve')
    expect(buttons).toContain('Reject')
    expect(buttons).toContain('Suspend')
  })

  it('emits rowClick when row is clicked', async () => {
    const wrapper = mount(CustomerTable, {
      props: { customers: [pendingCustomer], loading: false },
    })
    const row = wrapper.find('tbody tr')
    await row.trigger('click')
    expect(wrapper.emitted('rowClick')).toBeTruthy()
    expect(wrapper.emitted('rowClick')![0]).toEqual([pendingCustomer])
  })
})
