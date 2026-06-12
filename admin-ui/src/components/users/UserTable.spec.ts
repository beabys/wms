import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import UserTable from './UserTable.vue'
import type { UserResponse } from '@/api/types'

const sampleUsers: UserResponse[] = [
  { id: '1', email: 'alice@test.com', name: 'Alice', role: 'admin', created_at: '2024-01-15T10:30:00Z', updated_at: '' },
  { id: '2', email: 'bob@test.com', name: 'Bob', role: 'warehouse_staff', created_at: '2024-02-20T14:00:00Z', updated_at: '' },
]

describe('UserTable', () => {
  it('renders users in table rows', () => {
    const wrapper = mount(UserTable, {
      props: { users: sampleUsers, loading: false },
    })
    const rows = wrapper.findAll('tbody tr')
    expect(rows).toHaveLength(2)
    expect(rows[0].text()).toContain('alice@test.com')
    expect(rows[0].text()).toContain('Alice')
    expect(rows[1].text()).toContain('bob@test.com')
    expect(rows[1].text()).toContain('Bob')
  })

  it('shows loading state', () => {
    const wrapper = mount(UserTable, {
      props: { users: [], loading: true },
    })
    expect(wrapper.text()).toContain('Loading users...')
    expect(wrapper.find('table').exists()).toBe(false)
  })

  it('shows empty state', () => {
    const wrapper = mount(UserTable, {
      props: { users: [], loading: false },
    })
    expect(wrapper.text()).toContain('No users found.')
    expect(wrapper.find('table').exists()).toBe(false)
  })

  it('renders role badge for each user', () => {
    const wrapper = mount(UserTable, {
      props: { users: sampleUsers, loading: false },
    })
    const badges = wrapper.findAllComponents({ name: 'WmsBadge' })
    expect(badges).toHaveLength(2)
    expect(badges[0].text()).toContain('admin')
  })
})
