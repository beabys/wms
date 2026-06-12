import { request } from '@/composables/useApi'
import type {
  LoginRequest,
  LoginResponse,
  RefreshRequest,
  RefreshResponse,
  UserResponse,
  CreateUserRequest,
  UserListResponse,
  InviteUserResponse,
  InviteListResponse,
} from '@/api/types'

const BASE = '/api/v1'

export const authClient = {
  async login(data: LoginRequest): Promise<LoginResponse> {
    return request<LoginResponse>(`${BASE}/auth/login`, {
      method: 'POST',
      body: JSON.stringify(data),
    })
  },

  async refresh(refreshToken: string): Promise<RefreshResponse> {
    const data: RefreshRequest = { refresh_token: refreshToken }
    return request<RefreshResponse>(`${BASE}/auth/refresh`, {
      method: 'POST',
      body: JSON.stringify(data),
    })
  },

  async getMe(): Promise<UserResponse> {
    return request<UserResponse>(`${BASE}/auth/me`)
  },

  async createUser(data: CreateUserRequest): Promise<UserResponse> {
    return request<UserResponse>(`${BASE}/users`, {
      method: 'POST',
      body: JSON.stringify(data),
    })
  },

  async listUsers(params?: { role?: string; page?: number; page_size?: number }): Promise<UserListResponse> {
    const query = new URLSearchParams()
    if (params?.role) query.set('role', params.role)
    if (params?.page) query.set('page', String(params.page))
    if (params?.page_size) query.set('page_size', String(params.page_size))
    const qs = query.toString()
    return request<UserListResponse>(`${BASE}/users${qs ? '?' + qs : ''}`)
  },

  async getUser(id: string): Promise<UserResponse> {
    return request<UserResponse>(`${BASE}/users/${id}`)
  },

  async updateUser(id: string, data: Partial<CreateUserRequest>): Promise<UserResponse> {
    return request<UserResponse>(`${BASE}/users/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    })
  },

  async logout(): Promise<void> {
    await request(`${BASE}/auth/logout`, { method: 'POST' })
  },

  async inviteUser(email: string): Promise<InviteUserResponse> {
    return request<InviteUserResponse>(`${BASE}/auth/invite`, {
      method: 'POST',
      body: JSON.stringify({ email }),
    })
  },

  async listInvites(params?: {
    page?: number
    page_size?: number
    status?: string
    expired?: boolean
    created_after?: number
    created_before?: number
  }): Promise<InviteListResponse> {
    const query = new URLSearchParams()
    if (params?.page) query.set('page', String(params.page))
    if (params?.page_size) query.set('page_size', String(params.page_size))
    if (params?.status) query.set('status', params.status)
    if (params?.expired !== undefined) query.set('expired', String(params.expired))
    if (params?.created_after) query.set('created_after', String(params.created_after))
    if (params?.created_before) query.set('created_before', String(params.created_before))
    const qs = query.toString()
    return request<InviteListResponse>(`${BASE}/auth/invites${qs ? '?' + qs : ''}`)
  },

  async cancelInvite(token: string): Promise<void> {
    await request(`${BASE}/auth/invites/${token}/cancel`, {
      method: 'POST',
    })
  },
}
