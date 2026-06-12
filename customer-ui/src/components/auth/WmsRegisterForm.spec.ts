import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import WmsRegisterForm from './WmsRegisterForm.vue'

describe('WmsRegisterForm', () => {
  function createWrapper(props = {}) {
    setActivePinia(createPinia())
    return mount(WmsRegisterForm, {
      props,
      global: { stubs: { WmsButton: true, WmsInput: true, WmsPasswordInput: true } },
    })
  }

  it('renders title and fields', () => {
    const wrapper = createWrapper()
    expect(wrapper.text()).toContain('Create Your Account')
  })

  it('renders invite token when provided', () => {
    const wrapper = createWrapper({ inviteToken: 'invite-abc' })
    expect(wrapper.exists()).toBe(true)
  })

  it('emits submit event with form data', async () => {
    const wrapper = createWrapper({ inviteToken: 'invite-123' })
    ;(wrapper.vm as any).companyName = 'Acme Corp'
    ;(wrapper.vm as any).email = 'admin@acme.com'
    ;(wrapper.vm as any).password = 'strongPass1'
    await wrapper.vm.$nextTick()
    await wrapper.find('form').trigger('submit.prevent')
    expect(wrapper.emitted('submit')).toBeTruthy()
    const emitted = (wrapper.emitted('submit')![0][0] as any) as { company_name: string; token: string }
    expect(emitted.company_name).toBe('Acme Corp')
    expect(emitted.token).toBe('invite-123')
  })
})
