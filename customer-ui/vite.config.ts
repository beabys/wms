import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
    },
  },
  server: {
    port: 5174,
    proxy: {
      '/api/v1/auth': {
        target: 'http://localhost:8081',
        changeOrigin: true,
        rewrite: (path) => path.replace('/api/v1', '/v1'),
      },
      '/api/v1/users': {
        target: 'http://localhost:8081',
        changeOrigin: true,
        rewrite: (path) => path.replace('/api/v1', '/v1'),
      },
      '/api/v1/roles': {
        target: 'http://localhost:8081',
        changeOrigin: true,
        rewrite: (path) => path.replace('/api/v1', '/v1'),
      },
      '/api/v1/customers': {
        target: 'http://localhost:8082',
        changeOrigin: true,
        rewrite: (path) => path.replace('/api/v1', '/v1'),
      },
      '/api/v1/companies': {
        target: 'http://localhost:8082',
        changeOrigin: true,
        rewrite: (path) => path.replace('/api/v1', '/v1'),
      },
    },
  },
})
