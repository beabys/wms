/**
 * Customer Portal E2E — login, dashboard, submit inbound
 *
 * Prerequisites:
 * - Backend running: go run ./cmd/wms
 * - FE running: pnpm --filter @wms/customer-portal dev
 * - VITE_USE_MOCK=false
 */

import { test, expect, type Page } from '@playwright/test';
import { TEST_USERS, seedCustomerUser, setupConsoleMonitor, expectNoConsoleErrors, createNetworkMonitor } from './helpers';

test.describe('Customer Portal — Auth', () => {
  const consoleErrors: string[] = [];

  test.beforeEach(async ({ page }) => {
    consoleErrors.length = 0;
    setupConsoleMonitor(page, consoleErrors);
  });

  test('should load login page without console errors', async ({ page }) => {
    const net = createNetworkMonitor(page);
    page.on('response', net.track);

    await page.goto('/login');
    await expect(page.locator('h1, button[type="submit"]')).toBeVisible({ timeout: 10000 });
    net.expectNoServerErrors();
    expectNoConsoleErrors(consoleErrors);
  });

  test('should login as customer and redirect to dashboard', async ({ page }) => {
    const net = createNetworkMonitor(page);
    page.on('response', net.track);
    await seedCustomerUser(page.request);

    await page.goto('/login');
    await page.waitForSelector('input[type="email"], input[name="email"], input[placeholder*="email"]', { timeout: 10000 });

    // Fill login form
    const emailInput = page.locator('input[type="email"], input[name="email"]').first();
    const passwordInput = page.locator('input[type="password"]').first();
    const submitButton = page.locator('button[type="submit"]').first();

    await emailInput.fill(TEST_USERS.customer.email);
    await passwordInput.fill(TEST_USERS.customer.password);
    await submitButton.click();

    // Should redirect to dashboard
    await expect(page).toHaveURL(/\/dashboard/, { timeout: 15000 });
    net.expectNoServerErrors();
    expectNoConsoleErrors(consoleErrors);
  });
});

async function loginAsCustomer(page: Page, net: { track: (r: any) => void }) {
  await seedCustomerUser(page.request);
  await page.goto('/login');
  await page.waitForSelector('input[type="email"], input[name="email"]', { timeout: 10000 });
  await page.locator('input[type="email"], input[name="email"]').first().fill(TEST_USERS.customer.email);
  await page.locator('input[type="password"]').first().fill(TEST_USERS.customer.password);
  await page.locator('button[type="submit"]').first().click();
  await page.waitForURL(/\/dashboard/, { timeout: 15000 });
}

test.describe('Customer Portal — Dashboard', () => {
  const consoleErrors: string[] = [];

  test.beforeEach(async ({ page }) => {
    consoleErrors.length = 0;
    setupConsoleMonitor(page, consoleErrors);
  });

  test('should show dashboard stats without console errors', async ({ page }) => {
    const net = createNetworkMonitor(page);
    page.on('response', net.track);
    await loginAsCustomer(page, net);

    // Dashboard should show stat cards with numbers (not dashes)
    await expect(page.locator('text=Active Inbounds')).toBeVisible({ timeout: 10000 });
    await expect(page.locator('text=Total Inbounds')).toBeVisible();
    await expect(page.locator('text=Stock Entries')).toBeVisible();
    net.expectNoServerErrors();
    expectNoConsoleErrors(consoleErrors);
  });

  test('should navigate to New Inbound form', async ({ page }) => {
    const net = createNetworkMonitor(page);
    page.on('response', net.track);
    await loginAsCustomer(page, net);

    // Navigate directly to inbound form
    await page.goto('/inbound/new');
    await page.waitForURL(/\/inbound\/new/, { timeout: 10000 });
    await expect(page.locator('text=New Inbound Shipment')).toBeVisible({ timeout: 5000 });
    net.expectNoServerErrors();
    expectNoConsoleErrors(consoleErrors);
  });
});

test.describe('Customer Portal — Submit Inbound', () => {
  const consoleErrors: string[] = [];

  test.beforeEach(async ({ page }) => {
    consoleErrors.length = 0;
    setupConsoleMonitor(page, consoleErrors);
  });

  test('should submit an inbound shipment successfully', async ({ page }) => {
    const net = createNetworkMonitor(page);
    page.on('response', net.track);
    await loginAsCustomer(page, net);

    // Navigate to inbound form
    await page.goto('/inbound/new');
    await page.waitForSelector('text=New Inbound Shipment', { timeout: 10000 });

    // Add a product row — wait for the row to render
    await page.locator('button:has-text("Add Product")').click();
    await page.locator('input[placeholder*="SKU"]').first().waitFor({ state: 'visible', timeout: 5000 });

    // Fill SKU
    const skuInput = page.locator('input[placeholder*="SKU"]').first();
    await skuInput.fill('SKU-E2E-001');

    // Fill quantity
    const qtyInput = page.locator('input[type="number"]').first();
    await qtyInput.fill('10');

    // Submit
    await page.locator('button:has-text("Submit Shipment")').click();

    // Should see success toast notification
    await expect(page.locator('text=submitted successfully')).toBeVisible({ timeout: 15000 });
    net.expectNoServerErrors();
    expectNoConsoleErrors(consoleErrors);
  });
});
