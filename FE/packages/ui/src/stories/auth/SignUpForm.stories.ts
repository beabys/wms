import type { Meta, StoryObj } from '@storybook/vue3'
import SignUpForm from '../../components/auth/SignUpForm.vue'

const meta: Meta<typeof SignUpForm> = {
  title: 'Auth/SignUpForm',
  component: SignUpForm,
  tags: ['autodocs'],
  argTypes: {
    validating: { control: 'boolean' },
    success: { control: 'boolean' },
    error: { control: 'text' }
  }
}

export default meta
type Story = StoryObj<typeof SignUpForm>

export const Default: Story = {
  args: {
    validating: false,
    success: false,
    error: null
  }
}

export const Validating: Story = {
  args: {
    validating: true,
    success: false,
    error: null
  }
}

export const Success: Story = {
  args: {
    validating: false,
    success: true,
    error: null
  }
}

export const Error: Story = {
  args: {
    validating: false,
    success: false,
    error: 'This email is already registered.'
  }
}
