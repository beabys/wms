import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import WmsInput from './WmsInput.vue'

describe('WmsInput', () => {
  it('renders label', () => {
    const wrapper = mount(WmsInput, { props: { label: 'Email' } })
    expect(wrapper.find('.wms-input__label').text()).toBe('Email')
  })

  it('shows error', () => {
    const wrapper = mount(WmsInput, { props: { error: 'Required' } })
    expect(wrapper.find('.wms-input__error').text()).toBe('Required')
  })

  it('emits update on input', async () => {
    const wrapper = mount(WmsInput)
    await wrapper.find('input').setValue('test@test.com')
    expect(wrapper.emitted('update:modelValue')).toBeTruthy()
  })

  it('renders prefix slot', () => {
    const wrapper = mount(WmsInput, { slots: { prefix: '<span class="pfx">@</span>' } })
    expect(wrapper.find('.pfx').exists()).toBe(true)
  })

  it('renders suffix slot', () => {
    const wrapper = mount(WmsInput, { slots: { suffix: '<button>T</button>' } })
    expect(wrapper.find('button').exists()).toBe(true)
  })
})
