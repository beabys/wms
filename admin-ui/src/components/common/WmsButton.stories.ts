import type { Meta, StoryObj } from '@storybook/vue3'
import WmsButton from './WmsButton.vue'

const meta: Meta<typeof WmsButton> = {
  title: 'Common/WmsButton',
  component: WmsButton,
  tags: ['autodocs'],
  argTypes: {
    variant: { control: 'select', options: ['primary', 'secondary', 'ghost', 'danger'] },
    size: { control: 'select', options: ['sm', 'md', 'lg'] },
  },
}

export default meta
type Story = StoryObj<typeof WmsButton>

export const Primary: Story = {
  args: {
    variant: 'primary',
    size: 'md',
    default: 'Click me',
  },
}

export const Secondary: Story = {
  args: {
    variant: 'secondary',
    size: 'md',
    default: 'Cancel',
  },
}

export const Ghost: Story = {
  args: {
    variant: 'ghost',
    size: 'md',
    default: 'Cancel',
  },
}

export const Danger: Story = {
  args: {
    variant: 'danger',
    size: 'md',
    default: 'Delete',
  },
}

export const Loading: Story = {
  args: {
    variant: 'primary',
    size: 'md',
    loading: true,
    default: 'Loading...',
  },
}

export const Disabled: Story = {
  args: {
    variant: 'primary',
    size: 'md',
    disabled: true,
    default: 'Disabled',
  },
}

export const Small: Story = {
  args: {
    variant: 'primary',
    size: 'sm',
    default: 'Small',
  },
}

export const Large: Story = {
  args: {
    variant: 'primary',
    size: 'lg',
    default: 'Large',
  },
}
