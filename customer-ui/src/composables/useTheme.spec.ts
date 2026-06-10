import { describe, it, expect, beforeEach } from 'vitest'
import { useTheme } from './useTheme'

describe('useTheme', () => {
  beforeEach(() => {
    localStorage.clear()
    document.documentElement.removeAttribute('data-theme')
  })

  it('initializes with light theme when no saved preference', () => {
    const theme = useTheme()
    theme.init()
    expect(theme.isDark.value).toBe(false)
    expect(document.documentElement.getAttribute('data-theme')).toBe('light')
  })

  it('initializes with dark theme when saved preference is dark', () => {
    localStorage.setItem('wms_customer_theme', 'dark')
    const theme = useTheme()
    theme.init()
    expect(theme.isDark.value).toBe(true)
    expect(document.documentElement.getAttribute('data-theme')).toBe('dark')
  })

  it('toggles from light to dark', () => {
    const theme = useTheme()
    theme.init()
    theme.toggle()
    expect(theme.isDark.value).toBe(true)
    expect(document.documentElement.getAttribute('data-theme')).toBe('dark')
    expect(localStorage.getItem('wms_customer_theme')).toBe('dark')
  })

  it('setTheme sets specific theme', () => {
    const theme = useTheme()
    theme.init()
    theme.setTheme('dark')
    expect(theme.isDark.value).toBe(true)
    expect(document.documentElement.getAttribute('data-theme')).toBe('dark')
    theme.setTheme('light')
    expect(theme.isDark.value).toBe(false)
  })
})
