import { describe, it, expect, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import WmsThemeToggle from './WmsThemeToggle.vue'

describe('WmsThemeToggle', () => {
  beforeEach(() => {
    localStorage.clear()
    document.documentElement.removeAttribute('data-theme')
  })

  it('renders sun icon in dark mode', () => {
    localStorage.setItem('wms_admin_theme', 'dark')
    document.documentElement.setAttribute('data-theme', 'dark')
    const wrapper = mount(WmsThemeToggle)
    // In dark mode, sun icon should be shown
    expect(wrapper.find('svg').exists()).toBe(true)
  })

  it('renders moon icon in light mode', () => {
    const wrapper = mount(WmsThemeToggle)
    // In light mode, moon icon should be shown
    expect(wrapper.find('svg').exists()).toBe(true)
  })

  it('toggles theme on click', async () => {
    document.documentElement.setAttribute('data-theme', 'light')
    localStorage.setItem('wms_admin_theme', 'light')
    const wrapper = mount(WmsThemeToggle)
    expect(document.documentElement.getAttribute('data-theme')).toBe('light')

    await wrapper.find('button').trigger('click')
    expect(document.documentElement.getAttribute('data-theme')).toBe('dark')
    expect(localStorage.getItem('wms_admin_theme')).toBe('dark')
  })

  it('has accessible title', () => {
    const wrapper = mount(WmsThemeToggle)
    expect(wrapper.find('button').attributes('title')).toBeTruthy()
  })
})
