import type { Meta, StoryObj } from '@storybook/vue3'
import InviteCustomerModal from './InviteCustomerModal.vue'

const meta: Meta<typeof InviteCustomerModal> = {
  title: 'Views/InviteCustomerModal',
  component: InviteCustomerModal,
  tags: ['autodocs'],
}

export default meta
type Story = StoryObj<typeof InviteCustomerModal>

export const Default: Story = {
  args: {
    visible: true,
  },
}

export const DarkMode: Story = {
  args: {
    visible: true,
  },
  parameters: {
    themes: { theme: 'dark' },
  },
}
