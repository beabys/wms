import type { Meta, StoryObj } from '@storybook/vue3'
import WmsLayout from './WmsLayout.vue'

const meta: Meta<typeof WmsLayout> = {
  title: 'Layout/WmsLayout',
  component: WmsLayout,
  tags: ['autodocs'],
  parameters: {
    layout: 'fullscreen',
  },
}

export default meta
type Story = StoryObj<typeof WmsLayout>

export const Default: Story = {
  render: () => ({
    components: { WmsLayout },
    template: '<WmsLayout><div style="padding: 2rem"><h1>Dashboard Content</h1><p>Welcome to the customer portal.</p></div></WmsLayout>',
  }),
}
