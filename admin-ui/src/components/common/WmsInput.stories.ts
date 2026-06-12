import type { Meta, StoryObj } from '@storybook/vue3'
import WmsInput from './WmsInput.vue'

const meta: Meta<typeof WmsInput> = {
  title: 'Common/WmsInput',
  component: WmsInput,
  tags: ['autodocs'],
  argTypes: {
    type: { control: 'select', options: ['text', 'email', 'password'] },
  },
}

export default meta
type Story = StoryObj<typeof WmsInput>

export const Default: Story = {
  args: {
    label: 'Email',
    placeholder: 'you@example.com',
    type: 'email',
  },
}

export const WithValue: Story = {
  args: {
    label: 'Email',
    modelValue: 'user@example.com',
    type: 'email',
  },
}

export const WithError: Story = {
  args: {
    label: 'Email',
    modelValue: 'invalid',
    type: 'email',
    error: 'Invalid email format',
  },
}

export const Disabled: Story = {
  args: {
    label: 'Email',
    placeholder: 'you@example.com',
    disabled: true,
  },
}

export const WithPrefix: Story = {
  args: {
    label: 'Search',
    placeholder: 'Search...',
  },
  render: (args) => ({
    components: { WmsInput },
    setup: () => ({ args }),
    template: '<WmsInput v-bind="args"><template #prefix><span style="color: var(--color-text-muted)">🔍</span></template></WmsInput>',
  }),
}
