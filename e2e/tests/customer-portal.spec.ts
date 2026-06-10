import { test, expect } from '@playwright/test'

const CUSTOMER_UI = 'http://localhost:5174'

test.describe('Customer Portal', () => {
  test('customer login page renders', async ({ page }) => {
    await page.goto(`${CUSTOMER_UI}/login`)
    await page.waitForSelector('text=Welcome Back')

    await expect(page).toHaveTitle('WMS Customer Portal')
    await expect(page.getByRole('heading', { name: 'Welcome Back' })).toBeVisible()
    await expect(page.getByPlaceholder('you@example.com')).toBeVisible()
    await expect(page.getByPlaceholder('Enter your password')).toBeVisible()
    await expect(page.getByRole('button', { name: 'Sign In' })).toBeVisible()
    await expect(page.getByRole('link', { name: /register here/i })).toBeVisible()
  })

  test('customer register page renders with all fields', async ({ page }) => {
    await page.goto(`${CUSTOMER_UI}/register`)
    await page.waitForSelector('text=Create Your Account')

    await expect(page).toHaveTitle('WMS Customer Portal')
    await expect(page.getByRole('heading', { name: 'Create Your Account' })).toBeVisible()
    await expect(page.getByPlaceholder('Your company name')).toBeVisible()
    await expect(page.getByPlaceholder('you@example.com')).toBeVisible()
    await expect(page.getByPlaceholder('Create a password')).toBeVisible()
    await expect(page.getByRole('button', { name: 'Create Account' })).toBeVisible()
    await expect(page.getByRole('link', { name: /sign in/i })).toBeVisible()
  })

  test('customer portal shows correct page title', async ({ page }) => {
    await page.goto(`${CUSTOMER_UI}/login`)
    await expect(page).toHaveTitle('WMS Customer Portal')
  })
})
