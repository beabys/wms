import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import { useInviteStore } from '@/stores/invites'
import InviteListView from './InviteListView.vue'
import type { InviteEntry } from '@/api/types'

const teleportStub = { template: '<div><slot /></div>' }

const stubs = {
  InviteCustomerModal: true,
  Teleport: teleportStub,
}

function flushMicrotasks() {
  return new Promise(resolve => setTimeout(resolve, 0))
}

vi.mock('@/api/authClient', () => ({
  authClient: {
    listInvites: vi.fn().mockResolvedValue({
      invites: [],
      pagination: { page: 1, page_size: 10, total_items: 0 },
    }),
    inviteUser: vi.fn(),
    cancelInvite: vi.fn(),
  },
}))

function makeInvite(overrides: Partial<InviteEntry> = {}): InviteEntry {
  return {
    id: '1',
    email: 'test@test.com',
    token: 'tok1',
    invited_by: 'admin@test.com',
    status: 'pending',
    expires_at: 9999999999,
    created_at: 1710000000,
    ...overrides,
  }
}

describe('InviteListView', () => {
  async function createWrapper() {
    setActivePinia(createPinia())
    const wrapper = mount(InviteListView, {
      global: { stubs },
    })
    // wait for onMounted fetchInvites to settle
    await flushMicrotasks()
    await wrapper.vm.$nextTick()
    return wrapper
  }

  it('renders title and Send New Invite button', async () => {
    const wrapper = await createWrapper()
    expect(wrapper.text()).toContain('Sent Invitations')
    const buttons = wrapper.findAllComponents({ name: 'WmsButton' })
    const inviteBtn = buttons.find(b => b.text() === 'Send New Invite')
    expect(inviteBtn).toBeTruthy()
  })

  it('renders InviteListFilters', async () => {
    const wrapper = await createWrapper()
    const filters = wrapper.findComponent({ name: 'InviteListFilters' })
    expect(filters.exists()).toBe(true)
  })

  it('shows empty state when no invites', async () => {
    const wrapper = await createWrapper()
    await wrapper.vm.$nextTick()
    expect(wrapper.text()).toContain('No invitations found.')
  })

  it('shows loading state', async () => {
    const wrapper = await createWrapper()
    const store = useInviteStore()
    store.loading = true
    await wrapper.vm.$nextTick()
    expect(wrapper.text()).toContain('Loading invitations...')
  })

  it('shows error state', async () => {
    const wrapper = await createWrapper()
    const store = useInviteStore()
    store.error = 'Something broke'
    await wrapper.vm.$nextTick()
    expect(wrapper.text()).toContain('Something broke')
    expect(wrapper.text()).toContain('Retry')
  })

  it('renders invite rows in table', async () => {
    const wrapper = await createWrapper()
    const store = useInviteStore()
    store.invites = [
      makeInvite({ id: '1', email: 'alice@test.com' }),
      makeInvite({ id: '2', email: 'bob@test.com', status: 'used' }),
    ]
    store.pagination = { page: 1, page_size: 10, total_items: 2 }
    await wrapper.vm.$nextTick()
    expect(wrapper.text()).toContain('alice@test.com')
    expect(wrapper.text()).toContain('bob@test.com')
    expect(wrapper.text()).toContain('Copy Link')
  })

  it('shows pagination when items exist', async () => {
    const wrapper = await createWrapper()
    const store = useInviteStore()
    store.invites = [makeInvite()]
    store.pagination = { page: 1, page_size: 10, total_items: 25 }
    await wrapper.vm.$nextTick()
    expect(wrapper.text()).toContain('Page')
    expect(wrapper.text()).toContain('3')
  })

  it('opens invite modal on button click', async () => {
    const wrapper = await createWrapper()
    const buttons = wrapper.findAllComponents({ name: 'WmsButton' })
    const inviteBtn = buttons.find(b => b.text() === 'Send New Invite')
    await inviteBtn?.trigger('click')
    await wrapper.vm.$nextTick()
    const modal = wrapper.findComponent({ name: 'InviteCustomerModal' })
    expect(modal.props('visible')).toBe(true)
  })

  it('shows Cancel button only for pending invites', async () => {
    const wrapper = await createWrapper()
    const store = useInviteStore()
    store.invites = [
      makeInvite({ id: '1', email: 'pending@test.com', status: 'pending' }),
      makeInvite({ id: '2', email: 'used@test.com', status: 'used' }),
      makeInvite({ id: '3', email: 'cancelled@test.com', status: 'cancelled' }),
    ]
    store.pagination = { page: 1, page_size: 10, total_items: 3 }
    await wrapper.vm.$nextTick()

    const allButtons = wrapper.findAllComponents({ name: 'WmsButton' })
    const cancelButtons = allButtons.filter(b => b.text() === 'Cancel')
    // Only one cancel button (for the pending invite)
    expect(cancelButtons.length).toBe(1)
  })

  it('calls cancelInvite on Cancel button confirm', async () => {
    const wrapper = await createWrapper()
    const store = useInviteStore()
    store.invites = [makeInvite({ token: 'tok-cancel' })]
    store.pagination = { page: 1, page_size: 10, total_items: 1 }
    await wrapper.vm.$nextTick()

    const cancelBtn = wrapper.findAllComponents({ name: 'WmsButton' }).find(b => b.text() === 'Cancel')
    expect(cancelBtn).toBeTruthy()

    vi.spyOn(store, 'cancelInvite').mockResolvedValue(null)

    // Click Cancel row button → shows confirm dialog
    await cancelBtn?.trigger('click')
    await wrapper.vm.$nextTick()

    // Find the confirm dialog's confirm button (second WmsButton inside it)
    const dialogConfirmBtn = wrapper.findAllComponents({ name: 'WmsButton' }).find(
      b => b.text() === 'Cancel Invitation',
    )
    expect(dialogConfirmBtn).toBeTruthy()
    await dialogConfirmBtn?.trigger('click')
    await wrapper.vm.$nextTick()

    expect(store.cancelInvite).toHaveBeenCalledWith('tok-cancel')
  })

  it('does not call cancelInvite when confirm dialog is cancelled', async () => {
    const wrapper = await createWrapper()
    const store = useInviteStore()
    store.invites = [makeInvite({ token: 'tok-nocancel' })]
    store.pagination = { page: 1, page_size: 10, total_items: 1 }
    await wrapper.vm.$nextTick()

    const cancelBtn = wrapper.findAllComponents({ name: 'WmsButton' }).find(b => b.text() === 'Cancel')
    vi.spyOn(store, 'cancelInvite').mockResolvedValue(null)

    // Click Cancel row button → shows confirm dialog
    await cancelBtn?.trigger('click')
    await wrapper.vm.$nextTick()

    // Emit cancel on the dialog (simulating clicking cancel or overlay)
    const dialog = wrapper.findComponent({ name: 'WmsConfirmDialog' })
    expect(dialog.exists()).toBe(true)
    dialog.vm.$emit('cancel')
    await wrapper.vm.$nextTick()

    expect(store.cancelInvite).not.toHaveBeenCalled()
  })

  it('shows correct badge for Pending status', async () => {
    const wrapper = await createWrapper()
    const store = useInviteStore()
    store.invites = [makeInvite({ status: 'pending' })]
    store.pagination = { page: 1, page_size: 10, total_items: 1 }
    await wrapper.vm.$nextTick()
    expect(wrapper.text()).toContain('Pending')
  })

  it('shows correct badge for Used status', async () => {
    const wrapper = await createWrapper()
    const store = useInviteStore()
    store.invites = [makeInvite({ status: 'used' })]
    store.pagination = { page: 1, page_size: 10, total_items: 1 }
    await wrapper.vm.$nextTick()
    expect(wrapper.text()).toContain('Used')
  })

  it('shows correct badge for Cancelled status', async () => {
    const wrapper = await createWrapper()
    const store = useInviteStore()
    store.invites = [makeInvite({ status: 'cancelled' })]
    store.pagination = { page: 1, page_size: 10, total_items: 1 }
    await wrapper.vm.$nextTick()
    expect(wrapper.text()).toContain('Cancelled')
  })

  it('shows Expired badge when pending invite is past expiry', async () => {
    const wrapper = await createWrapper()
    const store = useInviteStore()
    store.invites = [makeInvite({ status: 'pending', expires_at: 1000000000 })]
    store.pagination = { page: 1, page_size: 10, total_items: 1 }
    await wrapper.vm.$nextTick()
    expect(wrapper.text()).toContain('Expired')
  })
})
