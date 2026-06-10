export interface ColorPalette {
  primary: string
  primaryHover: string
  primaryLight: string
  bgPrimary: string
  bgSecondary: string
  bgTertiary: string
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
  primary: '#059669',
  primaryHover: '#047857',
  primaryLight: '#d1fae5',
  bgPrimary: '#fafaf9',
  bgSecondary: '#f5f5f4',
  bgTertiary: '#e7e5e4',
  cardBg: '#ffffff',
  textPrimary: '#1c1917',
  textSecondary: '#78716c',
  textMuted: '#a8a29e',
  border: '#e7e5e4',
  borderLight: '#f5f5f4',
  success: '#16a34a',
  warning: '#d97706',
  error: '#dc2626',
  info: '#059669',
  inputBg: '#ffffff',
  inputBorder: '#d6d3d1',
  inputFocus: '#059669',
  overlay: 'rgba(0,0,0,0.5)',
}

export const darkPalette: ColorPalette = {
  primary: '#34d399',
  primaryHover: '#6ee7b7',
  primaryLight: '#064e3b',
  bgPrimary: '#1c1917',
  bgSecondary: '#292524',
  bgTertiary: '#44403c',
  cardBg: '#292524',
  textPrimary: '#fafaf9',
  textSecondary: '#a8a29e',
  textMuted: '#78716c',
  border: '#44403c',
  borderLight: '#292524',
  success: '#4ade80',
  warning: '#fbbf24',
  error: '#f87171',
  info: '#34d399',
  inputBg: '#292524',
  inputBorder: '#57534e',
  inputFocus: '#34d399',
  overlay: 'rgba(0,0,0,0.7)',
}
