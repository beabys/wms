import type { Meta, StoryObj } from '@storybook/vue3'
import RegisterView from './RegisterView.vue'

const meta: Meta<typeof RegisterView> = {
  title: 'Views/RegisterView',
  component: RegisterView,
  tags: ['autodocs'],
}

export default meta
type Story = StoryObj<typeof RegisterView>

export const WithToken: Story = {
  parameters: { vueRouter: { query: { token: 'invite-abc-123' } } },
}

export const NoToken: Story = {}

export const DarkMode: Story = {
  parameters: { themes: { theme: 'dark' } },
}
