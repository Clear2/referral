import {
  Coins,
  LayoutDashboard,
  ShieldCheck,
  UserCog,
  Users,
} from "lucide-react"
import { Link, useLocation, useRouteLoaderData } from "react-router"

import type { AdminAccess } from "./root"

const sections = [
  {
    label: "业务管理",
    items: [
      {
        name: "邀请概览",
        to: "/",
        icon: LayoutDashboard,
        permission: "referral:read",
      },
      {
        name: "普通用户",
        to: "/users",
        icon: Users,
        permission: "system:rbac",
      },
      {
        name: "管理账号",
        to: "/administrators",
        icon: UserCog,
        permission: "system:rbac",
      },
      {
        name: "权限管理",
        to: "/permissions",
        icon: ShieldCheck,
        permission: "system:rbac",
      },
    ],
  },
  {
    label: "数据审计",
    items: [
      {
        name: "邀请记录",
        to: "/?view=referrals",
        icon: Users,
        permission: "referral:read",
      },
      {
        name: "Credit 流水",
        to: "/?view=credits",
        icon: Coins,
        permission: "referral:read",
      },
    ],
  },
]

export function AdminNavigation({
  currentPath: _currentPath,
}: {
  currentPath: string
}) {
  const access = useRouteLoaderData("root") as AdminAccess | null
  const location = useLocation()
  const roles = access?.roles ?? []
  const permissions = access?.permissions ?? []
  const unrestricted =
    roles.includes("super_admin") || permissions.includes("system:*")

  return (
    <>
      {sections.map((section) => {
        const items = section.items.filter(
          (item) => unrestricted || permissions.includes(item.permission)
        )
        if (!items.length) return null
        return (
          <div className="admin-nav-section" key={section.label}>
            <span className="admin-nav-section-label">{section.label}</span>
            {items.map(({ name, to, icon: Icon }) => {
              const [path, query = ""] = to.split("?")
              const active =
                location.pathname === path &&
                (query
                  ? new URLSearchParams(location.search).get("view") ===
                    new URLSearchParams(query).get("view")
                  : path !== "/" ||
                    !new URLSearchParams(location.search).has("view"))
              return (
                <Link className={active ? "active" : ""} key={to} to={to}>
                  <Icon />
                  <span>{name}</span>
                </Link>
              )
            })}
          </div>
        )
      })}
    </>
  )
}
