import type { Meta, StoryObj } from '@storybook/vue3'
import ProductForm from '../../components/inventory/ProductForm.vue'

const meta: Meta<typeof ProductForm> = {
  title: 'Inventory/ProductForm',
  component: ProductForm,
  tags: ['autodocs']
}

export default meta
type Story = StoryObj<typeof ProductForm>

export const CreateMode: Story = {
  args: {
    mode: 'create',
    loading: false,
  }
}

export const EditMode: Story = {
  args: {
    mode: 'edit',
    initial: {
      sku: 'SKU-001',
      name: 'Widget Alpha',
      description: 'A sample widget',
      category: 'Electronics',
      unit: 'piece',
      weight_kg: 1.5,
      low_stock_threshold: 10,
    },
    loading: false,
  }
}
