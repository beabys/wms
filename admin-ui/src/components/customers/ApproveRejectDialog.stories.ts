import type { Meta, StoryObj } from '@storybook/vue3'
import ApproveRejectDialog from './ApproveRejectDialog.vue'

const customer = {
  id: '1',
  company_name: 'Acme Corp',
  email: 'acme@test.com',
  phone: '123456789',
  vat_number: 'VAT123',
  status: 'pending',
  created_at: '2024-01-01',
}

const meta: Meta<typeof ApproveRejectDialog> = {
  title: 'Customers/ApproveRejectDialog',
  component: ApproveRejectDialog,
  tags: ['autodocs'],
  argTypes: {
    action: { control: 'select', options: ['approve', 'reject', 'suspend'] },
  },
}

export default meta
type Story = StoryObj<typeof ApproveRejectDialog>

export const Approve: Story = {
  args: {
    visible: true,
    customer,
    action: 'approve',
  },
}

export const Reject: Story = {
  args: {
    visible: true,
    customer,
    action: 'reject',
  },
}

export const Suspend: Story = {
  args: {
    visible: true,
    customer,
    action: 'suspend',
  },
}

export const DarkMode: Story = {
  args: {
    visible: true,
    customer,
    action: 'reject',
  },
  parameters: {
    themes: { theme: 'dark' },
  },
}
