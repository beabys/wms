import { request } from '@/composables/useApi'
import type { CustomerResponse, CustomerListResponse, AuditLogListResponse } from '@/api/types'

const BASE = '/api/v1'

export const customerClient = {
  async listCustomers(params?: { status?: string; page?: number; page_size?: number }): Promise<CustomerListResponse> {
    const query = new URLSearchParams()
    if (params?.status) query.set('status', params.status)
    if (params?.page) query.set('page', String(params.page))
    if (params?.page_size) query.set('page_size', String(params.page_size))
    const qs = query.toString()
    return request<CustomerListResponse>(`${BASE}/customers${qs ? '?' + qs : ''}`)
  },

  async getCustomer(id: string): Promise<CustomerResponse> {
    return request<CustomerResponse>(`${BASE}/customers/${id}`)
  },

  async updateCustomer(id: string, data: Partial<CustomerResponse>): Promise<CustomerResponse> {
    return request<CustomerResponse>(`${BASE}/customers/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    })
  },

  async approveCustomer(id: string): Promise<CustomerResponse> {
    return request<CustomerResponse>(`${BASE}/customers/${id}/approve`, {
      method: 'POST',
    })
  },

  async rejectCustomer(id: string, reason?: string): Promise<CustomerResponse> {
    return request<CustomerResponse>(`${BASE}/customers/${id}/reject`, {
      method: 'POST',
      body: reason ? JSON.stringify({ reason }) : undefined,
    })
  },

  async suspendCustomer(id: string): Promise<void> {
    await request(`${BASE}/customers/${id}/suspend`, { method: 'POST' })
  },

  async restoreCustomer(id: string): Promise<CustomerResponse> {
    return request<CustomerResponse>(`${BASE}/customers/${id}/restore`, {
      method: 'POST',
    })
  },

  async listAuditLogs(id: string, params?: { page?: number; page_size?: number }): Promise<AuditLogListResponse> {
    const query = new URLSearchParams()
    if (params?.page) query.set('page', String(params.page))
    if (params?.page_size) query.set('page_size', String(params.page_size))
    const qs = query.toString()
    return request<AuditLogListResponse>(`${BASE}/customers/${id}/audit${qs ? '?' + qs : ''}`)
  },
}
