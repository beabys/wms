import { test, expect } from '@playwright/test'

test.describe('Admin Customer Management', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/login')
    await page.waitForSelector('text=Sign In')
    await page.getByPlaceholder('you@example.com').fill('admin@wms.local')
    await page.getByPlaceholder('Enter your password').fill('admin123')
    await page.getByRole('button', { name: 'Sign In' }).click()
    await page.waitForURL('**/dashboard')
  })

  test('customer list page loads with table columns', async ({ page }) => {
    await page.getByRole('navigation').getByRole('button', { name: 'Customers' }).click()
    await page.waitForURL('**/customers')
    await page.waitForSelector('text=COMPANY')

    await expect(page.getByRole('heading', { name: 'Customers' })).toBeVisible()
    await expect(page.getByText('COMPANY').first()).toBeVisible()
    await expect(page.getByText('EMAIL').first()).toBeVisible()
    await expect(page.getByText('STATUS').first()).toBeVisible()
    await expect(page.getByText('CREATED').first()).toBeVisible()
    await expect(page.getByText('ACTIONS').first()).toBeVisible()
  })

  test('status filter dropdown renders all options', async ({ page }) => {
    await page.getByRole('navigation').getByRole('button', { name: 'Customers' }).click()
    await page.waitForURL('**/customers')
    await page.waitForSelector('text=COMPANY')

    // Check select element exists with all expected options (PRD: All, Pending, Active, Rejected, Suspended)
    const statusSelect = page.getByRole('combobox', { name: /status/i })
    await expect(statusSelect).toBeVisible()

    // Use evaluate to check option texts in native select
    const options = await statusSelect.evaluate((select) => {
      const sel = select as HTMLSelectElement
      return Array.from(sel.options).map(o => o.text)
    })
    expect(options).toContain('All')
    expect(options).toContain('Pending')
    expect(options).toContain('Active')
    expect(options).toContain('Rejected')
    expect(options).toContain('Suspended')
  })

  test('Invite Customer button exists on customer list page', async ({ page }) => {
    await page.getByRole('navigation').getByRole('button', { name: 'Customers' }).click()
    await page.waitForURL('**/customers')

    await expect(page.getByRole('main').getByRole('button', { name: 'Invite Customer' })).toBeVisible()
  })
})
