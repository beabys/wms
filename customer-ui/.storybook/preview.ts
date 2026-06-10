import type { Preview } from '@storybook/vue3'
import { setup } from '@storybook/vue3'
import { createPinia, setActivePinia } from 'pinia'
import '@/assets/styles/reset.css'
import '@/assets/styles/variables.css'

// Initialize Pinia for all stories
setup(() => {
  setActivePinia(createPinia())
})

const preview: Preview = {
  parameters: {
    controls: {
      matchers: {
        color: /(background|color)$/i,
        date: /Date$/i,
      },
    },
  },
  globalTypes: {
    theme: {
      name: 'Theme',
      description: 'Global theme for components',
      defaultValue: 'light',
      toolbar: {
        icon: 'circlehollow',
        items: [
          { value: 'light', icon: 'sun', title: 'Light' },
          { value: 'dark', icon: 'moon', title: 'Dark' },
        ],
        dynamicTitle: true,
      },
    },
  },
  decorators: [
    (story, context) => {
      const theme = context.globals.theme || 'light'
      document.documentElement.setAttribute('data-theme', theme)
      return {
        components: { story },
        template: '<div :data-theme="theme" style="padding: 2rem"><story /></div>',
        data: () => ({ theme }),
      }
    },
  ],
}

export default preview
