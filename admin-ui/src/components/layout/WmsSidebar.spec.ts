import { describe, it, expect, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import { nextTick } from 'vue'
import WmsSidebar from './WmsSidebar.vue'
import { useAuthStore } from '@/stores/auth'

const routes = [
  { path: '/', redirect: '/login' },
  { path: '/login', name: 'Login', component: { template: '<div>Login</div>' } },
  { path: '/dashboard', name: 'Dashboard', component: { template: '<div>Dashboard</div>' } },
  { path: '/users', name: 'Users', component: { template: '<div>Users</div>' } },
  { path: '/users/create', name: 'CreateUser', component: { template: '<div>Create</div>' } },
  { path: '/customers', name: 'Customers', component: { template: '<div>Customers</div>' } },
  { path: '/customers/invite', name: 'InviteCustomer', component: { template: '<div>Invite</div>' } },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

describe('WmsSidebar', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('renders brand name', async () => {
    await router.push('/dashboard')
    await router.isReady()
    const wrapper = mount(WmsSidebar, {
      global: { plugins: [router] },
    })
    expect(wrapper.text()).toContain('WMS')
    expect(wrapper.text()).toContain('Admin')
  })

  it('renders all navigation links', async () => {
    await router.push('/dashboard')
    await router.isReady()
    const wrapper = mount(WmsSidebar, {
      global: { plugins: [router] },
    })
    const links = wrapper.findAll('.wms-sidebar__link')
    expect(links).toHaveLength(5)
    expect(links[0].text()).toBe('Dashboard')
    expect(links[1].text()).toBe('Users')
    expect(links[2].text()).toBe('Create User')
    expect(links[3].text()).toBe('Customers')
    expect(links[4].text()).toBe('Invite Customer')
  })

  it('highlights active link', async () => {
    await router.push('/users')
    await router.isReady()
    const wrapper = mount(WmsSidebar, {
      global: { plugins: [router] },
    })
    await nextTick()
    const activeLinks = wrapper.findAll('.wms-sidebar__link--active')
    expect(activeLinks).toHaveLength(1)
    expect(activeLinks[0].text()).toBe('Users')
  })

  it('shows user email when authenticated', async () => {
    await router.push('/dashboard')
    await router.isReady()
    const auth = useAuthStore()
    auth.setUser({ id: '1', email: 'admin@test.com', name: 'Admin', role: 'admin', created_at: '', updated_at: '' })
    const wrapper = mount(WmsSidebar, {
      global: { plugins: [router] },
    })
    expect(wrapper.text()).toContain('admin@test.com')
  })

  it('renders theme toggle', async () => {
    await router.push('/dashboard')
    await router.isReady()
    const wrapper = mount(WmsSidebar, {
      global: { plugins: [router] },
    })
    expect(wrapper.findComponent({ name: 'WmsThemeToggle' }).exists()).toBe(true)
  })

  it('emits logout event when logout clicked', async () => {
    await router.push('/dashboard')
    await router.isReady()
    const wrapper = mount(WmsSidebar, {
      global: { plugins: [router] },
    })
    await wrapper.find('.wms-sidebar__logout').trigger('click')
    expect(wrapper.emitted('logout')).toBeTruthy()
    expect(wrapper.emitted('logout')).toHaveLength(1)
  })
})
