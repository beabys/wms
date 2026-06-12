import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import UserListView from './UserListView.vue'
import { useUserStore } from '@/stores/users'
import type { UserResponse } from '@/api/types'

const sampleUsers: UserResponse[] = [
  { id: '1', email: 'alice@test.com', name: 'Alice', role: 'admin', created_at: '2024-01-01T00:00:00Z', updated_at: '' },
]

const routes = [
  { path: '/users', name: 'Users', component: { template: '<div>Users</div>' } },
  { path: '/users/create', name: 'CreateUser', component: { template: '<div>Create</div>' } },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

const emptyListResponse = { users: [], pagination: { page: 1, page_size: 10, total_items: 0 } }

describe('UserListView', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.restoreAllMocks()
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 200,
      ok: true,
      json: async () => ({ success: true, data: emptyListResponse }),
    }))
  })

  it('renders title and create button', async () => {
    const wrapper = mount(UserListView, {
      global: { plugins: [router] },
    })
    await wrapper.vm.$nextTick()

    expect(wrapper.find('h1').text()).toBe('Users')
    expect(wrapper.text()).toContain('Create User')
  })

  it('shows loading state', async () => {
    const store = useUserStore()
    store.loading = true

    const wrapper = mount(UserListView, {
      global: { plugins: [router] },
    })
    await wrapper.vm.$nextTick()

    expect(wrapper.text()).toContain('Loading users...')
  })

  it('shows empty state when no users', async () => {
    const wrapper = mount(UserListView, {
      global: { plugins: [router] },
    })
    await wrapper.vm.$nextTick()

    expect(wrapper.findComponent({ name: 'UserTable' }).exists()).toBe(true)
  })

  it('renders users when data loaded', async () => {
    const fetchMock = vi.fn().mockResolvedValue({
      status: 200,
      ok: true,
      json: async () => ({
        success: true,
        data: { users: sampleUsers, pagination: { page: 1, page_size: 10, total_items: 1 } },
      }),
    })
    vi.stubGlobal('fetch', fetchMock)

    const wrapper = mount(UserListView, {
      global: { plugins: [router] },
    })
    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 50))

    expect(wrapper.text()).toContain('alice@test.com')
  })

  it('navigates to create page on button click', async () => {
    const pushSpy = vi.spyOn(router, 'push')
    const wrapper = mount(UserListView, {
      global: { plugins: [router] },
    })
    await wrapper.vm.$nextTick()

    const buttons = wrapper.findAll('button')
    const createBtn = buttons.find(b => b.text().includes('Create User'))
    expect(createBtn).toBeTruthy()
    await createBtn!.trigger('click')
    expect(pushSpy).toHaveBeenCalledWith('/users/create')
  })
})
