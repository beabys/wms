import type { Preview } from '@storybook/vue3'
import { setup } from '@storybook/vue3'
import { createPinia } from 'pinia'
import { createRouter, createWebHashHistory } from 'vue-router'
import '../src/style.css'

const pinia = createPinia()
const router = createRouter({
  history: createWebHashHistory(),
  routes: []
})

setup((app) => {
  app.use(pinia)
  app.use(router)
})

// Dark mode decorator — reads globalType and toggles .dark on root
const darkModeDecorator = (story: any, context: any) => {
  const isDark = context.globals.theme === 'dark'
  document.documentElement.classList.toggle('dark', isDark)
  return story()
}

const preview: Preview = {
  decorators: [darkModeDecorator],
  parameters: {
    controls: {
      matchers: {
        color: /(background|color)$/i,
        date: /Date$/i
      }
    }
  },
  globalTypes: {
    theme: {
      name: 'Theme',
      description: 'Global theme for components',
      defaultValue: 'light',
      toolbar: {
        icon: 'circlehollow',
        items: [
          { value: 'light', icon: 'circlehollow', title: 'Light' },
          { value: 'dark', icon: 'circle', title: 'Dark' },
        ],
      },
    },
  },
}

export default preview
