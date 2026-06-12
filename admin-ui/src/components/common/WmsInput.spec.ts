import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import WmsInput from './WmsInput.vue'

describe('WmsInput', () => {
  it('renders label when provided', () => {
    const wrapper = mount(WmsInput, {
      props: { label: 'Email' },
    })
    expect(wrapper.find('.wms-input__label').text()).toBe('Email')
  })

  it('renders placeholder', () => {
    const wrapper = mount(WmsInput, {
      props: { placeholder: 'you@example.com' },
    })
    expect(wrapper.find('input').attributes('placeholder')).toBe('you@example.com')
  })

  it('shows error message', () => {
    const wrapper = mount(WmsInput, {
      props: { error: 'Invalid email' },
    })
    expect(wrapper.find('.wms-input__error').text()).toBe('Invalid email')
  })

  it('disables input when disabled prop is true', () => {
    const wrapper = mount(WmsInput, {
      props: { disabled: true },
    })
    expect(wrapper.find('input').attributes('disabled')).toBeDefined()
  })

  it('emits update:modelValue on input', async () => {
    const wrapper = mount(WmsInput, {
      props: { modelValue: '' },
    })
    await wrapper.find('input').setValue('test@test.com')
    expect(wrapper.emitted('update:modelValue')).toBeTruthy()
    expect(wrapper.emitted('update:modelValue')![0]).toEqual(['test@test.com'])
  })

  it('sets input type', () => {
    const wrapper = mount(WmsInput, {
      props: { type: 'email' },
    })
    expect(wrapper.find('input').attributes('type')).toBe('email')
  })

  it('renders prefix slot', () => {
    const wrapper = mount(WmsInput, {
      slots: { prefix: '<span class="prefix-icon">@</span>' },
    })
    expect(wrapper.find('.prefix-icon').exists()).toBe(true)
  })

  it('renders suffix slot', () => {
    const wrapper = mount(WmsInput, {
      slots: { suffix: '<button>Toggle</button>' },
    })
    expect(wrapper.find('button').exists()).toBe(true)
  })

  it('applies error class to wrapper', () => {
    const wrapper = mount(WmsInput, {
      props: { error: 'Error' },
    })
    expect(wrapper.find('.wms-input__wrapper').classes()).toContain('wms-input__wrapper--error')
  })
})
