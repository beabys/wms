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

    // ── Customer stock ─────────────────────────────────────────────────
    if (url.includes('/v1/stock/customer/') && method === 'GET') {
      return new Response(JSON.stringify({
        success: true,
        data: {
          entries: [
            { id: 'cst-1', product_id: 'p1', bin_location_id: 'loc-1', quantity: 100, reserved_quantity: 5, lot_number: 'LOT-001', expiry_date: '2027-01-01', status: 'available', last_updated: '2026-05-20T10:00:00Z' },
            { id: 'cst-2', product_id: 'p2', bin_location_id: 'loc-2', quantity: 5, reserved_quantity: 2, lot_number: 'LOT-002', expiry_date: '2026-12-01', status: 'available', last_updated: '2026-05-19T10:00:00Z' },
            { id: 'cst-3', product_id: 'p3', bin_location_id: 'loc-3', quantity: 0, reserved_quantity: 0, lot_number: '', expiry_date: null, status: 'inactive', last_updated: '2026-05-18T10:00:00Z' },
            { id: 'cst-4', product_id: 'p4', bin_location_id: 'loc-4', quantity: 15, reserved_quantity: 10, lot_number: 'LOT-004', expiry_date: '2026-06-15', status: 'reserved', last_updated: '2026-05-20T08:00:00Z' },
          ],
          page_size: 20,
        },
      }), { status: 200 })
    }

    // ── Products (read-only lookup) ─────────────────────────────────────
    if (url.includes('/v1/products') && method === 'GET') {
      return new Response(JSON.stringify({
        success: true,
        data: {
          products: [
            { id: 'p1', sku: 'SKU-001', name: 'Widget Alpha', description: 'Premium widget', category: 'Electronics', unit: 'piece', weight_kg: 1.5, is_active: true, low_stock_threshold: 10, created_at: '2026-01-01T00:00:00Z', updated_at: '2026-05-01T00:00:00Z' },
            { id: 'p2', sku: 'SKU-002', name: 'Gadget Beta', description: 'Portable gadget', category: 'Electronics', unit: 'piece', weight_kg: 0.8, is_active: true, low_stock_threshold: 10, created_at: '2026-01-15T00:00:00Z', updated_at: '2026-05-10T00:00:00Z' },
          ],
          total: 2,
          page: 1,
          page_size: 20,
        },
      }), { status: 200 })
    }

    // Fallback — pass through to real fetch
    return originalFetch(input, init)
  }
}

export {}
