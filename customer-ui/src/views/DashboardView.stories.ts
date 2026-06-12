import type { Meta, StoryObj } from '@storybook/vue3'
import { setActivePinia, createPinia } from 'pinia'
import DashboardView from './DashboardView.vue'
import { useCustomerStore } from '@/stores/customer'

const meta: Meta<typeof DashboardView> = {
  title: 'Views/Dashboard',
  component: DashboardView,
}

export default meta
type Story = StoryObj<typeof DashboardView>

export const Default: Story = {
  render: () => ({
    components: { DashboardView },
    template: '<DashboardView />',
  }),
}

export const Pending: Story = {
  render: () => {
    setActivePinia(createPinia())
    const store = useCustomerStore()
    store.customer = { id: '1', company_name: 'Acme', email: 'a@b.com', phone: '', vat_number: '', status: 'pending', created_at: '' }
    return {
      components: { DashboardView },
      template: '<DashboardView />',
    }
  },
}

export const Rejected: Story = {
  render: () => {
    setActivePinia(createPinia())
    const store = useCustomerStore()
    store.customer = { id: '1', company_name: 'Acme', email: 'a@b.com', phone: '', vat_number: '', status: 'rejected', created_at: '' }
    return {
      components: { DashboardView },
      template: '<DashboardView />',
    }
  },
}

export const Suspended: Story = {
  render: () => {
    setActivePinia(createPinia())
    const store = useCustomerStore()
    store.customer = { id: '1', company_name: 'Acme', email: 'a@b.com', phone: '', vat_number: '', status: 'suspended', created_at: '' }
    return {
      components: { DashboardView },
      template: '<DashboardView />',
    }
  },
}

export const Active: Story = {
  render: () => {
    setActivePinia(createPinia())
    const store = useCustomerStore()
    store.customer = { id: '1', company_name: 'Acme', email: 'a@b.com', phone: '', vat_number: '', status: 'active', created_at: '' }
    return {
      components: { DashboardView },
      template: '<DashboardView />',
    }
  },
}
