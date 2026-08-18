type APIResponse<T> = { code: number; message: string; payload?: T }

export const PASSWORD_REQUIREMENTS =
  "至少 8 位，并包含大写字母、小写字母、数字和特殊字符"

export function isSecurePassword(password: string) {
  return (
    password.length >= 8 &&
    password.length <= 72 &&
    /[A-Z]/.test(password) &&
    /[a-z]/.test(password) &&
    /\d/.test(password) &&
    /[^A-Za-z0-9\s]/.test(password) &&
    !/\s/.test(password)
  )
}

let refreshPromise: Promise<boolean> | null = null

async function refreshSession() {
  if (!refreshPromise) {
    refreshPromise = fetch("/api/v1/auth/refresh", {
      method: "POST",
      credentials: "include",
      cache: "no-store",
    })
      .then((response) => response.ok)
      .catch(() => false)
      .finally(() => {
        refreshPromise = null
      })
  }
  return refreshPromise
}

export async function apiRequest<T>(
  url: string,
  init?: RequestInit
): Promise<T> {
  const send = () => {
    const headers = new Headers(init?.headers)
    return fetch(url, {
      credentials: "include",
      cache: "no-store",
      ...init,
      headers,
    })
  }
  let response = await send()
  if (
    response.status === 401 &&
    url !== "/api/v1/auth/refresh" &&
    (await refreshSession())
  ) {
    response = await send()
  }
  const body = (await response.json().catch(() => ({}))) as APIResponse<T>
  if (!response.ok) {
    const error = new Error(body.message || "请求失败") as Error & {
      status?: number
    }
    error.status = response.status
    throw error
  }
  if (body.payload === undefined) throw new Error("服务返回的数据不完整")
  return body.payload
}

export async function logout() {
  await fetch("/api/v1/auth/logout", { method: "POST", credentials: "include" })
}
