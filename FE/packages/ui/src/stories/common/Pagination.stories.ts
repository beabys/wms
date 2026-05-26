import type { Meta, StoryObj } from '@storybook/vue3'
import Pagination from '../../components/common/Pagination.vue'

const meta: Meta<typeof Pagination> = {
  title: 'Common/Pagination',
  component: Pagination,
  tags: ['autodocs']
}

export default meta
type Story = StoryObj<typeof Pagination>

export const FirstPage: Story = {
  args: {
    page: 1,
    total: 50,
    limit: 10
  }
}

export const MiddlePage: Story = {
  args: {
    page: 3,
    total: 50,
    limit: 10
  }
}

export const LastPage: Story = {
  args: {
    page: 5,
    total: 50,
    limit: 10
  }
}

export const SinglePage: Story = {
  args: {
    page: 1,
    total: 3,
    limit: 10
  }
}
