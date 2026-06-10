import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import WmsSidebar from './WmsSidebar.vue'
import { useAuthStore } from '@/stores/auth'
import { useCustomerStore } from '@/stores/customer'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', redirect: '/login' },
    { path: '/dashboard', name: 'Dashboard', component: { template: '<div>Dashboard</div>' } },
    { path: '/login', name: 'Login', component: { template: '<div>Login</div>' } },
  ],
})

function createWrapper(options: { user?: { email: string }; customerStatus?: string } = {}) {
  const store = useAuthStore()
  if (options.user) {
    store.setUser(options.user as any)
  }
  const customerStore = useCustomerStore()
  if (options.customerStatus) {
    customerStore.customer = { id: '1', company_name: 'Acme', email: 'a@b.com', phone: '', vat_number: '', status: options.customerStatus, created_at: '' }
  }
  return mount(WmsSidebar, {
    global: {
      plugins: [router],
    },
  })
}

describe('WmsSidebar', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('renders brand WMS Customer', () => {
    const wrapper = createWrapper()
    expect(wrapper.text()).toContain('WMS')
    expect(wrapper.text()).toContain('Customer')
  })

  it('shows nav links when customer approved', () => {
    const wrapper = createWrapper({ customerStatus: 'active' })
    expect(wrapper.text()).toContain('Dashboard')
  })

  it('hides nav links when customer pending', () => {
    const wrapper = createWrapper({ customerStatus: 'pending' })
    expect(wrapper.text()).not.toContain('Dashboard')
  })

  it('hides nav links when customer rejected', () => {
    const wrapper = createWrapper({ customerStatus: 'rejected' })
    expect(wrapper.text()).not.toContain('Dashboard')
  })

  it('hides nav links when customer suspended', () => {
    const wrapper = createWrapper({ customerStatus: 'suspended' })
    expect(wrapper.text()).not.toContain('Dashboard')
  })

  it('hides nav links when customer null', () => {
    const wrapper = createWrapper()
    expect(wrapper.text()).not.toContain('Dashboard')
  })

  it('highlights active nav link', async () => {
    const store = useAuthStore()
    store.setTokens('token', 'refresh')
    const customerStore = useCustomerStore()
    customerStore.customer = { id: '1', company_name: 'Acme', email: 'a@b.com', phone: '', vat_number: '', status: 'active', created_at: '' }
    await router.push('/dashboard')
    const wrapper = createWrapper({ customerStatus: 'active' })
    const activeLink = wrapper.find('.wms-sidebar__link--active')
    expect(activeLink.exists()).toBe(true)
    expect(activeLink.text()).toBe('Dashboard')
  })

  it('shows user email when authenticated', () => {
    const wrapper = createWrapper({ user: { email: 'cust@test.com' } })
    expect(wrapper.text()).toContain('cust@test.com')
  })

  it('renders theme toggle', () => {
    const wrapper = createWrapper()
    expect(wrapper.findComponent({ name: 'WmsThemeToggle' }).exists()).toBe(true)
  })

  it('renders logout button with ghost class', () => {
    const wrapper = createWrapper()
    const logoutBtn = wrapper.find('.wms-sidebar__logout--ghost')
    expect(logoutBtn.exists()).toBe(true)
    expect(logoutBtn.text()).toBe('Logout')
  })

  it('emits logout when logout button clicked', async () => {
    const wrapper = createWrapper()
    await wrapper.find('.wms-sidebar__logout--ghost').trigger('click')
    expect(wrapper.emitted('logout')).toBeTruthy()
  })

  it('calls router.push with dashboard path on link click', async () => {
    const pushSpy = vi.spyOn(router, 'push')
    const wrapper = createWrapper({ customerStatus: 'active' })
    await wrapper.find('.wms-sidebar__link').trigger('click')
    expect(pushSpy).toHaveBeenCalledWith('/dashboard')
  })
})
