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
  args: { label: 'Password', placeholder: 'Enter your password' },
}

export const Filled: Story = {
  args: { label: 'Password', modelValue: 'mysecretpassword' },
}

export const WithStrength: Story = {
  args: { label: 'Password', modelValue: 'StrongP@ss1', showStrength: true, strength: 4, strengthLabel: 'Strong' },
}

export const WeakStrength: Story = {
  args: { label: 'Password', modelValue: 'weak', showStrength: true, strength: 1, strengthLabel: 'Weak' },
}
