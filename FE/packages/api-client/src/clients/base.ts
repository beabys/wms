/**
 * Base API client — thin fetch wrapper.
 * Unwraps { success, data, error } envelope.
 */
export class ApiClient {
  protected baseURL: string
  protected token: string | null = null

  constructor(baseURL = '') {
    this.baseURL = baseURL
  }

  setToken(token: string | null): void {
    this.token = token
  }

  protected async request<T>(path: string, options: RequestInit = {}): Promise<T> {
    const { body, ...rest } = options
    const headers: Record<string, string> = {
      ...(rest.headers as Record<string, string> | undefined),
    }

    // Default Content-Type for non-GET requests with body
    if (body && !headers['Content-Type']) {
      headers['Content-Type'] = 'application/json'
    }

    if (this.token) {
      headers['Authorization'] = `Bearer ${this.token}`
    }

    const response = await fetch(`${this.baseURL}${path}`, {
      ...rest,
      headers,
      body,
    })

    const json: { success: boolean; data?: T; error?: string | null } = await response.json()

    if (json.success === false) {
      throw new Error(json.error || 'Request failed')
    }

    return json.data as T
  }
}
