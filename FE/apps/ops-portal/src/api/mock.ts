// Mock API interceptor — intercepts fetch calls in dev mode.

import type { LoginData, CustomerQueueData, InboundQueueData, InboundItem } from '@wms/api-client/types/responses'
import type { ProductListData, StockListData, BinLocationListData, LowStockAlertData } from '@wms/api-client/types/inventory'

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

    // ── Products ──────────────────────────────────────────────────────
    if (url.includes('/v1/products') && method === 'GET') {
      const data: ProductListData = {
        products: [
          { id: 'p1', sku: 'SKU-001', name: 'Widget Alpha', description: 'Premium widget', category: 'Electronics', unit: 'piece', weight_kg: 1.5, is_active: true, low_stock_threshold: 10, created_at: '2026-01-01T00:00:00Z', updated_at: '2026-05-01T00:00:00Z' },
          { id: 'p2', sku: 'SKU-002', name: 'Gadget Beta', description: 'Portable gadget', category: 'Electronics', unit: 'piece', weight_kg: 0.8, is_active: true, low_stock_threshold: 10, created_at: '2026-01-15T00:00:00Z', updated_at: '2026-05-10T00:00:00Z' },
          { id: 'p3', sku: 'SKU-003', name: 'Component Gamma', description: 'Replacement component', category: 'Automotive', unit: 'piece', weight_kg: 2.0, is_active: true, low_stock_threshold: 20, created_at: '2026-02-01T00:00:00Z', updated_at: '2026-04-20T00:00:00Z' },
          { id: 'p4', sku: 'SKU-004', name: 'Device Delta', description: 'IoT sensor device', category: 'Electronics', unit: 'piece', weight_kg: 0.3, is_active: true, low_stock_threshold: 15, created_at: '2026-03-01T00:00:00Z', updated_at: '2026-05-15T00:00:00Z' },
          { id: 'p5', sku: 'SKU-005', name: 'Material Epsilon', description: 'Raw material', category: 'Other', unit: 'kg', weight_kg: 25.0, is_active: false, low_stock_threshold: 50, created_at: '2026-01-20T00:00:00Z', updated_at: '2026-03-01T00:00:00Z' },
        ],
        total: 5,
        page: 1,
        page_size: 20,
      }
      return new Response(JSON.stringify({ success: true, data }), { status: 200 })
    }

    if (url.includes('/v1/products') && method === 'POST') {
      return new Response(JSON.stringify({
        success: true,
        data: { id: 'p-new', sku: 'SKU-NEW', name: 'New Product', description: '', category: 'Other', unit: 'piece', weight_kg: 1.0, is_active: true, low_stock_threshold: 10, created_at: new Date().toISOString(), updated_at: new Date().toISOString() },
      }), { status: 200 })
    }

    if (url.match(/\/v1\/products\/[^/]+$/) && method === 'PUT') {
      return new Response(JSON.stringify({ success: true, data: { id: url.split('/').pop(), updated: true } }), { status: 200 })
    }

    if (url.match(/\/v1\/products\/[^/]+$/) && method === 'DELETE') {
      return new Response(JSON.stringify({ success: true, data: null }), { status: 200 })
    }

    // ── Stock ─────────────────────────────────────────────────────────
    if (url.includes('/v1/stock') && method === 'GET' && !url.includes('/alerts') && !url.includes('/customer/')) {
      const data: StockListData = {
        entries: [
          { id: 'st-1', product_id: 'p1', bin_location_id: 'loc-1', quantity: 100, reserved_quantity: 5, lot_number: 'LOT-001', expiry_date: '2027-01-01', status: 'available', last_updated: '2026-05-20T10:00:00Z' },
          { id: 'st-2', product_id: 'p2', bin_location_id: 'loc-2', quantity: 5, reserved_quantity: 2, lot_number: 'LOT-002', expiry_date: '2026-12-01', status: 'available', last_updated: '2026-05-19T10:00:00Z' },
          { id: 'st-3', product_id: 'p3', bin_location_id: 'loc-3', quantity: 0, reserved_quantity: 0, lot_number: '', expiry_date: null, status: 'inactive', last_updated: '2026-05-18T10:00:00Z' },
          { id: 'st-4', product_id: 'p4', bin_location_id: 'loc-4', quantity: 15, reserved_quantity: 10, lot_number: 'LOT-004', expiry_date: '2026-06-15', status: 'reserved', last_updated: '2026-05-20T08:00:00Z' },
          { id: 'st-5', product_id: 'p5', bin_location_id: 'loc-1', quantity: 200, reserved_quantity: 0, lot_number: 'LOT-005', expiry_date: null, status: 'available', last_updated: '2026-05-17T10:00:00Z' },
        ],
        page_size: 20,
      }
      return new Response(JSON.stringify({ success: true, data }), { status: 200 })
    }

    if (url.includes('/v1/stock') && method === 'POST') {
      return new Response(JSON.stringify({
        success: true,
        data: { id: 'st-new', product_id: 'p1', bin_location_id: 'loc-1', quantity: 50, reserved_quantity: 0, lot_number: 'LOT-NEW', expiry_date: null, status: 'available', last_updated: new Date().toISOString() },
      }), { status: 200 })
    }

    if (url.includes('/v1/stock/alerts/low') && method === 'GET') {
      const data: LowStockAlertData = {
        alerts: [
          { product: { id: 'p2', sku: 'SKU-002', name: 'Gadget Beta', description: '', category: 'Electronics', unit: 'piece', weight_kg: 0.8, is_active: true, low_stock_threshold: 10, created_at: '2026-01-15T00:00:00Z', updated_at: '2026-05-10T00:00:00Z' }, current_quantity: 5, threshold: 10, status: 'critical' },
          { product: { id: 'p3', sku: 'SKU-003', name: 'Component Gamma', description: '', category: 'Automotive', unit: 'piece', weight_kg: 2.0, is_active: true, low_stock_threshold: 20, created_at: '2026-02-01T00:00:00Z', updated_at: '2026-04-20T00:00:00Z' }, current_quantity: 0, threshold: 20, status: 'critical' },
        ],
      }
      return new Response(JSON.stringify({ success: true, data }), { status: 200 })
    }

    if (url.match(/\/v1\/stock\/[^/]+\/adjust$/) && method === 'POST') {
      return new Response(JSON.stringify({ success: true, data: { id: url.split('/').slice(-2, -1)[0], quantity: 95, status: 'available' } }), { status: 200 })
    }

    if (url.match(/\/v1\/stock\/[^/]+\/reserve$/) && method === 'POST') {
      return new Response(JSON.stringify({ success: true, data: { id: url.split('/').slice(-2, -1)[0], status: 'reserved' } }), { status: 200 })
    }

    if (url.match(/\/v1\/stock\/[^/]+\/release$/) && method === 'POST') {
      return new Response(JSON.stringify({ success: true, data: { id: url.split('/').slice(-2, -1)[0], status: 'available' } }), { status: 200 })
    }

    // ── Bin Locations ──────────────────────────────────────────────────
    if (url.includes('/v1/bin-locations') && method === 'GET') {
      const data: BinLocationListData = {
        locations: [
          { id: 'loc-1', warehouse_zone: 'A', aisle: 'A1', rack: 'R01', shelf: 'S1', is_active: true },
          { id: 'loc-2', warehouse_zone: 'A', aisle: 'A1', rack: 'R01', shelf: 'S2', is_active: true },
          { id: 'loc-3', warehouse_zone: 'B', aisle: 'B2', rack: 'R03', shelf: 'S1', is_active: true },
          { id: 'loc-4', warehouse_zone: 'C', aisle: 'C1', rack: 'R02', shelf: 'S3', is_active: true },
          { id: 'loc-5', warehouse_zone: 'A', aisle: 'A3', rack: 'R01', shelf: 'S2', is_active: false },
        ],
        page_size: 20,
      }
      return new Response(JSON.stringify({ success: true, data }), { status: 200 })
    }

    if (url.includes('/v1/bin-locations') && method === 'POST') {
      return new Response(JSON.stringify({
        success: true,
        data: { id: 'loc-new', warehouse_zone: 'A', aisle: 'A4', rack: 'R05', shelf: 'S1', is_active: true },
      }), { status: 200 })
    }

    // ── Customer stock ─────────────────────────────────────────────────
    if (url.includes('/v1/stock/customer/') && method === 'GET') {
      const customerId = url.split('/').pop()
      const data: StockListData = {
        entries: [
          { id: 'cst-1', product_id: 'p1', bin_location_id: 'loc-1', quantity: 100, reserved_quantity: 5, lot_number: 'LOT-001', expiry_date: '2027-01-01', status: 'available', last_updated: '2026-05-20T10:00:00Z' },
          { id: 'cst-2', product_id: 'p2', bin_location_id: 'loc-2', quantity: 5, reserved_quantity: 2, lot_number: 'LOT-002', expiry_date: '2026-12-01', status: 'available', last_updated: '2026-05-19T10:00:00Z' },
        ],
        page_size: 20,
      }
      return new Response(JSON.stringify({ success: true, data }), { status: 200 })
    }

    return originalFetch(input, init)
  }
}

export {}
