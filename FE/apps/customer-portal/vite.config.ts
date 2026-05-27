import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'
import { resolve } from 'path'

export default defineConfig({
  plugins: [vue(), tailwindcss()],
  server: {
    proxy: {
      '/v1/auth': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/v1/inbounds': {
        target: 'http://localhost:8083',
        changeOrigin: true,
      },
      '/v1/stock': {
        target: 'http://localhost:8084',
        changeOrigin: true,
      },
    },
  },
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
      '@wms/ui': resolve(__dirname, '../../packages/ui/src'),
      '@wms/composables': resolve(__dirname, '../../packages/composables/src'),
      '@wms/api-client': resolve(__dirname, '../../packages/api-client/src'),
    },
  },
})
