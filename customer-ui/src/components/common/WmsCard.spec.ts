import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import WmsCard from './WmsCard.vue'

describe('WmsCard', () => {
  it('renders default slot', () => {
    const wrapper = mount(WmsCard, { slots: { default: 'Content' } })
    expect(wrapper.text()).toContain('Content')
  })

  it('renders header and footer', () => {
    const wrapper = mount(WmsCard, { slots: { header: 'H', default: 'B', footer: 'F' } })
    expect(wrapper.find('.wms-card__header').text()).toBe('H')
    expect(wrapper.find('.wms-card__footer').text()).toBe('F')
  })

  it('applies padding class', () => {
    const wrapper = mount(WmsCard, { props: { padding: 'lg' } })
    expect(wrapper.classes()).toContain('wms-card--lg')
  })

  it('defaults to md padding', () => {
    const wrapper = mount(WmsCard, { slots: { default: 'C' } })
    expect(wrapper.classes()).toContain('wms-card--md')
  })
})
