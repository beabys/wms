import type { Meta, StoryObj } from '@storybook/vue3'
import { createPinia, setActivePinia } from 'pinia'
import UserListView from './UserListView.vue'
import { useUserStore } from '@/stores/users'

const meta: Meta<typeof UserListView> = {
  title: 'Views/UserListView',
  component: UserListView,
  tags: ['autodocs'],
}

export default meta
type Story = StoryObj<typeof UserListView>

export const Default: Story = {
  render: () => ({
    components: { UserListView },
    template: '<UserListView />',
  }),
}

export const Loading: Story = {
  render: () => ({
    components: { UserListView },
    setup: () => {
      const pinia = createPinia()
      setActivePinia(pinia)
      const store = useUserStore()
      store.loading = true
      return {}
    },
    template: '<UserListView />',
  }),
}

export const Empty: Story = {
  render: () => ({
    components: { UserListView },
    setup: () => {
      const pinia = createPinia()
      setActivePinia(pinia)
      return {}
    },
    template: '<UserListView />',
  }),
}

export const WithPagination: Story = {
  render: () => ({
    components: { UserListView },
    setup: () => {
      const pinia = createPinia()
      setActivePinia(pinia)
      const store = useUserStore()
      store.users = [
        { id: '1', email: 'alice@example.com', name: 'Alice Johnson', role: 'admin', created_at: '2024-01-15T10:30:00Z', updated_at: '' },
        { id: '2', email: 'bob@example.com', name: 'Bob Smith', role: 'warehouse_staff', created_at: '2024-02-20T14:00:00Z', updated_at: '' },
      ]
      store.pagination = { page: 1, page_size: 10, total_items: 25 }
      store.loading = false
      return {}
    },
    template: '<UserListView />',
  }),
}

export const DarkMode: Story = {
  render: () => ({
    components: { UserListView },
    setup: () => {
      const pinia = createPinia()
      setActivePinia(pinia)
      const store = useUserStore()
      store.users = [
        { id: '1', email: 'alice@example.com', name: 'Alice Johnson', role: 'admin', created_at: '2024-01-15T10:30:00Z', updated_at: '' },
      ]
      store.pagination = { page: 1, page_size: 10, total_items: 1 }
      store.loading = false
      return {}
    },
    template: '<UserListView />',
  }),
  parameters: {
    themes: { theme: 'dark' },
  },
}
