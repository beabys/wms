import type { Meta, StoryObj } from '@storybook/vue3'
import WmsThemeToggle from './WmsThemeToggle.vue'

const meta: Meta<typeof WmsThemeToggle> = {
  title: 'Layout/WmsThemeToggle',
  component: WmsThemeToggle,
  tags: ['autodocs'],
}

export default meta
type Story = StoryObj<typeof WmsThemeToggle>

export const Default: Story = {}
