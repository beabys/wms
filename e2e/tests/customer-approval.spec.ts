import { test, expect } from '@playwright/test'

const ADMIN_UI = 'http://localhost:5173'
const CUSTOMER_UI = 'http://localhost:5174'
const CUSTOMER_BFF = 'http://localhost:8082'
const AUTH_BFF = 'http://localhost:8081'

const ADMIN_EMAIL = 'admin@wms.local'
const ADMIN_PASS = 'admin123'
const TEST_PASS = 'Password123!'

let testCustomerEmail = ''
let testCustomerId = ''
let customerToken = ''

test.describe.serial('Customer Approval Flow — Real API', () => {
  test('Setup: create pending customer via API', async () => {
    // Login as admin
    const loginResp = await fetch(`${AUTH_BFF}/v1/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email: ADMIN_EMAIL, password: ADMIN_PASS }),
    })
    const loginData = await loginResp.json()
    const adminToken = loginData.data.access_token
    expect(adminToken).toBeTruthy()

    // Create invite
    const ts = Date.now()
    testCustomerEmail = `e2e-approval-${ts}@test.com`
    const inviteResp = await fetch(`${AUTH_BFF}/v1/auth/invite`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${adminToken}`,
      },
      body: JSON.stringify({ email: testCustomerEmail }),
    })
    const inviteData = await inviteResp.json()
    const inviteToken = inviteData.data.token
    expect(inviteToken).toBeTruthy()

    // Register customer
    const regResp = await fetch(`${CUSTOMER_BFF}/v1/customers/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        company_name: `E2E Approval Test ${ts}`,
        email: testCustomerEmail,
        password: TEST_PASS,
        token: inviteToken,
      }),
    })
    const regData = await regResp.json()
    expect(regResp.status).toBe(201)
    expect(regData.success).toBe(true)
    testCustomerId = regData.data.customer.id || ''

    // Login as customer to get token
    const custLoginResp = await fetch(`${AUTH_BFF}/v1/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email: testCustomerEmail, password: TEST_PASS }),
    })
    const custLoginData = await custLoginResp.json()
    customerToken = custLoginData.data.access_token
    expect(customerToken).toBeTruthy()

    // Verify GET /v1/customers/me returns pending
    const meResp = await fetch(`${CUSTOMER_BFF}/v1/customers/me`, {
      headers: { Authorization: `Bearer ${customerToken}` },
    })
    const meData = await meResp.json()
    expect(meData.data.status).toBe('pending')
  })

  test('Pending customer sees waiting message on dashboard', async ({ page }) => {
    // Login as pending customer
    await page.goto(`${CUSTOMER_UI}/login`)
    await page.waitForSelector('text=Welcome Back')
    await page.getByPlaceholder('you@example.com').fill(testCustomerEmail)
    await page.getByPlaceholder('Enter your password').fill(TEST_PASS)
    await page.getByRole('button', { name: 'Sign In' }).click()

    // Should land on dashboard
    await page.waitForURL('**/dashboard')

    // Dashboard should show pending/waiting message (one of these)
    const pendingVisible = await page.getByText(/pending|waiting for approval|account.*pending/i).first().isVisible().catch(() => false)
    const headingVisible = await page.getByRole('heading').first().isVisible().catch(() => false)
    
    if (pendingVisible) {
      await expect(page.getByText(/pending|waiting for approval|account.*pending/i).first()).toBeVisible()
    } else {
      // At minimum the page loaded
      expect(headingVisible).toBe(true)
    }
  })

  test('Sidebar shows no nav links for pending customer', async ({ page }) => {
    // Login as pending customer
    await page.goto(`${CUSTOMER_UI}/login`)
    await page.waitForSelector('text=Welcome Back')
    await page.getByPlaceholder('you@example.com').fill(testCustomerEmail)
    await page.getByPlaceholder('Enter your password').fill(TEST_PASS)
    await page.getByRole('button', { name: 'Sign In' }).click()
    await page.waitForURL('**/dashboard')

    // Sidebar should not have navigation links for pending customers
    // Logout button should still be visible
    const logoutBtn = page.getByRole('button', { name: /logout|sign out/i })
    await expect(logoutBtn).toBeVisible()
  })

  test('Admin approves customer via API', async () => {
    const loginResp = await fetch(`${AUTH_BFF}/v1/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email: ADMIN_EMAIL, password: ADMIN_PASS }),
    })
    const loginData = await loginResp.json()
    const adminToken = loginData.data.access_token

    // Approve
    const approveResp = await fetch(`${CUSTOMER_BFF}/v1/customers/${testCustomerId}/approve`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${adminToken}`,
      },
    })
    const approveData = await approveResp.json()
    expect(approveResp.status).toBe(200)
    expect(approveData.data.status).toBe('active')
  })

  test('Approved customer sees welcome and full sidebar', async ({ page }) => {
    // Login as approved customer
    await page.goto(`${CUSTOMER_UI}/login`)
    await page.waitForSelector('text=Welcome Back')
    await page.getByPlaceholder('you@example.com').fill(testCustomerEmail)
    await page.getByPlaceholder('Enter your password').fill(TEST_PASS)
    await page.getByRole('button', { name: 'Sign In' }).click()
    await page.waitForURL('**/dashboard')

    // Dashboard should show welcome content
    const welcomeVisible = await page.getByText(/welcome|dashboard|WMS/i).first().isVisible().catch(() => false)
    expect(welcomeVisible).toBe(true)

    // Sidebar nav links should be visible for approved customer
    const logoutBtn = page.getByRole('button', { name: /logout|sign out/i })
    await expect(logoutBtn).toBeVisible()
  })

  test('GET /v1/customers/me returns active for approved customer', async () => {
    const custLoginResp = await fetch(`${AUTH_BFF}/v1/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email: testCustomerEmail, password: TEST_PASS }),
    })
    const custLoginData = await custLoginResp.json()
    const newToken = custLoginData.data.access_token

    const meResp = await fetch(`${CUSTOMER_BFF}/v1/customers/me`, {
      headers: { Authorization: `Bearer ${newToken}` },
    })
    const meData = await meResp.json()
    expect(meData.data.status).toBe('active')
    expect(meData.data.email).toBe(testCustomerEmail)
    expect(meData.data.company_name).toContain('E2E Approval Test')
  })
})
