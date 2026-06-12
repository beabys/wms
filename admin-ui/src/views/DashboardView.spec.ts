import { describe, it, expect, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import DashboardView from './DashboardView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', redirect: '/login' },
    { path: '/dashboard', name: 'Dashboard', component: DashboardView },
  ],
})

describe('DashboardView', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('renders dashboard heading', async () => {
    router.push('/dashboard')
    await router.isReady()
    const wrapper = mount(DashboardView, {
      global: { plugins: [router] },
    })
    expect(wrapper.text()).toContain('Dashboard')
  })

  it('renders welcome message', async () => {
    router.push('/dashboard')
    await router.isReady()
    const wrapper = mount(DashboardView, {
      global: { plugins: [router] },
    })
    expect(wrapper.text()).toContain('Welcome')
  })
})
