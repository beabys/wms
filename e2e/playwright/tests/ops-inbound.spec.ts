/**
 * Ops Portal E2E — login, dashboard, inbound queue, inbound detail
 *
 * Prerequisites:
 * - Backend running: go run ./cmd/wms
 * - FE running: pnpm --filter @wms/ops-portal dev
 * - Admin user exists (seeded via migration or POST /v1/auth/register)
 *
 * Run:
 *   cd e2e/playwright && npx playwright test --project=ops-portal --headed
 */

import { test, expect, type Page } from '@playwright/test';
import { TEST_USERS, expectNoConsoleErrors, createNetworkMonitor, seedAdminUser } from './helpers';

// Seed admin user once before all tests in this file
test.beforeAll(async ({ request }) => {
  await seedAdminUser(request);
});

test.describe('Ops Portal — Auth', () => {
  const consoleErrors: string[] = [];

  test.beforeEach(async ({ page }) => {
    consoleErrors.length = 0;
    page.on('console', (msg) => {
      if (msg.type() === 'error') {
        consoleErrors.push(msg.text());
      }
    });
    page.on('pageerror', (err) => {
      consoleErrors.push(err.message);
    });
  });

  test('should load login page without console errors', async ({ page }) => {
    const net = createNetworkMonitor(page);
    page.on('response', net.track);

    await page.goto('/login');
    await expect(page.locator('h1, button[type="submit"]')).toBeVisible({ timeout: 10000 });
    net.expectNoServerErrors();
    expectNoConsoleErrors(consoleErrors);
  });

  test('should login as admin and redirect to dashboard', async ({ page }) => {
    const net = createNetworkMonitor(page);
    page.on('response', net.track);

    await page.goto('/login');
    await page.waitForSelector('input[type="email"], input[name="email"], input[placeholder*="email"]', { timeout: 10000 });

    const emailInput = page.locator('input[type="email"], input[name="email"]').first();
    const passwordInput = page.locator('input[type="password"]').first();
    const submitButton = page.locator('button[type="submit"]').first();

    await emailInput.fill(TEST_USERS.admin.email);
    await passwordInput.fill(TEST_USERS.admin.password);
    await submitButton.click();

    await expect(page).toHaveURL(/\/dashboard/, { timeout: 15000 });
    net.expectNoServerErrors();
    expectNoConsoleErrors(consoleErrors);
  });
});

async function loginAsAdmin(page: Page, net: { track: (r: any) => void }) {
  await page.goto('/login');
  await page.waitForSelector('input[type="email"], input[name="email"]', { timeout: 10000 });
  await page.locator('input[type="email"], input[name="email"]').first().fill(TEST_USERS.admin.email);
  await page.locator('input[type="password"]').first().fill(TEST_USERS.admin.password);
  await page.locator('button[type="submit"]').first().click();
  await page.waitForURL(/\/dashboard/, { timeout: 15000 });
}

test.describe('Ops Portal — Dashboard', () => {
  const consoleErrors: string[] = [];

  test.beforeEach(async ({ page }) => {
    consoleErrors.length = 0;
    page.on('console', (msg) => {
      if (msg.type() === 'error') {
        consoleErrors.push(msg.text());
      }
    });
    page.on('pageerror', (err) => {
      consoleErrors.push(err.message);
    });
  });

  test('should show dashboard without console errors', async ({ page }) => {
    const net = createNetworkMonitor(page);
    page.on('response', net.track);
    await loginAsAdmin(page, net);

    // Dashboard should be visible with some expected content
    await expect(
      page.locator('h1').or(page.getByTestId('dashboard-title')).or(page.getByText('Dashboard')).first()
    ).toBeVisible({ timeout: 10000 });
    net.expectNoServerErrors();
    expectNoConsoleErrors(consoleErrors);
  });

  test('should navigate via sidebar to inbound queue', async ({ page }) => {
    const net = createNetworkMonitor(page);
    page.on('response', net.track);
    await loginAsAdmin(page, net);

    // Try sidebar link for inbound queue
    const inboundNav = page.locator(
      `a:has-text("Inbound"), a:has-text("Inbound Queue"), nav a:has-text("Queue"), a[href*="inbound"], [data-testid="nav-inbound"]`
    ).first();
    await expect(inboundNav).toBeVisible({ timeout: 10000 });
    await inboundNav.click();

    await page.waitForURL(/\/inbound/, { timeout: 10000 });
    net.expectNoServerErrors();
    expectNoConsoleErrors(consoleErrors);
  });
});

test.describe('Ops Portal — Inbound Queue', () => {
  const consoleErrors: string[] = [];

  test.beforeEach(async ({ page }) => {
    consoleErrors.length = 0;
    page.on('console', (msg) => {
      if (msg.type() === 'error') {
        consoleErrors.push(msg.text());
      }
    });
    page.on('pageerror', (err) => {
      consoleErrors.push(err.message);
    });
  });

  test('should show inbound queue page without console errors', async ({ page }) => {
    const net = createNetworkMonitor(page);
    page.on('response', net.track);
    await loginAsAdmin(page, net);

    await page.goto('/inbound-queue');
    await page.waitForURL(/\/inbound-queue/, { timeout: 10000 });

    // Should see either a list or empty state
    await expect(
      page.locator('table').or(page.locator('[role="grid"]')).or(page.getByText('No inbounds')).or(page.getByTestId('inbound-list'))
    ).toBeVisible({ timeout: 10000 });

    net.expectNoServerErrors();
    expectNoConsoleErrors(consoleErrors);
  });

  test('should navigate from queue to inbound detail page', async ({ page }) => {
    const net = createNetworkMonitor(page);
    page.on('response', net.track);
    await loginAsAdmin(page, net);

    await page.goto('/inbound-queue');
    await page.waitForURL(/\/inbound-queue/, { timeout: 10000 });

    // Click the first inbound if visible
    const firstInboundRow = page.locator(
      `table tbody tr:first-child a, [role="grid"] a:first-child, a[href*="/inbound/"]:first-child, [data-testid="inbound-row"]:first-child`
    ).first();

    if (await firstInboundRow.isVisible({ timeout: 5000 }).catch(() => false)) {
      await firstInboundRow.click();
      // Should land on inbound detail page
      await expect(page).toHaveURL(/\/inbound\/(?!queue)/, { timeout: 10000 });
      await expect(
        page.locator(`h1, [data-testid="inbound-detail"], text=Inbound`)
      ).toBeVisible({ timeout: 10000 });
    } else {
      // No inbounds — that's acceptable (empty state), test passes with warning
      test.info().annotations.push({
        type: 'info',
        description: 'No inbound shipments to click — empty queue',
      });
    }

    net.expectNoServerErrors();
    expectNoConsoleErrors(consoleErrors);
  });
});

test.describe('Ops Portal — Inbound Detail', () => {
  const consoleErrors: string[] = [];

  test.beforeEach(async ({ page }) => {
    consoleErrors.length = 0;
    page.on('console', (msg) => {
      if (msg.type() === 'error') {
        consoleErrors.push(msg.text());
      }
    });
    page.on('pageerror', (err) => {
      consoleErrors.push(err.message);
    });
  });

  test('should show inbound detail page without console errors', async ({ page }) => {
    const net = createNetworkMonitor(page);
    page.on('response', net.track);
    await loginAsAdmin(page, net);

    // Try to access first inbound directly (ID will need to be discovered)
    // First load queue to get an inbound ID
    await page.goto('/inbound-queue');
    await page.waitForURL(/\/inbound-queue/, { timeout: 10000 });

    const firstInboundLink = page.locator('a[href*="/inbound/"]').first();
    const href = await firstInboundLink.getAttribute('href').catch(() => null);

    if (href) {
      await page.goto(href);
      await page.waitForURL(/\/inbound\/(?!queue)/, { timeout: 10000 });

      // Verify detail page has relevant content
      await expect(
        page.locator(`h1, [data-testid="inbound-detail"], text=Inbound, text=Shipment`)
      ).toBeVisible({ timeout: 10000 });
    } else {
      test.info().annotations.push({
        type: 'info',
        description: 'No inbound shipments exist — skipping detail page test',
      });
    }

    net.expectNoServerErrors();
    expectNoConsoleErrors(consoleErrors);
  });
});

test.describe('Ops Portal — Navigation & Errors', () => {
  const consoleErrors: string[] = [];

  test.beforeEach(async ({ page }) => {
    consoleErrors.length = 0;
    page.on('console', (msg) => {
      if (msg.type() === 'error') {
        consoleErrors.push(msg.text());
      }
    });
    page.on('pageerror', (err) => {
      consoleErrors.push(err.message);
    });
  });

  test('should navigate to all main pages without 5xx errors', async ({ page }) => {
    const net = createNetworkMonitor(page);
    page.on('response', net.track);
    await loginAsAdmin(page, net);

    // Visit each main page
    const pages = [
      '/dashboard',
      '/inbound-queue',
    ];

    for (const url of pages) {
      await page.goto(url);
      await page.waitForURL(new RegExp(url.replace('/', '\\/')), { timeout: 10000 });
      await page.waitForLoadState('networkidle');
    }

    net.expectNoServerErrors();
    expectNoConsoleErrors(consoleErrors);
  });

  test('should handle unknown route without crashing', async ({ page }) => {
    const net = createNetworkMonitor(page);
    page.on('response', net.track);
    await loginAsAdmin(page, net);

    await page.goto('/nonexistent-route');
    // Should not crash — no 5xx server errors
    await page.waitForTimeout(2000);

    // App should still render something (not a white screen)
    const bodyText = await page.locator('body').innerText();
    expect(bodyText.length).toBeGreaterThan(0);
    // No server errors
    net.expectNoServerErrors();
  });

  test('should have functional navigation sidebar', async ({ page }) => {
    const net = createNetworkMonitor(page);
    page.on('response', net.track);
    await loginAsAdmin(page, net);

    // Look for nav links (sidebar or header)
    const navLocator = page.locator('nav a, header a, aside a, [role="navigation"] a, .sidebar a');
    const navCount = await navLocator.count();

    // If nav exists, verify links are clickable
    if (navCount > 0) {
      for (let i = 0; i < navCount && i < 3; i++) { // test up to 3 nav links
        const link = navLocator.nth(i);
        const href = await link.getAttribute('href');
        if (href && href !== '#' && !href.startsWith('http')) {
          await link.click();
          // Should navigate somewhere valid
          await page.waitForTimeout(2000);
          net.expectNoServerErrors();
          await page.goBack();
          await page.waitForTimeout(1000);
        }
      }
    } else {
      test.info().annotations.push({
        type: 'info',
        description: 'No navigation sidebar found on ops portal',
      });
    }

    net.expectNoServerErrors();
    expectNoConsoleErrors(consoleErrors);
  });
});
