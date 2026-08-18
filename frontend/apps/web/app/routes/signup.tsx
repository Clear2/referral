import { useEffect, useState, type FormEvent } from "react"
import {
  ArrowRight,
  Eye,
  EyeOff,
  KeyRound,
  Link2,
  LoaderCircle,
  UserPlus,
  X,
} from "lucide-react"
import { Link, useNavigate, useSearchParams } from "react-router"
import { seoMeta } from "../seo"

import {
  apiRequest,
  isSecurePassword,
  PASSWORD_REQUIREMENTS,
} from "@referral/api"

type DemoCode = { demo: boolean; code: string; message: string }
type RegisterResult = { user_id: number; name: string }

export function meta() {
  return seoMeta({
    title: "注册 · Referral",
    description: "创建 Referral 账户，开始分享邀请链接并获得 Credit 奖励。",
    path: "/register",
    index: false,
  })
}

export default function Signup() {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const [name, setName] = useState("")
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [showPassword, setShowPassword] = useState(false)
  const [confirmPassword, setConfirmPassword] = useState("")
  const [showConfirmPassword, setShowConfirmPassword] = useState(false)
  const [code, setCode] = useState("")
  const [demoCode, setDemoCode] = useState("")
  const [modalOpen, setModalOpen] = useState(false)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState("")
  const requested = searchParams.get("next")
  const destination =
    requested?.startsWith("/") && !requested.startsWith("//") ? requested : "/"
  const loginURL =
    destination === "/"
      ? "/login"
      : `/login?next=${encodeURIComponent(destination)}`

  useEffect(() => {
    if (!modalOpen) return

    const previousOverflow = document.body.style.overflow
    document.body.style.overflow = "hidden"
    return () => {
      document.body.style.overflow = previousOverflow
    }
  }, [modalOpen])

  async function requestCode(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setError("")
    if (!isSecurePassword(password)) {
      setError(`密码需要满足：${PASSWORD_REQUIREMENTS}`)
      return
    }
    if (password !== confirmPassword) {
      setError("两次输入的密码不一致")
      return
    }
    setLoading(true)
    try {
      const result = await apiRequest<DemoCode>(
        "/api/v1/auth/registration-code",
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ email }),
        }
      )
      setDemoCode(result.code)
      setCode("")
      setModalOpen(true)
    } catch (reason) {
      setError(reason instanceof Error ? reason.message : "获取验证码失败")
    } finally {
      setLoading(false)
    }
  }

  async function completeRegistration(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setError("")
    if (password !== confirmPassword) {
      setError("两次输入的密码不一致")
      setModalOpen(false)
      return
    }
    setLoading(true)
    try {
      await apiRequest<RegisterResult>("/api/v1/auth/register", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          name,
          email,
          password,
          confirm_password: confirmPassword,
          code,
        }),
      })
      navigate(destination, { replace: true })
    } catch (reason) {
      setError(
        reason instanceof Error ? reason.message : "注册失败，请稍后重试"
      )
    } finally {
      setLoading(false)
    }
  }

  return (
    <main className="register-page login-page">
      <header className="register-header">
        <Link to="/" className="invite-brand">
          <span className="invite-brand-icon">
            <Link2 size={17} />
          </span>
          <span>REFERRAL</span>
          <small>邀请奖励系统</small>
        </Link>
      </header>
      <section className="login-shell">
        <div className="register-card login-card-new">
          <div className="register-heading">
            <span>
              <UserPlus size={15} /> 创建账户
            </span>
            <h2>加入 Referral</h2>
            <p>填写账户信息，并通过邮箱验证码完成注册。</p>
          </div>
          <div className="demo-environment-note">
            <strong>演示环境</strong>
            <span>暂不发送真实邮件，点击注册后将在弹窗中显示验证码。</span>
          </div>
          <form onSubmit={requestCode}>
            <label className="register-field">
              <span>姓名</span>
              <input
                value={name}
                onChange={(event) => setName(event.target.value)}
                autoComplete="name"
                maxLength={100}
                required
                disabled={loading}
              />
            </label>
            <label className="register-field">
              <span>邮箱</span>
              <input
                type="email"
                value={email}
                onChange={(event) => setEmail(event.target.value)}
                placeholder="name@example.com"
                autoComplete="email"
                maxLength={254}
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
                  placeholder="创建安全密码"
                  autoComplete="new-password"
                  minLength={8}
                  maxLength={72}
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
              <small
                className={
                  password && !isSecurePassword(password)
                    ? "password-hint is-invalid"
                    : "password-hint"
                }
              >
                {PASSWORD_REQUIREMENTS}
              </small>
            </label>
            <label className="register-field">
              <span>确认密码</span>
              <span className="register-password-control">
                <input
                  type={showConfirmPassword ? "text" : "password"}
                  value={confirmPassword}
                  onChange={(event) => setConfirmPassword(event.target.value)}
                  placeholder="再次输入密码"
                  autoComplete="new-password"
                  minLength={8}
                  maxLength={72}
                  required
                  disabled={loading}
                />
                <button
                  type="button"
                  onClick={() => setShowConfirmPassword((visible) => !visible)}
                  aria-label={
                    showConfirmPassword ? "隐藏确认密码" : "显示确认密码"
                  }
                  aria-pressed={showConfirmPassword}
                  disabled={loading}
                >
                  {showConfirmPassword ? <EyeOff /> : <Eye />}
                </button>
              </span>
              <small
                className={
                  confirmPassword && password !== confirmPassword
                    ? "password-hint is-invalid"
                    : "password-hint"
                }
              >
                {confirmPassword && password !== confirmPassword
                  ? "两次输入的密码不一致"
                  : "请再次输入密码"}
              </small>
            </label>
            <div
              className={`register-error ${error ? "is-visible" : ""}`}
              role="alert"
            >
              {error}
            </div>
            <button
              type="submit"
              disabled={
                loading ||
                !name.trim() ||
                !email.trim() ||
                !isSecurePassword(password) ||
                password !== confirmPassword
              }
            >
              {loading ? (
                <>
                  <LoaderCircle className="spin" size={17} /> 正在获取验证码…
                </>
              ) : (
                <>
                  注册 <ArrowRight size={17} />
                </>
              )}
            </button>
          </form>
          <p className="register-note">
            已有账户？
            <Link to={loginURL}>返回登录</Link>
          </p>
        </div>
      </section>

      {modalOpen && (
        <div className="verification-backdrop" role="presentation">
          <section
            className="verification-modal"
            role="dialog"
            aria-modal="true"
            aria-labelledby="verification-title"
          >
            <button
              className="verification-close"
              type="button"
              onClick={() => setModalOpen(false)}
              aria-label="关闭"
            >
              <X size={17} />
            </button>
            <span className="verification-icon">
              <KeyRound size={22} />
            </span>
            <h2 id="verification-title">输入邮箱验证码</h2>
            <p>
              验证码已发送至 <strong>{email}</strong>
            </p>
            <div className="demo-code-box">
              <span>演示环境验证码</span>
              <strong>{demoCode}</strong>
              <small>无需等待邮件，输入上方验证码即可继续</small>
            </div>
            <form onSubmit={completeRegistration}>
              <input
                className="verification-input"
                value={code}
                onChange={(event) =>
                  setCode(event.target.value.replace(/\D/g, "").slice(0, 6))
                }
                inputMode="numeric"
                autoComplete="one-time-code"
                placeholder="请输入 6 位验证码"
                autoFocus
                required
              />
              <div
                className={`register-error ${error ? "is-visible" : ""}`}
                role="alert"
              >
                {error}
              </div>
              <button type="submit" disabled={loading || code.length !== 6}>
                {loading ? (
                  <>
                    <LoaderCircle className="spin" size={17} /> 正在注册…
                  </>
                ) : (
                  <>
                    完成注册 <ArrowRight size={17} />
                  </>
                )}
              </button>
            </form>
          </section>
        </div>
      )}
    </main>
  )
}
