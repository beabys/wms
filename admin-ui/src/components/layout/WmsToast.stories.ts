import type { Meta, StoryObj } from '@storybook/vue3'
import WmsToast from './WmsToast.vue'
import { useToastStore } from '@/stores/toast'

const meta: Meta<typeof WmsToast> = {
  title: 'Layout/WmsToast',
  component: WmsToast,
  tags: ['autodocs'],
  decorators: [() => ({ template: '<div style="min-height: 200px"><story /></div>' })],
}

export default meta
type Story = StoryObj<typeof WmsToast>

export const Success: Story = {
  play: () => {
    const store = useToastStore()
    store.addToast('Operation completed successfully', 'success')
  },
}

export const Error: Story = {
  play: () => {
    const store = useToastStore()
    store.addToast('Something went wrong', 'error')
  },
}

export const Warning: Story = {
  play: () => {
    const store = useToastStore()
    store.addToast('This action cannot be undone', 'warning')
  },
}

export const Info: Story = {
  play: () => {
    const store = useToastStore()
    store.addToast('New updates available', 'info')
  },
}

export const Multiple: Story = {
  play: () => {
    const store = useToastStore()
    store.addToast('Success!', 'success')
    store.addToast('Warning message', 'warning')
    store.addToast('Info message', 'info')
  },
}

export const DarkMode: Story = {
  parameters: {
    themes: { default: 'dark' },
  },
  play: () => {
    const store = useToastStore()
    store.addToast('Operation completed successfully', 'success')
    store.addToast('Something went wrong', 'error')
  },
}
