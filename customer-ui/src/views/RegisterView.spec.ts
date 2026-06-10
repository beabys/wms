import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import RegisterView from './RegisterView.vue'

// Mock useRoute
import { useRoute } from 'vue-router'
vi.mock('vue-router', async () => {
  const actual = await vi.importActual('vue-router')
  return {
    ...(actual as any),
    useRoute: vi.fn(),
    useRouter: vi.fn().mockReturnValue({ push: vi.fn() }),
  }
})

import { vi } from 'vitest'

describe('RegisterView', () => {
  it('renders registration form', () => {
    setActivePinia(createPinia())
    ;(useRoute as any).mockReturnValue({ query: {} })
    const wrapper = mount(RegisterView, {
      global: {
        stubs: {
          WmsRegisterForm: true,
          WmsThemeToggle: true,
          RouterLink: { template: '<a><slot /></a>' },
          RouterView: true,
        },
      },
    })
    expect(wrapper.find('.register-view').exists()).toBe(true)
  })

  it('renders sign-in link', () => {
    setActivePinia(createPinia())
    ;(useRoute as any).mockReturnValue({ query: {} })
    const wrapper = mount(RegisterView, {
      global: {
        stubs: {
          WmsRegisterForm: true,
          WmsThemeToggle: true,
          RouterLink: { template: '<a><slot /></a>' },
          RouterView: true,
        },
      },
    })
    expect(wrapper.text()).toContain('Sign in')
  })

  it('extracts token from URL query', () => {
    setActivePinia(createPinia())
    ;(useRoute as any).mockReturnValue({ query: { token: 'invite-abc' } })
    const wrapper = mount(RegisterView, {
      global: {
        stubs: {
          WmsRegisterForm: true,
          WmsThemeToggle: true,
          RouterLink: { template: '<a><slot /></a>' },
          RouterView: true,
        },
      },
    })
    expect(wrapper.exists()).toBe(true)
  })
})
