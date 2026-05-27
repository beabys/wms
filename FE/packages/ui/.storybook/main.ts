import type { StorybookConfig } from '@storybook/vue3-vite'
import tailwindcss from '@tailwindcss/vite'

const config: StorybookConfig = {
  stories: ['../src/stories/**/*.stories.ts'],
  addons: [
    '@storybook/addon-essentials',
    '@storybook/addon-links'
  ],
  framework: {
    name: '@storybook/vue3-vite',
    options: {}
  },
  async viteFinal(config) {
    config.plugins = [...(config.plugins || []), tailwindcss()]
    return config
  }
}

export default config
