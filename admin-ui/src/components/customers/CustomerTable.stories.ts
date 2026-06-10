import type { Meta, StoryObj } from '@storybook/vue3'
import CustomerTable from './CustomerTable.vue'
import type { CustomerResponse } from '@/api/types'

const pendingCustomers: CustomerResponse[] = [
  { id: '1', company_name: 'Acme Corp', email: 'acme@test.com', phone: '', vat_number: 'VAT001', status: 'pending', created_at: '2024-01-15' },
  { id: '2', company_name: 'Globex Inc', email: 'globex@test.com', phone: '', vat_number: 'VAT002', status: 'pending', created_at: '2024-02-01' },
]

const mixedCustomers: CustomerResponse[] = [
  { id: '1', company_name: 'Acme Corp', email: 'acme@test.com', phone: '', vat_number: 'VAT001', status: 'active', created_at: '2024-01-15' },
  { id: '2', company_name: 'Globex Inc', email: 'globex@test.com', phone: '', vat_number: 'VAT002', status: 'pending', created_at: '2024-02-01' },
  { id: '3', company_name: 'Initech', email: 'initech@test.com', phone: '', vat_number: 'VAT003', status: 'rejected', created_at: '2024-03-10' },
  { id: '4', company_name: 'Umbrella Co', email: 'umbrella@test.com', phone: '', vat_number: 'VAT004', status: 'suspended', created_at: '2024-03-15' },
]

const meta: Meta<typeof CustomerTable> = {
  title: 'Customers/CustomerTable',
  component: CustomerTable,
  tags: ['autodocs'],
}

export default meta
type Story = StoryObj<typeof CustomerTable>

export const Default: Story = {
  args: {
    customers: mixedCustomers,
    loading: false,
  },
}

export const Loading: Story = {
  args: {
    customers: [],
    loading: true,
  },
}

export const Empty: Story = {
  args: {
    customers: [],
    loading: false,
  },
}

export const WithPendingCustomers: Story = {
  args: {
    customers: pendingCustomers,
    loading: false,
  },
}

export const DarkMode: Story = {
  args: {
    customers: mixedCustomers,
    loading: false,
  },
  parameters: {
    themes: { theme: 'dark' },
  },
}
