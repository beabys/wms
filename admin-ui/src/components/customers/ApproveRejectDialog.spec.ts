import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import ApproveRejectDialog from './ApproveRejectDialog.vue'

const customer = {
  id: '1',
  company_name: 'Acme Corp',
  email: 'acme@test.com',
  phone: '123456789',
  vat_number: 'VAT123',
  status: 'pending',
  created_at: '2024-01-01',
}

// happy-dom doesn't support <Teleport> natively, stub it to render inline
const teleportStub = { template: '<div><slot /></div>' }

describe('ApproveRejectDialog', () => {
  function mountDialog(props = {}, slots = {}) {
    return mount(ApproveRejectDialog, {
      props: { visible: true, customer, action: 'approve', ...props },
      slots,
      global: {
        stubs: { Teleport: teleportStub },
      },
    })
  }

  it('does not render when visible is false', () => {
    const wrapper = mount(ApproveRejectDialog, {
      props: { visible: false, customer, action: 'approve' },
      global: { stubs: { Teleport: teleportStub } },
    })
    expect(wrapper.find('.dialog-overlay').exists()).toBe(false)
  })

  it('renders when visible is true', () => {
    const wrapper = mountDialog()
    expect(wrapper.find('.dialog-overlay').exists()).toBe(true)
    expect(wrapper.text()).toContain('Approve Customer')
    expect(wrapper.text()).toContain('Acme Corp')
  })

  it('shows approve title and primary button for approve action', () => {
    const wrapper = mountDialog()
    expect(wrapper.text()).toContain('Approve')
    const buttons = wrapper.findAllComponents({ name: 'WmsButton' })
    const confirmBtn = buttons[1]
    expect(confirmBtn.props('variant')).toBe('primary')
  })

  it('shows danger variant for reject action', () => {
    const wrapper = mountDialog({ action: 'reject' })
    const buttons = wrapper.findAllComponents({ name: 'WmsButton' })
    const confirmBtn = buttons[1]
    expect(confirmBtn.props('variant')).toBe('danger')
  })

  it('shows danger variant for suspend action', () => {
    const wrapper = mountDialog({ action: 'suspend' })
    const buttons = wrapper.findAllComponents({ name: 'WmsButton' })
    const confirmBtn = buttons[1]
    expect(confirmBtn.props('variant')).toBe('danger')
  })

  it('emits confirm when confirm button clicked', async () => {
    const wrapper = mountDialog()
    const buttons = wrapper.findAllComponents({ name: 'WmsButton' })
    await buttons[1].trigger('click')
    expect(wrapper.emitted('confirm')).toBeTruthy()
  })

  it('emits cancel when cancel button clicked', async () => {
    const wrapper = mountDialog()
    const buttons = wrapper.findAllComponents({ name: 'WmsButton' })
    await buttons[0].trigger('click')
    expect(wrapper.emitted('cancel')).toBeTruthy()
  })

  it('emits cancel when overlay clicked', async () => {
    const wrapper = mountDialog()
    await wrapper.find('.dialog-overlay').trigger('click')
    expect(wrapper.emitted('cancel')).toBeTruthy()
  })

  it('renders default slot content', () => {
    const wrapper = mountDialog({ action: 'suspend' }, { default: 'Custom confirmation message' })
    expect(wrapper.text()).toContain('Custom confirmation message')
  })
})
