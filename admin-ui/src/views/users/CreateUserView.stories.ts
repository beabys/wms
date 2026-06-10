import type { Meta, StoryObj } from '@storybook/vue3'
import { createPinia, setActivePinia } from 'pinia'
import CreateUserView from './CreateUserView.vue'
import { useUserStore } from '@/stores/users'

const meta: Meta<typeof CreateUserView> = {
  title: 'Views/CreateUserView',
  component: CreateUserView,
  tags: ['autodocs'],
}

export default meta
type Story = StoryObj<typeof CreateUserView>

export const Default: Story = {
  render: () => ({
    components: { CreateUserView },
    template: '<CreateUserView />',
  }),
}

export const ValidationErrors: Story = {
  render: () => ({
    components: { CreateUserView },
    template: '<CreateUserView />',
  }),
}

export const Submitting: Story = {
  render: () => ({
    components: { CreateUserView },
    setup: () => {
      const pinia = createPinia()
      setActivePinia(pinia)
      const store = useUserStore()
      store.loading = true
      return {}
    },
    template: '<CreateUserView />',
  }),
}

export const DarkMode: Story = {
  render: () => ({
    components: { CreateUserView },
    template: '<CreateUserView />',
  }),
  parameters: {
    themes: { theme: 'dark' },
  },
}
