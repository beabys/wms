import { test, expect } from '@playwright/test'

const TEST_EMAIL = `e2e-cancel-${Date.now()}@test.com`

test.describe.serial('Invite Cancel Flow — Real API', () => {
  test('Step 1: Admin creates invite for test email', async ({ page }) => {
    // Login as admin
    await page.goto('/login')
    await page.waitForSelector('text=Sign In')
    await page.getByPlaceholder('you@example.com').fill('admin@wms.local')
    await page.getByPlaceholder('Enter your password').fill('admin123')
    await page.getByRole('button', { name: 'Sign In' }).click()
    await page.waitForURL('**/dashboard')

    // Navigate to invite page
    await page.getByRole('navigation').getByRole('button', { name: 'Invite Customer' }).click()
    await page.waitForURL('**/customers/invite')
    await page.waitForSelector('text=Sent Invitations')

    // Click "Send New Invite" to open modal
    await page.getByRole('button', { name: /send new invite/i }).click()
    await page.waitForSelector('text=Invite Customer')

    // Fill email and submit
    await page.getByPlaceholder('customer@company.com').fill(TEST_EMAIL)
    await page.getByRole('button', { name: 'Send Invitation' }).click()

    // Wait for success toast
    await expect(page.getByText('Invitation sent!')).toBeVisible({ timeout: 10000 })

    // Close modal
    await page.getByRole('button', { name: 'Done' }).click()
    await page.waitForTimeout(1000)
  })

  test('Step 2: Cancel the invite from invite list', async ({ page }) => {
    // Login as admin
    await page.goto('/login')
    await page.waitForSelector('text=Sign In')
    await page.getByPlaceholder('you@example.com').fill('admin@wms.local')
    await page.getByPlaceholder('Enter your password').fill('admin123')
    await page.getByRole('button', { name: 'Sign In' }).click()
    await page.waitForURL('**/dashboard')

    // Navigate to invite page
    await page.getByRole('navigation').getByRole('button', { name: 'Invite Customer' }).click()
    await page.waitForURL('**/customers/invite')
    await page.waitForSelector('text=Sent Invitations')

    // Filter to pending to find our invite
    const statusFilter = page.getByRole('combobox', { name: /status/i })
    await statusFilter.selectOption({ label: 'Pending' })
    await page.getByRole('button', { name: 'Apply' }).click()
    await page.waitForTimeout(2000)

    // Find the row with our test email
    const row = page.locator('table tbody tr').filter({ hasText: TEST_EMAIL })
    await expect(row).toBeVisible({ timeout: 10000 })

    // The row should have a Cancel button
    const cancelBtn = row.getByRole('button', { name: /cancel/i })
    await expect(cancelBtn).toBeVisible()

    // Click Cancel — opens WmsConfirmDialog modal
    await cancelBtn.click()

    // Wait for WmsConfirmDialog modal
    await expect(page.getByRole('heading', { name: 'Cancel Invitation' })).toBeVisible({ timeout: 5000 })
    await expect(page.getByRole('button', { name: 'Cancel Invitation' })).toBeVisible()

    // Click Confirm to cancel the invite
    await page.getByRole('button', { name: 'Cancel Invitation' }).click()

    // Wait for success toast
    await expect(page.getByText('Invitation cancelled')).toBeVisible({ timeout: 15000 })

    // List should refresh — verify status changed
    await page.waitForTimeout(2000)

    // Now filter to show cancelled
    await statusFilter.selectOption({ label: 'Cancelled' })
    await page.getByRole('button', { name: 'Apply' }).click()
    await page.waitForTimeout(2000)

    // Our invite should now appear with Cancelled status
    const cancelledRow = page.locator('table tbody tr').filter({ hasText: TEST_EMAIL })
    await expect(cancelledRow).toBeVisible({ timeout: 10000 })
  })

  test('Step 3: Can create new invite for same email after cancellation', async ({ page }) => {
    // Login as admin
    await page.goto('/login')
    await page.waitForSelector('text=Sign In')
    await page.getByPlaceholder('you@example.com').fill('admin@wms.local')
    await page.getByPlaceholder('Enter your password').fill('admin123')
    await page.getByRole('button', { name: 'Sign In' }).click()
    await page.waitForURL('**/dashboard')

    // Navigate to invite page
    await page.getByRole('navigation').getByRole('button', { name: 'Invite Customer' }).click()
    await page.waitForURL('**/customers/invite')
    await page.waitForSelector('text=Sent Invitations')

    // Send new invite for same email
    await page.getByRole('button', { name: /send new invite/i }).click()
    await page.waitForSelector('text=Invite Customer')

    await page.getByPlaceholder('customer@company.com').fill(TEST_EMAIL)
    await page.getByRole('button', { name: 'Send Invitation' }).click()

    // Should succeed — no duplicate pending invite error after cancellation
    await expect(page.getByText('Invitation sent!')).toBeVisible({ timeout: 10000 })

    // Close modal
    await page.getByRole('button', { name: 'Done' }).click()
    await page.waitForTimeout(500)
  })
})
