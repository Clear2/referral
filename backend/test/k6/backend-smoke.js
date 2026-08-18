import http from "k6/http"
import { check, fail, sleep } from "k6"

const baseURL = (__ENV.BASE_URL || "http://localhost:8999").replace(/\/$/, "")
const adminEmail = __ENV.ADMIN_EMAIL || ""
const adminPassword = __ENV.ADMIN_PASSWORD || ""
const userEmail = __ENV.USER_EMAIL || adminEmail
const userPassword = __ENV.USER_PASSWORD || adminPassword

const vus = Number(__ENV.VUS || 2)
const duration = __ENV.DURATION || "20s"

export const options = {
  noCookiesReset: true,
  scenarios: {
    user_api: {
      executor: "constant-vus",
      exec: "userAPI",
      vus,
      duration,
      gracefulStop: "5s",
      tags: { actor: "user" },
    },
    admin_api: {
      executor: "constant-vus",
      exec: "adminAPI",
      vus: Math.max(1, Math.ceil(vus / 2)),
      duration,
      gracefulStop: "5s",
      tags: { actor: "admin" },
    },
  },
  thresholds: {
    checks: ["rate>0.99"],
    http_req_failed: ["rate<0.01"],
    http_req_duration: ["p(95)<750", "p(99)<1500"],
    "http_req_duration{actor:user}": ["p(95)<750"],
    "http_req_duration{actor:admin}": ["p(95)<750"],
  },
}

let userID = 0
let userAuthenticated = false
let adminAuthenticated = false

function jsonHeaders() {
  return { headers: { "Content-Type": "application/json" } }
}

function payload(response) {
  try {
    return response.json("payload")
  } catch (_) {
    return null
  }
}

function login(email, password, actor) {
  if (!email || !password) {
    fail(`Missing ${actor.toUpperCase()}_EMAIL or ${actor.toUpperCase()}_PASSWORD`)
  }
  const response = http.post(
    `${baseURL}/api/v1/auth/login-with-account`,
    JSON.stringify({ email, password }),
    { ...jsonHeaders(), tags: { name: "POST /auth/login", actor } }
  )
  const body = payload(response)
  const ok = check(response, {
    [`${actor} login returns 200`]: (item) => item.status === 200,
    [`${actor} login returns user id`]: () => Number(body?.user_id) > 0,
  })
  if (!ok) fail(`${actor} login failed with HTTP ${response.status}`)
  return Number(body.user_id)
}

function get(path, name, actor, expected = 200) {
  const response = http.get(`${baseURL}${path}`, {
    tags: { name, actor },
  })
  check(response, {
    [`${name} returns ${expected}`]: (item) => item.status === expected,
    [`${name} returns API envelope`]: (item) => {
      try {
        return item.json("code") === expected
      } catch (_) {
        return false
      }
    },
  })
  return response
}

export function userAPI() {
  if (!userAuthenticated) {
    userID = login(userEmail, userPassword, "user")
    userAuthenticated = true
  }

  const me = get("/api/v1/users/me", "GET /users/me", "user")
  check(me, {
    "current user matches login": () => Number(payload(me)?.id) === userID,
  })

  const code = http.post(
    `${baseURL}/api/v1/users/${userID}/referral-code`,
    null,
    { tags: { name: "POST /users/:id/referral-code", actor: "user" } }
  )
  check(code, {
    "referral code returns 200": (item) => item.status === 200,
    "referral code is present": () => Boolean(payload(code)?.code),
  })

  const dashboard = get(
    `/api/v1/users/${userID}/referral-dashboard`,
    "GET /users/:id/referral-dashboard",
    "user"
  )
  check(dashboard, {
    "dashboard has numeric totals": () => {
      const body = payload(dashboard)
      return (
        Number.isFinite(body?.successfulReferrals) &&
        Number.isFinite(body?.totalCreditsEarned)
      )
    },
  })
  sleep(1)
}

export function adminAPI() {
  if (!adminAuthenticated) {
    login(adminEmail, adminPassword, "admin")
    adminAuthenticated = true
  }

  const session = get("/api/v1/admin/session", "GET /admin/session", "admin")
  check(session, {
    "admin session identifies the account": () => {
      const body = payload(session)
      return Number(body?.id ?? body?.user_id) > 0
    },
  })

  const access = get("/api/v1/access/me", "GET /access/me", "admin")
  check(access, {
    "RBAC access arrays are never null": () => {
      const body = payload(access)
      return (
        Array.isArray(body?.roles) &&
        Array.isArray(body?.permissions) &&
        Array.isArray(body?.menu_ids) &&
        Array.isArray(body?.menus)
      )
    },
  })

  get("/api/v1/admin/referral-stats", "GET /admin/referral-stats", "admin")
  get("/api/v1/admin/referrals?page=1&page_size=20", "GET /admin/referrals", "admin")
  get(
    "/api/v1/admin/credit-transactions?page=1&page_size=20",
    "GET /admin/credit-transactions",
    "admin"
  )
  get(
    "/api/v1/admin/users?page=1&page_size=20&account_type=customer",
    "GET /admin/users customers",
    "admin"
  )
  get(
    "/api/v1/admin/users?page=1&page_size=20&account_type=admin",
    "GET /admin/users administrators",
    "admin"
  )
  get("/api/v1/admin/rbac", "GET /admin/rbac", "admin")
  sleep(1)
}

export function setup() {
  if (!adminEmail || !adminPassword) {
    fail("Set ADMIN_EMAIL and ADMIN_PASSWORD before running this test")
  }
  const health = http.get(`${baseURL}/api/v1/healthz`, {
    tags: { name: "GET /healthz", actor: "health" },
  })
  check(health, {
    "backend health check returns 200": (response) => response.status === 200,
  })
}
