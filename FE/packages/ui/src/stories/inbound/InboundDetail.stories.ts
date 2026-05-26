import type { Meta, StoryObj } from '@storybook/vue3'
import InboundDetail, { type InboundDetailShipment } from '../../components/inbound/InboundDetail.vue'

const meta: Meta<typeof InboundDetail> = {
  title: 'Inbound/InboundDetail',
  component: InboundDetail,
  tags: ['autodocs']
}

export default meta
type Story = StoryObj<typeof InboundDetail>

const sampleShipment: InboundDetailShipment = {
  id: 'INB-001',
  customerName: 'Acme Corp',
  skuCount: 5,
  totalUnits: 250,
  status: 'pending',
  expectedDate: '2026-05-28'
}

export const Submitted: Story = {
  args: {
    shipment: sampleShipment,
    currentStep: 0,
    loading: false
  }
}

export const Inspection: Story = {
  args: {
    shipment: { ...sampleShipment, status: 'approved' },
    currentStep: 1,
    loading: false
  }
}

export const ApproveFlagHold: Story = {
  args: {
    shipment: { ...sampleShipment, status: 'flagged' },
    currentStep: 2,
    loading: false
  }
}

export const Release: Story = {
  args: {
    shipment: { ...sampleShipment, status: 'completed' },
    currentStep: 3,
    loading: false
  }
}
