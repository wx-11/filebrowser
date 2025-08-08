# Testing Guide - File Browser Frontend

## Test Types

### Integration Tests (Playwright)
Location: `./tests/`

These are end-to-end tests that require a running File Browser server.

**Local Development:**
```bash
# Terminal 1: Start the backend server
cd .. && go run . --port 3000 --username admin --password admin

# Terminal 2: Run tests
pnpm test
```

**CI Environment:**
Integration tests are currently skipped in CI because they require:
- Running File Browser backend server
- Database setup
- User authentication

### Unit Tests (Future)
Consider adding unit tests with Vitest for:
- Vue components
- Utility functions  
- Store logic
- API client functions

```bash
# Future unit test setup
pnpm add -D vitest @vue/test-utils jsdom
pnpm run test:unit
```

## CI Configuration

The CI workflow uses a smart testing strategy:

1. **Type Checking**: Always runs (`pnpm run typecheck`)
2. **Linting**: Always runs (`pnpm run lint`)
3. **Unit Tests**: Runs if `test:unit` script exists
4. **Integration Tests**: 
   - Skips in CI due to server dependency
   - Runs locally with `playwright.config.ts`
   - Has CI-ready config in `playwright.ci.config.ts`

## Playwright Configurations

- `playwright.config.ts` - Full config with webServer (local development)
- `playwright.ci.config.ts` - CI-optimized config without webServer

## Adding Tests

### For New Features:
1. Add unit tests for components and utilities
2. Add integration tests for user workflows
3. Update this README with test coverage info

### Best Practices:
- Use data-testid attributes for reliable element selection
- Mock external dependencies in unit tests
- Keep integration tests focused on critical user paths
- Use page objects for complex interactions

## Test Environment Setup

### Prerequisites:
- Node.js 22+
- pnpm 9+
- Running File Browser server (for integration tests)

### Commands:
```bash
# Install dependencies
pnpm install

# Install Playwright browsers
pnpm exec playwright install

# Run all tests (requires server)
pnpm test

# Run tests without server (CI mode)  
pnpm run test:ci

# Run with UI mode
pnpm exec playwright test --ui

# Generate test report
pnpm exec playwright show-report
```