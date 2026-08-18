import { useCallback, useEffect, useState, type FormEvent } from "react"
import {
  ArrowLeft,
  ArrowRight,
  Gift,
  LayoutDashboard,
  Link2,
  LoaderCircle,
  RefreshCw,
  Search,
  Users,
} from "lucide-react"
import { useNavigate, useSearchParams } from "react-router"

import { apiRequest } from "@referral/api"
import { getLocale } from "@referral/i18n"
import { AdminGlobalHeader } from "../components/admin-global-header"
import { AdminNavigation } from "../admin-navigation"

type User = { id: number; name: string; email: string; creditBalance: number }
type Pagination = {
  page: number
  pageSize: number
  total: number
  totalPages: number
}
type Referral = {
  id: number
  inviter: User
  invitee: User
  reward: number
  createdAt: string
}
type CreditTransaction = {
  id: number
  user: User
  amount: number
  reason: string
  referralId: number
  createdAt: string
}
type Page<T> = { items: T[]; pagination: Pagination }
type Stats = {
  totalUsers: number
  totalInviters: number
  totalReferrals: number
  totalCreditsAwarded: number
}
type View = "referrals" | "credits"

export function meta() {
  return [{ title: "Referral 管理中心" }]
}

const formatTime = (value: string) =>
  new Intl.DateTimeFormat(getLocale(), {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
  }).format(new Date(value))

export default function ReferralAdmin() {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const view: View =
    searchParams.get("view") === "credits" ? "credits" : "referrals"
  const [stats, setStats] = useState<Stats | null>(null)
  const [referrals, setReferrals] = useState<Page<Referral> | null>(null)
  const [credits, setCredits] = useState<Page<CreditTransaction> | null>(null)
  const [email, setEmail] = useState("")
  const [query, setQuery] = useState("")
  const [from, setFrom] = useState("")
  const [to, setTo] = useState("")
  const [page, setPage] = useState(1)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState("")

  const handleError = useCallback(
    (reason: unknown) => {
      const typed = reason as Error & { status?: number }
      if (typed.status === 401) {
        navigate("/login", { replace: true })
        return
      }
      setError(
        typed.status === 403
          ? "当前账户没有查看 Referral 管理数据的权限。"
          : typed.message || "数据加载失败"
      )
    },
    [navigate]
  )

  const loadStats = useCallback(async () => {
    try {
      setStats(await apiRequest<Stats>("/api/v1/admin/referral-stats"))
    } catch (reason) {
      handleError(reason)
    }
  }, [handleError])

  const loadRecords = useCallback(async () => {
    setLoading(true)
    setError("")
    const params = new URLSearchParams({ page: String(page), page_size: "20" })
    if (query) params.set("email", query)
    if (from) params.set("created_at_from", from)
    if (to) params.set("created_at_to", to)
    try {
      if (view === "referrals")
        setReferrals(
          await apiRequest<Page<Referral>>(`/api/v1/admin/referrals?${params}`)
        )
      else
        setCredits(
          await apiRequest<Page<CreditTransaction>>(
            `/api/v1/admin/credit-transactions?${params}`
          )
        )
    } catch (reason) {
      handleError(reason)
    } finally {
      setLoading(false)
    }
  }, [from, handleError, page, query, to, view])

  useEffect(() => {
    void loadStats()
  }, [loadStats])
  useEffect(() => {
    void loadRecords()
  }, [loadRecords])

  function search(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setPage(1)
    setQuery(email.trim())
  }

  const currentPagination =
    view === "referrals" ? referrals?.pagination : credits?.pagination

  return (
    <main className="ref-admin">
      <AdminGlobalHeader />
      <aside className="ref-admin-rail">
        <nav>
          <AdminNavigation currentPath="/" />
        </nav>
      </aside>

      <section className="ref-admin-main">
        <header className="ref-admin-header">
          <div>
            <span>REFERRAL OPERATIONS</span>
            <h1>邀请奖励概览</h1>
            <p>追踪每一段邀请关系和每一笔 Credit 奖励。</p>
          </div>
          <button
            onClick={() => {
              void loadStats()
              void loadRecords()
            }}
          >
            <RefreshCw />
            刷新数据
          </button>
        </header>

        {error && (
          <div className="ref-admin-error" role="alert">
            {error}
          </div>
        )}

        <section className="ref-admin-stats" aria-label="全局统计">
          <article>
            <span>
              <Users />
            </span>
            <div>
              <small>用户总数</small>
              <strong>{stats?.totalUsers ?? "--"}</strong>
            </div>
          </article>
          <article>
            <span>
              <Link2 />
            </span>
            <div>
              <small>发起邀请用户</small>
              <strong>{stats?.totalInviters ?? "--"}</strong>
            </div>
          </article>
          <article>
            <span>
              <LayoutDashboard />
            </span>
            <div>
              <small>成功邀请</small>
              <strong>{stats?.totalReferrals ?? "--"}</strong>
            </div>
          </article>
          <article className="credit-stat">
            <span>
              <Gift />
            </span>
            <div>
              <small>已发放 Credit</small>
              <strong>{stats?.totalCreditsAwarded ?? "--"}</strong>
            </div>
          </article>
        </section>

        <section className="ref-admin-trail" aria-label="奖励规则">
          <div>
            <i>A</i>
            <span>邀请人</span>
          </div>
          <b />
          <div>
            <i>B</i>
            <span>新用户注册</span>
          </div>
          <b />
          <div className="reward">
            <i>+100</i>
            <span>Credit 入账</span>
          </div>
        </section>

        <section className="ref-admin-panel">
          <div className="ref-admin-panel-head">
            <div>
              <h2>{view === "referrals" ? "邀请记录" : "Credit 变更流水"}</h2>
              <p>共 {currentPagination?.total ?? 0} 条记录</p>
            </div>
            <form onSubmit={search} className="ref-admin-filters">
              <label>
                <Search />
                <input
                  value={email}
                  onChange={(event) => setEmail(event.target.value)}
                  placeholder="按邮箱搜索"
                />
              </label>
              <input
                type="date"
                value={from}
                onChange={(event) => {
                  setFrom(event.target.value)
                  setPage(1)
                }}
                aria-label="开始日期"
              />
              <span>至</span>
              <input
                type="date"
                value={to}
                onChange={(event) => {
                  setTo(event.target.value)
                  setPage(1)
                }}
                aria-label="结束日期"
              />
              <button type="submit">查询</button>
            </form>
          </div>

          {loading ? (
            <div className="ref-admin-state">
              <LoaderCircle className="spin" />
              正在读取数据
            </div>
          ) : view === "referrals" ? (
            <div className="ref-admin-table referrals-table">
              <div className="table-head">
                <span>邀请人</span>
                <span>关系</span>
                <span>被邀请人</span>
                <span>奖励</span>
                <span>创建时间</span>
              </div>
              {referrals?.items.map((item) => (
                <div className="ref-admin-data-row" key={item.id}>
                  <span>
                    <strong>{item.inviter.name}</strong>
                    <small>{item.inviter.email}</small>
                  </span>
                  <span className="relation">
                    <i />→<i />
                  </span>
                  <span>
                    <strong>{item.invitee.name}</strong>
                    <small>{item.invitee.email}</small>
                  </span>
                  <span className="credit-value">+{item.reward} Credit</span>
                  <time>{formatTime(item.createdAt)}</time>
                </div>
              ))}
              {!referrals?.items.length && (
                <div className="ref-admin-state">没有符合条件的邀请记录</div>
              )}
            </div>
          ) : (
            <div className="ref-admin-table credits-table">
              <div className="table-head">
                <span>获得用户</span>
                <span>变更原因</span>
                <span>关联邀请</span>
                <span>金额</span>
                <span>入账时间</span>
              </div>
              {credits?.items.map((item) => (
                <div className="ref-admin-data-row" key={item.id}>
                  <span>
                    <strong>{item.user.name}</strong>
                    <small>{item.user.email}</small>
                  </span>
                  <span>邀请奖励</span>
                  <span>#{item.referralId}</span>
                  <span className="credit-value">+{item.amount} Credit</span>
                  <time>{formatTime(item.createdAt)}</time>
                </div>
              ))}
              {!credits?.items.length && (
                <div className="ref-admin-state">
                  没有符合条件的 Credit 流水
                </div>
              )}
            </div>
          )}

          {(currentPagination?.totalPages ?? 0) > 1 && (
            <footer className="ref-admin-pagination">
              <button
                disabled={page <= 1}
                onClick={() => setPage((value) => value - 1)}
              >
                <ArrowLeft />
                上一页
              </button>
              <span>
                第 {page} / {currentPagination?.totalPages} 页
              </span>
              <button
                disabled={page >= (currentPagination?.totalPages ?? 1)}
                onClick={() => setPage((value) => value + 1)}
              >
                下一页
                <ArrowRight />
              </button>
            </footer>
          )}
        </section>
      </section>
    </main>
  )
}
