import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import WmsCard from './WmsCard.vue'

describe('WmsCard', () => {
  it('renders default slot content', () => {
    const wrapper = mount(WmsCard, {
      slots: { default: 'Card content' },
    })
    expect(wrapper.text()).toContain('Card content')
  })

  it('renders header slot', () => {
    const wrapper = mount(WmsCard, {
      slots: { header: 'Header' },
    })
    expect(wrapper.find('.wms-card__header').text()).toBe('Header')
  })

  it('renders footer slot', () => {
    const wrapper = mount(WmsCard, {
      slots: { footer: 'Footer' },
    })
    expect(wrapper.find('.wms-card__footer').text()).toBe('Footer')
  })

  it('does not render header when no slot', () => {
    const wrapper = mount(WmsCard, {
      slots: { default: 'Content' },
    })
    expect(wrapper.find('.wms-card__header').exists()).toBe(false)
  })

  it('does not render footer when no slot', () => {
    const wrapper = mount(WmsCard, {
      slots: { default: 'Content' },
    })
    expect(wrapper.find('.wms-card__footer').exists()).toBe(false)
  })

  it('applies padding class', () => {
    const wrapper = mount(WmsCard, {
      props: { padding: 'sm' },
    })
    expect(wrapper.classes()).toContain('wms-card--sm')
  })

  it('defaults to md padding', () => {
    const wrapper = mount(WmsCard, {
      slots: { default: 'Content' },
    })
    expect(wrapper.classes()).toContain('wms-card--md')
  })
})
