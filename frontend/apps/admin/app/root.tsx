import {
  Links,
  Meta,
  Outlet,
  Scripts,
  ScrollRestoration,
  isRouteErrorResponse,
  redirect,
} from "react-router"
import type { Route } from "./+types/root"
import { I18nProvider } from "@referral/i18n"
import { apiRequest } from "@referral/api"
import "@referral/ui/styles.css"

export type AdminMenu = {
  id: number
  name: string
  path: string
  icon: string
  type: string
  parent_id?: number
  sort_order: number
  enabled: boolean
}
export type AdminAccess = {
  user_id: number
  roles: string[]
  permissions: string[]
  menu_ids: number[]
  menus: AdminMenu[]
}

export async function clientLoader({ request }: Route.ClientLoaderArgs) {
  const url = new URL(request.url)
  if (url.pathname === "/admin/login" || url.pathname === "/login") return null
  try {
    await apiRequest("/api/v1/admin/session")
    const response = await apiRequest<AdminAccess>("/api/v1/access/me")
    const access: AdminAccess = {
      ...response,
      roles: response.roles ?? [],
      permissions: response.permissions ?? [],
      menu_ids: response.menu_ids ?? [],
      menus: response.menus ?? [],
    }
    const path = url.pathname.replace(/^\/admin/, "") || "/"
    const isSuperAdmin = access.roles.includes("super_admin")
    const canOpen =
      isSuperAdmin ||
      access.menus.some(
        (menu) =>
          menu.enabled &&
          ["MENU", "CATALOG"].includes(menu.type) &&
          (menu.path === path ||
            (menu.path !== "/" && path.startsWith(`${menu.path}/`)))
      )
    if (path !== "/" && !isSuperAdmin && !canOpen)
      throw new Response("Forbidden", { status: 403 })
    return access
  } catch (error) {
    if (error instanceof Response && error.status === 403) throw error
    if ((error as Error & { status?: number }).status === 403)
      throw new Response("Forbidden", { status: 403 })
    const next = `${url.pathname.replace(/^\/admin/, "") || "/"}${url.search}`
    throw redirect(
      next === "/" ? "/login" : `/login?next=${encodeURIComponent(next)}`
    )
  }
}

clientLoader.hydrate = true as const

export function HydrateFallback() {
  return <main className="admin-auth-check">正在验证管理权限…</main>
}

export function Layout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="zh-CN">
      <head>
        <meta charSet="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <meta name="robots" content="noindex, nofollow, noarchive" />
        <meta name="theme-color" content="#111827" />
        <Meta />
        <Links />
      </head>
      <body>
        {children}
        <ScrollRestoration />
        <Scripts />
      </body>
    </html>
  )
}

export default function App() {
  return (
    <I18nProvider defaultSwitcher={false}>
      <Outlet />
    </I18nProvider>
  )
}

export function ErrorBoundary({ error }: Route.ErrorBoundaryProps) {
  const notFound = isRouteErrorResponse(error) && error.status === 404
  const forbidden = isRouteErrorResponse(error) && error.status === 403
  return (
    <main className="container mx-auto p-8">
      <h1>{notFound ? "404" : forbidden ? "无权访问" : "页面出错"}</h1>
      <p>
        {notFound
          ? "找不到这个管理页面。"
          : forbidden
            ? "当前角色没有访问此管理页面的菜单权限。"
            : "无法加载管理中心，请稍后重试。"}
      </p>
    </main>
  )
}
