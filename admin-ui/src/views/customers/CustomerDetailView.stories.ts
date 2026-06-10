import type { Meta, StoryObj } from '@storybook/vue3'
import { createPinia, setActivePinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import { useCustomerDetailStore } from '@/stores/customerDetail'
import CustomerDetailView from './CustomerDetailView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/customers/:id', name: 'CustomerDetail', component: CustomerDetailView },
  ],
})

function createStore(overrides: Record<string, any> = {}) {
  const pinia = createPinia()
  setActivePinia(pinia)
  const store = useCustomerDetailStore()
  Object.assign(store, {
    customer: null,
    auditLogs: [],
    loading: false,
    error: '',
    ...overrides,
  })
  return pinia
}

const baseCustomer = {
  id: '1',
  company_name: 'Acme Corp',
  email: 'acme@test.com',
  phone: '123-456-7890',
  vat_number: 'VAT-US-12345',
  address: '123 Main Street',
  city: 'Springfield',
  postal_code: '12345',
  country: 'United States',
  created_at: '2024-01-15',
}

const auditEntries = [
  { id: 'a1', customer_id: '1', action: 'approved', performed_by: 'admin@wms.com', details: 'Approved by admin', created_at: 1700000000 },
  { id: 'a2', customer_id: '1', action: 'suspended', performed_by: 'admin@wms.com', details: 'Violation of terms', created_at: 1700100000 },
  { id: 'a3', customer_id: '1', action: 'restored', performed_by: 'admin@wms.com', details: 'Issue resolved', created_at: 1700200000 },
]

const meta: Meta<typeof CustomerDetailView> = {
  title: 'Customers/CustomerDetailView',
  component: CustomerDetailView,
  tags: ['autodocs'],
  parameters: {
    vue: {
      setup() {
        return {}
      },
    },
  },
}

export default meta
type Story = StoryObj<typeof CustomerDetailView>

export const ActiveCustomer: Story = {
  decorators: [
    () => ({
      template: '<story />',
      setup() {
        const pinia = createStore({
          customer: { ...baseCustomer, status: 'active' },
          auditLogs: auditEntries,
        })
        router.push('/customers/1')
        return { pinia }
      },
    }),
  ],
  parameters: {
    vueRouter: { router },
  },
}

export const PendingCustomer: Story = {
  decorators: [
    () => ({
      template: '<story />',
      setup() {
        const pinia = createStore({
          customer: { ...baseCustomer, status: 'pending' },
          auditLogs: [],
        })
        router.push('/customers/1')
        return { pinia }
      },
    }),
  ],
  parameters: {
    vueRouter: { router },
  },
}

export const SuspendedCustomer: Story = {
  decorators: [
    () => ({
      template: '<story />',
      setup() {
        const pinia = createStore({
          customer: { ...baseCustomer, status: 'suspended' },
          auditLogs: auditEntries,
        })
        router.push('/customers/1')
        return { pinia }
      },
    }),
  ],
  parameters: {
    vueRouter: { router },
  },
}

export const LoadingState: Story = {
  decorators: [
    () => ({
      template: '<story />',
      setup() {
        const pinia = createStore({
          loading: true,
        })
        router.push('/customers/1')
        return { pinia }
      },
    }),
  ],
  parameters: {
    vueRouter: { router },
  },
}

export const ErrorState: Story = {
  decorators: [
    () => ({
      template: '<story />',
      setup() {
        const pinia = createStore({
          error: 'Failed to load customer details',
        })
        router.push('/customers/1')
        return { pinia }
      },
    }),
  ],
  parameters: {
    vueRouter: { router },
  },
}
