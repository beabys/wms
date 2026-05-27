import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
  testDir: './tests',
  fullyParallel: false,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: 1,
  reporter: [
    ['html', { outputFolder: 'playwright-report' }],
    ['list'],
  ],
  use: {
    baseURL: 'http://localhost:5173',
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
    /* Collect console errors */
    actionTimeout: 15000,
  },

  projects: [
    {
      name: 'customer-portal',
      use: { ...devices['Desktop Chrome'], baseURL: 'http://localhost:5173' },
      testMatch: /.*customer.*\.spec\.ts/,
    },
    {
      name: 'ops-portal',
      use: { ...devices['Desktop Chrome'], baseURL: 'http://localhost:5174' },
      testMatch: /.*ops.*\.spec\.ts/,
    },
  ],

  /* Run local dev servers before tests */
  webServer: [
    {
      command: 'cd /Users/beabys/go/src/github.com/beabys/wms && go run ./cmd/wms',
      port: 8084,
      timeout: 120 * 1000,
      reuseExistingServer: !process.env.CI,
      cwd: '/Users/beabys/go/src/github.com/beabys/wms',
    },
    {
      command: 'cd /Users/beabys/go/src/github.com/beabys/wms/FE && pnpm --filter @wms/customer-portal dev',
      port: 5173,
      timeout: 60 * 1000,
      reuseExistingServer: !process.env.CI,
      cwd: '/Users/beabys/go/src/github.com/beabys/wms/FE',
      env: { VITE_USE_MOCK: 'false' },
    },
    {
      command: 'cd /Users/beabys/go/src/github.com/beabys/wms/FE && pnpm --filter @wms/ops-portal dev',
      port: 5174,
      timeout: 60 * 1000,
      reuseExistingServer: !process.env.CI,
      cwd: '/Users/beabys/go/src/github.com/beabys/wms/FE',
      env: { VITE_USE_MOCK: 'false' },
    },
  ],
});
