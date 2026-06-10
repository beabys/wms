import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import UserForm from './UserForm.vue'
import type { CreateUserRequest } from '@/api/types'

const defaultData: CreateUserRequest = { email: '', password: '', name: '', role: '' }

describe('UserForm', () => {
  it('renders all form fields', () => {
    const wrapper = mount(UserForm, {
      props: { modelValue: { ...defaultData } },
    })
    expect(wrapper.find('input[type="email"]').exists()).toBe(true)
    expect(wrapper.find('input[type="password"]').exists()).toBe(true)
    expect(wrapper.findAll('input')).toHaveLength(3)
    expect(wrapper.find('select').exists()).toBe(true)
    expect(wrapper.find('button[type="submit"]').exists()).toBe(true)
  })

  it('shows validation errors for empty fields', () => {
    const wrapper = mount(UserForm, {
      props: { modelValue: { ...defaultData } },
    })
    // Submit should be disabled (isValid false)
    const btn = wrapper.find('button[type="submit"]')
    expect(btn.attributes('disabled')).toBeDefined()
    // Field errors should show
    expect(wrapper.text()).toContain('Email is required')
    expect(wrapper.text()).toContain('Password is required')
    expect(wrapper.text()).toContain('Name is required')
  })

  it('shows invalid email error', () => {
    const wrapper = mount(UserForm, {
      props: { modelValue: { email: 'bad', password: 'pass1234', name: 'Test', role: 'admin' } },
    })
    expect(wrapper.text()).toContain('Invalid email format')
  })

  it('shows short password error', () => {
    const wrapper = mount(UserForm, {
      props: { modelValue: { email: 'test@test.com', password: 'short', name: 'Test', role: 'admin' } },
    })
    expect(wrapper.text()).toContain('Password must be at least 8 characters')
  })

  it('shows server error when provided', () => {
    const wrapper = mount(UserForm, {
      props: { modelValue: { email: 'test@test.com', password: 'password123', name: 'Test', role: 'admin' }, error: 'Email already exists' },
    })
    expect(wrapper.text()).toContain('Email already exists')
  })

  it('emits update:modelValue on field change', async () => {
    const wrapper = mount(UserForm, {
      props: { modelValue: { ...defaultData } },
    })
    const emailInput = wrapper.find('input[type="email"]')
    await emailInput.setValue('new@test.com')
    expect(wrapper.emitted('update:modelValue')).toBeTruthy()
    const emitted = wrapper.emitted('update:modelValue') as any[][]
    expect(emitted[0][0].email).toBe('new@test.com')
  })

  it('emits submit when form is submitted with valid data', async () => {
    const wrapper = mount(UserForm, {
      props: { modelValue: { email: 'test@test.com', password: 'password123', name: 'Test', role: 'admin' } },
    })
    await wrapper.find('form').trigger('submit.prevent')
    expect(wrapper.emitted('submit')).toBeTruthy()
  })

  it('shows loading state on submit button', () => {
    const wrapper = mount(UserForm, {
      props: { modelValue: { ...defaultData }, loading: true },
    })
    const btn = wrapper.find('button[type="submit"]')
    expect(btn.attributes('disabled')).toBeDefined()
  })
})
