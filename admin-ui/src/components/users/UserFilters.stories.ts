import type { Meta, StoryObj } from '@storybook/vue3'
import { ref } from 'vue'
import UserFilters from './UserFilters.vue'

const meta: Meta<typeof UserFilters> = {
  title: 'Users/UserFilters',
  component: UserFilters,
  tags: ['autodocs'],
}

export default meta
type Story = StoryObj<typeof UserFilters>

export const Default: Story = {
  render: () => ({
    components: { UserFilters },
    setup: () => {
      const filters = ref<{ role?: string }>({})
      return { filters }
    },
    template: '<UserFilters v-model="filters" />',
  }),
}

export const WithRoleSelected: Story = {
  render: () => ({
    components: { UserFilters },
    setup: () => {
      const filters = ref<{ role?: string }>({ role: 'admin' })
      return { filters }
    },
    template: '<UserFilters v-model="filters" />',
  }),
}

export const DarkMode: Story = {
  render: () => ({
    components: { UserFilters },
    setup: () => {
      const filters = ref<{ role?: string }>({})
      return { filters }
    },
    template: '<UserFilters v-model="filters" />',
  }),
  parameters: {
    themes: { theme: 'dark' },
  },
}
