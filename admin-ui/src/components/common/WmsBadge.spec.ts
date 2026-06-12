import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import WmsBadge from './WmsBadge.vue'

describe('WmsBadge', () => {
  it('renders slot content', () => {
    const wrapper = mount(WmsBadge, {
      slots: { default: 'Active' },
    })
    expect(wrapper.text()).toBe('Active')
  })

  it('applies variant class', () => {
    const wrapper = mount(WmsBadge, {
      props: { variant: 'success' },
    })
    expect(wrapper.classes()).toContain('wms-badge--success')
  })

  it('applies size class', () => {
    const wrapper = mount(WmsBadge, {
      props: { size: 'md' },
    })
    expect(wrapper.classes()).toContain('wms-badge--md')
  })

  it('defaults to sm size', () => {
    const wrapper = mount(WmsBadge)
    expect(wrapper.classes()).toContain('wms-badge--sm')
  })

  it('defaults to neutral variant', () => {
    const wrapper = mount(WmsBadge)
    expect(wrapper.classes()).toContain('wms-badge--neutral')
  })

  it('renders all variants without error', () => {
    const variants = ['success', 'warning', 'error', 'info', 'neutral'] as const
    for (const variant of variants) {
      const wrapper = mount(WmsBadge, {
        props: { variant },
        slots: { default: variant },
      })
      expect(wrapper.classes()).toContain(`wms-badge--${variant}`)
    }
  })
})
