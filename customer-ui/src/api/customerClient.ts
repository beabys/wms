import { request } from '@/composables/useApi'
import type { RegisterRequest, RegisterResponse, CustomerResponse } from '@/api/types'

const BASE = '/api/v1'

export const customerClient = {
  async register(data: RegisterRequest): Promise<RegisterResponse> {
    return request<RegisterResponse>(`${BASE}/customers/register`, {
      method: 'POST',
      body: JSON.stringify(data),
    })
  },

  async getCustomer(id: string): Promise<CustomerResponse> {
    return request<CustomerResponse>(`${BASE}/customers/${id}`)
  },

  async getMyCustomer(): Promise<CustomerResponse> {
    return request<CustomerResponse>(`${BASE}/customers/me`)
  },
}
