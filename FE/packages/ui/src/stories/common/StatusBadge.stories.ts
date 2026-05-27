import type { Meta, StoryObj } from '@storybook/vue3'
import StatusBadge from '../../components/common/StatusBadge.vue'

const meta: Meta<typeof StatusBadge> = {
  title: 'Common/StatusBadge',
  component: StatusBadge,
  tags: ['autodocs'],
  argTypes: {
    status: {
      control: 'select',
      options: ['pending', 'approved', 'flagged', 'held', 'shipped', 'completed']
    }
  }
}

export default meta
type Story = StoryObj<typeof StatusBadge>

export const Pending: Story = { args: { status: 'pending' } }
export const Approved: Story = { args: { status: 'approved' } }
export const Flagged: Story = { args: { status: 'flagged' } }
export const Held: Story = { args: { status: 'held' } }
export const Shipped: Story = { args: { status: 'shipped' } }
export const Completed: Story = { args: { status: 'completed' } }
