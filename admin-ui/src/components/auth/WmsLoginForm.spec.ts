import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import WmsLoginForm from './WmsLoginForm.vue'

describe('WmsLoginForm', () => {
  function createWrapper(props = {}) {
    setActivePinia(createPinia())
    return mount(WmsLoginForm, {
      props,
      global: {
        stubs: {
          WmsButton: true,
          WmsInput: true,
          WmsPasswordInput: true,
        },
      },
    })
  }

  it('renders title and subtitle', () => {
    const wrapper = createWrapper()
    expect(wrapper.text()).toContain('Sign In')
    expect(wrapper.text()).toContain('Welcome back!')
  })

  it('renders error message when error prop is provided', () => {
    const wrapper = createWrapper({ error: 'Invalid credentials' })
    expect(wrapper.text()).toContain('Invalid credentials')
  })

  it('emits submit event with email and password', async () => {
    const wrapper = createWrapper()
    // Set internal reactive state
    ;(wrapper.vm as any).email = 'test@test.com'
    ;(wrapper.vm as any).password = 'password123'
    await wrapper.vm.$nextTick()
    // Trigger submit
    await wrapper.find('form').trigger('submit.prevent')
    expect(wrapper.emitted('submit')).toBeTruthy()
  })

  it('shows button with loading state', () => {
    const wrapper = createWrapper({ loading: true })
    const button = wrapper.findComponent({ name: 'WmsButton' })
    expect(button.props('loading')).toBe(true)
  })
})
