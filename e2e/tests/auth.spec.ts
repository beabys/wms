import { test, expect } from '@playwright/test'

test.describe('Authentication', () => {
  test('login with valid admin credentials redirects to dashboard', async ({ page }) => {
    await page.goto('/login')
    await page.waitForSelector('text=Sign In')

    await page.getByPlaceholder('you@example.com').fill('admin@wms.local')
    await page.getByPlaceholder('Enter your password').fill('admin123')
    await page.getByRole('button', { name: 'Sign In' }).click()

    await page.waitForURL('**/dashboard')
    await expect(page.getByRole('heading', { name: 'Dashboard' })).toBeVisible()
    await expect(page.getByText('Welcome, System Admin!')).toBeVisible()
  })

  test('login with invalid credentials shows error', async ({ page }) => {
    await page.goto('/login')
    await page.waitForSelector('text=Sign In')

    await page.getByPlaceholder('you@example.com').fill('wrong@email.com')
    await page.getByPlaceholder('Enter your password').fill('wrongpass')
    await page.getByRole('button', { name: 'Sign In' }).click()

    // Error appears (toast or inline message)
    await expect(page.getByText('unauthenticated').first()).toBeVisible({ timeout: 8000 })
    await expect(page).toHaveURL(/\/login/)
  })

  test('login with empty fields stays on login page', async ({ page }) => {
    await page.goto('/login')
    await page.waitForSelector('text=Sign In')

    // Click Sign In with empty fields — form validation should block or stay on page
    await page.getByRole('button', { name: 'Sign In' }).click()
    await expect(page).toHaveURL(/\/login/)
    await expect(page.getByRole('heading', { name: 'Sign In' })).toBeVisible()
  })

  test('logout clears session and redirects to login', async ({ page }) => {
    // Login first
    await page.goto('/login')
    await page.waitForSelector('text=Sign In')
    await page.getByPlaceholder('you@example.com').fill('admin@wms.local')
    await page.getByPlaceholder('Enter your password').fill('admin123')
    await page.getByRole('button', { name: 'Sign In' }).click()
    await page.waitForURL('**/dashboard')

    // Click Logout
    await page.getByRole('button', { name: 'Logout' }).click()
    await page.waitForURL('**/login')
    await expect(page.getByRole('heading', { name: 'Sign In' })).toBeVisible()
  })

  test('unauthenticated user redirected to login from dashboard', async ({ page }) => {
    await page.goto('/dashboard')
    await page.waitForURL('**/login')
    await expect(page.getByRole('heading', { name: 'Sign In' })).toBeVisible()
  })

  test('unauthenticated user redirected to login from /users', async ({ page }) => {
    await page.goto('/users')
    await page.waitForURL('**/login')
    await expect(page.getByRole('heading', { name: 'Sign In' })).toBeVisible()
  })

  test('unauthenticated user redirected to login from /customers', async ({ page }) => {
    await page.goto('/customers')
    await page.waitForURL('**/login')
    await expect(page.getByRole('heading', { name: 'Sign In' })).toBeVisible()
  })
})
