import type { Meta, StoryObj } from '@storybook/vue3'
import LowStockAlert, { type LowStockAlertItem } from '../../components/inventory/LowStockAlert.vue'

const meta: Meta<typeof LowStockAlert> = {
  title: 'Inventory/LowStockAlert',
  component: LowStockAlert,
  tags: ['autodocs']
}

export default meta
type Story = StoryObj<typeof LowStockAlert>

const sampleAlerts: LowStockAlertItem[] = [
  { product_id: 'p1', product_sku: 'SKU-002', product_name: 'Gadget Beta', current_quantity: 5, threshold: 10, status: 'critical' },
  { product_id: 'p2', product_sku: 'SKU-003', product_name: 'Component Gamma', current_quantity: 0, threshold: 10, status: 'critical' },
  { product_id: 'p3', product_sku: 'SKU-006', product_name: 'Nutrient Zeta', current_quantity: 14, threshold: 10, status: 'warning' },
  { product_id: 'p4', product_sku: 'SKU-007', product_name: 'Powder Eta', current_quantity: 18, threshold: 10, status: 'warning' },
]

export const Empty: Story = {
  args: {
    alerts: [],
  }
}

export const WithAlerts: Story = {
  args: {
    alerts: sampleAlerts,
  }
}
