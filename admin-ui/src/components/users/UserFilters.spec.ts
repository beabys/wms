import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import UserFilters from './UserFilters.vue'

describe('UserFilters', () => {
  it('renders role select with all options', () => {
    const wrapper = mount(UserFilters, {
      props: { modelValue: {} },
    })
    const select = wrapper.find('select')
    expect(select.exists()).toBe(true)
    const options = select.findAll('option')
    expect(options).toHaveLength(4)
    expect(options[0].text()).toBe('All Roles')
    expect(options[1].text()).toBe('Admin')
    expect(options[2].text()).toBe('Warehouse Staff')
    expect(options[3].text()).toBe('Billing Manager')
  })

  it('selects current role value', () => {
    const wrapper = mount(UserFilters, {
      props: { modelValue: { role: 'warehouse_staff' } },
    })
    const select = wrapper.find('select').element as HTMLSelectElement
    expect(select.value).toBe('warehouse_staff')
  })

  it('emits update:modelValue on role change', async () => {
    const wrapper = mount(UserFilters, {
      props: { modelValue: {} },
    })
    const select = wrapper.find('select')
    await select.setValue('admin')
    expect(wrapper.emitted('update:modelValue')).toBeTruthy()
  })

  it('disables select when loading', () => {
    const wrapper = mount(UserFilters, {
      props: { modelValue: {}, loading: true },
    })
    expect(wrapper.find('select').attributes('disabled')).toBeDefined()
  })
})
