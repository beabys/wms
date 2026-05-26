import { ApiClient } from './base'
import type {
  LoginData,
  RegisterData,
  RefreshData,
  MeData,
} from '../types/responses'

export class AuthClient extends ApiClient {
  /** POST /v1/auth/login */
  async login(email: string, password: string): Promise<LoginData> {
    return this.request<LoginData>('/v1/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    })
  }

  /** POST /v1/auth/register */
  async register(email: string, password: string, company_name?: string): Promise<RegisterData> {
    return this.request<RegisterData>('/v1/auth/register', {
      method: 'POST',
      body: JSON.stringify({ email, password, company_name }),
    })
  }

  /** POST /v1/auth/refresh */
  async refresh(refresh_token: string): Promise<RefreshData> {
    return this.request<RefreshData>('/v1/auth/refresh', {
      method: 'POST',
      body: JSON.stringify({ refresh_token }),
    })
  }

  /** POST /v1/auth/logout */
  async logout(refresh_token: string): Promise<void> {
    await this.request<undefined>('/v1/auth/logout', {
      method: 'POST',
      body: JSON.stringify({ refresh_token }),
    })
  }

  /** GET /v1/auth/me */
  async me(): Promise<MeData> {
    return this.request<MeData>('/v1/auth/me', {
      method: 'GET',
    })
  }
}
