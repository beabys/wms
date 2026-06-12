import type { Meta, StoryObj } from '@storybook/vue3'
import WmsCard from './WmsCard.vue'

const meta: Meta<typeof WmsCard> = {
  title: 'Common/WmsCard',
  component: WmsCard,
  tags: ['autodocs'],
  argTypes: {
    padding: { control: 'select', options: ['none', 'sm', 'md', 'lg'] },
  },
}

export default meta
type Story = StoryObj<typeof WmsCard>

export const Default: Story = {
  args: {
    padding: 'md',
  },
  render: (args) => ({
    components: { WmsCard },
    setup: () => ({ args }),
    template: '<WmsCard v-bind="args">This is a basic card with some content.</WmsCard>',
  }),
}

export const WithHeaderAndFooter: Story = {
  args: {
    padding: 'md',
  },
  render: (args) => ({
    components: { WmsCard },
    setup: () => ({ args }),
    template: `
      <WmsCard v-bind="args">
        <template #header><strong style="font-size: 1.125rem">Card Header</strong></template>
        Main card content goes here.
        <template #footer><span style="color: var(--color-text-muted)">Last updated: today</span></template>
      </WmsCard>
    `,
  }),
}

export const Compact: Story = {
  args: {
    padding: 'sm',
  },
  render: (args) => ({
    components: { WmsCard },
    setup: () => ({ args }),
    template: '<WmsCard v-bind="args">Compact card content.</WmsCard>',
  }),
}

export const LargePadding: Story = {
  args: {
    padding: 'lg',
  },
  render: (args) => ({
    components: { WmsCard },
    setup: () => ({ args }),
    template: '<WmsCard v-bind="args">Spacious card with extra padding.</WmsCard>',
  }),
}
