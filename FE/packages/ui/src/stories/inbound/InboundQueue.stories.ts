import type { Meta, StoryObj } from '@storybook/vue3'
import InboundQueue, { type InboundShipment } from '../../components/inbound/InboundQueue.vue'

const meta: Meta<typeof InboundQueue> = {
  title: 'Inbound/InboundQueue',
  component: InboundQueue,
  tags: ['autodocs']
}

export default meta
type Story = StoryObj<typeof InboundQueue>

const sampleShipments: InboundShipment[] = [
  { id: 'INB-001', customerName: 'Acme Corp', productCount: 5, status: 'pending', expectedDate: '2026-05-28', receivedDate: undefined },
  { id: 'INB-002', customerName: 'Globex Inc', productCount: 3, status: 'approved', expectedDate: '2026-05-27', receivedDate: '2026-05-26' },
  { id: 'INB-003', customerName: 'Initech', productCount: 12, status: 'flagged', expectedDate: '2026-05-29', receivedDate: undefined },
  { id: 'INB-004', customerName: 'Umbrella Corp', productCount: 8, status: 'held', expectedDate: '2026-05-25', receivedDate: '2026-05-25' },
  { id: 'INB-005', customerName: 'Hooli LLC', productCount: 2, status: 'pending', expectedDate: '2026-06-01', receivedDate: undefined },
  { id: 'INB-006', customerName: 'Stark Industries', productCount: 20, status: 'shipped', expectedDate: '2026-05-26', receivedDate: '2026-05-25' },
  { id: 'INB-007', customerName: 'Wayne Enterprises', productCount: 7, status: 'completed', expectedDate: '2026-05-24', receivedDate: '2026-05-24' },
  { id: 'INB-008', customerName: 'Oscorp', productCount: 15, status: 'pending', expectedDate: '2026-05-30', receivedDate: undefined }
]

export const Empty: Story = {
  args: {
    shipments: [],
    loading: false
  }
}

export const EightRows: Story = {
  args: {
    shipments: sampleShipments,
    loading: false
  }
}
