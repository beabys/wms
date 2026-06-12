export interface ColorPalette {
  primary: string
  primaryHover: string
  primaryLight: string
  bgPrimary: string
  bgSecondary: string
  bgTertiary: string
  sidebarBg: string
  sidebarHover: string
  cardBg: string
  textPrimary: string
  textSecondary: string
  textMuted: string
  border: string
  borderLight: string
  success: string
  warning: string
  error: string
  info: string
  inputBg: string
  inputBorder: string
  inputFocus: string
  overlay: string
}

export const lightPalette: ColorPalette = {
  primary: '#2563eb',
  primaryHover: '#1d4ed8',
  primaryLight: '#dbeafe',
  bgPrimary: '#ffffff',
  bgSecondary: '#f8fafc',
  bgTertiary: '#f1f5f9',
  sidebarBg: '#f1f5f9',
  sidebarHover: '#e2e8f0',
  cardBg: '#ffffff',
  textPrimary: '#0f172a',
  textSecondary: '#64748b',
  textMuted: '#94a3b8',
  border: '#e2e8f0',
  borderLight: '#f1f5f9',
  success: '#16a34a',
  warning: '#d97706',
  error: '#dc2626',
  info: '#2563eb',
  inputBg: '#ffffff',
  inputBorder: '#cbd5e1',
  inputFocus: '#2563eb',
  overlay: 'rgba(0,0,0,0.5)',
}

export const darkPalette: ColorPalette = {
  primary: '#60a5fa',
  primaryHover: '#93c5fd',
  primaryLight: '#1e3a5f',
  bgPrimary: '#0f172a',
  bgSecondary: '#1e293b',
  bgTertiary: '#334155',
  sidebarBg: '#020617',
  sidebarHover: '#1e293b',
  cardBg: '#1e293b',
  textPrimary: '#f8fafc',
  textSecondary: '#94a3b8',
  textMuted: '#64748b',
  border: '#334155',
  borderLight: '#1e293b',
  success: '#4ade80',
  warning: '#fbbf24',
  error: '#f87171',
  info: '#60a5fa',
  inputBg: '#1e293b',
  inputBorder: '#475569',
  inputFocus: '#60a5fa',
  overlay: 'rgba(0,0,0,0.7)',
}
