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

// ── Clients ─────────────────────────────────────────
export { AuthClient } from './clients/auth-client'
export { CustomerClient } from './clients/customer-client'
export { InboundClient } from './clients/inbound-client'
export { ApiClient } from './clients/base'
