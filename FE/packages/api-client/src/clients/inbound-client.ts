import { ApiClient } from './base'
import type {
  CreateInboundData,
  InboundQueueData,
  InboundItem,
  InspectData,
  FlagData,
} from '../types/responses'

export class InboundClient extends ApiClient {
  /** POST /v1/inbounds */
  async create(data: {
    customer_id: string
    expected_date: string
    notes?: string
    items: Array<{
      sku: string
      quantity_declared: number
      quantity_received?: number
      dimensions?: string
      weight?: number
    }>
  }): Promise<CreateInboundData> {
    return this.request<CreateInboundData>('/v1/inbounds', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  /** GET /v1/inbounds/queue */
  async getQueue(params?: {
    page_size?: number
    page_token?: string
    status?: string
  }): Promise<InboundQueueData> {
    const query = new URLSearchParams()
    if (params?.page_size) query.set('page_size', String(params.page_size))
    if (params?.page_token) query.set('page_token', params.page_token)
    if (params?.status) query.set('status', params.status)
    const qs = query.toString()
    return this.request<InboundQueueData>(
      `/v1/inbounds/queue${qs ? `?${qs}` : ''}`,
      { method: 'GET' }
    )
  }

  /** GET /v1/inbounds/{id} */
  async getById(id: string): Promise<InboundItem> {
    return this.request<InboundItem>(`/v1/inbounds/${id}`, {
      method: 'GET',
    })
  }

  /** POST /v1/inbounds/{id}/inspect */
  async inspect(id: string, data: {
    inspector_id: string
    notes?: string
    passed?: boolean
  }): Promise<InspectData> {
    return this.request<InspectData>(`/v1/inbounds/${id}/inspect`, {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  /** POST /v1/inbounds/{id}/approve */
  async approve(id: string): Promise<void> {
    await this.request<undefined>(`/v1/inbounds/${id}/approve`, {
      method: 'POST',
    })
  }

  /** POST /v1/inbounds/{id}/flag */
  async flag(id: string, reason: string): Promise<FlagData> {
    return this.request<FlagData>(`/v1/inbounds/${id}/flag`, {
      method: 'POST',
      body: JSON.stringify({ reason }),
    })
  }

  /** POST /v1/inbounds/{id}/hold */
  async hold(id: string, reason: string): Promise<void> {
    await this.request<undefined>(`/v1/inbounds/${id}/hold`, {
      method: 'POST',
      body: JSON.stringify({ reason }),
    })
  }

  /** POST /v1/inbounds/{id}/release */
  async release(id: string): Promise<void> {
    await this.request<undefined>(`/v1/inbounds/${id}/release`, {
      method: 'POST',
    })
  }

  /** GET /v1/inbounds/customer/{id} */
  async getCustomerInbounds(
    id: string,
    params?: { page_size?: number; page_token?: string }
  ): Promise<InboundQueueData> {
    const query = new URLSearchParams()
    if (params?.page_size) query.set('page_size', String(params.page_size))
    if (params?.page_token) query.set('page_token', params.page_token)
    const qs = query.toString()
    return this.request<InboundQueueData>(
      `/v1/inbounds/customer/${id}${qs ? `?${qs}` : ''}`,
      { method: 'GET' }
    )
  }
}
