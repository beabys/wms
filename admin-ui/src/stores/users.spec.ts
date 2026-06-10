import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useUserStore } from './users'
import type { UserResponse, Pagination, UserListResponse, InviteUserResponse } from '@/api/types'

const mockUsers: UserResponse[] = [
  { id: '1', email: 'alice@test.com', name: 'Alice', role: 'admin', created_at: '2024-01-01T00:00:00Z', updated_at: '2024-01-01T00:00:00Z' },
  { id: '2', email: 'bob@test.com', name: 'Bob', role: 'warehouse_staff', created_at: '2024-02-01T00:00:00Z', updated_at: '2024-02-01T00:00:00Z' },
]

const mockPagination: Pagination = { page: 1, page_size: 10, total_items: 2 }
const mockListResponse: UserListResponse = { users: mockUsers, pagination: mockPagination }

describe('useUserStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.restoreAllMocks()
  })

  it('starts with empty state', () => {
    const store = useUserStore()
    expect(store.users).toEqual([])
    expect(store.pagination.page).toBe(1)
    expect(store.loading).toBe(false)
    expect(store.error).toBe('')
  })

  it('fetchUsers populates users and pagination', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true, data: mockListResponse }),
    }))

    const store = useUserStore()
    await store.fetchUsers()

    expect(store.users).toEqual(mockUsers)
    expect(store.pagination).toEqual(mockPagination)
    expect(store.loading).toBe(false)
    expect(store.error).toBe('')
  })

  it('fetchUsers with params builds query string', async () => {
    const fetchFn = vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true, data: mockListResponse }),
    })
    vi.stubGlobal('fetch', fetchFn)

    const store = useUserStore()
    await store.fetchUsers({ role: 'admin', page: 2, page_size: 5 })

    expect(fetchFn).toHaveBeenCalledWith(
      '/api/v1/users?role=admin&page=2&page_size=5',
      expect.anything(),
    )
  })

  it('fetchUsers sets error on failure', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 500,
      json: async () => ({ success: false, error: 'Server error' }),
    }))

    const store = useUserStore()
    await expect(store.fetchUsers()).rejects.toThrow('Server error')
    expect(store.error).toBe('Server error')
    expect(store.loading).toBe(false)
  })

  it('createUser calls authClient and returns user', async () => {
    const newUser: UserResponse = { id: '3', email: 'charlie@test.com', name: 'Charlie', role: 'billing_manager', created_at: '2024-03-01T00:00:00Z', updated_at: '2024-03-01T00:00:00Z' }
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true, data: newUser }),
    }))

    const store = useUserStore()
    const result = await store.createUser({ email: 'charlie@test.com', password: 'pass1234', name: 'Charlie', role: 'billing_manager' })

    expect(result).toEqual(newUser)
    expect(store.loading).toBe(false)
  })

  it('createUser sets error on failure', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 400,
      json: async () => ({ success: false, error: 'Email already exists' }),
    }))

    const store = useUserStore()
    await expect(store.createUser({ email: 'dupe@test.com', password: 'pass1234', name: 'Dupe', role: 'admin' })).rejects.toThrow('Email already exists')
    expect(store.error).toBe('Email already exists')
    expect(store.loading).toBe(false)
  })

  it('inviteUser calls authClient and returns invite response', async () => {
    const inviteResponse: InviteUserResponse = { token: 'invite-token', invite_link: 'http://example.com/invite/token', expires_at: 9999999999 }
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 200,
      json: async () => ({ success: true, data: inviteResponse }),
    }))

    const store = useUserStore()
    const result = await store.inviteUser('invite@test.com')

    expect(result).toEqual(inviteResponse)
    expect(store.loading).toBe(false)
  })

  it('inviteUser sets error on failure', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      status: 500,
      json: async () => ({ success: false, error: 'Invite failed' }),
    }))

    const store = useUserStore()
    await expect(store.inviteUser('bad@test.com')).rejects.toThrow('Invite failed')
    expect(store.error).toBe('Invite failed')
    expect(store.loading).toBe(false)
  })

  it('$reset clears all state', () => {
    const store = useUserStore()
    store.$reset()
    expect(store.users).toEqual([])
    expect(store.pagination.page).toBe(1)
    expect(store.loading).toBe(false)
    expect(store.error).toBe('')
  })
})
