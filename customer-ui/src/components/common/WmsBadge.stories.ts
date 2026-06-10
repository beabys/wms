import type { Meta, StoryObj } from '@storybook/vue3'
import WmsBadge from './WmsBadge.vue'

const meta: Meta<typeof WmsBadge> = {
  title: 'Common/WmsBadge',
  component: WmsBadge,
  tags: ['autodocs'],
  argTypes: {
    variant: { control: 'select', options: ['success', 'warning', 'error', 'info', 'neutral'] },
    size: { control: 'select', options: ['sm', 'md'] },
  },
}

export default meta
type Story = StoryObj<typeof WmsBadge>

export const Success: Story = { args: { variant: 'success', default: 'Active' } }
export const Warning: Story = { args: { variant: 'warning', default: 'Pending' } }
export const Error: Story = { args: { variant: 'error', default: 'Blocked' } }
export const Info: Story = { args: { variant: 'info', default: 'New' } }
export const Neutral: Story = { args: { variant: 'neutral', default: 'Draft' } }
export const MediumSize: Story = { args: { variant: 'success', size: 'md', default: 'Verified' } }
