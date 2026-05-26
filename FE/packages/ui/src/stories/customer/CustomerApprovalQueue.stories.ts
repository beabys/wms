import type { Meta, StoryObj } from '@storybook/vue3'
import CustomerApprovalQueue, { type CustomerApproval } from '../../components/customer/CustomerApprovalQueue.vue'

const meta: Meta<typeof CustomerApprovalQueue> = {
  title: 'Customer/CustomerApprovalQueue',
  component: CustomerApprovalQueue,
  tags: ['autodocs']
}

export default meta
type Story = StoryObj<typeof CustomerApprovalQueue>

const sampleCustomers: CustomerApproval[] = [
  { id: '1', name: 'Alice Johnson', email: 'alice@acme.com', company: 'Acme Corp', registeredAt: '2026-05-20', status: 'pending' },
  { id: '2', name: 'Bob Smith', email: 'bob@globex.com', company: 'Globex Inc', registeredAt: '2026-05-21', status: 'pending' },
  { id: '3', name: 'Carol White', email: 'carol@initech.com', company: 'Initech', registeredAt: '2026-05-22', status: 'pending' },
  { id: '4', name: 'Dave Brown', email: 'dave@umbrella.com', company: 'Umbrella Corp', registeredAt: '2026-05-23', status: 'pending' },
  { id: '5', name: 'Eve Davis', email: 'eve@hooli.com', company: 'Hooli LLC', registeredAt: '2026-05-24', status: 'pending' }
]

export const Empty: Story = {
  args: {
    customers: [],
    loading: false
  }
}

export const FiveRows: Story = {
  args: {
    customers: sampleCustomers,
    loading: false
  }
}
