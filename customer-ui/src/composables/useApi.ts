export class ApiError extends Error {
  status: number
  constructor(message: string, status: number) {
    super(message)
    this.name = 'ApiError'
    this.status = status
  }
}

interface ApiResponse<T> {
  success: boolean
  data?: T
  error?: string
}

let getAccessToken: () => string | null = () => null
let getRefreshToken: () => string | null = () => null
let onRefresh: () => Promise<string | null> = async () => null
let onLogout: () => void = () => {}

export function configureApi(config: {
  getAccessToken: () => string | null
  getRefreshToken: () => string | null
  onRefresh: () => Promise<string | null>
  onLogout: () => void
}) {
  getAccessToken = config.getAccessToken
  getRefreshToken = config.getRefreshToken
  onRefresh = config.onRefresh
  onLogout = config.onLogout
}

export async function request<T>(
  url: string,
  options: RequestInit = {},
): Promise<T> {
  const token = getAccessToken()
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options.headers as Record<string, string>),
  }
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }

  let res = await fetch(url, { ...options, headers })

  // Auto-refresh on 401
  if (res.status === 401 && getRefreshToken()) {
    const newToken = await onRefresh()
    if (newToken) {
      headers['Authorization'] = `Bearer ${newToken}`
      res = await fetch(url, { ...options, headers })
    } else {
      onLogout()
      throw new ApiError('Session expired', 401)
    }
  }

  const body: ApiResponse<T> = await res.json()
  if (!body.success) {
    throw new ApiError(body.error || 'Unknown error', res.status)
  }
  return body.data as T
}
