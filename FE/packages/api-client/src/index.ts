// ── Types (handcrafted response types) ──────────────
export type {
  LoginData,
  RegisterData,
  RefreshData,
  MeData,
  CustomerItem,
  CustomerQueueData,
  CustomerDetailData,
  InviteData,
  RegisterCustomerData,
  DashboardData,
  InboundItem,
  InboundItemDetail,
  InboundQueueData,
  CreateInboundData,
  InspectData,
  FlagData,
} from './types/responses'

// ── Inventory types ─────────────────────────────────
export type {
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
} from './types/inventory'

// ── Clients ─────────────────────────────────────────
export { AuthClient } from './clients/auth-client'
export { CustomerClient } from './clients/customer-client'
export { InboundClient } from './clients/inbound-client'
export { InventoryClient } from './clients/inventory-client'
export { ApiClient } from './clients/base'
