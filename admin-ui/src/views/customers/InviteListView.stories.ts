import type { Meta, StoryObj } from '@storybook/vue3'
import InviteListView from './InviteListView.vue'
import { useInviteStore } from '@/stores/invites'
import type { InviteEntry, Pagination } from '@/api/types'

const meta: Meta<typeof InviteListView> = {
  title: 'Views/InviteListView',
  component: InviteListView,
  tags: ['autodocs'],
}

export default meta
type Story = StoryObj<typeof InviteListView>

interface StoreState {
  invites?: InviteEntry[]
  pagination?: Pagination
  loading?: boolean
  error?: string
}

function createRender(state: StoreState) {
  return () => {
    const store = useInviteStore()
    store.$reset()
    if (state.invites) store.invites = state.invites
    if (state.pagination) store.pagination = state.pagination
    if (state.loading !== undefined) store.loading = state.loading
    if (state.error !== undefined) store.error = state.error
    return { components: { InviteListView }, template: '<InviteListView />' }
  }
}

export const Default: Story = {}

export const Loading: Story = {
  render: createRender({ loading: true }),
}

export const Empty: Story = {
  render: createRender({}),
}

export const WithData: Story = {
  render: createRender({
    invites: [
      {
        id: '1',
        email: 'alice@example.com',
        token: 'abc123',
        invited_by: 'admin@wms.com',
        status: 'pending',
        expires_at: 9999999999,
        created_at: 1710000000,
      },
      {
        id: '2',
        email: 'bob@example.com',
        token: 'def456',
        invited_by: 'admin@wms.com',
        status: 'used',
        expires_at: 1000000000,
        created_at: 1700000000,
      },
      {
        id: '3',
        email: 'carol@example.com',
        token: 'ghi789',
        invited_by: 'admin@wms.com',
        status: 'cancelled',
        expires_at: 1000000000,
        created_at: 1690000000,
      },
    ],
    pagination: { page: 1, page_size: 10, total_items: 3 },
  }),
}

export const Filtered: Story = {
  render: createRender({
    invites: [
      {
        id: '2',
        email: 'bob@example.com',
        token: 'def456',
        invited_by: 'admin@wms.com',
        status: 'used',
        expires_at: 1000000000,
        created_at: 1700000000,
      },
    ],
    pagination: { page: 1, page_size: 10, total_items: 1 },
  }),
}

export const Paginated: Story = {
  render: createRender({
    invites: [
      {
        id: '1',
        email: 'alice@example.com',
        token: 'abc123',
        invited_by: 'admin@wms.com',
        status: 'pending',
        expires_at: 9999999999,
        created_at: 1710000000,
      },
    ],
    pagination: { page: 1, page_size: 10, total_items: 25 },
  }),
}

export const ErrorState: Story = {
  render: createRender({
    error: 'Failed to fetch invitations',
  }),
}

export const DarkMode: Story = {
  render: createRender({
    invites: [
      {
        id: '1',
        email: 'alice@example.com',
        token: 'abc123',
        invited_by: 'admin@wms.com',
        status: 'pending',
        expires_at: 9999999999,
        created_at: 1710000000,
      },
    ],
    pagination: { page: 1, page_size: 10, total_items: 1 },
  }),
  parameters: {
    themes: { theme: 'dark' },
  },
}
