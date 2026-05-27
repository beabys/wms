import { defineConfig } from 'vitest/config'
import { resolve } from 'path'

export default defineConfig({
  define: {
    'import.meta.env.VITE_API_BASE_URL': '""',
  },
  resolve: {
    alias: {
      '@wms/api-client': resolve(__dirname, '../../packages/api-client/src'),
    },
  },
  test: {
    environment: 'node',
  },
})
