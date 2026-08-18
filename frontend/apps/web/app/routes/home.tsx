import {
  useCallback,
  useEffect,
  useMemo,
  useState,
  type FormEvent,
} from "react"
import {
  ArrowUpRight,
  Check,
  Copy,
  Gift,
  Link2,
  LogOut,
  Mail,
  RefreshCw,
  Send,
  Share2,
  Sparkles,
  Users,
} from "lucide-react"
import type { Route } from "./+types/home"
import { apiRequest, logout } from "@referral/api"
import { getLocale, LanguageSwitcher, useLocale } from "@referral/i18n"
import { useNavigate } from "react-router"
import { seoMeta } from "../seo"

const defaultMessageZh =
  "嗨，我正在使用 Referral。通过我的邀请链接注册，我们可以一起获得奖励。"
const defaultMessageEn =
  "Hi! I'm using Referral. Register through my invitation link so we can earn rewards together."

type ReferralRecord = {
  id: number
  invitee: { id: number; name: string; email: string }
  reward: number
  createdAt: string
}

type ReferralDashboard = {
  successfulReferrals: number
  totalCreditsEarned: number
  user: { name: string; email: string; creditBalance: number }
  referrals: ReferralRecord[]
}

const formatReferralTime = (value: string) =>
  new Intl.DateTimeFormat(getLocale(), {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
  }).format(new Date(value))

export function meta({}: Route.MetaArgs) {
  return seoMeta({
    title: "Referral | 邀请好友，追踪推荐并获得 Credit 奖励",
    description:
      "Referral 是支持中英文的邀请奖励平台，帮助你分享专属邀请链接、追踪成功推荐并管理 Credit 奖励。",
    path: "/",
  })
}

export default function ReferralHome() {
  const navigate = useNavigate()
  const locale = useLocale()
  const defaultMessage =
    locale === "en-US" ? defaultMessageEn : defaultMessageZh
  const inviteSubject =
    locale === "en-US" ? "You're invited to Referral" : "邀请你加入 Referral"
  const [inviteCode, setInviteCode] = useState("")
  const [summary, setSummary] = useState<ReferralDashboard | null>(null)
  const [userID, setUserID] = useState<number | null>(null)
  const [historyRefreshing, setHistoryRefreshing] = useState(false)
  const [dataError, setDataError] = useState("")
  const inviteURL = inviteCode
    ? `${typeof window === "undefined" ? "" : window.location.origin}/ref/${inviteCode}`
    : "正在生成邀请链接…"
  const [copied, setCopied] = useState(false)
  const [emails, setEmails] = useState("")
  const [message, setMessage] = useState(defaultMessageZh)
  const [formError, setFormError] = useState("")
  const [loggingOut, setLoggingOut] = useState(false)
  const [accountMenuOpen, setAccountMenuOpen] = useState(false)
  const [logoutConfirmOpen, setLogoutConfirmOpen] = useState(false)
  const emailList = useMemo(
    () =>
      emails
        .split(/[，,\s]+/)
        .map((item) => item.trim())
        .filter(Boolean),
    [emails]
  )

  useEffect(() => {
    setMessage((current) =>
      current === defaultMessageZh || current === defaultMessageEn
        ? defaultMessage
        : current
    )
  }, [defaultMessage])

  const handleDataError = useCallback(
    (reason: unknown) => {
      const status = (reason as Error & { status?: number }).status
      if (status === 401 || status === 403) {
        navigate("/login", { replace: true })
        return
      }
      setDataError(
        reason instanceof Error ? reason.message : "无法载入邀请数据"
      )
    },
    [navigate]
  )

  const refreshDashboard = useCallback(
    async (id: number) => {
      setHistoryRefreshing(true)
      try {
        setSummary(
          await apiRequest<ReferralDashboard>(
            `/api/v1/users/${id}/referral-dashboard`
          )
        )
        setDataError("")
      } catch (reason) {
        handleDataError(reason)
      } finally {
        setHistoryRefreshing(false)
      }
    },
    [handleDataError]
  )

  useEffect(() => {
    apiRequest<{ id: number }>("/api/v1/users/me")
      .then((session) => {
        setUserID(session.id)
        return Promise.all([
          apiRequest<{ code: string }>(
            `/api/v1/users/${session.id}/referral-code`,
            {
              method: "POST",
            }
          ),
          refreshDashboard(session.id),
        ])
      })
      .then(([invitation]) => {
        setInviteCode(invitation.code)
      })
      .catch(handleDataError)
  }, [handleDataError, refreshDashboard])

  useEffect(() => {
    if (!userID) return
    const refreshWhenVisible = () => {
      if (document.visibilityState === "visible") void refreshDashboard(userID)
    }
    window.addEventListener("focus", refreshWhenVisible)
    document.addEventListener("visibilitychange", refreshWhenVisible)
    return () => {
      window.removeEventListener("focus", refreshWhenVisible)
      document.removeEventListener("visibilitychange", refreshWhenVisible)
    }
  }, [refreshDashboard, userID])

  useEffect(() => {
    if (!accountMenuOpen && !logoutConfirmOpen) return
    const previousOverflow = document.body.style.overflow
    const closeOnEscape = (event: KeyboardEvent) => {
      if (event.key !== "Escape") return
      setAccountMenuOpen(false)
      setLogoutConfirmOpen(false)
    }
    if (logoutConfirmOpen) document.body.style.overflow = "hidden"
    window.addEventListener("keydown", closeOnEscape)
    return () => {
      document.body.style.overflow = previousOverflow
      window.removeEventListener("keydown", closeOnEscape)
    }
  }, [accountMenuOpen, logoutConfirmOpen])

  async function copyInvite() {
    if (!inviteCode) return
    await navigator.clipboard?.writeText(inviteURL)
    setCopied(true)
    window.setTimeout(() => setCopied(false), 1600)
  }

  async function shareInvite() {
    if (!inviteCode) return
    if (navigator.share) {
      await navigator.share({
        title: locale === "en-US" ? "Join Referral" : "加入 Referral",
        text: defaultMessage,
        url: inviteURL,
      })
    } else {
      await copyInvite()
    }
  }

  async function handleLogout() {
    setLoggingOut(true)
    try {
      await logout()
    } finally {
      window.location.assign("/login")
    }
  }

  function sendInvites(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const invalid = emailList.find(
      (email) => !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)
    )
    if (!emailList.length) return setFormError("请至少输入一个邮箱地址")
    if (invalid) return setFormError(`邮箱格式不正确：${invalid}`)
    setFormError("")
    const subject = encodeURIComponent(inviteSubject)
    const body = encodeURIComponent(`${message}\n\n${inviteURL}`)
    window.location.href = `mailto:${emailList.join(",")}?subject=${subject}&body=${body}`
  }

  const mailHref = `mailto:?subject=${encodeURIComponent(inviteSubject)}&body=${encodeURIComponent(`${defaultMessage}\n\n${inviteURL}`)}`

  return (
    <main className="invite-page">
      <header className="invite-header">
        <div className="invite-header-inner">
          <a className="invite-brand" href="/" aria-label="Referral 首页">
            <span className="invite-brand-icon">
              <Link2 size={17} />
            </span>
            <span className="invite-brand-copy">
              <strong>REFERRAL</strong>
              <small>邀请奖励</small>
            </span>
          </a>
          <nav className="invite-nav" aria-label="首页导航">
            <a href="#invite-link">邀请链接</a>
            <a href="#referral-history">邀请记录</a>
            <a href="#email-invite">邮件邀请</a>
          </nav>
          <div className="invite-account">
            <span className="account-credit" title="当前 Credit 余额">
              <Gift size={14} />
              <span>{summary?.user.creditBalance ?? "--"}</span>
              <small>Credit</small>
            </span>
            <LanguageSwitcher
              className="invite-language-switcher"
              variant="toggle"
            />
            <div className="account-menu-anchor">
              <button
                className="account-user"
                type="button"
                onClick={() => setAccountMenuOpen((open) => !open)}
                aria-haspopup="menu"
                aria-expanded={accountMenuOpen}
              >
                <span className="account-avatar">
                  {summary?.user.name.trim().charAt(0).toUpperCase() || "?"}
                </span>
                <span className="account-name">
                  <small>当前账户</small>
                  <strong>{summary?.user.name || "加载中"}</strong>
                </span>
              </button>
              {accountMenuOpen && (
                <>
                  <button
                    className="account-menu-dismiss"
                    type="button"
                    aria-label="关闭账户菜单"
                    onClick={() => setAccountMenuOpen(false)}
                  />
                  <div className="account-dropdown" role="menu">
                    <div className="account-dropdown-profile">
                      <span className="account-dropdown-avatar">
                        {summary?.user.name.trim().charAt(0).toUpperCase() ||
                          "?"}
                      </span>
                      <span>
                        <strong>{summary?.user.name || "加载中"}</strong>
                        <small>
                          {summary?.user.email || "正在读取账户信息"}
                        </small>
                      </span>
                    </div>
                    <div className="account-dropdown-separator" />
                    <button
                      className="account-dropdown-logout"
                      type="button"
                      role="menuitem"
                      onClick={() => {
                        setAccountMenuOpen(false)
                        setLogoutConfirmOpen(true)
                      }}
                    >
                      <LogOut size={16} /> 退出登录
                    </button>
                  </div>
                </>
              )}
            </div>
          </div>
        </div>
      </header>

      {logoutConfirmOpen && (
        <div className="logout-modal-backdrop" role="presentation">
          <section
            className="logout-modal"
            role="dialog"
            aria-modal="true"
            aria-labelledby="logout-modal-title"
          >
            <h2 id="logout-modal-title">提示</h2>
            <p>是否退出登录？</p>
            <footer>
              <button
                type="button"
                onClick={() => setLogoutConfirmOpen(false)}
                disabled={loggingOut}
              >
                取消
              </button>
              <button
                className="logout-modal-confirm"
                type="button"
                onClick={handleLogout}
                disabled={loggingOut}
              >
                {loggingOut ? "正在退出…" : "确认"}
              </button>
            </footer>
          </section>
        </div>
      )}

      <div className="invite-content">
        <section className="invite-welcome">
          <div>
            <p className="invite-kicker">
              <Sparkles size={14} /> 邀请好友
            </p>
            <h1>
              分享你的链接，
              <br />
              让奖励自然发生。
            </h1>
            <p>好友通过你的链接注册后，你将获得 100 Credit。</p>
          </div>
          <div className="invite-summary" aria-label="邀请概览">
            <div>
              <Users size={17} />
              <span>
                <strong>{summary?.successfulReferrals ?? "--"}</strong>
                <small>成功邀请</small>
              </span>
            </div>
            <div>
              <Gift size={17} />
              <span>
                <strong>{summary?.totalCreditsEarned ?? "--"}</strong>
                <small>累计 Credit</small>
              </span>
            </div>
          </div>
          {dataError && (
            <p className="invite-data-error" role="alert">
              邀请数据加载失败：{dataError}
            </p>
          )}
        </section>

        <section className="invite-panel link-panel" id="invite-link">
          <div className="panel-heading">
            <span className="panel-icon">
              <Link2 size={19} />
            </span>
            <div>
              <h2>你的邀请链接</h2>
              <p>每个链接唯一对应你的账户。</p>
            </div>
          </div>
          <div className="link-control">
            <code>{inviteURL}</code>
            <button
              type="button"
              onClick={copyInvite}
              disabled={!inviteCode}
              className={copied ? "is-copied" : ""}
            >
              {copied ? <Check size={18} /> : <Copy size={18} />}
              {copied ? "已复制" : "复制链接"}
            </button>
          </div>
          <div className="share-actions">
            <button type="button" onClick={shareInvite}>
              <Share2 size={16} /> 分享链接
            </button>
            <a href={mailHref}>
              <Mail size={16} /> 在邮件中打开
            </a>
          </div>
        </section>

        <section
          className="invite-panel referral-history-panel"
          id="referral-history"
        >
          <div className="panel-heading">
            <span className="panel-icon">
              <Users size={19} />
            </span>
            <div>
              <h2>我邀请的用户</h2>
              <p>好友通过你的邀请链接完成注册后会显示在这里。</p>
            </div>
            <button
              className="referral-history-refresh"
              type="button"
              disabled={!userID || historyRefreshing}
              onClick={() => userID && void refreshDashboard(userID)}
            >
              <RefreshCw
                className={historyRefreshing ? "spin" : ""}
                size={15}
              />
              {historyRefreshing ? "刷新中…" : "刷新记录"}
            </button>
          </div>
          {!summary ? (
            <p className="referral-history-empty">正在加载邀请记录…</p>
          ) : summary.referrals.length ? (
            <div className="referral-history-list">
              {summary.referrals.map((record) => (
                <article key={record.id}>
                  <span className="referral-history-avatar">
                    {record.invitee.name.trim().charAt(0).toUpperCase() || "?"}
                  </span>
                  <span className="referral-history-user">
                    <strong>{record.invitee.name}</strong>
                    <small>{record.invitee.email}</small>
                  </span>
                  <time dateTime={record.createdAt}>
                    {formatReferralTime(record.createdAt)}
                  </time>
                  <em>+{record.reward} Credit</em>
                </article>
              ))}
            </div>
          ) : (
            <p className="referral-history-empty">
              还没有好友通过你的链接完成注册，分享邀请链接开始邀请吧。
            </p>
          )}
        </section>

        <section className="invite-panel email-panel" id="email-invite">
          <div className="panel-heading">
            <span className="panel-icon">
              <Mail size={19} />
            </span>
            <div>
              <h2>通过邮件邀请</h2>
              <p>一次可以邀请多位好友。</p>
            </div>
          </div>
          <form onSubmit={sendInvites}>
            <label className="invite-field">
              <span>邮箱地址</span>
              <input
                value={emails}
                onChange={(event) => setEmails(event.target.value)}
                placeholder="alice@example.com, bob@example.com"
                aria-describedby="email-help email-error"
              />
              <small id="email-help">多个邮箱请用逗号分隔。</small>
            </label>
            <label className="invite-field">
              <span>邀请留言</span>
              <textarea
                value={message}
                onChange={(event) => setMessage(event.target.value)}
                rows={4}
              />
            </label>
            <div className="invite-form-footer">
              <p id="email-error" role="alert">
                {formError || "邀请链接会自动添加到邮件末尾。"}
              </p>
              <button type="submit">
                <Send size={17} /> 生成邀请邮件 <ArrowUpRight size={16} />
              </button>
            </div>
          </form>
        </section>
      </div>
    </main>
  )
}
