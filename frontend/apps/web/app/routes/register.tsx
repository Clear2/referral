import { useState, type FormEvent } from "react"
import { ArrowRight, Check, Gift, Link2, LoaderCircle, UserPlus } from "lucide-react"
import { Link, useParams } from "react-router"

import { apiRequest } from "@referral/api"
import { seoMeta } from "../seo"

type Registration = {
  invitee: { id: number; name: string; email: string; creditBalance: number }
  reward: number
  inviterCreditBalance: number
}

export function meta() {
  return seoMeta({
    title: "接受邀请 · Referral",
    description: "通过专属邀请链接加入 Referral。",
    path: "/ref/invitation",
    index: false,
  })
}

export default function ReferralRegister() {
  const { code = "" } = useParams()
  const [name, setName] = useState("")
  const [email, setEmail] = useState("")
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState("")
  const [result, setResult] = useState<Registration | null>(null)

  async function register(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setError("")
    setLoading(true)
    try {
      setResult(
        await apiRequest<Registration>("/api/v1/referrals/register", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ code, name, email }),
        })
      )
    } catch (reason) {
      setError(reason instanceof Error ? reason.message : "注册失败，请稍后重试")
    } finally {
      setLoading(false)
    }
  }

  return (
    <main className="register-page">
      <header className="register-header">
        <Link to="/" className="invite-brand">
          <span className="invite-brand-icon"><Link2 size={17} /></span>
          <span>REFERRAL</span>
          <small>邀请奖励系统</small>
        </Link>
      </header>

      <section className="register-shell">
        <div className="register-story">
          <span className="register-story-icon"><Gift /></span>
          <p>你收到了一份邀请</p>
          <h1>加入 Referral，<br />让连接产生价值。</h1>
          <div className="register-reward">
            <span>注册成功后</span>
            <strong>邀请人获得 100 Credit</strong>
          </div>
          <code>邀请码 · {code.toUpperCase()}</code>
        </div>

        <div className="register-card">
          {result ? (
            <div className="register-success" role="status">
              <span><Check size={28} /></span>
              <p>注册成功</p>
              <h2>欢迎加入，{result.invitee.name}</h2>
              <p>邀请关系已建立，{result.reward} Credit 已发放给邀请人。</p>
              <Link to="/">进入 Referral <ArrowRight size={17} /></Link>
            </div>
          ) : (
            <>
              <div className="register-heading">
                <span><UserPlus size={15} /> 接受邀请</span>
                <h2>加入 Referral</h2>
                <p>按题目要求，只需填写姓名和邮箱即可完成注册。</p>
              </div>
              <form onSubmit={register}>
                <label className="register-field">
                  <span>姓名</span>
                  <input value={name} onChange={(event) => setName(event.target.value)} placeholder="你的姓名" autoComplete="name" maxLength={100} required disabled={loading} />
                </label>
                <label className="register-field">
                  <span>邮箱</span>
                  <input type="email" value={email} onChange={(event) => setEmail(event.target.value)} placeholder="name@example.com" autoComplete="email" maxLength={254} required disabled={loading} />
                </label>
                <div className={`register-error ${error ? "is-visible" : ""}`} role="alert">{error}</div>
                <button type="submit" disabled={loading || !name.trim() || !email.trim()}>
                  {loading ? <><LoaderCircle className="spin" size={17} /> 正在注册…</> : <>接受邀请 <ArrowRight size={17} /></>}
                </button>
              </form>
              <p className="register-note">继续即表示你接受本系统的账户与隐私条款。</p>
            </>
          )}
        </div>
      </section>
    </main>
  )
}
