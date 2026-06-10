import type { Meta, StoryObj } from '@storybook/vue3'
import WmsSidebar from './WmsSidebar.vue'
import { useAuthStore } from '@/stores/auth'
import { createRouter, createWebHistory } from 'vue-router'
import { h } from 'vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', redirect: '/dashboard' },
    { path: '/dashboard', name: 'Dashboard', component: { template: '<div>Dashboard</div>' } },
    { path: '/users', name: 'Users', component: { template: '<div>Users</div>' } },
    { path: '/users/create', name: 'CreateUser', component: { template: '<div>Create</div>' } },
    { path: '/customers', name: 'Customers', component: { template: '<div>Customers</div>' } },
    { path: '/customers/invite', name: 'InviteCustomer', component: { template: '<div>Invite</div>' } },
  ],
})

const meta: Meta<typeof WmsSidebar> = {
  title: 'Layout/WmsSidebar',
  component: WmsSidebar,
  tags: ['autodocs'],
  decorators: [
    () => ({
      router,
      render() {
        return h('div', { style: 'min-height: 400px' }, [h(WmsSidebar)])
      },
    }),
  ],
}

export default meta
type Story = StoryObj<typeof WmsSidebar>

export const Default: Story = {}

export const WithActiveLink: Story = {
  play: async () => {
    await router.push('/users')
  },
}

export const DarkMode: Story = {
  parameters: {
    themes: { default: 'dark' },
  },
}

export const WithUser: Story = {
  play: () => {
    const auth = useAuthStore()
    auth.setUser({ id: '1', email: 'admin@example.com', name: 'Admin', role: 'admin', created_at: '', updated_at: '' })
  },
}
