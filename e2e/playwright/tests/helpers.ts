/**
 * Shared test helpers for WMS E2E tests.
 */

export const TEST_USERS = {
  admin: { email: 'ops-e2e@test.com', password: 'TestPass123!' },
  customer: { email: 'customer-e2e@test.com', password: 'TestPass123!' },
};

/**
 * Seed a customer user by calling the register API directly.
 * Run once before customer portal tests.
 */
export async function seedCustomerUser(request: any): Promise<boolean> {
  const res = await request.post('http://localhost:8080/v1/auth/register', {
    data: {
      email: TEST_USERS.customer.email,
      password: TEST_USERS.customer.password,
      company_name: 'E2E Test Corp',
    },
  });
  // 409 = already exists (fine), 200 = created
  return res.ok() || res.status() === 409;
}

/**
 * Seed an admin user by calling the register API directly.
 * Run once before ops portal tests.
 */
export async function seedAdminUser(request: any): Promise<boolean> {
  const res = await request.post('http://localhost:8080/v1/auth/register', {
    data: {
      email: TEST_USERS.admin.email,
      password: TEST_USERS.admin.password,
      company_name: 'Ops E2E',
    },
  });
  // 409 = already exists (fine), 200 = created
  return res.ok() || res.status() === 409;
}

/**
 * Set up console error monitoring on a page.
 * Collects error messages including their source URL for filtering.
 */
export function setupConsoleMonitor(page: any, errors: string[]): void {
  page.on('console', (msg: any) => {
    if (msg.type() === 'error') {
      const loc = msg.location();
      const url = loc?.url || '';
      errors.push(url.includes('favicon.ico')
        ? `${msg.text()} [${url}]`
        : msg.text());
    }
  });
  page.on('pageerror', (err: Error) => {
    errors.push(err.message);
  });
}

/**
 * Check if there are console errors in the page.
 * Fails the test if any error-level message appears.
 */
export function expectNoConsoleErrors(errors: string[]) {
    const filtered = errors.filter(
      (msg) =>
        !msg.includes('bootstrap-autofill-overlay') && // browser extension noise
        !msg.includes('chrome-extension') &&
        !msg.includes('favicon.ico')
    );
  if (filtered.length > 0) {
    throw new Error(`Console errors found:\n${filtered.join('\n')}`);
  }
}

/**
 * Creates a network response monitor that tracks 5xx responses.
 * Call returned checker after each navigation/action.
 */
export function createNetworkMonitor(page: any): {
  track: (response: any) => void;
  expectNoServerErrors: () => void;
  reset: () => void;
} {
  const serverErrors: string[] = [];

  function track(response: any) {
    const status = response.status();
    if (status >= 500) {
      serverErrors.push(`${response.status()} ${response.url()}`);
    }
  }

  function expectNoServerErrors() {
    if (serverErrors.length > 0) {
      throw new Error(`Server errors (5xx) found:\n${serverErrors.join('\n')}`);
    }
  }

  function reset() {
    serverErrors.length = 0;
  }

  return { track, expectNoServerErrors, reset };
}
