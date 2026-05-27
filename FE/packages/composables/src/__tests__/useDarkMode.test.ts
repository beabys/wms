import { describe, it, expect, beforeEach, vi } from 'vitest'

describe('useDarkMode', () => {
  beforeEach(() => {
    localStorage.clear()
    document.documentElement.classList.remove('dark')
    vi.resetModules()
  })

  it('defaults to light mode when no preference stored', async () => {
    const { useDarkMode } = await import('../useDarkMode')
    const { isDark } = useDarkMode()
    expect(isDark.value).toBe(false)
  })

  it('reads stored dark preference from localStorage', async () => {
    localStorage.setItem('wms_dark_mode', 'true')
    const { useDarkMode } = await import('../useDarkMode')
    const { isDark } = useDarkMode()
    expect(isDark.value).toBe(true)
  })

  it('toggle switches the value', async () => {
    const { useDarkMode } = await import('../useDarkMode')
    const { isDark, toggle } = useDarkMode()
    const before = isDark.value
    toggle()
    expect(isDark.value).toBe(!before)
  })

  it('toggle persists to localStorage', async () => {
    const { useDarkMode } = await import('../useDarkMode')
    const { toggle } = useDarkMode()
    toggle()
    const val = localStorage.getItem('wms_dark_mode')
    expect(['true', 'false']).toContain(val)
  })

  it('toggle adds/removes dark class on html', async () => {
    const { useDarkMode } = await import('../useDarkMode')
    const { isDark, toggle } = useDarkMode()
    const beforeClass = document.documentElement.classList.contains('dark')
    expect(beforeClass).toBe(isDark.value)
    toggle()
    expect(document.documentElement.classList.contains('dark')).toBe(isDark.value)
  })
})
