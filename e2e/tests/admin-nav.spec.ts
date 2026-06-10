import { test, expect } from '@playwright/test'

test.describe('Admin Navigation / Sidebar', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/login')
    await page.waitForSelector('text=Sign In')
    await page.getByPlaceholder('you@example.com').fill('admin@wms.local')
    await page.getByPlaceholder('Enter your password').fill('admin123')
    await page.getByRole('button', { name: 'Sign In' }).click()
    await page.waitForURL('**/dashboard')
    await page.waitForSelector('text=Dashboard')
  })

  test('sidebar shows WMS Admin brand', async ({ page }) => {
    await expect(page.getByText('WMS').first()).toBeVisible()
    await expect(page.getByText('Admin').first()).toBeVisible()
  })

  test('all nav links render in sidebar', async ({ page }) => {
    const nav = page.getByRole('navigation')
    await expect(nav.getByRole('button', { name: 'Dashboard' })).toBeVisible()
    await expect(nav.getByRole('button', { name: 'Users' })).toBeVisible()
    await expect(nav.getByRole('button', { name: 'Create User' })).toBeVisible()
    await expect(nav.getByRole('button', { name: 'Customers' })).toBeVisible()
    await expect(nav.getByRole('button', { name: 'Invite Customer' })).toBeVisible()
  })

  test('nav links navigate to correct routes', async ({ page }) => {
    // Test Users link
    await page.getByRole('navigation').getByRole('button', { name: 'Users' }).click()
    await page.waitForURL('**/users')
    await expect(page.getByRole('heading', { name: 'Users' })).toBeVisible()

    // Test Create User link
    await page.getByRole('navigation').getByRole('button', { name: 'Create User' }).click()
    await page.waitForURL('**/users/create')
    await expect(page.getByRole('heading', { name: 'Create User' })).toBeVisible()

    // Test Dashboard link
    await page.getByRole('navigation').getByRole('button', { name: 'Dashboard' }).click()
    await page.waitForURL('**/dashboard')
    await expect(page.getByRole('heading', { name: 'Dashboard' })).toBeVisible()

    // Test Customers link
    await page.getByRole('navigation').getByRole('button', { name: 'Customers' }).click()
    await page.waitForURL('**/customers')
    await expect(page.getByRole('heading', { name: 'Customers' })).toBeVisible()
  })

  test('theme toggle exists in sidebar', async ({ page }) => {
    await expect(page.getByRole('button', { name: /switch to (dark|light) mode/i })).toBeVisible()
  })

  test('logout button exists in sidebar', async ({ page }) => {
    await expect(page.getByRole('button', { name: 'Logout' })).toBeVisible()
  })

  test('user email displayed in sidebar', async ({ page }) => {
    await expect(page.getByText('admin@wms.local').first()).toBeVisible()
  })
})
