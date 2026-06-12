import type { Meta, StoryObj } from '@storybook/vue3'
import InviteListFilters from './InviteListFilters.vue'

const meta: Meta<typeof InviteListFilters> = {
  title: 'Components/Customers/InviteListFilters',
  component: InviteListFilters,
  tags: ['autodocs'],
}

export default meta
type Story = StoryObj<typeof InviteListFilters>

export const Default: Story = {
  args: {
    modelValue: {},
  },
}

export const WithStatusFilter: Story = {
  args: {
    modelValue: { status: 'pending' },
  },
}

export const WithExpiredFilter: Story = {
  args: {
    modelValue: { expired: true },
  },
}

export const WithDateRange: Story = {
  args: {
    modelValue: {
      created_after: 1700000000,
      created_before: 1710000000,
    },
  },
}

export const DarkMode: Story = {
  args: {
    modelValue: {},
  },
  parameters: {
    themes: { theme: 'dark' },
  },
}
