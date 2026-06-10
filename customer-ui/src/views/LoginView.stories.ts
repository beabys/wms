import type { Meta, StoryObj } from '@storybook/vue3'
import LoginView from './LoginView.vue'

const meta: Meta<typeof LoginView> = {
  title: 'Views/LoginView',
  component: LoginView,
  tags: ['autodocs'],
}

export default meta
type Story = StoryObj<typeof LoginView>

export const Default: Story = {}
export const DarkMode: Story = {
  parameters: { themes: { theme: 'dark' } },
}
