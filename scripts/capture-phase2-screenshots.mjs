import { chromium } from '../apps/frontend/node_modules/playwright/index.mjs';
import path from 'node:path';
import fs from 'node:fs';

const root = 'C:/Users/Utkarsh Shukla/OneDrive/Documents/CAPSTONE NAPL/smart-pharmacy';
const screenshots = path.join(root, 'docs/screenshots');
const samplePrescription = path.join(root, 'tests/fixtures/prescriptions/valid_prescription.png');
const baseUrl = 'http://localhost:5173';

const browser = await chromium.launch({ headless: true });
const context = await browser.newContext({ viewport: { width: 1440, height: 1000 }, deviceScaleFactor: 1 });
const page = await context.newPage();

async function shot(name) {
  if (fs.existsSync(path.join(screenshots, name))) {
    console.log('preserving existing', name);
    return;
  }
  console.log('capturing', name);
  await page.screenshot({ path: path.join(screenshots, name), fullPage: true });
}

async function waitForApp() {
  await page.waitForLoadState('domcontentloaded');
  await page.waitForTimeout(700);
}

async function login(role) {
  console.log('logging in', role);
  await page.goto(`${baseUrl}/#/login`, { waitUntil: 'domcontentloaded', timeout: 15000 });
  await waitForApp();
  await page.getByRole('button', { name: new RegExp(`Login as Demo ${role}`, 'i') }).click({ timeout: 10000 });
  await page.waitForTimeout(900);
  console.log('logged in', role, page.url());
}

await page.goto(`${baseUrl}/#/login`, { waitUntil: 'domcontentloaded', timeout: 15000 });
await waitForApp();
await shot('01-login.png');

await page.goto(`${baseUrl}/#/medicines`, { waitUntil: 'domcontentloaded', timeout: 15000 });
await waitForApp();
await shot('02-medicine-catalogue.png');

await login('Customer');
await page.goto(`${baseUrl}/#/prescriptions`, { waitUntil: 'domcontentloaded', timeout: 15000 });
await waitForApp();
await page.locator('#file-input').setInputFiles(samplePrescription);
await page.locator('#preview-area').waitFor({ state: 'visible' });
await page.waitForTimeout(500);
await shot('03-prescription-upload.png');

await page.getByRole('button', { name: 'Submit Prescription for Verification' }).click();
await page.getByText('Analysis complete!', { exact: true }).waitFor({ state: 'visible', timeout: 45000 });
await page.waitForTimeout(500);
await shot('04-prescription-ocr-result.png');

await page.goto(`${baseUrl}/#/orders`, { waitUntil: 'domcontentloaded', timeout: 15000 });
await waitForApp();
await shot('05-customer-orders.png');

await login('Pharmacist');
await page.goto(`${baseUrl}/#/staff`, { waitUntil: 'domcontentloaded', timeout: 15000 });
await waitForApp();
await page.getByRole('button', { name: 'Prescription Queue', exact: true }).waitFor({ state: 'visible', timeout: 10000 });
await page.waitForTimeout(700);
await shot('06-pharmacist-dashboard.png');

const note = page.locator('textarea.rx-note-input').first();
if (await note.count()) {
  await note.fill('Reviewed for Phase-II demonstration; final pharmacist decision pending.');
}
await page.waitForTimeout(300);
await shot('07-prescription-review.png');

await page.getByRole('button', { name: 'Order Fulfillment', exact: true }).click();
await page.waitForTimeout(900);
await shot('08-order-fulfilment.png');

await login('Admin');
await page.goto(`${baseUrl}/#/admin`, { waitUntil: 'domcontentloaded', timeout: 15000 });
await waitForApp();
await shot('09-admin-dashboard.png');

await browser.close();
console.log('Captured Phase-II screenshots to', screenshots);
