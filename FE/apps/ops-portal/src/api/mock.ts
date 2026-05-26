// Mock API interceptor — intercepts fetch calls in dev mode.

import type { LoginData, CustomerQueueData, InboundQueueData, InboundItem } from '@wms/api-client/types/responses'

const originalFetch = window.fetch

if (import.meta.env.DEV && import.meta.env.VITE_USE_MOCK !== 'false') {
  window.fetch = async (input: RequestInfo | URL, init?: RequestInit): Promise<Response> => {
    const url = typeof input === 'string' ? input : input instanceof URL ? input.href : input.url
    const method = init?.method || 'GET'

    await new Promise(r => setTimeout(r, 300))

    // ── Auth ──────────────────────────────────────────────────────────
    if (url.includes('/v1/auth/login') && method === 'POST') {
      const data: LoginData = {
        access_token: 'mock-ops-jwt-' + Date.now(),
        refresh_token: 'mock-refresh-' + Date.now(),
        user: { id: 'ops-1', email: 'ops@wms.com', role: 'ops', customer_id: '', created_at: '2026-01-01T00:00:00Z' },
      }
      return new Response(JSON.stringify({ success: true, data }), { status: 200 })
    }

    if (url.includes('/v1/auth/me') && method === 'GET') {
      return new Response(JSON.stringify({
        success: true,
        data: { user: { id: 'ops-1', email: 'ops@wms.com', role: 'ops', customer_id: '', created_at: '2026-01-01T00:00:00Z' } },
      }), { status: 200 })
    }

    // ── Customer approvals ────────────────────────────────────────────
    if (url.includes('/v1/customers/queue') && method === 'GET') {
      const data: CustomerQueueData = {
        customers: [
          { id: 'c1', company_name: 'Acme Corp', email: 'alice@acme.com', vat_number: 'VAT-001', line1: '123 Main St', city: 'New York', country: 'US', status: 'pending', created_at: '2026-05-20T08:00:00Z' },
          { id: 'c2', company_name: 'Globex Inc', email: 'bob@globex.com', vat_number: 'VAT-002', line1: '456 Oak Ave', city: 'San Francisco', country: 'US', status: 'pending', created_at: '2026-05-21T08:00:00Z' },
          { id: 'c3', company_name: 'Initech', email: 'carol@initech.com', vat_number: 'VAT-003', line1: '789 Pine Rd', city: 'Austin', country: 'US', status: 'pending', created_at: '2026-05-22T08:00:00Z' },
          { id: 'c4', company_name: 'Umbrella Corp', email: 'dave@umbrella.com', vat_number: 'VAT-004', line1: '321 Elm St', city: 'Raccoon City', country: 'US', status: 'pending', created_at: '2026-05-23T08:00:00Z' },
          { id: 'c5', company_name: 'Hooli LLC', email: 'eve@hooli.com', vat_number: 'VAT-005', line1: '654 Birch Ln', city: 'Palo Alto', country: 'US', status: 'pending', created_at: '2026-05-24T08:00:00Z' },
        ],
        total: 5,
        page: 1,
        page_size: 20,
      }
      return new Response(JSON.stringify({ success: true, data }), { status: 200 })
    }

    if (url.includes('/v1/customers/') && method === 'POST' && url.includes('/approve')) {
      return new Response(JSON.stringify({ success: true, data: null }), { status: 200 })
    }

    if (url.includes('/v1/customers/') && method === 'POST' && url.includes('/suspend')) {
      return new Response(JSON.stringify({ success: true, data: null }), { status: 200 })
    }

    if (url.match(/\/v1\/customers\/[^/]+$/) && method === 'GET') {
      return new Response(JSON.stringify({
        success: true,
        data: {
          customer: { id: 'c1', company_name: 'Acme Corp', email: 'alice@acme.com', vat_number: 'VAT-001', line1: '123 Main St', city: 'New York', country: 'US', status: 'pending', created_at: '2026-05-20T08:00:00Z' },
        },
      }), { status: 200 })
    }

    // ── Inbound queue ────────────────────────────────────────────────
    if (url.includes('/v1/inbounds/queue') && method === 'GET') {
      const inbounds: InboundItem[] = [
        { id: 'INB-001', customer_id: 'c-1', customer_name: 'Acme Corp', status: 'pending', expected_date: '2026-05-28', created_at: '2026-05-20T08:00:00Z', items: [{ sku: 'SKU-001', quantity_declared: 100 }] },
        { id: 'INB-002', customer_id: 'c-2', customer_name: 'Globex Inc', status: 'approved', expected_date: '2026-05-27', notes: 'Urgent', created_at: '2026-05-19T08:00:00Z', items: [{ sku: 'SKU-002', quantity_declared: 50 }] },
        { id: 'INB-003', customer_id: 'c-3', customer_name: 'Initech', status: 'flagged', expected_date: '2026-05-29', notes: 'Damaged', created_at: '2026-05-18T08:00:00Z', items: [{ sku: 'SKU-003', quantity_declared: 200 }] },
        { id: 'INB-004', customer_id: 'c-4', customer_name: 'Umbrella Corp', status: 'held', expected_date: '2026-05-25', created_at: '2026-05-17T08:00:00Z', items: [{ sku: 'SKU-004', quantity_declared: 75 }] },
        { id: 'INB-005', customer_id: 'c-5', customer_name: 'Hooli LLC', status: 'pending', expected_date: '2026-06-01', created_at: '2026-05-22T08:00:00Z', items: [{ sku: 'SKU-005', quantity_declared: 30 }] },
        { id: 'INB-006', customer_id: 'c-1', customer_name: 'Stark Industries', status: 'shipped', expected_date: '2026-05-26', notes: 'Rush', created_at: '2026-05-16T08:00:00Z', items: [{ sku: 'SKU-006', quantity_declared: 500 }] },
        { id: 'INB-007', customer_id: 'c-2', customer_name: 'Wayne Enterprises', status: 'completed', expected_date: '2026-05-24', created_at: '2026-05-15T08:00:00Z', items: [{ sku: 'SKU-007', quantity_declared: 350 }] },
        { id: 'INB-008', customer_id: 'c-3', customer_name: 'Oscorp', status: 'pending', expected_date: '2026-05-30', created_at: '2026-05-23T08:00:00Z', items: [{ sku: 'SKU-008', quantity_declared: 750 }] },
      ]
      const data: InboundQueueData = { inbounds, page_size: 20 }
      return new Response(JSON.stringify({ success: true, data }), { status: 200 })
    }

    if (url.match(/\/v1\/inbounds\/(?!queue|customer\/)[^/]+$/) && method === 'GET') {
      const id = url.split('/').pop() || ''
      const data: InboundItem = {
        id,
        customer_id: 'c-1',
        customer_name: 'Acme Corp',
        status: 'pending',
        expected_date: '2026-05-28',
        notes: 'Sample inbound detail',
        created_at: '2026-05-20T08:00:00Z',
        items: [
          { sku: 'SKU-001', quantity_declared: 100, quantity_received: 95, dimensions: '30x20x15 cm', weight: 2.5 },
        ],
      }
      return new Response(JSON.stringify({ success: true, data }), { status: 200 })
    }

    if (url.includes('/v1/inbounds/') && method === 'POST' && url.includes('/approve')) {
      return new Response(JSON.stringify({ success: true, data: null }), { status: 200 })
    }

    return originalFetch(input, init)
  }
}

export {}
