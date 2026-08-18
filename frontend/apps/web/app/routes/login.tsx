import { useState, type FormEvent } from "react"
import {
  ArrowRight,
  Eye,
  EyeOff,
  Link2,
  LoaderCircle,
  LockKeyhole,
} from "lucide-react"
import { Link, useNavigate, useSearchParams } from "react-router"
import { seoMeta } from "../seo"

export function meta() {
  return seoMeta({
    title: "登录 · Referral",
    description: "登录 Referral，查看邀请进度与 Credit 奖励。",
    path: "/login",
    index: false,
  })
}

type LoginResponse = {
  message?: string
  payload?: { user_id: number; name: string }
}

export default function Login() {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [showPassword, setShowPassword] = useState(false)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState("")

  const requested = searchParams.get("next")
  const destination =
    requested?.startsWith("/") && !requested.startsWith("//") ? requested : "/"
  const registerURL =
    destination === "/"
      ? "/register"
      : `/register?next=${encodeURIComponent(destination)}`

  function continueWithGoogle() {
    window.location.assign(
      `/api/v1/auth/google/login?next=${encodeURIComponent(destination)}&origin=${encodeURIComponent(window.location.origin)}`
    )
  }

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setError("")
    setLoading(true)
    try {
      const response = await fetch("/api/v1/auth/login-with-account", {
        method: "POST",
        credentials: "include",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, password }),
      })
      const body = (await response.json().catch(() => ({}))) as LoginResponse
      if (!response.ok) throw new Error(body.message || "邮箱或密码错误")
      if (!body.payload?.user_id) throw new Error("登录响应缺少用户信息")
      navigate(destination, { replace: true })
    } catch (reason) {
      setError(
        reason instanceof Error ? reason.message : "登录失败，请稍后重试"
      )
    } finally {
      setLoading(false)
    }
  }

  return (
    <main className="register-page login-page">
      <header className="register-header">
        <a href="/login" className="invite-brand">
          <span className="invite-brand-icon">
            <Link2 size={17} />
          </span>
          <span>REFERRAL</span>
          <small>邀请奖励系统</small>
        </a>
      </header>
      <section className="login-shell">
        <div className="register-card login-card-new">
          <div className="register-heading">
            <span>
              <LockKeyhole size={15} /> 安全登录
            </span>
            <h2>欢迎回来</h2>
            <p>登录后查看邀请进度与 Credit 奖励。</p>
          </div>
          <button
            type="button"
            className="google-auth-button"
            onClick={continueWithGoogle}
            disabled={loading}
          >
            <svg viewBox="0 0 24 24" aria-hidden="true">
              <path
                fill="#4285f4"
                d="M21.6 12.23c0-.71-.06-1.4-.18-2.07H12v3.92h5.38a4.6 4.6 0 0 1-2 3.02v2.54h3.24c1.9-1.75 2.98-4.33 2.98-7.41Z"
              />
              <path
                fill="#34a853"
                d="M12 22c2.7 0 4.98-.9 6.63-2.43l-3.24-2.53c-.9.6-2.05.96-3.39.96-2.61 0-4.82-1.76-5.61-4.13H3.04v2.61A10 10 0 0 0 12 22Z"
              />
              <path
                fill="#fbbc05"
                d="M6.39 13.87A6 6 0 0 1 6.08 12c0-.65.11-1.28.31-1.87V7.52H3.04A10 10 0 0 0 2 12c0 1.61.38 3.14 1.04 4.48l3.35-2.61Z"
              />
              <path
                fill="#ea4335"
                d="M12 6c1.47 0 2.79.51 3.83 1.5l2.87-2.88A9.64 9.64 0 0 0 12 2a10 10 0 0 0-8.96 5.52l3.35 2.61C7.18 7.76 9.39 6 12 6Z"
              />
            </svg>
            使用 Google 登录或注册
          </button>
          <div className="login-divider">
            <span>或使用邮箱登录</span>
          </div>
          <form onSubmit={submit}>
            <label className="register-field">
              <span>邮箱</span>
              <input
                type="email"
                value={email}
                onChange={(event) => setEmail(event.target.value)}
                placeholder="name@example.com"
                autoComplete="username"
                required
                disabled={loading}
              />
            </label>
            <label className="register-field">
              <span>密码</span>
              <span className="register-password-control">
                <input
                  type={showPassword ? "text" : "password"}
                  value={password}
                  onChange={(event) => setPassword(event.target.value)}
                  placeholder="输入登录密码"
                  autoComplete="current-password"
                  required
                  disabled={loading}
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
            <button type="submit" disabled={loading || !email || !password}>
              {loading ? (
                <>
                  <LoaderCircle className="spin" size={17} /> 正在登录…
                </>
              ) : (
                <>
                  登录 <ArrowRight size={17} />
                </>
              )}
            </button>
          </form>
          <p className="register-note">
            没有账户？
            <Link to={registerURL}>创建新账户</Link>
          </p>
        </div>
      </section>
    </main>
  )
}
