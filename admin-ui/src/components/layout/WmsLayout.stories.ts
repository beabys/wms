import type { Meta, StoryObj } from '@storybook/vue3'
import WmsLayout from './WmsLayout.vue'
import { createRouter, createWebHistory } from 'vue-router'
import { h } from 'vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', redirect: '/dashboard' },
    { path: '/dashboard', name: 'Dashboard', component: { template: '<div>Dashboard</div>' } },
    { path: '/login', name: 'Login', component: { template: '<div>Login</div>' } },
  ],
})

const meta: Meta<typeof WmsLayout> = {
  title: 'Layout/WmsLayout',
  component: WmsLayout,
  tags: ['autodocs'],
  decorators: [
    () => ({
      router,
      render() {
        return h('div', { style: 'min-height: 400px' }, [
          h(WmsLayout, {}, {
            default: () => h('div', { style: 'padding: 2rem' }, 'Page content goes here'),
          }),
        ])
      },
    }),
  ],
}

export default meta
type Story = StoryObj<typeof WmsLayout>

export const Default: Story = {}

export const DarkMode: Story = {
  parameters: {
    themes: { default: 'dark' },
  },
}
