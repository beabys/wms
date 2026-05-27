import { ApiClient } from './base'
import type {
  InviteData,
  RegisterCustomerData,
  CustomerQueueData,
  CustomerDetailData,
} from '../types/responses'

export class CustomerClient extends ApiClient {
  /** POST /v1/customers/invite */
  async invite(email: string): Promise<InviteData> {
    return this.request<InviteData>('/v1/customers/invite', {
      method: 'POST',
      body: JSON.stringify({ email }),
    })
  }

  /** POST /v1/customers/register */
  async register(data: {
    company_name: string
    vat_number: string
    line1: string
    line2?: string
    city: string
    postal_code?: string
    country: string
    rate_card_id?: string
  }): Promise<RegisterCustomerData> {
    return this.request<RegisterCustomerData>('/v1/customers/register', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  /** GET /v1/customers/queue */
  async getQueue(page = 1, page_size = 20): Promise<CustomerQueueData> {
    return this.request<CustomerQueueData>(
      `/v1/customers/queue?page=${page}&page_size=${page_size}`,
      { method: 'GET' }
    )
  }

  /** POST /v1/customers/{id}/approve */
  async approve(id: string): Promise<void> {
    await this.request<undefined>(`/v1/customers/${id}/approve`, {
      method: 'POST',
    })
  }

  /** POST /v1/customers/{id}/suspend */
  async suspend(id: string, reason?: string): Promise<void> {
    await this.request<undefined>(`/v1/customers/${id}/suspend`, {
      method: 'POST',
      body: reason ? JSON.stringify({ reason }) : undefined,
    })
  }

  /** GET /v1/customers/{id} */
  async getById(id: string): Promise<CustomerDetailData> {
    return this.request<CustomerDetailData>(`/v1/customers/${id}`, {
      method: 'GET',
    })
  }
}
