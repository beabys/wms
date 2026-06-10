import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useInviteStore } from './invites'
import type { InviteEntry, Pagination } from '@/api/types'

const mockInvites: InviteEntry[] = [
  {
    id: '1',
    email: 'alice@test.com',
    token: 'tok1',
    invited_by: 'admin@test.com',
    status: 'pending',
    expires_at: 9999999999,
    created_at: 1710000000,
  },
  {
    id: '2',
    email: 'bob@test.com',
    token: 'tok2',
    invited_by: 'admin@test.com',
    status: 'used',
    expires_at: 1000000000,
    created_at: 1700000000,
  },
]

const mockPagination: Pagination = { page: 1, page_size: 10, total_items: 2 }

describe('useInviteStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.restoreAllMocks()
  })

  it('starts with default state', () => {
    const store = useInviteStore()
    expect(store.invites).toEqual([])
    expect(store.pagination).toEqual({ page: 1, page_size: 10, total_items: 0 })
    expect(store.loading).toBe(false)
    expect(store.error).toBe('')
  })

  it('fetchInvites populates invites and pagination', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true, data: { invites: mockInvites, pagination: mockPagination } }),
    }))

    const store = useInviteStore()
    await store.fetchInvites()

    expect(store.invites).toEqual(mockInvites)
    expect(store.pagination).toEqual(mockPagination)
    expect(store.loading).toBe(false)
    expect(store.error).toBe('')
  })

  it('fetchInvites with params builds query string', async () => {
    const fetchFn = vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true, data: { invites: [], pagination: { page: 1, page_size: 10, total_items: 0 } } }),
    })
    vi.stubGlobal('fetch', fetchFn)

    const store = useInviteStore()
    await store.fetchInvites({ status: 'pending', expired: true, page: 2, page_size: 5 })

    expect(fetchFn).toHaveBeenCalledWith(
      '/api/v1/auth/invites?page=2&page_size=5&status=pending&expired=true',
      expect.anything(),
    )
  })

  it('fetchInvites with date range builds query string', async () => {
    const fetchFn = vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true, data: { invites: [], pagination: { page: 1, page_size: 10, total_items: 0 } } }),
    })
    vi.stubGlobal('fetch', fetchFn)

    const store = useInviteStore()
    await store.fetchInvites({ created_after: 1700000000, created_before: 1710000000 })

    expect(fetchFn).toHaveBeenCalledWith(
      '/api/v1/auth/invites?created_after=1700000000&created_before=1710000000',
      expect.anything(),
    )
  })

  it('fetchInvites sets error on API failure response', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 500,
      json: async () => ({ success: false, error: 'Server error' }),
    }))

    const store = useInviteStore()
    await store.fetchInvites()
    expect(store.invites).toEqual([])
    expect(store.error).toBe('Server error')
    expect(store.loading).toBe(false)
  })

  it('fetchInvites sets error on network failure', async () => {
    vi.stubGlobal('fetch', vi.fn().mockRejectedValue(new Error('Network error')))

    const store = useInviteStore()
    await store.fetchInvites()
    expect(store.invites).toEqual([])
    expect(store.error).toBe('Network error')
    expect(store.loading).toBe(false)
  })

  it('cancelInvite calls API and refreshes list on success', async () => {
    const fetchFn = vi.fn()
    fetchFn
      .mockResolvedValueOnce({ status: 200, json: async () => ({ success: true }) }) // POST cancel
      .mockResolvedValueOnce({ // refresh
        status: 200,
        json: async () => ({ success: true, data: { invites: mockInvites, pagination: mockPagination } }),
      })
    vi.stubGlobal('fetch', fetchFn)

    const store = useInviteStore()
    const err = await store.cancelInvite('tok1')

    expect(err).toBeNull()
    // First call should be POST to cancel
    expect(fetchFn.mock.calls[0][0]).toBe('/api/v1/auth/invites/tok1/cancel')
    expect(fetchFn.mock.calls[0][1]?.method).toBe('POST')
    // Second call should refresh invites
    expect(fetchFn.mock.calls[1][0]).toBe('/api/v1/auth/invites')
    expect(store.invites).toEqual(mockInvites)
  })

  it('cancelInvite returns error on failure', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 400,
      json: async () => ({ success: false, error: 'Cannot cancel used invite' }),
    }))

    const store = useInviteStore()
    const err = await store.cancelInvite('tok2')

    expect(err).toBe('Cannot cancel used invite')
  })

  it('inviteUser returns result on success', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true, data: { token: 'new-tok', invite_link: 'http://...', expires_at: 9999999999 } }),
    }))

    const store = useInviteStore()
    const result = await store.inviteUser('new@test.com')

    expect(typeof result).not.toBe('string')
    if (typeof result !== 'string') {
      expect(result.token).toBe('new-tok')
    }
  })

  it('inviteUser returns duplicate email error message', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 409,
      json: async () => ({ success: false, error: 'Email already has a pending invite' }),
    }))

    const store = useInviteStore()
    const result = await store.inviteUser('dup@test.com')

    expect(result).toBe('This email already has a pending invitation')
  })

  it('inviteUser returns generic error for other failures', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 500,
      json: async () => ({ success: false, error: 'Server error' }),
    }))

    const store = useInviteStore()
    const result = await store.inviteUser('fail@test.com')

    expect(result).toBe('Server error')
  })

  it('$reset restores default state', () => {
    const store = useInviteStore()
    store.invites = mockInvites
    store.pagination = { page: 2, page_size: 20, total_items: 5 }
    store.loading = true
    store.error = 'some error'
    store.$reset()
    expect(store.invites).toEqual([])
    expect(store.pagination).toEqual({ page: 1, page_size: 10, total_items: 0 })
    expect(store.loading).toBe(false)
    expect(store.error).toBe('')
  })
})
