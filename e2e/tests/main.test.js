// Jetbrains IDE doesn't respect env and globals from .eslintrc. This is a workaround.
// eslint-disable-next-line no-redeclare
/* global describe, expect, it, page */

const { hostUrl, timeout } = require('../args');

describe('The app error handling', () => {
  it('should show 404 for resources that are not found', async () => {
    const res = await page.goto(`${hostUrl}/non-exist`, { waitUntil: 'networkidle2' });

    // Test error page.
    await expect(res.status()).toBe(404);
    await expect(page.title()).resolves.toMatch('404');
    await expect(page).toMatchElement('p#status', { text: /404/ });
  }, timeout);

  it('should show 401 for unauthorized access to protected pages', async () => {
    // No cookies so launch a new page as incognito.
    const ctx = await browser.createIncognitoBrowserContext();
    const incogPage = await ctx.newPage();
    const res = await incogPage.goto(`${hostUrl}/userinfo`, { waitUntil: 'networkidle2' });

    // Test error page.
    await expect(res.status()).toBe(401);
    await expect(incogPage.title()).resolves.toMatch('401');
    await expect(incogPage).toMatchElement('p#status', { text: /401/ });

    // Clear history.
    await ctx.close();
  }, timeout);
});
