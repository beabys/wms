import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import WmsLoginForm from './WmsLoginForm.vue'

describe('WmsLoginForm', () => {
  function createWrapper(props = {}) {
    setActivePinia(createPinia())
    return mount(WmsLoginForm, {
      props,
      global: { stubs: { WmsButton: true, WmsInput: true, WmsPasswordInput: true } },
    })
  }

  it('renders title', () => {
    const wrapper = createWrapper()
    expect(wrapper.text()).toContain('Welcome Back')
  })

  it('renders error message', () => {
    const wrapper = createWrapper({ error: 'Invalid credentials' })
    expect(wrapper.text()).toContain('Invalid credentials')
  })

  it('emits submit event', async () => {
    const wrapper = createWrapper()
    ;(wrapper.vm as any).email = 'test@test.com'
    ;(wrapper.vm as any).password = 'password123'
    await wrapper.vm.$nextTick()
    await wrapper.find('form').trigger('submit.prevent')
    expect(wrapper.emitted('submit')).toBeTruthy()
  })
})
