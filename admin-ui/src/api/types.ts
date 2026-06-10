export interface ApiResponse<T> {
  success: boolean
  data?: T
  error?: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface LoginResponse {
  access_token: string
  refresh_token: string
  expires_in: number
}

export interface RefreshRequest {
  refresh_token: string
}

export interface RefreshResponse {
  access_token: string
  refresh_token: string
  expires_in: number
}

export interface UserResponse {
  id: string
  email: string
  name: string
  role: string
  created_at: string
  updated_at: string
}

export interface CreateUserRequest {
  email: string
  password: string
  name: string
  role: string
}

export interface CustomerResponse {
  id: string
  company_name: string
  email: string
  phone: string
  vat_number: string
  status: string
  address?: string
  city?: string
  postal_code?: string
  country?: string
  created_at: string
  updated_at?: string
}

export interface AuditEntry {
  id: string
  customer_id: string
  action: string
  performed_by: string
  details: string
  created_at: number
}

export interface AuditLogListResponse {
  entries: AuditEntry[]
  pagination: Pagination
}

export interface Pagination {
  page: number
  page_size: number
  total_items: number
}

export interface UserListResponse {
  users: UserResponse[]
  pagination: Pagination
}

export interface InviteUserResponse {
  token: string
  invite_link: string
  expires_at: number
}

export interface CustomerListResponse {
  customers: CustomerResponse[]
  pagination: Pagination
}

export interface InviteEntry {
  id: string
  email: string
  token: string
  invited_by: string
  status: string
  expires_at: number
  created_at: number
}

export interface InviteListResponse {
  invites: InviteEntry[]
  pagination: Pagination
}
