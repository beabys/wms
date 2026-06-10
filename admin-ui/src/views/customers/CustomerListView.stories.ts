import type { Meta, StoryObj } from '@storybook/vue3'
import { createPinia, setActivePinia } from 'pinia'
import CustomerListView from './CustomerListView.vue'
import { useCustomerStore } from '@/stores/customers'
import type { CustomerResponse } from '@/api/types'

const meta: Meta<typeof CustomerListView> = {
  title: 'Views/CustomerListView',
  component: CustomerListView,
  tags: ['autodocs'],
}

export default meta
type Story = StoryObj<typeof CustomerListView>

const sampleCustomers: CustomerResponse[] = [
  { id: 'c1', company_name: 'ACME Corp', email: 'admin@acme.com', status: 'active', phone: '', vat_number: '', address: '', city: '', postal_code: '', country: '', created_at: '' },
  { id: 'c2', company_name: 'Globex Inc', email: 'admin@globex.com', status: 'pending', phone: '', vat_number: '', address: '', city: '', postal_code: '', country: '', created_at: '' },
  { id: 'c3', company_name: 'Initech', email: 'admin@initech.com', status: 'suspended', phone: '', vat_number: '', address: '', city: '', postal_code: '', country: '', created_at: '' },
]

export const Default: Story = {
  render: () => ({
    components: { CustomerListView },
    template: '<CustomerListView />',
  }),
}

export const Loading: Story = {
  render: () => ({
    components: { CustomerListView },
    setup: () => {
      const pinia = createPinia()
      setActivePinia(pinia)
      const store = useCustomerStore()
      store.loading = true
      store.customers = []
      return {}
    },
    template: '<CustomerListView />',
  }),
}

export const Empty: Story = {
  render: () => ({
    components: { CustomerListView },
    setup: () => {
      const pinia = createPinia()
      setActivePinia(pinia)
      const store = useCustomerStore()
      store.customers = []
      store.pagination = { page: 1, page_size: 10, total_items: 0 }
      store.loading = false
      return {}
    },
    template: '<CustomerListView />',
  }),
}

export const Error: Story = {
  render: () => ({
    components: { CustomerListView },
    setup: () => {
      const pinia = createPinia()
      setActivePinia(pinia)
      const store = useCustomerStore()
      store.customers = []
      store.error = 'Failed to fetch customers. Please try again.'
      store.loading = false
      return {}
    },
    template: '<CustomerListView />',
  }),
}

export const WithPagination: Story = {
  render: () => ({
    components: { CustomerListView },
    setup: () => {
      const pinia = createPinia()
      setActivePinia(pinia)
      const store = useCustomerStore()
      store.customers = sampleCustomers
      store.pagination = { page: 1, page_size: 10, total_items: 3 }
      store.loading = false
      return {}
    },
    template: '<CustomerListView />',
  }),
}

export const DarkMode: Story = {
  render: () => ({
    components: { CustomerListView },
    setup: () => {
      const pinia = createPinia()
      setActivePinia(pinia)
      const store = useCustomerStore()
      store.customers = sampleCustomers
      store.pagination = { page: 1, page_size: 10, total_items: 3 }
      store.loading = false
      return {}
    },
    template: '<CustomerListView />',
  }),
  parameters: {
    themes: { theme: 'dark' },
  },
}
