import { test, expect } from '@playwright/test'

test.describe('Toast Notifications', () => {
  test('error toast appears on failed login (admin portal)', async ({ page }) => {
    await page.goto('http://localhost:5173/login')
    await page.waitForSelector('text=Sign In')

    // Submit invalid credentials
    await page.getByPlaceholder('you@example.com').fill('invalid@test.com')
    await page.getByPlaceholder('Enter your password').fill('wrongpass')
    await page.getByRole('button', { name: 'Sign In' }).click()

    // Wait for error toast to appear
    await expect(page.getByText('unauthenticated').first()).toBeVisible({ timeout: 8000 })

    // Verify no raw API error text outside toast container
    const body = page.locator('body')
    await expect(body).not.toContainText('"error"')
    await expect(body).not.toContainText('"success":false')
  })

  test('protected page redirects without showing raw error', async ({ page }) => {
    await page.goto('http://localhost:5173/dashboard')
    await page.waitForURL('**/login')
    await expect(page.getByRole('heading', { name: 'Sign In' })).toBeVisible()

    // No raw API error on the page
    const body = page.locator('body')
    await expect(body).not.toContainText('"error"')
    await expect(body).not.toContainText('"success":false')
  })
})
