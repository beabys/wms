import type { Meta, StoryObj } from '@storybook/vue3'
import { ref } from 'vue'
import UserForm from './UserForm.vue'
import type { CreateUserRequest } from '@/api/types'

const defaultData: CreateUserRequest = { email: '', password: '', name: '', role: '' }

const meta: Meta<typeof UserForm> = {
  title: 'Users/UserForm',
  component: UserForm,
  tags: ['autodocs'],
}

export default meta
type Story = StoryObj<typeof UserForm>

export const Default: Story = {
  render: () => ({
    components: { UserForm },
    setup: () => {
      const data = ref<CreateUserRequest>({ ...defaultData })
      return { data }
    },
    template: '<UserForm v-model="data" />',
  }),
}

export const ValidationErrors: Story = {
  render: () => ({
    components: { UserForm },
    setup: () => {
      const data = ref<CreateUserRequest>({ email: 'bad', password: 'short', name: '', role: '' })
      return { data }
    },
    template: '<UserForm v-model="data" />',
  }),
}

export const Submitting: Story = {
  render: () => ({
    components: { UserForm },
    setup: () => {
      const data = ref<CreateUserRequest>({ email: 'test@test.com', password: 'password123', name: 'Test User', role: 'admin' })
      return { data }
    },
    template: '<UserForm v-model="data" :loading="true" />',
  }),
}

export const ServerError: Story = {
  render: () => ({
    components: { UserForm },
    setup: () => {
      const data = ref<CreateUserRequest>({ email: 'test@test.com', password: 'password123', name: 'Test User', role: 'admin' })
      return { data }
    },
    template: '<UserForm v-model="data" error="Email already exists" />',
  }),
}

export const DarkMode: Story = {
  render: () => ({
    components: { UserForm },
    setup: () => {
      const data = ref<CreateUserRequest>({ ...defaultData })
      return { data }
    },
    template: '<UserForm v-model="data" />',
  }),
  parameters: {
    themes: { theme: 'dark' },
  },
}
