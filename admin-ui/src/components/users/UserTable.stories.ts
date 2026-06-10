import type { Meta, StoryObj } from '@storybook/vue3'
import UserTable from './UserTable.vue'
import type { UserResponse } from '@/api/types'

const sampleUsers: UserResponse[] = [
  { id: '1', email: 'alice@example.com', name: 'Alice Johnson', role: 'admin', created_at: '2024-01-15T10:30:00Z', updated_at: '2024-01-15T10:30:00Z' },
  { id: '2', email: 'bob@example.com', name: 'Bob Smith', role: 'warehouse_staff', created_at: '2024-02-20T14:00:00Z', updated_at: '2024-02-20T14:00:00Z' },
  { id: '3', email: 'carol@example.com', name: 'Carol Davis', role: 'billing_manager', created_at: '2024-03-10T09:15:00Z', updated_at: '2024-03-10T09:15:00Z' },
]

const meta: Meta<typeof UserTable> = {
  title: 'Users/UserTable',
  component: UserTable,
  tags: ['autodocs'],
}

export default meta
type Story = StoryObj<typeof UserTable>

export const Default: Story = {
  args: {
    users: sampleUsers,
    loading: false,
  },
}

export const Loading: Story = {
  args: {
    users: [],
    loading: true,
  },
}

export const Empty: Story = {
  args: {
    users: [],
    loading: false,
  },
}

export const WithData: Story = {
  args: {
    users: [sampleUsers[0]],
    loading: false,
  },
}

export const DarkMode: Story = {
  args: {
    users: sampleUsers,
    loading: false,
  },
  parameters: {
    themes: { theme: 'dark' },
  },
}
