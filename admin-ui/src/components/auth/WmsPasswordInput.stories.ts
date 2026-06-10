import type { Meta, StoryObj } from '@storybook/vue3'
import WmsPasswordInput from './WmsPasswordInput.vue'

const meta: Meta<typeof WmsPasswordInput> = {
  title: 'Auth/WmsPasswordInput',
  component: WmsPasswordInput,
  tags: ['autodocs'],
}

export default meta
type Story = StoryObj<typeof WmsPasswordInput>

export const Default: Story = {
  args: {
    label: 'Password',
    placeholder: 'Enter your password',
  },
}

export const Filled: Story = {
  args: {
    label: 'Password',
    modelValue: 'mysecretpassword',
  },
}
