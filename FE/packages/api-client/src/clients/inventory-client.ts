import { ApiClient } from './base'
import type {
  Product,
  StockEntry,
  BinLocation,
  CreateProductRequest,
  UpdateProductRequest,
  AddStockRequest,
  AdjustStockRequest,
  ReserveStockRequest,
  CreateBinLocationRequest,
  LowStockAlert,
  ProductListData,
  StockListData,
  BinLocationListData,
  LowStockAlertData,
} from '../types/inventory'

export class InventoryClient extends ApiClient {
  // ── Products ───────────────────────────────────────

  /** GET /v1/products */
  async listProducts(params?: {
    page?: number
    page_size?: number
    search?: string
    category?: string
  }): Promise<ProductListData> {
    const query = new URLSearchParams()
    if (params?.page) query.set('page', String(params.page))
    if (params?.page_size) query.set('page_size', String(params.page_size))
    if (params?.search) query.set('search', params.search)
    if (params?.category) query.set('category', params.category)
    const qs = query.toString()
    return this.request<ProductListData>(
      `/v1/products${qs ? `?${qs}` : ''}`,
      { method: 'GET' }
    )
  }

  /** POST /v1/products */
  async createProduct(data: CreateProductRequest): Promise<Product> {
    return this.request<Product>('/v1/products', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  /** GET /v1/products/{id} */
  async getProduct(id: string): Promise<Product> {
    return this.request<Product>(`/v1/products/${id}`, { method: 'GET' })
  }

  /** PUT /v1/products/{id} */
  async updateProduct(id: string, data: UpdateProductRequest): Promise<Product> {
    return this.request<Product>(`/v1/products/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    })
  }

  /** DELETE /v1/products/{id} (archive) */
  async archiveProduct(id: string): Promise<void> {
    await this.request<undefined>(`/v1/products/${id}`, { method: 'DELETE' })
  }

  // ── Stock ──────────────────────────────────────────

  /** GET /v1/stock */
  async listStock(params?: {
    page_size?: number
    page_token?: string
    product_id?: string
    bin_location_id?: string
    status?: string
  }): Promise<StockListData> {
    const query = new URLSearchParams()
    if (params?.page_size) query.set('page_size', String(params.page_size))
    if (params?.page_token) query.set('page_token', params.page_token)
    if (params?.product_id) query.set('product_id', params.product_id)
    if (params?.bin_location_id) query.set('bin_location_id', params.bin_location_id)
    if (params?.status) query.set('status', params.status)
    const qs = query.toString()
    return this.request<StockListData>(
      `/v1/stock${qs ? `?${qs}` : ''}`,
      { method: 'GET' }
    )
  }

  /** POST /v1/stock */
  async addStock(data: AddStockRequest): Promise<StockEntry> {
    return this.request<StockEntry>('/v1/stock', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  /** GET /v1/stock/{id} */
  async getStock(id: string): Promise<StockEntry> {
    return this.request<StockEntry>(`/v1/stock/${id}`, { method: 'GET' })
  }

  /** POST /v1/stock/{id}/adjust */
  async adjustStock(id: string, data: AdjustStockRequest): Promise<StockEntry> {
    return this.request<StockEntry>(`/v1/stock/${id}/adjust`, {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  /** POST /v1/stock/{id}/reserve */
  async reserveStock(id: string, data: ReserveStockRequest): Promise<StockEntry> {
    return this.request<StockEntry>(`/v1/stock/${id}/reserve`, {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  /** POST /v1/stock/{id}/release */
  async releaseStock(id: string, quantity?: number): Promise<StockEntry> {
    return this.request<StockEntry>(`/v1/stock/${id}/release`, {
      method: 'POST',
      body: quantity !== undefined ? JSON.stringify({ quantity }) : undefined,
    })
  }

  // ── Bin Locations ──────────────────────────────────

  /** GET /v1/bin-locations */
  async listBinLocations(params?: {
    page_size?: number
    page_token?: string
    warehouse_zone?: string
  }): Promise<BinLocationListData> {
    const query = new URLSearchParams()
    if (params?.page_size) query.set('page_size', String(params.page_size))
    if (params?.page_token) query.set('page_token', params.page_token)
    if (params?.warehouse_zone) query.set('warehouse_zone', params.warehouse_zone)
    const qs = query.toString()
    return this.request<BinLocationListData>(
      `/v1/bin-locations${qs ? `?${qs}` : ''}`,
      { method: 'GET' }
    )
  }

  /** POST /v1/bin-locations */
  async createBinLocation(data: CreateBinLocationRequest): Promise<BinLocation> {
    return this.request<BinLocation>('/v1/bin-locations', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  // ── Low Stock Alerts ───────────────────────────────

  /** GET /v1/stock/alerts/low */
  async getLowStockAlerts(): Promise<LowStockAlertData> {
    return this.request<LowStockAlertData>('/v1/stock/alerts/low', {
      method: 'GET',
    })
  }

  // ── Customer stock (read-only) ─────────────────────

  /** GET /v1/stock/customer/{customerId} */
  async getCustomerStock(
    customerId: string,
    params?: { page_size?: number; page_token?: string }
  ): Promise<StockListData> {
    const query = new URLSearchParams()
    if (params?.page_size) query.set('page_size', String(params.page_size))
    if (params?.page_token) query.set('page_token', params.page_token)
    const qs = query.toString()
    return this.request<StockListData>(
      `/v1/stock/customer/${customerId}${qs ? `?${qs}` : ''}`,
      { method: 'GET' }
    )
  }
}
