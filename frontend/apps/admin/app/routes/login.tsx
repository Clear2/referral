import { useState, type FormEvent } from "react"
import {
  ArrowRight,
  Eye,
  EyeOff,
  Link2,
  LoaderCircle,
  LockKeyhole,
} from "lucide-react"
import { useSearchParams } from "react-router"
import { LanguageSwitcher } from "@referral/i18n"

type LoginResponse = {
  message?: string
  payload?: { user_id: number; name: string }
}
export function meta() {
  return [{ title: "管理登录 · Referral" }]
}

export default function AdminLogin() {
  const [searchParams] = useSearchParams()
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [showPassword, setShowPassword] = useState(false)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState("")
  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const form = new FormData(event.currentTarget)
    const submittedEmail = String(form.get("email") ?? email)
    const submittedPassword = String(form.get("password") ?? password)
    setError("")
    setLoading(true)
    try {
      const response = await fetch("/api/v1/auth/login-with-account", {
        method: "POST",
        credentials: "include",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          email: submittedEmail,
          password: submittedPassword,
        }),
      })
      const body = (await response.json().catch(() => ({}))) as LoginResponse
      if (!response.ok) throw new Error(body.message || "邮箱或密码错误")
      if (!body.payload?.user_id) throw new Error("登录响应缺少用户信息")
      const session = await fetch("/api/v1/admin/session", {
        credentials: "include",
        cache: "no-store",
      })
      if (!session.ok) {
        await fetch("/api/v1/auth/logout", {
          method: "POST",
          credentials: "include",
        })
        throw new Error(
          session.status === 403
            ? "当前账户没有管理端访问权限"
            : "登录状态验证失败"
        )
      }
      const requested = searchParams.get("next")
      const next =
        requested?.startsWith("/") && !requested.startsWith("//")
          ? requested
          : "/"
      window.location.replace(`/admin${next === "/" ? "/" : next}`)
    } catch (reason) {
      setError(reason instanceof Error ? reason.message : "登录失败")
    } finally {
      setLoading(false)
    }
  }
  return (
    <main className="register-page login-page">
      <header className="register-header">
        <a href="/admin/login" className="invite-brand">
          <span className="invite-brand-icon">
            <Link2 size={17} />
          </span>
          <span>REFERRAL</span>
          <small>管理中心</small>
        </a>
        <LanguageSwitcher
          className="invite-language-switcher"
          variant="toggle"
        />
      </header>
      <section className="login-shell">
        <div className="register-card login-card-new">
          <div className="register-heading">
            <span>
              <LockKeyhole size={15} /> 管理端登录
            </span>
            <h2>管理中心</h2>
            <p>登录后查看邀请数据与 Credit 流水。</p>
          </div>
          <form onSubmit={submit}>
            <label className="register-field">
              <span>邮箱</span>
              <input
                name="email"
                type="email"
                value={email}
                onChange={(event) => setEmail(event.target.value)}
                autoComplete="username"
                required
              />
            </label>
            <label className="register-field">
              <span>密码</span>
              <span className="register-password-control">
                <input
                  name="password"
                  type={showPassword ? "text" : "password"}
                  value={password}
                  onChange={(event) => setPassword(event.target.value)}
                  autoComplete="current-password"
                  required
                />
                <button
                  type="button"
                  onClick={() => setShowPassword((visible) => !visible)}
                  aria-label={showPassword ? "隐藏密码" : "显示密码"}
                  aria-pressed={showPassword}
                  disabled={loading}
                >
                  {showPassword ? <EyeOff /> : <Eye />}
                </button>
              </span>
            </label>
            <div
              className={`register-error ${error ? "is-visible" : ""}`}
              role="alert"
            >
              {error}
            </div>
            <button type="submit" disabled={loading}>
              {loading ? (
                <>
                  <LoaderCircle className="spin" size={17} />
                  正在登录…
                </>
              ) : (
                <>
                  进入管理中心
                  <ArrowRight size={17} />
                </>
              )}
            </button>
          </form>
        </div>
      </section>
    </main>
  )
}
