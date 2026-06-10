import type { Meta, StoryObj } from '@storybook/vue3'
import WmsRegisterForm from './WmsRegisterForm.vue'

const meta: Meta<typeof WmsRegisterForm> = {
  title: 'Auth/WmsRegisterForm',
  component: WmsRegisterForm,
  tags: ['autodocs'],
}

export default meta
type Story = StoryObj<typeof WmsRegisterForm>

export const Default: Story = {
  args: { inviteToken: 'invite-abc-123' },
}

export const WithNoToken: Story = {}
export const Loading: Story = { args: { loading: true } }
export const WithInviteToken: Story = {
  args: { inviteToken: 'invite-xyz-789' },
}
export const ValidationErrors: Story = {
  play: async ({ canvasElement }) => {
    const form = canvasElement.querySelector('form')
    form?.requestSubmit()
  },
}
