import type { Meta, StoryObj } from '@storybook/vue3'
import WmsLoginForm from './WmsLoginForm.vue'

const meta: Meta<typeof WmsLoginForm> = {
  title: 'Auth/WmsLoginForm',
  component: WmsLoginForm,
  tags: ['autodocs'],
}

export default meta
type Story = StoryObj<typeof WmsLoginForm>

export const Default: Story = {}
export const Loading: Story = { args: { loading: true } }
export const WithError: Story = { args: { error: 'Invalid email or password.' } }
export const ValidationErrors: Story = {
  play: async ({ canvasElement }) => {
    const form = canvasElement.querySelector('form')
    form?.requestSubmit()
  },
}
