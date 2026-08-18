import { expect, test } from "@playwright/test"

test("invitee registers with only name and email", async ({ page }) => {
  let submitted: Record<string, string> | undefined
  await page.route("**/api/v1/referrals/register", async (route) => {
    submitted = route.request().postDataJSON()
    await route.fulfill({
      contentType: "application/json",
      body: JSON.stringify({
        code: 200,
        message: "success",
        payload: {
          invitee: { id: 2, name: "Bob", email: "bob@example.com", creditBalance: 0 },
          reward: 100,
          inviterCreditBalance: 100,
        },
      }),
    })
  })

  await page.goto("/ref/ABCD2345")
  await expect(page.locator('input[type="password"]')).toHaveCount(0)
  await expect(page.locator(".verification-modal")).toHaveCount(0)
  await page.locator('input[autocomplete="name"]').fill("Bob")
  await page.locator('input[autocomplete="email"]').fill("bob@example.com")
  await page.locator('button[type="submit"]').click()

  await expect(page.getByRole("status")).toContainText("100 Credit")
  expect(submitted).toEqual({ code: "ABCD2345", name: "Bob", email: "bob@example.com" })
})

for (const scenario of [
  { status: 404, message: "邀请码不存在" },
  { status: 409, message: "邮箱已注册" },
]) {
  test(`shows API error: ${scenario.message}`, async ({ page }) => {
    await page.route("**/api/v1/referrals/register", (route) =>
      route.fulfill({
        status: scenario.status,
        contentType: "application/json",
        body: JSON.stringify({ code: scenario.status, message: scenario.message }),
      })
    )
    await page.goto("/ref/ABCD2345")
    await page.locator('input[autocomplete="name"]').fill("Bob")
    await page.locator('input[autocomplete="email"]').fill("bob@example.com")
    await page.locator('button[type="submit"]').click()
    await expect(page.getByRole("alert")).toHaveText(scenario.message)
  })
}
