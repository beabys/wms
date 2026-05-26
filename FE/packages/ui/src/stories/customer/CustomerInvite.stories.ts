import type { Meta, StoryObj } from '@storybook/vue3'
import CustomerInvite from '../../components/customer/CustomerInvite.vue'

const meta: Meta<typeof CustomerInvite> = {
  title: 'Customer/CustomerInvite',
  component: CustomerInvite,
  tags: ['autodocs']
}

export default meta
type Story = StoryObj<typeof CustomerInvite>

export const Default: Story = {
  args: {
    inviteLink: 'https://wms.example.com/sign-up?token=abc123xyz',
    expiresAt: '2026-06-26'
  }
}
