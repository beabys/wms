import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import WmsPasswordInput from './WmsPasswordInput.vue'

describe('WmsPasswordInput', () => {
  it('renders label and toggle button', () => {
    const wrapper = mount(WmsPasswordInput, { props: { label: 'Password' } })
    expect(wrapper.text()).toContain('Password')
    expect(wrapper.find('button').exists()).toBe(true)
  })

  it('toggles password visibility', async () => {
    const wrapper = mount(WmsPasswordInput)
    const input = wrapper.find('input')
    expect(input.attributes('type')).toBe('password')
    await wrapper.find('button').trigger('click')
    expect(input.attributes('type')).toBe('text')
  })

  it('shows strength indicator when showStrength is true', () => {
    const wrapper = mount(WmsPasswordInput, {
      props: { modelValue: 'StrongP@ss1', showStrength: true, strength: 4, strengthLabel: 'Strong' },
    })
    expect(wrapper.find('.wms-password-strength').exists()).toBe(true)
    expect(wrapper.text()).toContain('Strong')
  })

  it('emits update:modelValue on input', async () => {
    const wrapper = mount(WmsPasswordInput, { props: { modelValue: '' } })
    await wrapper.find('input').setValue('secret123')
    expect(wrapper.emitted('update:modelValue')).toBeTruthy()
  })
})
