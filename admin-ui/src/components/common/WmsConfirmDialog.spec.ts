import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import WmsConfirmDialog from './WmsConfirmDialog.vue'

// happy-dom doesn't support <Teleport> natively, stub it to render inline
const teleportStub = { template: '<div><slot /></div>' }

describe('WmsConfirmDialog', () => {
  function mountDialog(props = {}, slots = {}) {
    return mount(WmsConfirmDialog, {
      props: { visible: true, ...props },
      slots,
      global: {
        stubs: { Teleport: teleportStub },
      },
    })
  }

  it('renders when visible is true', () => {
    const wrapper = mountDialog({ title: 'Test Title', message: 'Test message' })
    expect(wrapper.find('.wms-confirm-overlay').exists()).toBe(true)
    expect(wrapper.text()).toContain('Test Title')
    expect(wrapper.text()).toContain('Test message')
  })

  it('does not render when visible is false', () => {
    const wrapper = mount(WmsConfirmDialog, {
      props: { visible: false },
      global: { stubs: { Teleport: teleportStub } },
    })
    expect(wrapper.find('.wms-confirm-overlay').exists()).toBe(false)
  })

  it('renders default slot content instead of message prop', () => {
    const wrapper = mountDialog(
      { message: 'prop message' },
      { default: '<strong>slot content</strong>' },
    )
    expect(wrapper.text()).toContain('slot content')
    expect(wrapper.text()).not.toContain('prop message')
  })

  it('emits confirm when confirm button clicked', async () => {
    const wrapper = mountDialog()
    const confirmBtn = wrapper.findAllComponents({ name: 'WmsButton' })[1]
    await confirmBtn.trigger('click')
    expect(wrapper.emitted('confirm')).toBeTruthy()
  })

  it('emits cancel when cancel button clicked', async () => {
    const wrapper = mountDialog()
    const cancelBtn = wrapper.findAllComponents({ name: 'WmsButton' })[0]
    await cancelBtn.trigger('click')
    expect(wrapper.emitted('cancel')).toBeTruthy()
  })

  it('emits cancel when overlay clicked', async () => {
    const wrapper = mountDialog()
    await wrapper.find('.wms-confirm-overlay').trigger('click')
    expect(wrapper.emitted('cancel')).toBeTruthy()
  })

  it('renders custom button text', () => {
    const wrapper = mountDialog({ confirmText: 'Yes', cancelText: 'No' })
    expect(wrapper.text()).toContain('Yes')
    expect(wrapper.text()).toContain('No')
  })

  it('applies danger variant to confirm button', () => {
    const wrapper = mountDialog({ variant: 'danger' })
    const confirmBtn = wrapper.findAllComponents({ name: 'WmsButton' })[1]
    expect(confirmBtn.props('variant')).toBe('danger')
  })
})
