import type { Meta, StoryObj } from '@storybook/vue3'
import DataTable, { type ColumnDef } from '../../components/common/DataTable.vue'

const meta: Meta<typeof DataTable> = {
  title: 'Common/DataTable',
  component: DataTable,
  tags: ['autodocs']
}

export default meta
type Story = StoryObj<typeof DataTable>

const columns: ColumnDef[] = [
  { key: 'id', label: 'ID', sortable: true },
  { key: 'name', label: 'Name', sortable: true },
  { key: 'email', label: 'Email', sortable: true },
  { key: 'role', label: 'Role' }
]

const sampleData = [
  { id: '001', name: 'Alice Johnson', email: 'alice@example.com', role: 'Admin' },
  { id: '002', name: 'Bob Smith', email: 'bob@example.com', role: 'User' },
  { id: '003', name: 'Carol White', email: 'carol@example.com', role: 'Editor' },
  { id: '004', name: 'Dave Brown', email: 'dave@example.com', role: 'User' },
  { id: '005', name: 'Eve Davis', email: 'eve@example.com', role: 'Admin' }
]

export const HasData: Story = {
  args: {
    columns,
    data: sampleData,
    loading: false
  }
}

export const Empty: Story = {
  args: {
    columns,
    data: [],
    loading: false,
    emptyText: 'No records found'
  }
}
