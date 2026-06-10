import { test, expect } from '@playwright/test'

test.describe('Admin User Management', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/login')
    await page.waitForSelector('text=Sign In')
    await page.getByPlaceholder('you@example.com').fill('admin@wms.local')
    await page.getByPlaceholder('Enter your password').fill('admin123')
    await page.getByRole('button', { name: 'Sign In' }).click()
    await page.waitForURL('**/dashboard')
  })

  test('user list page loads with table columns', async ({ page }) => {
    await page.getByRole('navigation').getByRole('button', { name: 'Users' }).click()
    await page.waitForURL('**/users')
    await page.waitForSelector('text=EMAIL')

    await expect(page.getByRole('heading', { name: 'Users' })).toBeVisible()
    await expect(page.getByText('EMAIL').first()).toBeVisible()
    await expect(page.getByText('NAME').first()).toBeVisible()
    await expect(page.getByText('ROLE').first()).toBeVisible()
    await expect(page.getByText('CREATED').first()).toBeVisible()
  })

  test('Create User button navigates to create page', async ({ page }) => {
    await page.getByRole('navigation').getByRole('button', { name: 'Users' }).click()
    await page.waitForURL('**/users')

    await page.getByRole('main').getByRole('button', { name: 'Create User' }).click()
    await page.waitForURL('**/users/create')
    await expect(page.getByRole('heading', { name: 'Create User' })).toBeVisible()
  })

  test('create user form renders with all fields', async ({ page }) => {
    await page.getByRole('navigation').getByRole('button', { name: 'Create User' }).click()
    await page.waitForURL('**/users/create')
    await page.waitForSelector('text=Create User')

    await expect(page.getByRole('heading', { name: 'Create User' })).toBeVisible()
    await expect(page.getByPlaceholder('user@example.com')).toBeVisible()
    await expect(page.getByPlaceholder('Min. 8 characters')).toBeVisible()
    await expect(page.getByPlaceholder('Full name')).toBeVisible()
    await expect(page.getByRole('combobox', { name: 'Role' })).toBeVisible()
  })

  test('create user form shows validation errors for empty fields', async ({ page }) => {
    await page.getByRole('navigation').getByRole('button', { name: 'Create User' }).click()
    await page.waitForURL('**/users/create')

    // Per PRD: validation errors visible inline for empty/untouched fields
    // Submit button is disabled when form invalid
    const submitBtn = page.getByRole('main').getByRole('button', { name: 'Create User' })
    await expect(submitBtn).toBeDisabled()

    // Validation errors are shown inline per PRD
    await expect(page.getByText('Email is required').first()).toBeVisible()
    await expect(page.getByText('Password is required').first()).toBeVisible()
    await expect(page.getByText('Name is required').first()).toBeVisible()
  })

  test('create user form validates email format', async ({ page }) => {
    await page.getByRole('navigation').getByRole('button', { name: 'Create User' }).click()
    await page.waitForURL('**/users/create')

    // Fill invalid email — "Email is required" should change to "Invalid email format"
    await page.getByPlaceholder('user@example.com').fill('not-an-email')
    await page.getByPlaceholder('Min. 8 characters').fill('password123')
    await page.getByPlaceholder('Full name').fill('Test User')

    // Submit button stays disabled with invalid email format
    const submitBtn = page.getByRole('main').getByRole('button', { name: 'Create User' })
    await expect(submitBtn).toBeDisabled()

    // Email format error should appear instead of "Email is required"
    await expect(page.getByText('Invalid email format').first()).toBeVisible({ timeout: 3000 })
  })

  test('role filter dropdown renders options', async ({ page }) => {
    await page.getByRole('navigation').getByRole('button', { name: 'Users' }).click()
    await page.waitForURL('**/users')
    await page.waitForSelector('text=EMAIL')

    const roleSelect = page.getByRole('combobox', { name: 'Role' })
    await expect(roleSelect).toBeVisible()

    const options = await roleSelect.evaluate((select) => {
      const sel = select as HTMLSelectElement
      return Array.from(sel.options).map(o => o.text)
    })
    expect(options).toContain('All Roles')
    expect(options).toContain('Admin')
    expect(options).toContain('Warehouse Staff')
    expect(options).toContain('Billing Manager')
  })
})
