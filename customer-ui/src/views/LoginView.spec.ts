import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import LoginView from './LoginView.vue'

describe('LoginView', () => {
  it('renders login form', () => {
    setActivePinia(createPinia())
    const wrapper = mount(LoginView, {
      global: {
        stubs: {
          WmsLoginForm: true,
          WmsThemeToggle: true,
          RouterLink: { template: '<a><slot /></a>' },
          RouterView: true,
        },
      },
    })
    expect(wrapper.find('.login-view').exists()).toBe(true)
  })

  it('renders register link', () => {
    setActivePinia(createPinia())
    const wrapper = mount(LoginView, {
      global: {
        stubs: {
          WmsLoginForm: true,
          WmsThemeToggle: true,
          RouterLink: { template: '<a><slot /></a>' },
          RouterView: true,
        },
      },
    })
    expect(wrapper.text()).toContain('Register here')
  })
})
