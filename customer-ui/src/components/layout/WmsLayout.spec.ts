import { describe, it, expect, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import WmsLayout from './WmsLayout.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', redirect: '/login' },
    { path: '/dashboard', name: 'Dashboard', component: { template: '<div>Dashboard</div>' } },
    { path: '/login', name: 'Login', component: { template: '<div>Login</div>' } },
  ],
})

describe('WmsLayout', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('renders sidebar and content slot', () => {
    const wrapper = mount(WmsLayout, {
      slots: { default: '<div class="page-content">Page</div>' },
      global: { plugins: [router] },
    })
    expect(wrapper.findComponent({ name: 'WmsSidebar' }).exists()).toBe(true)
    expect(wrapper.text()).toContain('Page')
  })

  it('renders slot content', () => {
    const wrapper = mount(WmsLayout, {
      slots: { default: '<p>Hello Customer</p>' },
      global: { plugins: [router] },
    })
    expect(wrapper.text()).toContain('Hello Customer')
  })
})
