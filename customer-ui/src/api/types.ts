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

export interface RegisterRequest {
  token?: string
  company_name: string
  email: string
  password: string
  phone?: string
  vat_number?: string
}

export interface RegisterResponse {
  id: string
  company_name: string
  email: string
  status: string
}

export interface CustomerResponse {
  id: string
  company_name: string
  email: string
  phone: string
  vat_number: string
  status: string
  created_at: string
}
