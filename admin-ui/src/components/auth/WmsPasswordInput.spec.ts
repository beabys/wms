import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import WmsPasswordInput from './WmsPasswordInput.vue'

describe('WmsPasswordInput', () => {
  it('renders label and placeholder', () => {
    const wrapper = mount(WmsPasswordInput, {
      props: { label: 'Password', placeholder: 'Enter password' },
      global: {
        stubs: { WmsInput: false },
      },
    })
    expect(wrapper.text()).toContain('Password')
  })

  it('renders show/hide toggle button', () => {
    const wrapper = mount(WmsPasswordInput)
    expect(wrapper.find('button').exists()).toBe(true)
  })

  it('toggles password visibility on click', async () => {
    const wrapper = mount(WmsPasswordInput)
    const input = wrapper.find('input')
    expect(input.attributes('type')).toBe('password')

    await wrapper.find('button').trigger('click')
    expect(input.attributes('type')).toBe('text')

    await wrapper.find('button').trigger('click')
    expect(input.attributes('type')).toBe('password')
  })

  it('emits update:modelValue on input', async () => {
    const wrapper = mount(WmsPasswordInput, {
      props: { modelValue: '' },
    })
    const input = wrapper.find('input')
    await input.setValue('secret123')
    expect(wrapper.emitted('update:modelValue')).toBeTruthy()
  })
})
