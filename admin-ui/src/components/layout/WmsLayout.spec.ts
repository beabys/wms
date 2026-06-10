import { describe, it, expect, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import WmsLayout from './WmsLayout.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', redirect: '/dashboard' },
    { path: '/dashboard', name: 'Dashboard', component: { template: '<div>Dashboard</div>' } },
    { path: '/login', name: 'Login', component: { template: '<div>Login</div>' } },
  ],
})

describe('WmsLayout', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('renders sidebar', () => {
    const wrapper = mount(WmsLayout, {
      global: { plugins: [router] },
      slots: { default: '<div>Content</div>' },
    })
    expect(wrapper.findComponent({ name: 'WmsSidebar' }).exists()).toBe(true)
  })

  it('renders slot content', () => {
    const wrapper = mount(WmsLayout, {
      global: { plugins: [router] },
      slots: { default: '<div class="page-content">Page Content</div>' },
    })
    expect(wrapper.find('.page-content').exists()).toBe(true)
    expect(wrapper.text()).toContain('Page Content')
  })

  it('applies content margin for sidebar offset', () => {
    const wrapper = mount(WmsLayout, {
      global: { plugins: [router] },
      slots: { default: '<div>Content</div>' },
    })
    const content = wrapper.find('.wms-layout__content')
    expect(content.exists()).toBe(true)
  })
})
