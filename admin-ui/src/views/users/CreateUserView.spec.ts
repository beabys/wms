import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import CreateUserView from './CreateUserView.vue'

const routes = [
  { path: '/users', name: 'Users', component: { template: '<div>Users</div>' } },
  { path: '/users/create', name: 'CreateUser', component: { template: '<div>Create</div>' } },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

describe('CreateUserView', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.restoreAllMocks()
  })

  it('renders title and form', () => {
    const wrapper = mount(CreateUserView, {
      global: { plugins: [router] },
    })
    expect(wrapper.find('h1').text()).toBe('Create User')
    expect(wrapper.findComponent({ name: 'UserForm' }).exists()).toBe(true)
  })

  it('renders cancel button', () => {
    const wrapper = mount(CreateUserView, {
      global: { plugins: [router] },
    })
    expect(wrapper.text()).toContain('Cancel')
  })

  it('navigates back on cancel click', async () => {
    const pushSpy = vi.spyOn(router, 'push')
    const wrapper = mount(CreateUserView, {
      global: { plugins: [router] },
    })
    await wrapper.vm.$nextTick()

    const cancelBtn = wrapper.findAll('button').find(b => b.text().includes('Cancel'))
    expect(cancelBtn).toBeTruthy()
    await cancelBtn!.trigger('click')
    expect(pushSpy).toHaveBeenCalledWith('/users')
  })

  it('calls store.createUser and navigates on submit', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({
        success: true,
        data: { id: '3', email: 'new@test.com', name: 'New', role: 'admin', created_at: '', updated_at: '' },
      }),
    }))

    const pushSpy = vi.spyOn(router, 'push')
    const wrapper = mount(CreateUserView, {
      global: { plugins: [router] },
    })
    await wrapper.vm.$nextTick()

    const form = wrapper.findComponent({ name: 'UserForm' })
    expect(form.exists()).toBe(true)

    form.vm.$emit('submit')
    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 50))

    expect(pushSpy).toHaveBeenCalledWith('/users')
  })
})
