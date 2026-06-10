import { test, expect } from '@playwright/test'

test.describe('Invite List Page — Real API', () => {
  test.beforeEach(async ({ page }) => {
    // Login as admin for each test
    await page.goto('/login')
    await page.waitForSelector('text=Sign In')
    await page.getByPlaceholder('you@example.com').fill('admin@wms.local')
    await page.getByPlaceholder('Enter your password').fill('admin123')
    await page.getByRole('button', { name: 'Sign In' }).click()
    await page.waitForURL('**/dashboard')
  })

  test('invite list page loads with table headers', async ({ page }) => {
    await page.getByRole('navigation').getByRole('button', { name: 'Invite Customer' }).click()
    await page.waitForURL('**/customers/invite')
    await page.waitForSelector('text=Sent Invitations')

    // Table headers per plan: email, status, invited_by, created, actions
    await expect(page.getByText('EMAIL').first()).toBeVisible()
    await expect(page.getByText('STATUS').first()).toBeVisible()
    await expect(page.getByText('INVITED BY').first()).toBeVisible()
    await expect(page.getByText('CREATED').first()).toBeVisible()
    await expect(page.getByText('ACTIONS').first()).toBeVisible()
  })

  test('invite list renders invite data from backend', async ({ page }) => {
    await page.getByRole('navigation').getByRole('button', { name: 'Invite Customer' }).click()
    await page.waitForURL('**/customers/invite')
    await page.waitForSelector('text=Sent Invitations')

    // Verify at least one invite email appears in the table (from existing data)
    const emailCells = page.locator('table tbody tr td').first()
    await expect(emailCells).toBeVisible()

    // Verify status badges render (Pending/Used/Cancelled)
    const statusCells = page.locator('table tbody tr td:nth-child(2)')
    const statusText = await statusCells.first().textContent()
    expect(['Pending', 'Used', 'Cancelled']).toContain(statusText?.trim())
  })

  test('Status filter dropdown renders with options', async ({ page }) => {
    await page.getByRole('navigation').getByRole('button', { name: 'Invite Customer' }).click()
    await page.waitForURL('**/customers/invite')
    await page.waitForSelector('text=Sent Invitations')

    const statusFilter = page.getByRole('combobox', { name: /status/i })
    await expect(statusFilter).toBeVisible()

    const options = await statusFilter.evaluate((select) => {
      const sel = select as HTMLSelectElement
      return Array.from(sel.options).map(o => o.text)
    })
    expect(options).toContain('All')
    expect(options).toContain('Pending')
    expect(options).toContain('Used')
    expect(options).toContain('Cancelled')
  })

  test('Expired/Unexpired filter renders with options', async ({ page }) => {
    await page.getByRole('navigation').getByRole('button', { name: 'Invite Customer' }).click()
    await page.waitForURL('**/customers/invite')
    await page.waitForSelector('text=Sent Invitations')

    const expiredFilter = page.getByRole('combobox', { name: /expired/i })
    await expect(expiredFilter).toBeVisible()

    const options = await expiredFilter.evaluate((select) => {
      const sel = select as HTMLSelectElement
      return Array.from(sel.options).map(o => o.text)
    })
    expect(options).toContain('All')
    expect(options).toContain('Expired')
    expect(options).toContain('Unexpired')
  })

  test('date range inputs render', async ({ page }) => {
    await page.getByRole('navigation').getByRole('button', { name: 'Invite Customer' }).click()
    await page.waitForURL('**/customers/invite')
    await page.waitForSelector('text=Sent Invitations')

    // Date range picker per plan
    await expect(page.getByText('From:').first()).toBeVisible()
    await expect(page.getByText('To:').first()).toBeVisible()
  })

  test('Apply and Reset filter buttons render', async ({ page }) => {
    await page.getByRole('navigation').getByRole('button', { name: 'Invite Customer' }).click()
    await page.waitForURL('**/customers/invite')
    await page.waitForSelector('text=Sent Invitations')

    await expect(page.getByRole('button', { name: 'Apply' })).toBeVisible()
    await expect(page.getByRole('button', { name: 'Reset' })).toBeVisible()
  })

  test('pagination controls render', async ({ page }) => {
    await page.getByRole('navigation').getByRole('button', { name: 'Invite Customer' }).click()
    await page.waitForURL('**/customers/invite')
    await page.waitForSelector('text=Sent Invitations')

    // Pagination controls — Previous/Next buttons
    const prevBtn = page.getByRole('button', { name: /previous/i })
    const nextBtn = page.getByRole('button', { name: /next/i })
    await expect(prevBtn).toBeVisible()
    await expect(nextBtn).toBeVisible()
  })

  test('"Send New Invite" button opens invite modal', async ({ page }) => {
    await page.getByRole('navigation').getByRole('button', { name: 'Invite Customer' }).click()
    await page.waitForURL('**/customers/invite')
    await page.waitForSelector('text=Sent Invitations')

    const sendBtn = page.getByRole('button', { name: /send new invite/i })
    await expect(sendBtn).toBeVisible()

    // Click — should open modal with email input and send button
    await sendBtn.click()
    await expect(page.getByText('Invite Customer').first()).toBeVisible({ timeout: 5000 })
    await expect(page.getByPlaceholder('customer@company.com')).toBeVisible()
    await expect(page.getByRole('button', { name: /send invitation/i })).toBeVisible()
  })

  test('Copy Link button exists on each invite row', async ({ page }) => {
    await page.getByRole('navigation').getByRole('button', { name: 'Invite Customer' }).click()
    await page.waitForURL('**/customers/invite')
    await page.waitForSelector('text=Sent Invitations')

    // Each row should have a "Copy Link" button
    const copyButtons = page.getByRole('button', { name: /copy link/i })
    const count = await copyButtons.count()
    expect(count).toBeGreaterThanOrEqual(1)
  })

  test('Cancel button visible on pending invite rows', async ({ page }) => {
    await page.getByRole('navigation').getByRole('button', { name: 'Invite Customer' }).click()
    await page.waitForURL('**/customers/invite')
    await page.waitForSelector('text=Sent Invitations')

    // Check cancel buttons exist on currently visible rows
    const cancelButtons = page.getByRole('button', { name: /cancel/i })
    const count = await cancelButtons.count()

    // Verify at least one pending invite row has a Cancel button
    const tableRows = page.locator('table tbody tr')
    const rowCount = await tableRows.count()

    let pendingWithCancelFound = false
    for (let i = 0; i < rowCount; i++) {
      const row = tableRows.nth(i)
      const statusCell = row.locator('td:nth-child(2)')
      const statusText = await statusCell.textContent()
      if (statusText?.trim() === 'Pending') {
        const cancelInRow = row.getByRole('button', { name: /cancel/i })
        const cancelCount = await cancelInRow.count()
        if (cancelCount === 1) {
          pendingWithCancelFound = true
        }
      }
    }
    expect(pendingWithCancelFound).toBe(true)
  })

  test('Cancel button not visible on non-pending invite rows', async ({ page }) => {
    await page.getByRole('navigation').getByRole('button', { name: 'Invite Customer' }).click()
    await page.waitForURL('**/customers/invite')
    await page.waitForSelector('text=Sent Invitations')

    // Scan all visible rows for Used/Cancelled — verify no Cancel button
    const tableRows = page.locator('table tbody tr')
    const rowCount = await tableRows.count()

    let nonPendingFound = false
    for (let i = 0; i < rowCount; i++) {
      const row = tableRows.nth(i)
      const statusCell = row.locator('td:nth-child(2)')
      const statusText = await statusCell.textContent()
      const trimmed = statusText?.trim() || ''
      if (trimmed === 'Used' || trimmed === 'Cancelled') {
        nonPendingFound = true
        const cancelBtn = row.getByRole('button', { name: /cancel/i })
        await expect(cancelBtn).toHaveCount(0)
      }
    }
    // If no non-pending rows, test still validates the logic with pending rows
    if (!nonPendingFound) {
      // At least verify Cancel button only appears in Pending rows
      for (let i = 0; i < rowCount; i++) {
        const row = tableRows.nth(i)
        const statusCell = row.locator('td:nth-child(2)')
        const statusText = await statusCell.textContent()
        const trimmed = statusText?.trim() || ''
        const cancelBtn = row.getByRole('button', { name: /cancel/i })
        if (trimmed === 'Pending') {
          await expect(cancelBtn).toHaveCount(1)
        } else {
          await expect(cancelBtn).toHaveCount(0)
        }
      }
    }
  })

  test('Cancel button opens WmsConfirmDialog modal', async ({ page }) => {
    await page.getByRole('navigation').getByRole('button', { name: 'Invite Customer' }).click()
    await page.waitForURL('**/customers/invite')
    await page.waitForSelector('text=Sent Invitations')

    // Find a row with Cancel button
    const cancelButtons = page.getByRole('button', { name: /cancel/i })
    const count = await cancelButtons.count()
    expect(count).toBeGreaterThanOrEqual(1)

    // Click Cancel
    await cancelButtons.first().click()

    // Verify WmsConfirmDialog modal appears with title, message, and buttons
    await expect(page.getByRole('heading', { name: 'Cancel Invitation' })).toBeVisible({ timeout: 5000 })
    await expect(page.getByText(/cancel invitation for/i)).toBeVisible()
    await expect(page.getByRole('button', { name: 'Cancel Invitation' })).toBeVisible()
    await expect(page.getByRole('button', { name: 'Go Back' })).toBeVisible()

    // Dismiss by clicking Go Back
    await page.getByRole('button', { name: 'Go Back' }).click()

    // Modal should disappear
    await expect(page.getByRole('heading', { name: 'Cancel Invitation' })).not.toBeVisible()
  })

  test('Confirming cancel sends API and shows success toast', async ({ page }) => {
    await page.getByRole('navigation').getByRole('button', { name: 'Invite Customer' }).click()
    await page.waitForURL('**/customers/invite')
    await page.waitForSelector('text=Sent Invitations')

    // Find first Cancel button
    const cancelButtons = page.getByRole('button', { name: /cancel/i })
    const count = await cancelButtons.count()
    expect(count).toBeGreaterThanOrEqual(1)

    // Click Cancel → opens confirm modal
    await cancelButtons.first().click()

    // Verify modal and click Confirm
    await expect(page.getByRole('heading', { name: 'Cancel Invitation' })).toBeVisible({ timeout: 5000 })
    await page.getByRole('button', { name: 'Cancel Invitation' }).click()

    // Wait for success toast after cancellation
    await expect(page.getByText('Invitation cancelled')).toBeVisible({ timeout: 15000 })
  })
})
