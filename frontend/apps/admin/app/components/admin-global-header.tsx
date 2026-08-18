import { useEffect, useState } from "react"
import { ExternalLink, Link2, LogOut } from "lucide-react"
import { useNavigate } from "react-router"

import { logout } from "@referral/api"
import { LanguageSwitcher } from "@referral/i18n"

type Session = { id: number; name: string; email: string }

export function AdminGlobalHeader() {
  const navigate = useNavigate()
  const [session, setSession] = useState<Session | null>(null)
  const [menuOpen, setMenuOpen] = useState(false)
  const [confirmOpen, setConfirmOpen] = useState(false)
  const [loggingOut, setLoggingOut] = useState(false)
  const webURL =
    typeof window !== "undefined" && window.location.port === "5174"
      ? `${window.location.protocol}//${window.location.hostname}:5173/`
      : "/"

  useEffect(() => {
    fetch("/api/v1/admin/session", {
      credentials: "include",
      cache: "no-store",
    })
      .then((response) => (response.ok ? response.json() : null))
      .then((body) => setSession(body?.payload ?? null))
      .catch(() => setSession(null))
  }, [])

  useEffect(() => {
    if (!menuOpen && !confirmOpen) return
    const close = (event: KeyboardEvent) => {
      if (event.key !== "Escape") return
      setMenuOpen(false)
      setConfirmOpen(false)
    }
    window.addEventListener("keydown", close)
    return () => window.removeEventListener("keydown", close)
  }, [confirmOpen, menuOpen])

  async function signOut() {
    setLoggingOut(true)
    try {
      await logout()
    } finally {
      navigate("/login", { replace: true })
    }
  }

  return (
    <>
      <header className="admin-global-header">
        <a className="admin-global-brand" href="/admin/">
          <span>
            <Link2 />
          </span>
          <strong>Referral Admin</strong>
        </a>
        <div className="admin-global-actions">
          <a href={webURL} title="用户邀请页">
            <ExternalLink />
            <span>用户邀请页</span>
          </a>
          <LanguageSwitcher
            className="invite-language-switcher"
            variant="toggle"
          />
          <div className="admin-account-anchor">
            <button
              className="admin-account-trigger"
              type="button"
              onClick={() => setMenuOpen((open) => !open)}
              aria-haspopup="menu"
              aria-expanded={menuOpen}
            >
              <span className="admin-account-avatar">
                {session?.name.trim().charAt(0).toUpperCase() || "A"}
              </span>
              <span>{session?.name || "Administrator"}</span>
            </button>
            {menuOpen && (
              <>
                <button
                  className="admin-account-dismiss"
                  type="button"
                  aria-label="关闭账户菜单"
                  onClick={() => setMenuOpen(false)}
                />
                <div className="admin-account-menu" role="menu">
                  <div>
                    <span className="admin-account-avatar large">
                      {session?.name.trim().charAt(0).toUpperCase() || "A"}
                    </span>
                    <span>
                      <strong>{session?.name || "Administrator"}</strong>
                      <small>{session?.email || ""}</small>
                    </span>
                  </div>
                  <hr />
                  <button
                    type="button"
                    role="menuitem"
                    onClick={() => {
                      setMenuOpen(false)
                      setConfirmOpen(true)
                    }}
                  >
                    <LogOut /> 退出登录
                  </button>
                </div>
              </>
            )}
          </div>
        </div>
      </header>
      {confirmOpen && (
        <div className="admin-logout-backdrop" role="presentation">
          <section
            className="admin-logout-modal"
            role="dialog"
            aria-modal="true"
            aria-labelledby="admin-logout-title"
          >
            <h2 id="admin-logout-title">提示</h2>
            <p>是否退出登录？</p>
            <footer>
              <button
                type="button"
                onClick={() => setConfirmOpen(false)}
                disabled={loggingOut}
              >
                取消
              </button>
              <button
                className="primary"
                type="button"
                onClick={() => void signOut()}
                disabled={loggingOut}
              >
                {loggingOut ? "正在退出…" : "确认"}
              </button>
            </footer>
          </section>
        </div>
      )}
    </>
  )
}
