import type { Meta, StoryObj } from '@storybook/vue3'
import BinLocationForm from '../../components/inventory/BinLocationForm.vue'

const meta: Meta<typeof BinLocationForm> = {
  title: 'Inventory/BinLocationForm',
  component: BinLocationForm,
  tags: ['autodocs']
}

export default meta
type Story = StoryObj<typeof BinLocationForm>

export const Default: Story = {
  args: {
    loading: false,
  }
}
