import type { Meta, StoryObj } from '@storybook/vue3'
import WmsSidebar from './WmsSidebar.vue'
import { useAuthStore } from '@/stores/auth'
import { useCustomerStore } from '@/stores/customer'
import { setActivePinia, createPinia } from 'pinia'

const meta: Meta<typeof WmsSidebar> = {
  title: 'Layout/WmsSidebar',
  component: WmsSidebar,
  tags: ['autodocs'],
  decorators: [
    () => ({
      template: '<div style="display: flex;"><story /></div>',
    }),
  ],
  parameters: {
    layout: 'fullscreen',
  },
}

export default meta
type Story = StoryObj<typeof WmsSidebar>

export const Default: Story = {
  render: () => ({
    components: { WmsSidebar },
    template: '<WmsSidebar @logout="() => {}" />',
  }),
}

export const Authenticated: Story = {
  render: () => {
    setActivePinia(createPinia())
    const store = useAuthStore()
    store.setUser({ id: '1', email: 'customer@example.com', name: 'Test User', role: 'customer', created_at: '', updated_at: '' })
    const customerStore = useCustomerStore()
    customerStore.customer = { id: '1', company_name: 'Acme', email: 'customer@example.com', phone: '', vat_number: '', status: 'active', created_at: '' }
    return {
      components: { WmsSidebar },
      template: '<WmsSidebar @logout="() => {}" />',
    }
  },
}

export const PendingApproval: Story = {
  render: () => {
    setActivePinia(createPinia())
    const store = useAuthStore()
    store.setUser({ id: '1', email: 'pending@example.com', name: 'Pending User', role: 'customer', created_at: '', updated_at: '' })
    const customerStore = useCustomerStore()
    customerStore.customer = { id: '1', company_name: 'Acme', email: 'pending@example.com', phone: '', vat_number: '', status: 'pending', created_at: '' }
    return {
      components: { WmsSidebar },
      template: '<WmsSidebar @logout="() => {}" />',
    }
  },
}
