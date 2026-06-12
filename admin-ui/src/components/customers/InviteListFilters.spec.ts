import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import InviteListFilters from './InviteListFilters.vue'

describe('InviteListFilters', () => {
  function createWrapper(modelValue = {}) {
    return mount(InviteListFilters, {
      props: { modelValue },
    })
  }

  it('renders all filter controls', () => {
    const wrapper = createWrapper()
    expect(wrapper.find('select').exists()).toBe(true)
    const selects = wrapper.findAll('select')
    expect(selects.length).toBe(2)
    const dateInputs = wrapper.findAll('input[type="date"]')
    expect(dateInputs.length).toBe(2)
    expect(wrapper.text()).toContain('Apply')
    expect(wrapper.text()).toContain('Reset')
  })

  it('emits status filter on Apply click', async () => {
    const wrapper = createWrapper()
    const selects = wrapper.findAll('select')
    // Set status filter to "pending"
    await selects[0].setValue('pending')
    const buttons = wrapper.findAllComponents({ name: 'WmsButton' })
    const applyBtn = buttons.find(b => b.text() === 'Apply')
    await applyBtn?.trigger('click')
    expect(wrapper.emitted('update:modelValue')).toBeTruthy()
    const emitted = wrapper.emitted('update:modelValue')![0][0] as any
    expect(emitted.status).toBe('pending')
  })

  it('emits empty filters on Reset click', async () => {
    const wrapper = createWrapper()
    const buttons = wrapper.findAllComponents({ name: 'WmsButton' })
    const resetBtn = buttons.find(b => b.text() === 'Reset')
    await resetBtn?.trigger('click')
    const emitted = wrapper.emitted('update:modelValue')![0][0] as any
    expect(emitted).toEqual({})
  })

  it('syncs selects from modelValue prop', async () => {
    const wrapper = createWrapper({ status: 'used', expired: false })
    const selects = wrapper.findAll('select')
    expect((selects[0].element as HTMLSelectElement).value).toBe('used')
    expect((selects[1].element as HTMLSelectElement).value).toBe('false')
  })

  it('emits status=used and expired=true properly', async () => {
    const wrapper = createWrapper()
    const selects = wrapper.findAll('select')
    await selects[0].setValue('used')
    await selects[1].setValue('true')
    const buttons = wrapper.findAllComponents({ name: 'WmsButton' })
    const applyBtn = buttons.find(b => b.text() === 'Apply')
    await applyBtn?.trigger('click')
    const emitted = wrapper.emitted('update:modelValue')![0][0] as any
    expect(emitted.status).toBe('used')
    expect(emitted.expired).toBe(true)
  })

  it('status dropdown has All/Pending/Used/Cancelled options', () => {
    const wrapper = createWrapper()
    const selects = wrapper.findAll('select')
    const statusOptions = selects[0].findAll('option')
    const optionTexts = statusOptions.map(o => o.text())
    expect(optionTexts).toEqual(['All', 'Pending', 'Used', 'Cancelled'])
  })
})
