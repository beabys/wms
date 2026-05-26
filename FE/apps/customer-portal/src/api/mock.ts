// Mock API interceptor — intercepts fetch calls in dev mode.
// In prod, VITE_USE_MOCK=false or not set — this file is still imported but condition is false.

import type { LoginData, RegisterData } from '@wms/api-client/types/responses'

const originalFetch = window.fetch

if (import.meta.env.DEV && import.meta.env.VITE_USE_MOCK !== 'false') {
  window.fetch = async (input: RequestInfo | URL, init?: RequestInit): Promise<Response> => {
    const url = typeof input === 'string' ? input : input instanceof URL ? input.href : input.url
    const method = init?.method || 'GET'

    // Simulate network delay
    await new Promise(r => setTimeout(r, 300))

    // ── Auth ──────────────────────────────────────────────────────────
    if (url.includes('/v1/auth/login') && method === 'POST') {
      const data: LoginData = {
        access_token: 'mock-customer-jwt-' + Date.now(),
        refresh_token: 'mock-refresh-' + Date.now(),
        user: { id: 'cust-1', email: 'customer@example.com', role: 'customer', customer_id: 'c-1', created_at: '2026-01-15T10:00:00Z' },
      }
      return new Response(JSON.stringify({ success: true, data }), { status: 200 })
    }

    if (url.includes('/v1/auth/register') && method === 'POST') {
      const data: RegisterData = {
        user: { id: 'cust-2', email: 'new@example.com', role: 'customer', customer_id: 'c-2', created_at: '2026-05-26T10:00:00Z' },
      }
      return new Response(JSON.stringify({ success: true, data }), { status: 200 })
    }

    if (url.includes('/v1/auth/me') && method === 'GET') {
      return new Response(JSON.stringify({
        success: true,
        data: { user: { id: 'cust-1', email: 'customer@example.com', role: 'customer', customer_id: 'c-1', created_at: '2026-01-15T10:00:00Z' } },
      }), { status: 200 })
    }

    // ── Inbound ───────────────────────────────────────────────────────
    if (url.includes('/v1/inbounds') && method === 'POST' && !url.includes('/customer/')) {
      return new Response(JSON.stringify({
        success: true,
        data: { id: 'INB-MOCK-' + Date.now(), message: 'Inbound shipment created successfully.' },
      }), { status: 200 })
    }

    if (url.includes('/v1/inbounds/customer/') && method === 'GET') {
      return new Response(JSON.stringify({
        success: true,
        data: {
          inbounds: [
            { id: 'INB-001', customer_id: 'c-1', customer_name: 'Acme Corp', status: 'pending', expected_date: '2026-05-28', created_at: '2026-05-20T08:00:00Z', items: [{ sku: 'SKU-001', quantity_declared: 100 }] },
            { id: 'INB-002', customer_id: 'c-1', customer_name: 'Acme Corp', status: 'approved', expected_date: '2026-05-27', notes: 'Urgent', created_at: '2026-05-19T08:00:00Z', items: [{ sku: 'SKU-002', quantity_declared: 50 }] },
          ],
          page_size: 20,
        },
      }), { status: 200 })
    }

    // Fallback — pass through to real fetch
    return originalFetch(input, init)
  }
}

export {}
