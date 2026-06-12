import { request } from '@/composables/useApi'
import type {
  LoginRequest,
  LoginResponse,
  RefreshRequest,
  RefreshResponse,
  UserResponse,
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

  async logout(): Promise<void> {
    await request(`${BASE}/auth/logout`, { method: 'POST' })
  },
}
