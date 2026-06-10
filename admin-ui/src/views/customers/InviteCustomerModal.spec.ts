import { describe, it, expect, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import InviteCustomerModal from './InviteCustomerModal.vue'

const teleportStub = { template: '<div><slot /></div>' }

describe('InviteCustomerModal', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  function mountModal(props = {}) {
    return mount(InviteCustomerModal, {
      props: { visible: true, ...props },
      global: {
        stubs: { Teleport: teleportStub },
      },
    })
  }

  it('does not render when visible is false', () => {
    const wrapper = mount(InviteCustomerModal, {
      props: { visible: false },
      global: { stubs: { Teleport: teleportStub } },
    })
    expect(wrapper.find('.invite-overlay').exists()).toBe(false)
  })

  it('renders when visible is true', () => {
    const wrapper = mountModal()
    expect(wrapper.find('.invite-overlay').exists()).toBe(true)
    expect(wrapper.text()).toContain('Invite Customer')
  })

  it('renders email input', () => {
    const wrapper = mountModal()
    const input = wrapper.findComponent({ name: 'WmsInput' })
    expect(input.exists()).toBe(true)
    expect(input.props('label')).toBe('Email address')
  })

  it('emits close when overlay clicked', async () => {
    const wrapper = mountModal()
    await wrapper.find('.invite-overlay').trigger('click')
    expect(wrapper.emitted('close')).toBeTruthy()
  })

  it('emits close when cancel button clicked', async () => {
    const wrapper = mountModal()
    const buttons = wrapper.findAllComponents({ name: 'WmsButton' })
    const cancelBtn = buttons.find(b => b.text() === 'Cancel')
    await cancelBtn?.trigger('click')
    expect(wrapper.emitted('close')).toBeTruthy()
  })

  it('disables submit button when email is empty', () => {
    const wrapper = mountModal()
    const buttons = wrapper.findAllComponents({ name: 'WmsButton' })
    const submitBtn = buttons.find(b => b.text() === 'Send Invitation')
    expect(submitBtn?.props('disabled')).toBe(true)
  })
})
