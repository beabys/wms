import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import CustomerStatusBadge from './CustomerStatusBadge.vue'

describe('CustomerStatusBadge', () => {
  it('renders pending status with warning variant', () => {
    const wrapper = mount(CustomerStatusBadge, {
      props: { status: 'pending' },
    })
    expect(wrapper.text()).toBe('Pending')
    const badge = wrapper.findComponent({ name: 'WmsBadge' })
    expect(badge.props('variant')).toBe('warning')
  })

  it('renders active status with success variant', () => {
    const wrapper = mount(CustomerStatusBadge, {
      props: { status: 'active' },
    })
    expect(wrapper.text()).toBe('Active')
    const badge = wrapper.findComponent({ name: 'WmsBadge' })
    expect(badge.props('variant')).toBe('success')
  })

  it('renders rejected status with error variant', () => {
    const wrapper = mount(CustomerStatusBadge, {
      props: { status: 'rejected' },
    })
    expect(wrapper.text()).toBe('Rejected')
    const badge = wrapper.findComponent({ name: 'WmsBadge' })
    expect(badge.props('variant')).toBe('error')
  })

  it('renders suspended status with neutral variant', () => {
    const wrapper = mount(CustomerStatusBadge, {
      props: { status: 'suspended' },
    })
    expect(wrapper.text()).toBe('Suspended')
    const badge = wrapper.findComponent({ name: 'WmsBadge' })
    expect(badge.props('variant')).toBe('neutral')
  })

  it('falls back to neutral for unknown status', () => {
    const wrapper = mount(CustomerStatusBadge, {
      props: { status: 'unknown' },
    })
    expect(wrapper.text()).toBe('Unknown')
    const badge = wrapper.findComponent({ name: 'WmsBadge' })
    expect(badge.props('variant')).toBe('neutral')
  })
})
