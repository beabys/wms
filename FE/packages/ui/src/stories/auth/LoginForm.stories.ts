import type { Meta, StoryObj } from '@storybook/vue3'
import LoginForm from '../../components/auth/LoginForm.vue'

const meta: Meta<typeof LoginForm> = {
  title: 'Auth/LoginForm',
  component: LoginForm,
  tags: ['autodocs'],
  argTypes: {
    loading: { control: 'boolean' },
    error: { control: 'text' }
  }
}

export default meta
type Story = StoryObj<typeof LoginForm>

export const Default: Story = {
  args: {
    loading: false,
    error: null
  }
}

export const Loading: Story = {
  args: {
    loading: true,
    error: null
  }
}

export const Error: Story = {
  args: {
    loading: false,
    error: 'Invalid email or password. Please try again.'
  }
}
