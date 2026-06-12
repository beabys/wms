import type { Meta, StoryObj } from '@storybook/vue3'
import WmsConfirmDialog from './WmsConfirmDialog.vue'

const meta: Meta<typeof WmsConfirmDialog> = {
  title: 'Common/WmsConfirmDialog',
  component: WmsConfirmDialog,
  tags: ['autodocs'],
  argTypes: {
    variant: { control: 'select', options: ['primary', 'danger'] },
    visible: { control: 'boolean' },
  },
  args: {
    visible: true,
  },
}

export default meta
type Story = StoryObj<typeof WmsConfirmDialog>

export const Default: Story = {
  args: {
    title: 'Confirm Action',
    message: 'Are you sure you want to proceed with this action?',
    confirmText: 'Confirm',
    cancelText: 'Cancel',
    variant: 'primary',
  },
}

export const Danger: Story = {
  args: {
    title: 'Delete Item',
    message: 'This action cannot be undone. Are you sure you want to delete this item?',
    confirmText: 'Delete',
    cancelText: 'Keep',
    variant: 'danger',
  },
}

export const LongMessage: Story = {
  args: {
    title: 'Terms and Conditions',
    message: 'By confirming you agree to the following terms: you will be billed monthly, you can cancel anytime, and your data will be processed according to our privacy policy. Please review carefully before proceeding.',
    confirmText: 'I Agree',
    cancelText: 'Cancel',
    variant: 'primary',
  },
}

export const CustomSlot: Story = {
  args: {
    title: 'Warning',
    confirmText: 'OK',
    cancelText: 'Cancel',
    variant: 'danger',
  },
  render: (args) => ({
    components: { WmsConfirmDialog },
    setup: () => ({ args }),
    template: `
      <WmsConfirmDialog v-bind="args">
        <p>This order includes <strong>3 items</strong> that are <em>out of stock</em>.</p>
      </WmsConfirmDialog>
    `,
  }),
}
