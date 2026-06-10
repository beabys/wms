import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import WmsButton from './WmsButton.vue'

describe('WmsButton', () => {
  it('renders slot content', () => {
    const wrapper = mount(WmsButton, {
      slots: { default: 'Click me' },
    })
    expect(wrapper.text()).toContain('Click me')
  })

  it('shows spinner when loading', () => {
    const wrapper = mount(WmsButton, {
      props: { loading: true },
      slots: { default: 'Loading' },
    })
    expect(wrapper.find('.wms-btn__spinner').exists()).toBe(true)
    expect(wrapper.attributes('disabled')).toBeDefined()
  })

  it('disables button when disabled prop is true', () => {
    const wrapper = mount(WmsButton, {
      props: { disabled: true },
      slots: { default: 'Disabled' },
    })
    expect(wrapper.attributes('disabled')).toBeDefined()
  })

  it('emits click event when clicked', async () => {
    const wrapper = mount(WmsButton, {
      slots: { default: 'Click' },
    })
    await wrapper.trigger('click')
    expect(wrapper.emitted('click')).toBeTruthy()
  })

  it('does not emit click when disabled', async () => {
    const wrapper = mount(WmsButton, {
      props: { disabled: true },
      slots: { default: 'Click' },
    })
    await wrapper.trigger('click')
    expect(wrapper.emitted('click')).toBeFalsy()
  })

  it('does not emit click when loading', async () => {
    const wrapper = mount(WmsButton, {
      props: { loading: true },
      slots: { default: 'Loading' },
    })
    await wrapper.trigger('click')
    expect(wrapper.emitted('click')).toBeFalsy()
  })

  it('applies variant class', () => {
    const wrapper = mount(WmsButton, {
      props: { variant: 'danger' },
    })
    expect(wrapper.classes()).toContain('wms-btn--danger')
  })

  it('applies size class', () => {
    const wrapper = mount(WmsButton, {
      props: { size: 'lg' },
    })
    expect(wrapper.classes()).toContain('wms-btn--lg')
  })

  it('sets type attribute', () => {
    const wrapper = mount(WmsButton, {
      props: { type: 'submit' },
    })
    expect(wrapper.attributes('type')).toBe('submit')
  })
})
