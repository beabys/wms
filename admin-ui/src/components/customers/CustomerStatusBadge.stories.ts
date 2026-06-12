import type { Meta, StoryObj } from '@storybook/vue3'
import CustomerStatusBadge from './CustomerStatusBadge.vue'

const meta: Meta<typeof CustomerStatusBadge> = {
  title: 'Customers/CustomerStatusBadge',
  component: CustomerStatusBadge,
  tags: ['autodocs'],
  argTypes: {
    status: { control: 'select', options: ['pending', 'active', 'rejected', 'suspended'] },
  },
}

export default meta
type Story = StoryObj<typeof CustomerStatusBadge>

export const Pending: Story = {
  args: { status: 'pending' },
}

export const Active: Story = {
  args: { status: 'active' },
}

export const Rejected: Story = {
  args: { status: 'rejected' },
}

export const Suspended: Story = {
  args: { status: 'suspended' },
}

export const DarkMode: Story = {
  args: { status: 'pending' },
  parameters: {
    themes: { theme: 'dark' },
  },
}
