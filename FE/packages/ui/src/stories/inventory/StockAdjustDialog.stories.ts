import type { Meta, StoryObj } from '@storybook/vue3'
import StockAdjustDialog from '../../components/inventory/StockAdjustDialog.vue'

const meta: Meta<typeof StockAdjustDialog> = {
  title: 'Inventory/StockAdjustDialog',
  component: StockAdjustDialog,
  tags: ['autodocs']
}

export default meta
type Story = StoryObj<typeof StockAdjustDialog>

export const Default: Story = {
  args: {
    productSku: 'SKU-001',
    productName: 'Widget Alpha',
    currentQuantity: 100,
    loading: false,
  }
}
