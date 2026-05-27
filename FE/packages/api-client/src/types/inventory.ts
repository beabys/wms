export interface Product {
  id: string
  sku: string
  name: string
  description: string
  category: string
  unit: string
  weight_kg: number
  is_active: boolean
  low_stock_threshold: number
  created_at: string
  updated_at: string
}

export interface StockEntry {
  id: string
  product_id: string
  bin_location_id: string
  quantity: number
  reserved_quantity: number
  lot_number: string
  expiry_date: string | null
  status: string
  last_updated: string
}

export interface BinLocation {
  id: string
  warehouse_zone: string
  aisle: string
  rack: string
  shelf: string
  is_active: boolean
}

export interface CreateProductRequest {
  sku: string
  name: string
  description?: string
  category?: string
  unit?: string
  weight_kg?: number
  low_stock_threshold?: number
}

export interface UpdateProductRequest {
  name?: string
  description?: string
  category?: string
  unit?: string
  weight_kg?: number
  low_stock_threshold?: number
}

export interface AddStockRequest {
  product_id: string
  bin_location_id?: string
  quantity: number
  lot_number?: string
  expiry_date?: string
}

export interface AdjustStockRequest {
  delta: number
  reason: string
  notes?: string
}

export interface ReserveStockRequest {
  quantity: number
  order_id?: string
}

export interface CreateBinLocationRequest {
  warehouse_zone: string
  aisle: string
  rack: string
  shelf: string
}

export interface LowStockAlert {
  product: Product
  current_quantity: number
  threshold: number
  status: 'critical' | 'warning'
}

export interface ProductListData {
  products: Product[]
  total: number
  page: number
  page_size: number
}

export interface StockListData {
  entries: StockEntry[]
  page_size: number
  page_token?: string
}

export interface BinLocationListData {
  locations: BinLocation[]
  page_size: number
  page_token?: string
}

export interface LowStockAlertData {
  alerts: LowStockAlert[]
}
