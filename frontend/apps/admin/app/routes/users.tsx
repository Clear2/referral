import { useCallback, useEffect, useState, type FormEvent } from "react"
import {
  ArrowLeft,
  Check,
  ChevronLeft,
  ChevronRight,
  CirclePlus,
  Coins,
  Eye,
  KeyRound,
  Link2,
  LoaderCircle,
  Pencil,
  Search,
  ShieldCheck,
  UserCog,
  UserRound,
  Users,
  X,
} from "lucide-react"
import { Link, useNavigate } from "react-router"
import { apiRequest } from "@referral/api"
import { AdminNavigation } from "../admin-navigation"
import { getLocale } from "@referral/i18n"
import { AdminGlobalHeader } from "../components/admin-global-header"

type Role = { id: number; name: string; code: string }
type User = {
  id: number
  name: string
  email: string
  enabled: boolean
  referral_code?: string
  credit_balance: number
  successful_referrals: number
  credit_transactions: number
  role_ids: number[]
  roles: Role[]
  created_at: string
  updated_at: string
}
type Pagination = {
  page: number
  page_size: number
  total: number
  total_pages: number
}
type UserPage = { items: User[]; pagination: Pagination }
type Referral = {
  id: number
  invitee_id: number
  name: string
  email: string
  reward: number
  created_at: string
}
type Credit = {
  id: number
  referral_id: number
  amount: number
  reason: string
  created_at: string
}
type Detail = {
  user: User
  referrals: Referral[]
  credit_transactions: Credit[]
}
type ManagementRole = Role & {
  description: string
  enabled: boolean
  is_system: boolean
  permission_ids: number[]
  menu_ids: number[]
}
type ManagementPermission = {
  id: number
  name: string
  code: string
  module: string
  enabled: boolean
}
type ManagementMenu = {
  id: number
  name: string
  path: string
  enabled: boolean
}
type RBACSnapshot = {
  roles: ManagementRole[]
  permissions: ManagementPermission[]
  menus: ManagementMenu[]
}
type AdminSession = { id: number; name: string; email: string }
type Editor = { kind: "create" } | { kind: "detail"; id: number } | null

export function meta() {
  return [{ title: "用户管理 · Referral Admin" }]
}
const formatTime = (value: string) =>
  new Intl.DateTimeFormat(getLocale(), {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
  }).format(new Date(value))

export default function UsersPage() {
  return <UserAdmin accountType="customer" />
}

export function UserAdmin({
  accountType,
}: {
  accountType: "customer" | "admin"
}) {
  const isAdmin = accountType === "admin"
  const navigate = useNavigate()
  const [data, setData] = useState<UserPage | null>(null)
  const [page, setPage] = useState(1)
  const [queryInput, setQueryInput] = useState("")
  const [query, setQuery] = useState("")
  const [enabled, setEnabled] = useState("")
  const [loading, setLoading] = useState(true)
  const [notice, setNotice] = useState("")
  const [editor, setEditor] = useState<Editor>(null)

  const handleError = useCallback(
    (reason: unknown) => {
      const error = reason as Error & { status?: number }
      if (error.status === 401)
        navigate(`/login?next=${isAdmin ? "/administrators" : "/users"}`, {
          replace: true,
        })
      else
        setNotice(
          error.status === 403
            ? "当前账户没有用户管理权限。"
            : error.message || "操作失败"
        )
    },
    [isAdmin, navigate]
  )

  const load = useCallback(async () => {
    setLoading(true)
    const params = new URLSearchParams({ page: String(page), page_size: "20" })
    params.set("account_type", accountType)
    if (query) params.set("query", query)
    if (enabled) params.set("enabled", enabled)
    try {
      setData(await apiRequest<UserPage>(`/api/v1/admin/users?${params}`))
    } catch (reason) {
      handleError(reason)
    } finally {
      setLoading(false)
    }
  }, [accountType, enabled, handleError, page, query])
  useEffect(() => {
    setPage(1)
    setData(null)
  }, [accountType])
  useEffect(() => {
    void load()
  }, [load])

  function search(event: FormEvent) {
    event.preventDefault()
    setPage(1)
    setQuery(queryInput.trim())
  }
  async function setStatus(user: User) {
    try {
      await apiRequest(`/api/v1/admin/users/${user.id}/status`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ enabled: !user.enabled }),
      })
      setNotice(user.enabled ? "账户已停用" : "账户已启用")
      await load()
    } catch (reason) {
      handleError(reason)
    }
  }

  return (
    <main className="user-admin-page">
      <AdminGlobalHeader />
      <aside className="user-admin-rail">
        <nav>
          <AdminNavigation
            currentPath={isAdmin ? "/administrators" : "/users"}
          />
        </nav>
      </aside>
      <section className="user-admin-main">
        <header className="user-admin-header">
          <div>
            <span>
              {isAdmin ? "ADMINISTRATOR ACCESS" : "CUSTOMER ACCOUNTS"}
            </span>
            <h1>{isAdmin ? "管理账号" : "普通用户"}</h1>
            <p>
              {isAdmin
                ? "集中维护可进入管理端的账号与访问状态。"
                : "维护用户资料、访问状态与邀请资产。"}
            </p>
          </div>
          {isAdmin ? (
            <button onClick={() => navigate("/permissions")}>
              <ShieldCheck />
              配置角色权限
            </button>
          ) : (
            <button onClick={() => setEditor({ kind: "create" })}>
              <CirclePlus />
              创建普通用户
            </button>
          )}
        </header>
        {notice && (
          <div className="user-admin-notice">
            <Check />
            {notice}
            <button onClick={() => setNotice("")}>
              <X />
            </button>
          </div>
        )}
        <section className="user-admin-card">
          <div className="user-admin-toolbar">
            <form onSubmit={search}>
              <label>
                <Search />
                <input
                  value={queryInput}
                  onChange={(event) => setQueryInput(event.target.value)}
                  placeholder="搜索姓名或邮箱"
                />
              </label>
              <select
                value={enabled}
                onChange={(event) => {
                  setEnabled(event.target.value)
                  setPage(1)
                }}
              >
                <option value="">全部状态</option>
                <option value="true">已启用</option>
                <option value="false">已停用</option>
              </select>
              <button type="submit">查询</button>
            </form>
            <span>
              共 {data?.pagination.total ?? 0} 个
              {isAdmin ? "管理账号" : "普通用户"}
            </span>
          </div>
          {loading ? (
            <div className="user-admin-state">
              <LoaderCircle className="spin" />
              正在读取用户
            </div>
          ) : (
            <div className="user-admin-table-wrap">
              <table className="user-admin-table">
                <thead>
                  <tr>
                    <th>用户</th>
                    <th>角色</th>
                    <th>{isAdmin ? "管理角色" : "邀请 / Credit"}</th>
                    <th>状态</th>
                    <th>注册时间</th>
                    <th>操作</th>
                  </tr>
                </thead>
                <tbody>
                  {data?.items.map((user) => (
                    <tr key={user.id}>
                      <td>
                        <div className="user-cell">
                          <i>{user.name.charAt(0).toUpperCase()}</i>
                          <span>
                            <strong>{user.name}</strong>
                            <small>{user.email}</small>
                          </span>
                        </div>
                      </td>
                      <td>
                        <div className="role-tags">
                          {user.roles?.length ? (
                            user.roles.map((role) => (
                              <em key={role.id}>{role.name}</em>
                            ))
                          ) : (
                            <small>{isAdmin ? "未配置角色" : "普通用户"}</small>
                          )}
                        </div>
                      </td>
                      <td>
                        {isAdmin ? (
                          <div className="role-tags">
                            {user.roles?.map((role) => (
                              <em key={role.id}>{role.name}</em>
                            ))}
                          </div>
                        ) : (
                          <>
                            <strong>{user.successful_referrals} 次</strong>
                            <small>{user.credit_balance} Credit</small>
                          </>
                        )}
                      </td>
                      <td>
                        <button
                          className={`status-switch ${user.enabled ? "on" : ""}`}
                          onClick={() => void setStatus(user)}
                        >
                          <i />
                          {user.enabled ? "启用" : "停用"}
                        </button>
                      </td>
                      <td>
                        <time>{formatTime(user.created_at)}</time>
                      </td>
                      <td>
                        <button
                          className="row-action"
                          onClick={() =>
                            setEditor({ kind: "detail", id: user.id })
                          }
                        >
                          <Eye />
                          查看与编辑
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
              {!data?.items.length && (
                <div className="user-admin-state">
                  没有符合条件的{isAdmin ? "管理账号" : "普通用户"}
                </div>
              )}
            </div>
          )}
          {(data?.pagination.total_pages ?? 0) > 1 && (
            <footer className="user-admin-pages">
              <button
                disabled={page <= 1}
                onClick={() => setPage((value) => value - 1)}
              >
                <ChevronLeft />
                上一页
              </button>
              <span>
                {page} / {data?.pagination.total_pages}
              </span>
              <button
                disabled={page >= (data?.pagination.total_pages ?? 1)}
                onClick={() => setPage((value) => value + 1)}
              >
                下一页
                <ChevronRight />
              </button>
            </footer>
          )}
        </section>
      </section>
      {!isAdmin && editor?.kind === "create" && (
        <CreateDrawer
          close={() => setEditor(null)}
          saved={async () => {
            setEditor(null)
            setNotice("普通用户已创建")
            await load()
          }}
          error={handleError}
        />
      )}
      {isAdmin && editor?.kind === "detail" && (
        <AdministratorDrawer
          id={editor.id}
          close={() => setEditor(null)}
          changed={async (message) => {
            setEditor(null)
            setNotice(message)
            await load()
          }}
          error={handleError}
        />
      )}
      {!isAdmin && editor?.kind === "detail" && (
        <DetailDrawer
          id={editor.id}
          close={() => setEditor(null)}
          changed={async (message) => {
            setNotice(message)
            await load()
          }}
          error={handleError}
        />
      )}
    </main>
  )
}

function AdministratorDrawer({
  id,
  close,
  changed,
  error,
}: {
  id: number
  close: () => void
  changed: (message: string) => void
  error: (reason: unknown) => void
}) {
  const [detail, setDetail] = useState<Detail | null>(null)
  const [rbac, setRBAC] = useState<RBACSnapshot | null>(null)
  const [roleIDs, setRoleIDs] = useState<number[]>([])
  const [session, setSession] = useState<AdminSession | null>(null)
  const [saving, setSaving] = useState(false)

  useEffect(() => {
    void Promise.all([
      apiRequest<Detail>(`/api/v1/admin/users/${id}`),
      apiRequest<RBACSnapshot>("/api/v1/admin/rbac"),
      apiRequest<AdminSession>("/api/v1/admin/session"),
    ])
      .then(([userDetail, snapshot, currentSession]) => {
        setDetail(userDetail)
        setRoleIDs(userDetail.user.role_ids ?? [])
        setRBAC({
          roles: (snapshot.roles ?? []).map((role) => ({
            ...role,
            permission_ids: role.permission_ids ?? [],
            menu_ids: role.menu_ids ?? [],
          })),
          permissions: snapshot.permissions ?? [],
          menus: snapshot.menus ?? [],
        })
        setSession(currentSession)
      })
      .catch(error)
  }, [error, id])

  const assignedRoles = (rbac?.roles ?? []).filter((role) =>
    roleIDs.includes(role.id)
  )
  const permissionIDs = new Set(
    assignedRoles.flatMap((role) => role.permission_ids)
  )
  const menuIDs = new Set(assignedRoles.flatMap((role) => role.menu_ids))
  const hasGlobalAccess = assignedRoles.some(
    (role) => role.code === "super_admin"
  )
  const permissions = hasGlobalAccess
    ? (rbac?.permissions ?? [])
    : (rbac?.permissions ?? []).filter((permission) =>
        permissionIDs.has(permission.id)
      )
  const menus = hasGlobalAccess
    ? (rbac?.menus ?? [])
    : (rbac?.menus ?? []).filter((menu) => menuIDs.has(menu.id))
  const isCurrentAccount = session?.id === id

  function toggleRole(roleID: number) {
    setRoleIDs((current) =>
      current.includes(roleID)
        ? current.filter((value) => value !== roleID)
        : [...current, roleID]
    )
  }

  async function saveRoles() {
    setSaving(true)
    try {
      await apiRequest(`/api/v1/admin/users/${id}/roles`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ role_ids: roleIDs }),
      })
      await changed(
        roleIDs.length
          ? "管理账号角色已更新"
          : "管理权限已移除，账号已转入普通用户"
      )
    } catch (reason) {
      error(reason)
    } finally {
      setSaving(false)
    }
  }

  return (
    <Drawer title="管理账号详情" close={close}>
      {!detail || !rbac ? (
        <div className="user-admin-state">
          <LoaderCircle className="spin" />
          正在读取真实权限
        </div>
      ) : (
        <div className="administrator-detail">
          <div className="user-detail-head">
            <i>{detail.user.name.charAt(0).toUpperCase()}</i>
            <div>
              <strong>{detail.user.name}</strong>
              <small>{detail.user.email}</small>
            </div>
            <StatusBadge enabled={detail.user.enabled} />
          </div>

          <section className="administrator-summary">
            <article>
              <small>已分配角色</small>
              <strong>{assignedRoles.length}</strong>
            </article>
            <article>
              <small>有效权限</small>
              <strong>{hasGlobalAccess ? "全部" : permissions.length}</strong>
            </article>
            <article>
              <small>可见菜单</small>
              <strong>{menus.length}</strong>
            </article>
          </section>

          <section className="administrator-section">
            <header>
              <div>
                <h3>角色分配</h3>
                <p>角色决定该账号能够进入哪些管理功能。</p>
              </div>
            </header>
            <div className="administrator-role-grid">
              {rbac.roles.map((role) => (
                <label
                  key={role.id}
                  className={roleIDs.includes(role.id) ? "selected" : ""}
                >
                  <input
                    type="checkbox"
                    checked={roleIDs.includes(role.id)}
                    onChange={() => toggleRole(role.id)}
                    disabled={isCurrentAccount}
                  />
                  <span>
                    <strong>{role.name}</strong>
                    <small>{role.code}</small>
                  </span>
                  {role.is_system && <em>系统角色</em>}
                </label>
              ))}
            </div>
            {isCurrentAccount && (
              <p className="administrator-self-note">
                当前登录账号不能修改自己的角色，避免意外失去管理权限。
              </p>
            )}
            <button
              className="drawer-primary"
              onClick={() => void saveRoles()}
              disabled={saving || isCurrentAccount}
            >
              <ShieldCheck />
              {saving ? "正在保存…" : "保存角色分配"}
            </button>
          </section>

          <section className="administrator-section">
            <header>
              <div>
                <h3>有效权限</h3>
                <p>
                  {hasGlobalAccess
                    ? "超级管理员拥有全部管理权限。"
                    : "根据当前选中的角色实时计算。"}
                </p>
              </div>
              <span>{hasGlobalAccess ? "全部" : permissions.length}</span>
            </header>
            <div className="administrator-access-list">
              {permissions.map((permission) => (
                <span key={permission.id}>
                  <strong>{permission.name}</strong>
                  <code>{permission.code}</code>
                </span>
              ))}
              {!permissions.length && <p>当前角色没有授予管理权限。</p>}
            </div>
          </section>

          <section className="administrator-section">
            <header>
              <div>
                <h3>菜单范围</h3>
                <p>角色关联的后台菜单资源。</p>
              </div>
              <span>{menus.length}</span>
            </header>
            <div className="administrator-access-list">
              {menus.map((menu) => (
                <span key={menu.id}>
                  <strong>{menu.name}</strong>
                  <code>{menu.path || "未配置路径"}</code>
                </span>
              ))}
              {!menus.length && <p>当前角色没有关联菜单。</p>}
            </div>
          </section>
        </div>
      )}
    </Drawer>
  )
}

function StatusBadge({ enabled }: { enabled: boolean }) {
  return (
    <span className={`administrator-status ${enabled ? "on" : ""}`}>
      <i />
      {enabled ? "已启用" : "已停用"}
    </span>
  )
}

function Drawer({
  title,
  close,
  children,
}: {
  title: string
  close: () => void
  children: React.ReactNode
}) {
  return (
    <div
      className="user-drawer-backdrop"
      onMouseDown={(event) => {
        if (event.target === event.currentTarget) close()
      }}
    >
      <aside className="user-drawer">
        <header>
          <h2>{title}</h2>
          <button onClick={close}>
            <X />
          </button>
        </header>
        {children}
      </aside>
    </div>
  )
}

function CreateDrawer({
  close,
  saved,
  error,
}: {
  close: () => void
  saved: () => void
  error: (reason: unknown) => void
}) {
  const [name, setName] = useState("")
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [confirmPassword, setConfirmPassword] = useState("")
  const [enabled, setEnabled] = useState(true)
  const [saving, setSaving] = useState(false)
  async function submit(event: FormEvent) {
    event.preventDefault()
    setSaving(true)
    try {
      await apiRequest("/api/v1/admin/users", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          name,
          email,
          password,
          confirm_password: confirmPassword,
          enabled,
        }),
      })
      await saved()
    } catch (reason) {
      error(reason)
    } finally {
      setSaving(false)
    }
  }
  return (
    <Drawer title="创建用户" close={close}>
      <form className="user-editor-form" onSubmit={submit}>
        <label>
          <span>姓名</span>
          <input
            value={name}
            onChange={(event) => setName(event.target.value)}
            required
          />
        </label>
        <label>
          <span>邮箱</span>
          <input
            type="email"
            value={email}
            onChange={(event) => setEmail(event.target.value)}
            required
          />
        </label>
        <label>
          <span>初始密码</span>
          <input
            type="password"
            minLength={8}
            maxLength={72}
            value={password}
            onChange={(event) => setPassword(event.target.value)}
            required
          />
          <small>至少 8 位，创建后可单独重置。</small>
        </label>
        <label>
          <span>确认初始密码</span>
          <input
            type="password"
            minLength={8}
            maxLength={72}
            value={confirmPassword}
            onChange={(event) => setConfirmPassword(event.target.value)}
            required
          />
          {confirmPassword && password !== confirmPassword && (
            <small>两次输入的密码不一致。</small>
          )}
        </label>
        <label className="check-field">
          <input
            type="checkbox"
            checked={enabled}
            onChange={(event) => setEnabled(event.target.checked)}
          />
          账户立即启用
        </label>
        <button
          className="drawer-primary"
          disabled={saving || password !== confirmPassword}
        >
          {saving ? "正在创建…" : "创建用户"}
        </button>
      </form>
    </Drawer>
  )
}

function DetailDrawer({
  id,
  close,
  changed,
  error,
}: {
  id: number
  close: () => void
  changed: (message: string) => void
  error: (reason: unknown) => void
}) {
  const [detail, setDetail] = useState<Detail | null>(null)
  const [name, setName] = useState("")
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [confirmPassword, setConfirmPassword] = useState("")
  const [tab, setTab] = useState<"profile" | "referrals" | "credits">("profile")
  const [saving, setSaving] = useState(false)
  const load = useCallback(async () => {
    try {
      const result = await apiRequest<Detail>(`/api/v1/admin/users/${id}`)
      setDetail(result)
      setName(result.user.name)
      setEmail(result.user.email)
    } catch (reason) {
      error(reason)
    }
  }, [error, id])
  useEffect(() => {
    void load()
  }, [load])
  async function saveProfile(event: FormEvent) {
    event.preventDefault()
    setSaving(true)
    try {
      await apiRequest(`/api/v1/admin/users/${id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ name, email }),
      })
      await load()
      await changed("用户资料已更新")
    } catch (reason) {
      error(reason)
    } finally {
      setSaving(false)
    }
  }
  async function resetPassword(event: FormEvent) {
    event.preventDefault()
    setSaving(true)
    try {
      await apiRequest(`/api/v1/admin/users/${id}/reset-password`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ password, confirm_password: confirmPassword }),
      })
      setPassword("")
      setConfirmPassword("")
      await changed("密码已重置")
    } catch (reason) {
      error(reason)
    } finally {
      setSaving(false)
    }
  }
  return (
    <Drawer title="用户详情" close={close}>
      {!detail ? (
        <div className="user-admin-state">
          <LoaderCircle className="spin" />
          正在加载
        </div>
      ) : (
        <>
          <div className="user-detail-head">
            <i>{detail.user.name.charAt(0).toUpperCase()}</i>
            <div>
              <strong>{detail.user.name}</strong>
              <small>{detail.user.email}</small>
            </div>
            <span>{detail.user.credit_balance} Credit</span>
          </div>
          <nav className="user-detail-tabs">
            <button
              className={tab === "profile" ? "active" : ""}
              onClick={() => setTab("profile")}
            >
              <UserRound />
              资料
            </button>
            <button
              className={tab === "referrals" ? "active" : ""}
              onClick={() => setTab("referrals")}
            >
              <Link2 />
              邀请 {detail.referrals.length}
            </button>
            <button
              className={tab === "credits" ? "active" : ""}
              onClick={() => setTab("credits")}
            >
              <Coins />
              流水 {detail.credit_transactions.length}
            </button>
          </nav>
          {tab === "profile" && (
            <div>
              <form className="user-editor-form" onSubmit={saveProfile}>
                <label>
                  <span>姓名</span>
                  <input
                    value={name}
                    onChange={(event) => setName(event.target.value)}
                    required
                  />
                </label>
                <label>
                  <span>邮箱</span>
                  <input
                    type="email"
                    value={email}
                    onChange={(event) => setEmail(event.target.value)}
                    required
                  />
                </label>
                <div className="user-info-grid">
                  <span>
                    邀请码
                    <strong>{detail.user.referral_code || "尚未生成"}</strong>
                  </span>
                  <span>
                    角色
                    <strong>
                      {detail.user.roles.map((role) => role.name).join("、") ||
                        "普通用户"}
                    </strong>
                  </span>
                </div>
                <button className="drawer-primary" disabled={saving}>
                  <Pencil />
                  保存资料
                </button>
              </form>
              <form className="reset-password-form" onSubmit={resetPassword}>
                <h3>
                  <KeyRound />
                  重置密码
                </h3>
                <div>
                  <input
                    type="password"
                    minLength={8}
                    maxLength={72}
                    value={password}
                    onChange={(event) => setPassword(event.target.value)}
                    placeholder="输入至少 8 位的新密码"
                    required
                  />
                  <input
                    type="password"
                    minLength={8}
                    maxLength={72}
                    value={confirmPassword}
                    onChange={(event) => setConfirmPassword(event.target.value)}
                    placeholder="再次输入新密码"
                    required
                  />
                  <button disabled={saving || password !== confirmPassword}>
                    重置
                  </button>
                </div>
                {confirmPassword && password !== confirmPassword && (
                  <small>两次输入的密码不一致。</small>
                )}
              </form>
              <Link className="role-manage-link" to="/permissions">
                管理用户角色 <ChevronRight />
              </Link>
            </div>
          )}
          {tab === "referrals" && (
            <div className="user-history">
              {detail.referrals.map((item) => (
                <article key={item.id}>
                  <span>
                    <strong>{item.name}</strong>
                    <small>{item.email}</small>
                  </span>
                  <em>+{item.reward} Credit</em>
                  <time>{formatTime(item.created_at)}</time>
                </article>
              ))}
              {!detail.referrals.length && <p>该用户还没有成功邀请。</p>}
            </div>
          )}
          {tab === "credits" && (
            <div className="user-history">
              {detail.credit_transactions.map((item) => (
                <article key={item.id}>
                  <span>
                    <strong>邀请奖励</strong>
                    <small>Referral #{item.referral_id}</small>
                  </span>
                  <em>+{item.amount} Credit</em>
                  <time>{formatTime(item.created_at)}</time>
                </article>
              ))}
              {!detail.credit_transactions.length && (
                <p>该用户还没有 Credit 流水。</p>
              )}
            </div>
          )}
        </>
      )}
    </Drawer>
  )
}
