import { describe, it, expect, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import WmsToast from './WmsToast.vue'
import { useToastStore } from '@/stores/toast'

describe('WmsToast', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('renders nothing when no toasts', () => {
    const wrapper = mount(WmsToast)
    expect(wrapper.find('.wms-toast__item').exists()).toBe(false)
  })

  it('renders toasts from store', () => {
    const store = useToastStore()
    store.addToast('Test message', 'info')
    const wrapper = mount(WmsToast)
    expect(wrapper.text()).toContain('Test message')
  })

  it('renders multiple toasts', () => {
    const store = useToastStore()
    store.addToast('First', 'info')
    store.addToast('Second', 'success')
    const wrapper = mount(WmsToast)
    const items = wrapper.findAll('.wms-toast__item')
    expect(items).toHaveLength(2)
  })

  it('removes toast when close button clicked', async () => {
    const store = useToastStore()
    store.addToast('Dismiss me', 'info')
    const wrapper = mount(WmsToast)
    expect(wrapper.text()).toContain('Dismiss me')

    await wrapper.find('.wms-toast__close').trigger('click')
    expect(wrapper.find('.wms-toast__item').exists()).toBe(false)
  })

  it('applies correct class for success type', () => {
    const store = useToastStore()
    store.addToast('OK', 'success')
    const wrapper = mount(WmsToast)
    expect(wrapper.find('.wms-toast__item--success').exists()).toBe(true)
  })

  it('applies correct class for error type', () => {
    const store = useToastStore()
    store.addToast('Fail', 'error')
    const wrapper = mount(WmsToast)
    expect(wrapper.find('.wms-toast__item--error').exists()).toBe(true)
  })

  it('applies correct class for warning type', () => {
    const store = useToastStore()
    store.addToast('Caution', 'warning')
    const wrapper = mount(WmsToast)
    expect(wrapper.find('.wms-toast__item--warning').exists()).toBe(true)
  })

  it('applies correct class for info type', () => {
    const store = useToastStore()
    store.addToast('Info', 'info')
    const wrapper = mount(WmsToast)
    expect(wrapper.find('.wms-toast__item--info').exists()).toBe(true)
  })
})
