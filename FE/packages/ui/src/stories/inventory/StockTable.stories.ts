import type { Meta, StoryObj } from '@storybook/vue3'
import StockTable, { type StockTableEntry } from '../../components/inventory/StockTable.vue'

const meta: Meta<typeof StockTable> = {
  title: 'Inventory/StockTable',
  component: StockTable,
  tags: ['autodocs']
}

export default meta
type Story = StoryObj<typeof StockTable>

const sampleEntries: StockTableEntry[] = [
  { id: 'st-1', product_sku: 'SKU-001', product_name: 'Widget Alpha', bin_location: 'A-01-R01-S1', quantity: 100, reserved_quantity: 5, status: 'available', lot_number: 'LOT-001', expiry_date: '2027-01-01', low_stock_threshold: 10 },
  { id: 'st-2', product_sku: 'SKU-002', product_name: 'Gadget Beta', bin_location: 'A-01-R01-S2', quantity: 5, reserved_quantity: 2, status: 'available', lot_number: 'LOT-002', expiry_date: '2026-12-01', low_stock_threshold: 10 },
  { id: 'st-3', product_sku: 'SKU-003', product_name: 'Component Gamma', bin_location: 'B-02-R03-S1', quantity: 0, reserved_quantity: 0, status: 'inactive', lot_number: '', expiry_date: null, low_stock_threshold: 10 },
  { id: 'st-4', product_sku: 'SKU-004', product_name: 'Device Delta', bin_location: 'C-01-R02-S3', quantity: 15, reserved_quantity: 10, status: 'reserved', lot_number: 'LOT-004', expiry_date: '2026-06-15', low_stock_threshold: 10 },
  { id: 'st-5', product_sku: 'SKU-005', product_name: 'Material Epsilon', bin_location: 'A-03-R01-S2', quantity: 200, reserved_quantity: 0, status: 'available', lot_number: 'LOT-005', expiry_date: null, low_stock_threshold: 50 },
]

export const Empty: Story = {
  args: {
    stockEntries: [],
    loading: false,
  }
}

export const WithData: Story = {
  args: {
    stockEntries: sampleEntries,
    loading: false,
  }
}

export const Loading: Story = {
  args: {
    stockEntries: [],
    loading: true,
  }
}
