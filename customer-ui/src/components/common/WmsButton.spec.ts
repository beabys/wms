import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import WmsButton from './WmsButton.vue'

describe('WmsButton', () => {
  it('renders slot content', () => {
    const wrapper = mount(WmsButton, { slots: { default: 'Get Started' } })
    expect(wrapper.text()).toContain('Get Started')
  })

  it('shows spinner when loading', () => {
    const wrapper = mount(WmsButton, { props: { loading: true }, slots: { default: 'Loading' } })
    expect(wrapper.find('.wms-btn__spinner').exists()).toBe(true)
    expect(wrapper.attributes('disabled')).toBeDefined()
  })

  it('disables when disabled', () => {
    const wrapper = mount(WmsButton, { props: { disabled: true }, slots: { default: 'Disabled' } })
    expect(wrapper.attributes('disabled')).toBeDefined()
  })

  it('emits click', async () => {
    const wrapper = mount(WmsButton, { slots: { default: 'Click' } })
    await wrapper.trigger('click')
    expect(wrapper.emitted('click')).toBeTruthy()
  })

  it('does not emit when disabled', async () => {
    const wrapper = mount(WmsButton, { props: { disabled: true }, slots: { default: 'Click' } })
    await wrapper.trigger('click')
    expect(wrapper.emitted('click')).toBeFalsy()
  })

  it('applies variant and size classes', () => {
    const wrapper = mount(WmsButton, { props: { variant: 'danger', size: 'lg' } })
    expect(wrapper.classes()).toContain('wms-btn--danger')
    expect(wrapper.classes()).toContain('wms-btn--lg')
  })
})
