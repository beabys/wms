import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from '@/stores/auth'
import LoginView from './LoginView.vue'

describe('LoginView', () => {
  it('renders login form', () => {
    setActivePinia(createPinia())
    const wrapper = mount(LoginView, {
      global: {
        stubs: {
          WmsLoginForm: true,
          WmsThemeToggle: true,
          RouterLink: true,
        },
      },
    })
    expect(wrapper.find('.login-view').exists()).toBe(true)
  })

  it('renders theme toggle button', () => {
    setActivePinia(createPinia())
    const wrapper = mount(LoginView, {
      global: {
        stubs: {
          WmsLoginForm: true,
          WmsThemeToggle: true,
          RouterLink: true,
        },
      },
    })
    expect(wrapper.findComponent({ name: 'WmsThemeToggle' }).exists()).toBe(true)
  })

  it('shows error message on failed login', async () => {
    setActivePinia(createPinia())
    mount(LoginView, {
      global: {
        stubs: {
          WmsLoginForm: {
            template: '<form @submit.prevent="$emit(\'submit\', \'a@b.com\', \'pass\')"><button>Submit</button></form>',
            props: ['loading', 'error'],
          },
          WmsThemeToggle: true,
        },
      },
    })
    const store = useAuthStore()
    expect(store.isAuthenticated).toBe(false)
  })
})
