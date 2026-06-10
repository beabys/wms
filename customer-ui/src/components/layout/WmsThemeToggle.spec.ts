import { describe, it, expect, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import WmsThemeToggle from './WmsThemeToggle.vue'

describe('WmsThemeToggle', () => {
  beforeEach(() => {
    localStorage.clear()
    document.documentElement.removeAttribute('data-theme')
  })

  it('renders icon', () => {
    const wrapper = mount(WmsThemeToggle)
    expect(wrapper.find('svg').exists()).toBe(true)
  })

  it('toggles theme on click', async () => {
    document.documentElement.setAttribute('data-theme', 'light')
    localStorage.setItem('wms_customer_theme', 'light')
    const wrapper = mount(WmsThemeToggle)
    await wrapper.find('button').trigger('click')
    expect(document.documentElement.getAttribute('data-theme')).toBe('dark')
    expect(localStorage.getItem('wms_customer_theme')).toBe('dark')
  })

  it('has accessible title', () => {
    const wrapper = mount(WmsThemeToggle)
    expect(wrapper.find('button').attributes('title')).toBeTruthy()
  })
})
