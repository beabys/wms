# E2E Tests — WMS

Independent end-to-end tests for all WMS services, BFFs, and frontend.

## Structure

```
e2e/
├── services/          # Go service E2E (via gRPC)
│   ├── inbound/       inbound-service full flow
│   ├── inventory/     inventory-service full flow
│   ├── login/         login-service auth flow
│   ├── customer/      customer-service flow
│   └── notification/  notification-service flow
├── bffs/              # BFF E2E (via HTTP)
│   ├── login-bff/
│   ├── customer-bff/
│   ├── inbound-bff/
│   └── inventory-bff/
├── playwright/        # Frontend E2E (browser automation)
│   ├── tests/
│   ├── playwright.config.ts
│   └── package.json
└── README.md
```

## Running

### Service E2E
```bash
go test ./e2e/services/...
```

### BFF E2E
```bash
go test ./e2e/bffs/...
```

### Playwright
```bash
cd e2e/playwright && npx playwright test
```
