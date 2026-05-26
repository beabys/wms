import type { Meta, StoryObj } from '@storybook/vue3'
import InboundForm, { type InboundProduct } from '../../components/inbound/InboundForm.vue'

const meta: Meta<typeof InboundForm> = {
  title: 'Inbound/InboundForm',
  component: InboundForm,
  tags: ['autodocs']
}

export default meta
type Story = StoryObj<typeof InboundForm>

export const Empty: Story = {
  args: {
    products: [],
    errors: undefined,
    loading: false
  }
}

const twoProducts: InboundProduct[] = [
  { id: 'p1', sku: 'SKU-001', name: 'Widget Alpha', quantity: 100, length: 30, width: 20, height: 15, weight: 2.5, packaging: 'box' },
  { id: 'p2', sku: 'SKU-002', name: 'Gadget Beta', quantity: 50, length: 40, width: 30, height: 25, weight: 5.0, packaging: 'pallet' }
]

export const TwoProductsFilled: Story = {
  args: {
    products: twoProducts,
    errors: undefined,
    loading: false
  }
}

export const ValidationErrors: Story = {
  args: {
    products: twoProducts,
    errors: {
      _general: 'Please fix the errors below before submitting.',
      p1: 'SKU already exists in another active shipment.',
      p2: 'Quantity exceeds available storage capacity.'
    },
    loading: false
  }
}
