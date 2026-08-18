import {
  useCallback,
  useEffect,
  useMemo,
  useState,
  type FormEvent,
} from "react"
import {
  Braces,
  CirclePlus,
  FilePenLine,
  FolderTree,
  KeyRound,
  LayoutDashboard,
  LoaderCircle,
  Menu as MenuIcon,
  RefreshCw,
  Route,
  Search,
  ShieldCheck,
  Trash2,
  UserRoundCog,
  Users,
  X,
} from "lucide-react"
import { useNavigate } from "react-router"
import { apiRequest } from "@referral/api"
import { AdminNavigation } from "../admin-navigation"
import { AdminGlobalHeader } from "../components/admin-global-header"

type Role = {
  id: number
  name: string
  code: string
  description: string
  enabled: boolean
  is_system: boolean
  permission_ids: number[]
  menu_ids: number[]
}
type Permission = {
  id: number
  name: string
  code: string
  module: string
  description: string
  enabled: boolean
  group_id?: number
  api_ids: number[]
  menu_ids: number[]
  role_count: number
}
type Group = {
  id: number
  name: string
  module: string
  description: string
  parent_id?: number
  sort_order: number
  enabled: boolean
}
type Menu = {
  id: number
  name: string
  path: string
  icon: string
  component: string
  redirect: string
  type: string
  parent_id?: number
  sort_order: number
  enabled: boolean
}
type API = {
  id: number
  name: string
  method: string
  path: string
  description: string
  enabled: boolean
}
type User = {
  id: number
  name: string
  email: string
  enabled: boolean
  role_ids: number[]
}
type Snapshot = {
  roles: Role[]
  permissions: Permission[]
  groups: Group[]
  menus: Menu[]
  apis: API[]
  users: User[]
}
type UserList = { items: User[] }
type Page = "permissions" | "roles" | "menus" | "apis" | "users"
type Editor = {
  kind: "role" | "permission" | "group" | "menu" | "api"
  item?: Role | Permission | Group | Menu | API
} | null
type Notice = { type: "success" | "error"; text: string } | null

const apiPageSize = 10

const empty: Snapshot = {
  roles: [],
  permissions: [],
  groups: [],
  menus: [],
  apis: [],
  users: [],
}

let apiSyncPromise: Promise<unknown> | null = null

function syncAPIsOnce() {
  if (!apiSyncPromise) {
    apiSyncPromise = apiRequest("/api/v1/admin/rbac/apis/sync", {
      method: "POST",
    }).finally(() => {
      apiSyncPromise = null
    })
  }
  return apiSyncPromise
}
export function meta() {
  return [{ title: "权限管理 · Referral Admin" }]
}

export default function PermissionAdmin() {
  const navigate = useNavigate()
  const [data, setData] = useState<Snapshot>(empty)
  const [page, setPage] = useState<Page>("permissions")
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [query, setQuery] = useState("")
  const [status, setStatus] = useState<"all" | "enabled" | "disabled">("all")
  const [selectedGroup, setSelectedGroup] = useState<number | undefined>()
  const [apiPage, setAPIPage] = useState(1)
  const [editor, setEditor] = useState<Editor>(null)
  const [notice, setNotice] = useState<Notice>(null)

  const load = useCallback(async () => {
    setLoading(true)
    try {
      let snapshot =
        await apiRequest<Omit<Snapshot, "users">>("/api/v1/admin/rbac")
      if ((snapshot.apis ?? []).length === 0) {
        await syncAPIsOnce()
        snapshot =
          await apiRequest<Omit<Snapshot, "users">>("/api/v1/admin/rbac")
      }
      const users = await apiRequest<UserList>(
        "/api/v1/admin/users?page_size=100"
      )
      setData({
        ...snapshot,
        roles: (snapshot.roles ?? []).map((role) => ({
          ...role,
          permission_ids: role.permission_ids ?? [],
          menu_ids: role.menu_ids ?? [],
        })),
        permissions: (snapshot.permissions ?? []).map((permission) => ({
          ...permission,
          api_ids: permission.api_ids ?? [],
          menu_ids: permission.menu_ids ?? [],
        })),
        groups: snapshot.groups ?? [],
        menus: snapshot.menus ?? [],
        apis: snapshot.apis ?? [],
        users: (users.items ?? []).map((user) => ({
          ...user,
          role_ids: user.role_ids ?? [],
        })),
      })
    } catch (reason) {
      const error = reason as Error & { status?: number }
      if (error.status === 401) navigate("/login", { replace: true })
      else setNotice({ type: "error", text: error.message })
    } finally {
      setLoading(false)
    }
  }, [navigate])
  useEffect(() => {
    void load()
  }, [load])
  useEffect(() => {
    setQuery("")
    setStatus("all")
    setAPIPage(1)
  }, [page])
  useEffect(() => setAPIPage(1), [query, status])

  async function mutate(url: string, init: RequestInit, message: string) {
    setSaving(true)
    try {
      await apiRequest(url, init)
      setEditor(null)
      setNotice({ type: "success", text: message })
      await load()
    } catch (reason) {
      setNotice({
        type: "error",
        text: reason instanceof Error ? reason.message : "操作失败",
      })
    } finally {
      setSaving(false)
    }
  }
  const json = (method: string, body: unknown): RequestInit => ({
    method,
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  })
  const remove = (kind: string, id: number) => {
    if (!window.confirm("确认删除这条记录？删除后无法恢复。")) return
    void mutate(
      `/api/v1/admin/rbac/${kind}/${id}`,
      { method: "DELETE" },
      "删除成功"
    )
  }
  const normalized = query.trim().toLowerCase()
  const matchesStatus = (enabled: boolean) =>
    status === "all" || (status === "enabled" ? enabled : !enabled)
  const permissions = data.permissions.filter(
    (item) =>
      matchesStatus(item.enabled) &&
      (!selectedGroup || item.group_id === selectedGroup) &&
      (!normalized ||
        `${item.name} ${item.code}`.toLowerCase().includes(normalized))
  )
  const roles = data.roles.filter(
    (item) =>
      matchesStatus(item.enabled) &&
      (!normalized ||
        `${item.name} ${item.code}`.toLowerCase().includes(normalized))
  )
  const menus = data.menus.filter(
    (item) =>
      matchesStatus(item.enabled) &&
      (!normalized ||
        `${item.name} ${item.path}`.toLowerCase().includes(normalized))
  )
  const apis = data.apis.filter(
    (item) =>
      matchesStatus(item.enabled) &&
      (!normalized ||
        `${item.name} ${item.method} ${item.path}`
          .toLowerCase()
          .includes(normalized))
  )
  const apiTotalPages = Math.max(1, Math.ceil(apis.length / apiPageSize))
  const visibleAPIs = apis.slice(
    (apiPage - 1) * apiPageSize,
    apiPage * apiPageSize
  )
  useEffect(() => {
    if (apiPage > apiTotalPages) setAPIPage(apiTotalPages)
  }, [apiPage, apiTotalPages])
  const users = data.users.filter(
    (item) =>
      matchesStatus(item.enabled) &&
      (!normalized ||
        `${item.name} ${item.email}`.toLowerCase().includes(normalized))
  )

  const nav = [
    {
      key: "permissions" as const,
      label: "权限点",
      icon: KeyRound,
      count: data.permissions.length,
    },
    {
      key: "roles" as const,
      label: "角色",
      icon: ShieldCheck,
      count: data.roles.length,
    },
    {
      key: "menus" as const,
      label: "菜单",
      icon: MenuIcon,
      count: data.menus.length,
    },
    {
      key: "apis" as const,
      label: "API",
      icon: Route,
      count: data.apis.length,
    },
    {
      key: "users" as const,
      label: "用户角色",
      icon: UserRoundCog,
      count: data.users.length,
    },
  ]
  const activeNav = nav.find((item) => item.key === page) ?? nav[0]
  return (
    <main className="vben-shell">
      <AdminGlobalHeader />
      <aside className="vben-sidebar">
        <nav>
          <AdminNavigation currentPath="/permissions" />
        </nav>
      </aside>
      <section className="vben-workspace">
        <div className="vben-content">
          <div className="vben-page-title">
            <div>
              <span>ACCESS CONTROL</span>
              <h1>权限管理</h1>
              <p>配置角色、权限点和可访问资源，保持管理边界清晰可审计。</p>
            </div>
            <button onClick={() => void load()} disabled={loading}>
              <RefreshCw />
              刷新数据
            </button>
          </div>
          {notice && (
            <div className={`vben-notice ${notice.type}`}>
              {notice.text}
              <button onClick={() => setNotice(null)}>
                <X />
              </button>
            </div>
          )}
          <nav className="vben-resource-tabs" aria-label="权限资源类型">
            {nav.map(({ key, label, icon: Icon, count }) => (
              <button
                key={key}
                className={page === key ? "active" : ""}
                onClick={() => setPage(key)}
              >
                <Icon />
                <span>{label}</span>
                <em>{count}</em>
              </button>
            ))}
          </nav>
          <section className="vben-filter">
            <label>
              <span>关键词</span>
              <div>
                <Search />
                <input
                  value={query}
                  onChange={(e) => setQuery(e.target.value)}
                  placeholder={`搜索${activeNav.label}名称或编码`}
                />
              </div>
            </label>
            <label>
              <span>状态</span>
              <select
                value={status}
                onChange={(event) =>
                  setStatus(event.target.value as typeof status)
                }
              >
                <option value="all">全部状态</option>
                <option value="enabled">仅启用</option>
                <option value="disabled">仅停用</option>
              </select>
            </label>
            <span className="vben-filter-count">共 {activeNav.count} 项</span>
            <button
              onClick={() => {
                setQuery("")
                setStatus("all")
              }}
            >
              <RefreshCw />
              重置
            </button>
          </section>
          {loading ? (
            <div className="vben-loading">
              <LoaderCircle />
              正在加载权限数据
            </div>
          ) : (
            <>
              {page === "permissions" && (
                <div className="permission-split">
                  <section className="vben-card group-tree">
                    <Toolbar
                      title="权限分组"
                      onCreate={() => setEditor({ kind: "group" })}
                    />
                    <button
                      className={!selectedGroup ? "selected" : ""}
                      onClick={() => setSelectedGroup(undefined)}
                    >
                      <FolderTree />
                      全部权限<span>{data.permissions.length}</span>
                    </button>
                    {data.groups.map((group) => (
                      <div
                        className={
                          selectedGroup === group.id
                            ? "tree-row selected"
                            : "tree-row"
                        }
                        key={group.id}
                      >
                        <button onClick={() => setSelectedGroup(group.id)}>
                          <Braces />
                          <span>
                            <strong>{group.name}</strong>
                            <small>{group.module}</small>
                          </span>
                          <em>
                            {
                              data.permissions.filter(
                                (p) => p.group_id === group.id
                              ).length
                            }
                          </em>
                        </button>
                        <button
                          onClick={() =>
                            setEditor({ kind: "group", item: group })
                          }
                        >
                          <FilePenLine />
                        </button>
                      </div>
                    ))}
                  </section>
                  <section className="vben-card">
                    <Toolbar
                      title="权限点列表"
                      onCreate={() => setEditor({ kind: "permission" })}
                    />
                    <DataTable
                      headers={[
                        "权限名称",
                        "权限编码",
                        "所属模块",
                        "关联角色",
                        "状态",
                        "操作",
                      ]}
                    >
                      {permissions.map((item) => (
                        <tr key={item.id}>
                          <td>
                            <strong>{item.name}</strong>
                            <small>{item.description || "暂无描述"}</small>
                          </td>
                          <td>
                            <code>{item.code}</code>
                          </td>
                          <td>
                            <Tag>{item.module}</Tag>
                          </td>
                          <td>{item.role_count}</td>
                          <td>
                            <Status enabled={item.enabled} />
                          </td>
                          <Actions
                            onEdit={() =>
                              setEditor({ kind: "permission", item })
                            }
                            onDelete={() => remove("permissions", item.id)}
                          />
                        </tr>
                      ))}
                    </DataTable>
                  </section>
                </div>
              )}
              {page === "roles" && (
                <section className="vben-card">
                  <Toolbar
                    title="角色列表"
                    onCreate={() => setEditor({ kind: "role" })}
                  />
                  <DataTable
                    headers={[
                      "角色名称",
                      "角色编码",
                      "权限数量",
                      "菜单数量",
                      "状态",
                      "操作",
                    ]}
                  >
                    {roles.map((item) => (
                      <tr key={item.id}>
                        <td>
                          <strong>{item.name}</strong>
                          <small>{item.description || "暂无描述"}</small>
                        </td>
                        <td>
                          <code>{item.code}</code>
                        </td>
                        <td>{item.permission_ids?.length ?? 0}</td>
                        <td>{item.menu_ids?.length ?? 0}</td>
                        <td>
                          <Status enabled={item.enabled} />
                          {item.is_system && <Tag>系统</Tag>}
                        </td>
                        <Actions
                          onEdit={() => setEditor({ kind: "role", item })}
                          onDelete={
                            item.is_system
                              ? undefined
                              : () => remove("roles", item.id)
                          }
                        />
                      </tr>
                    ))}
                  </DataTable>
                </section>
              )}
              {page === "menus" && (
                <section className="vben-card">
                  <Toolbar
                    title="菜单列表"
                    onCreate={() => setEditor({ kind: "menu" })}
                  />
                  <DataTable
                    headers={[
                      "菜单名称",
                      "类型",
                      "路由地址",
                      "组件",
                      "排序",
                      "状态",
                      "操作",
                    ]}
                  >
                    {menus.map((item) => (
                      <tr key={item.id}>
                        <td
                          style={{
                            paddingLeft: 16 + (item.parent_id ? 22 : 0),
                          }}
                        >
                          <strong>
                            {item.parent_id && "└ "}
                            {item.name}
                          </strong>
                        </td>
                        <td>
                          <Tag>{item.type}</Tag>
                        </td>
                        <td>
                          <code>{item.path || "-"}</code>
                        </td>
                        <td>{item.component || "-"}</td>
                        <td>{item.sort_order}</td>
                        <td>
                          <Status enabled={item.enabled} />
                        </td>
                        <Actions
                          onEdit={() => setEditor({ kind: "menu", item })}
                          onDelete={() => remove("menus", item.id)}
                        />
                      </tr>
                    ))}
                  </DataTable>
                </section>
              )}
              {page === "apis" && (
                <section className="vben-card">
                  <Toolbar
                    title="API 列表"
                    onCreate={() => setEditor({ kind: "api" })}
                    extra={
                      <button
                        onClick={() =>
                          void mutate(
                            "/api/v1/admin/rbac/apis/sync",
                            { method: "POST" },
                            "API 同步完成"
                          )
                        }
                      >
                        <RefreshCw />
                        同步 API
                      </button>
                    }
                  />
                  <DataTable
                    headers={[
                      "API 名称",
                      "请求方式",
                      "访问路径",
                      "说明",
                      "状态",
                      "操作",
                    ]}
                  >
                    {visibleAPIs.map((item) => (
                      <tr key={item.id}>
                        <td>
                          <strong>{item.name}</strong>
                        </td>
                        <td>
                          <Tag>{item.method}</Tag>
                        </td>
                        <td>
                          <code>{item.path}</code>
                        </td>
                        <td>{item.description || "-"}</td>
                        <td>
                          <Status enabled={item.enabled} />
                        </td>
                        <Actions
                          onEdit={() => setEditor({ kind: "api", item })}
                          onDelete={() => remove("apis", item.id)}
                        />
                      </tr>
                    ))}
                  </DataTable>
                  <Pagination
                    page={apiPage}
                    totalPages={apiTotalPages}
                    total={apis.length}
                    onChange={setAPIPage}
                  />
                </section>
              )}
              {page === "users" && (
                <section className="vben-card">
                  <Toolbar title="用户角色分配" />
                  <DataTable
                    headers={["用户", "邮箱", "已分配角色", "状态", "操作"]}
                  >
                    {users.map((user) => (
                      <UserRoleRow
                        key={user.id}
                        user={user}
                        roles={data.roles}
                        save={(ids) =>
                          mutate(
                            `/api/v1/admin/users/${user.id}/roles`,
                            json("PUT", { role_ids: ids }),
                            "角色分配已保存"
                          )
                        }
                      />
                    ))}
                  </DataTable>
                </section>
              )}
            </>
          )}
        </div>
      </section>
      {editor && (
        <EditorDrawer
          editor={editor}
          data={data}
          saving={saving}
          close={() => setEditor(null)}
          save={(kind, id, body) =>
            mutate(
              `/api/v1/admin/rbac/${kind}${id ? `/${id}` : ""}`,
              json(id ? "PUT" : "POST", body),
              id ? "更新成功" : "创建成功"
            )
          }
          saveGrants={(role, permissionIDs, menuIDs) =>
            mutate(
              `/api/v1/admin/rbac/roles/${role.id}/grants`,
              json("PUT", { permission_ids: permissionIDs, menu_ids: menuIDs }),
              "角色授权已保存"
            )
          }
          saveResources={(permission, apiIDs, menuIDs) =>
            mutate(
              `/api/v1/admin/rbac/permissions/${permission.id}/resources`,
              json("PUT", { api_ids: apiIDs, menu_ids: menuIDs }),
              "权限资源已保存"
            )
          }
        />
      )}
    </main>
  )
}

function Toolbar({
  title,
  onCreate,
  extra,
}: {
  title: string
  onCreate?: () => void
  extra?: React.ReactNode
}) {
  return (
    <header className="vben-toolbar">
      <h2>{title}</h2>
      <div>
        {extra}
        {onCreate && (
          <button className="primary" onClick={onCreate}>
            <CirclePlus />
            新增
          </button>
        )}
      </div>
    </header>
  )
}
function DataTable({
  headers,
  children,
}: {
  headers: string[]
  children: React.ReactNode
}) {
  return (
    <div className="vben-table-wrap">
      <table className="vben-table">
        <thead>
          <tr>
            {headers.map((h) => (
              <th key={h}>{h}</th>
            ))}
          </tr>
        </thead>
        <tbody>{children}</tbody>
      </table>
    </div>
  )
}
function Pagination({
  page,
  totalPages,
  total,
  onChange,
}: {
  page: number
  totalPages: number
  total: number
  onChange: (page: number) => void
}) {
  return (
    <footer className="vben-pagination">
      <span>
        共 {total} 项，每页 {apiPageSize} 项
      </span>
      <div>
        <button disabled={page <= 1} onClick={() => onChange(page - 1)}>
          上一页
        </button>
        <strong>
          {page} / {totalPages}
        </strong>
        <button
          disabled={page >= totalPages}
          onClick={() => onChange(page + 1)}
        >
          下一页
        </button>
      </div>
    </footer>
  )
}
function Tag({ children }: { children: React.ReactNode }) {
  return <span className="vben-tag">{children}</span>
}
function Status({ enabled }: { enabled: boolean }) {
  return (
    <span className={`vben-status ${enabled ? "on" : "off"}`}>
      <i />
      {enabled ? "启用" : "停用"}
    </span>
  )
}
function Actions({
  onEdit,
  onDelete,
}: {
  onEdit: () => void
  onDelete?: () => void
}) {
  return (
    <td className="vben-actions">
      <button onClick={onEdit} title="编辑">
        <FilePenLine />
      </button>
      {onDelete && (
        <button className="danger" onClick={onDelete} title="删除">
          <Trash2 />
        </button>
      )}
    </td>
  )
}

function UserRoleRow({
  user,
  roles,
  save,
}: {
  user: User
  roles: Role[]
  save: (ids: number[]) => void
}) {
  const [ids, setIDs] = useState(user.role_ids ?? [])
  return (
    <tr>
      <td>
        <strong>{user.name}</strong>
      </td>
      <td>{user.email}</td>
      <td>
        <div className="role-checks">
          {roles.map((role) => (
            <label key={role.id}>
              <input
                type="checkbox"
                checked={ids.includes(role.id)}
                onChange={() =>
                  setIDs((v) =>
                    v.includes(role.id)
                      ? v.filter((id) => id !== role.id)
                      : [...v, role.id]
                  )
                }
              />
              {role.name}
            </label>
          ))}
        </div>
      </td>
      <td>
        <Status enabled={user.enabled} />
      </td>
      <td>
        <button className="text-button" onClick={() => save(ids)}>
          保存分配
        </button>
      </td>
    </tr>
  )
}

function EditorDrawer({
  editor,
  data,
  saving,
  close,
  save,
  saveGrants,
  saveResources,
}: {
  editor: NonNullable<Editor>
  data: Snapshot
  saving: boolean
  close: () => void
  save: (kind: string, id: number | undefined, body: unknown) => void
  saveGrants: (role: Role, p: number[], m: number[]) => void
  saveResources: (permission: Permission, a: number[], m: number[]) => void
}) {
  const item = editor.item as any
  const [permissionIDs, setPermissionIDs] = useState<number[]>(
    item?.permission_ids || []
  )
  const [apiIDs, setAPIIDs] = useState<number[]>(item?.api_ids || [])
  const [menuIDs, setMenuIDs] = useState<number[]>(item?.menu_ids || [])
  function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const f = new FormData(event.currentTarget)
    const enabled = f.get("enabled") === "on"
    let body: any = {
      name: f.get("name"),
      description: f.get("description") || "",
      enabled,
    }
    if (editor.kind === "role") body = { ...body, code: f.get("code") }
    if (editor.kind === "permission")
      body = {
        ...body,
        code: f.get("code"),
        module: f.get("module"),
        group_id: Number(f.get("group_id")) || undefined,
      }
    if (editor.kind === "group")
      body = {
        ...body,
        module: f.get("module"),
        parent_id: Number(f.get("parent_id")) || undefined,
        sort_order: Number(f.get("sort_order")),
      }
    if (editor.kind === "menu")
      body = {
        ...body,
        path: f.get("path"),
        component: f.get("component"),
        redirect: f.get("redirect"),
        icon: f.get("icon") || "",
        type: f.get("type"),
        parent_id: Number(f.get("parent_id")) || undefined,
        sort_order: Number(f.get("sort_order")),
      }
    if (editor.kind === "api")
      body = { ...body, method: f.get("method"), path: f.get("path") }
    save(
      `${editor.kind === "group" ? "groups" : editor.kind === "api" ? "apis" : `${editor.kind}s`}`,
      item?.id,
      body
    )
  }
  const toggle = (list: number[], set: (v: number[]) => void, id: number) =>
    set(list.includes(id) ? list.filter((v) => v !== id) : [...list, id])
  return (
    <div
      className="vben-drawer-mask"
      onMouseDown={(e) => {
        if (e.target === e.currentTarget) close()
      }}
    >
      <aside className="vben-drawer">
        <header>
          <div>
            <h2>
              {item ? "编辑" : "新增"}
              {
                (
                  {
                    role: "角色",
                    permission: "权限点",
                    group: "权限分组",
                    menu: "菜单",
                    api: "API",
                  } as const
                )[editor.kind]
              }
            </h2>
            <p>填写信息后保存更改</p>
          </div>
          <button onClick={close}>
            <X />
          </button>
        </header>
        <form onSubmit={submit}>
          <label>
            名称
            <input name="name" defaultValue={item?.name} required />
          </label>
          {["role", "permission"].includes(editor.kind) && (
            <label>
              编码
              <input
                name="code"
                defaultValue={item?.code}
                required
                disabled={item?.is_system}
              />
            </label>
          )}
          {["permission", "group"].includes(editor.kind) && (
            <label>
              模块
              <input name="module" defaultValue={item?.module} required />
            </label>
          )}
          {editor.kind === "permission" && (
            <label>
              权限分组
              <select name="group_id" defaultValue={item?.group_id || ""}>
                <option value="">未分组</option>
                {data.groups.map((g) => (
                  <option key={g.id} value={g.id}>
                    {g.name}
                  </option>
                ))}
              </select>
            </label>
          )}
          {editor.kind === "group" && (
            <>
              <label>
                上级分组
                <select name="parent_id" defaultValue={item?.parent_id || ""}>
                  <option value="">顶级分组</option>
                  {data.groups
                    .filter((g) => g.id !== item?.id)
                    .map((g) => (
                      <option key={g.id} value={g.id}>
                        {g.name}
                      </option>
                    ))}
                </select>
              </label>
              <label>
                排序
                <input
                  name="sort_order"
                  type="number"
                  defaultValue={item?.sort_order || 0}
                />
              </label>
            </>
          )}
          {editor.kind === "menu" && (
            <>
              <label>
                菜单类型
                <select name="type" defaultValue={item?.type || "MENU"}>
                  {["CATALOG", "MENU", "BUTTON", "EMBEDDED", "LINK"].map(
                    (v) => (
                      <option key={v}>{v}</option>
                    )
                  )}
                </select>
              </label>
              <label>
                上级菜单
                <select name="parent_id" defaultValue={item?.parent_id || ""}>
                  <option value="">顶级菜单</option>
                  {data.menus
                    .filter((m) => m.id !== item?.id)
                    .map((m) => (
                      <option key={m.id} value={m.id}>
                        {m.name}
                      </option>
                    ))}
                </select>
              </label>
              <label>
                路由地址
                <input name="path" defaultValue={item?.path} />
              </label>
              <label>
                组件路径
                <input name="component" defaultValue={item?.component} />
              </label>
              <label>
                重定向
                <input name="redirect" defaultValue={item?.redirect} />
              </label>
              <label>
                图标
                <input name="icon" defaultValue={item?.icon} />
              </label>
              <label>
                排序
                <input
                  name="sort_order"
                  type="number"
                  defaultValue={item?.sort_order || 0}
                />
              </label>
            </>
          )}
          {editor.kind === "api" && (
            <>
              <label>
                请求方式
                <select name="method" defaultValue={item?.method || "GET"}>
                  {["GET", "POST", "PUT", "PATCH", "DELETE"].map((v) => (
                    <option key={v}>{v}</option>
                  ))}
                </select>
              </label>
              <label>
                访问路径
                <input name="path" defaultValue={item?.path} required />
              </label>
            </>
          )}
          <label>
            说明
            <textarea name="description" defaultValue={item?.description} />
          </label>
          <label className="switch-label">
            <input
              name="enabled"
              type="checkbox"
              defaultChecked={item?.enabled ?? true}
            />
            启用
          </label>
          {editor.kind === "role" && item && !item.is_system && (
            <section className="drawer-grants">
              <h3>权限与菜单授权</h3>
              <div>
                {data.permissions.map((p) => (
                  <label key={p.id}>
                    <input
                      type="checkbox"
                      checked={permissionIDs.includes(p.id)}
                      onChange={() =>
                        toggle(permissionIDs, setPermissionIDs, p.id)
                      }
                    />
                    {p.name}
                    <code>{p.code}</code>
                  </label>
                ))}
              </div>
              <h3>菜单</h3>
              <div>
                {data.menus.map((m) => (
                  <label key={m.id}>
                    <input
                      type="checkbox"
                      checked={menuIDs.includes(m.id)}
                      onChange={() => toggle(menuIDs, setMenuIDs, m.id)}
                    />
                    {m.name}
                  </label>
                ))}
              </div>
              <button
                type="button"
                onClick={() => saveGrants(item, permissionIDs, menuIDs)}
              >
                保存授权
              </button>
            </section>
          )}
          {editor.kind === "permission" && item && (
            <section className="drawer-grants">
              <h3>API 资源</h3>
              <div>
                {data.apis.map((api) => (
                  <label key={api.id}>
                    <input
                      type="checkbox"
                      checked={apiIDs.includes(api.id)}
                      onChange={() => toggle(apiIDs, setAPIIDs, api.id)}
                    />
                    {api.method} {api.path}
                  </label>
                ))}
              </div>
              <h3>可见菜单</h3>
              <div>
                {data.menus.map((menu) => (
                  <label key={menu.id}>
                    <input
                      type="checkbox"
                      checked={menuIDs.includes(menu.id)}
                      onChange={() => toggle(menuIDs, setMenuIDs, menu.id)}
                    />
                    {menu.name}
                  </label>
                ))}
              </div>
              <button
                type="button"
                onClick={() => saveResources(item, apiIDs, menuIDs)}
              >
                保存资源关联
              </button>
            </section>
          )}
          <footer>
            <button type="button" onClick={close}>
              取消
            </button>
            <button className="primary" disabled={saving}>
              {saving ? "保存中…" : "保存"}
            </button>
          </footer>
        </form>
      </aside>
    </div>
  )
}
