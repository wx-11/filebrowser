import { defineConfig, devices } from "@playwright/test";

/**
 * CI-specific Playwright configuration
 * This config is optimized for GitHub Actions and avoids webServer dependency
 */
export default defineConfig({
  testDir: "./tests",
  /* Run tests in files in parallel */
  fullyParallel: false,
  /* Fail the build on CI if you accidentally left test.only in the source code. */
  forbidOnly: !!process.env.CI,
  /* Retry on CI only */
  retries: process.env.CI ? 1 : 0,
  /* Use single worker in CI for stability */
  workers: 1,
  /* Reporter to use in CI */
  reporter: [
    ["list"],
    ["junit", { outputFile: "test-results/junit.xml" }]
  ],
  
  /* Shared settings for CI */
  use: {
    /* No base URL - tests should be self-contained or mock */
    // baseURL: "http://127.0.0.1:5173",
    
    /* Collect trace when retrying the failed test */
    trace: "retain-on-failure",
    
    /* Set default locale to English (US) */
    locale: "en-US",
    
    /* CI-optimized settings */
    actionTimeout: 10000,
    navigationTimeout: 30000,
    
    /* Screenshot settings */
    screenshot: "only-on-failure",
    video: "retain-on-failure",
  },

  /* Configure projects for CI (limited to Chrome for speed) */
  projects: [
    {
      name: "chromium",
      use: { ...devices["Desktop Chrome"] },
    },
  ],

  /* Timeout settings for CI */
  timeout: 30000,
  expect: {
    timeout: 10000,
  },

  /* No webServer in CI - tests should be unit/component tests or mock the backend */
  // webServer: {
  //   command: "npm run dev",
  //   url: "http://127.0.0.1:5173",
  //   reuseExistingServer: !process.env.CI,
  //   timeout: 120000,
  // },
});