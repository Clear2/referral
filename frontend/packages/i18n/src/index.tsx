import {
  createContext,
  useContext,
  useEffect,
  useState,
  type ReactNode,
} from "react"

export type Locale = "zh-CN" | "en-US"
const STORAGE_KEY = "referral.locale"
const LocaleContext = createContext<{
  locale: Locale
  select: (locale: Locale) => void
} | null>(null)

const en: Record<string, string> = {
  "API 资源": "API resources",
  保存资源关联: "Save resource mappings",
  权限资源已保存: "Permission resources saved",
  无权访问: "Access denied",
  "当前角色没有访问此管理页面的菜单权限。":
    "Your current role does not have menu access to this admin page.",
  "Referral 管理中心": "Referral admin console",
  "Referral 首页": "Referral home",
  关联角色: "Related roles",
  关闭账户中心: "Close account center",
  "刷新中…": "Refreshing…",
  刷新数据: "Refresh data",
  刷新记录: "Refresh records",
  刷新邀请记录: "Refresh referral history",
  正在刷新邀请记录: "Refreshing referral history",
  "嗨，我正在使用 Referral。通过我的邀请链接注册，我们可以一起获得奖励。":
    "Hi! I'm using Referral. Register through my invitation link so we can earn rewards together.",
  未配置路径: "No path configured",
  "正在生成邀请链接…": "Generating invitation link…",
  正在读取账户信息: "Loading account information",
  "正在退出…": "Signing out…",
  "邀请 / Credit": "Referrals / Credit",
  "邀请你加入 Referral": "You're invited to Referral",
  "Credit 入账": "Credit issued",
  "Credit 流水": "Credit transactions",
  用户邀请页: "User referral page",
  "REFERRAL OPERATIONS": "REFERRAL OPERATIONS",
  "ADMINISTRATOR ACCESS": "ADMINISTRATOR ACCESS",
  "CUSTOMER ACCOUNTS": "CUSTOMER ACCOUNTS",
  "ACCESS CONTROL": "ACCESS CONTROL",
  关闭账户菜单: "Close account menu",
  提示: "Confirmation",
  "是否退出登录？": "Sign out of this account?",
  取消: "Cancel",
  确认: "Confirm",
  菜单范围: "Menu scope",
  账户中心: "Account center",
  "这是你的普通用户账户，只能查看自己的邀请记录和奖励。":
    "This is your user account. You can only view your own referrals and rewards.",
  邀请奖励概览: "Referral rewards overview",
  首页导航: "Home navigation",
  接受邀请: "Accept invitation",
  "加入 Referral，": "Join Referral,",
  "让连接产生价值。": "and make connections rewarding.",
  注册成功后: "After registration",
  邀请码: "Invitation code",
  "邀请码 ·": "Invitation code ·",
  "填写账户信息并设置安全密码，即可完成注册。":
    "Enter your account details and create a secure password to register.",
  "填写账户信息，并通过邮箱验证码完成注册。":
    "Enter your account details and verify your email to register.",
  "至少 8 位，并包含大写字母、小写字母、数字和特殊字符":
    "At least 8 characters with uppercase, lowercase, number, and special character",
  "至少 8 位，创建后可单独重置。":
    "At least 8 characters. It can be reset later.",
  请再次输入密码: "Enter your password again",
  获取验证码: "Get verification code",
  "继续即表示你接受本系统的账户与隐私条款。":
    "By continuing, you accept the account and privacy terms.",
  管理端登录: "Admin sign in",
  "登录后查看邀请数据与 Credit 流水。":
    "Sign in to view referral data and Credit transactions.",
  创建账户: "Create account",
  "已有账户？": "Already have an account? ",
  返回登录: "Back to sign in",
  进入管理中心: "Open admin console",
  后台管理: "Admin console",
  演示环境: "Demo environment",
  演示环境验证码: "Demo verification code",
  "无需等待邮件，输入上方验证码即可继续":
    "Enter the code above to continue; no email wait is needed.",
  "暂不发送真实邮件，点击注册后将在弹窗中显示验证码。":
    "No real email is sent in demo mode. The verification code will appear in a dialog.",
  再次输入新密码: "Enter the new password again",
  创建安全密码: "Create a secure password",
  "输入至少 8 位的新密码": "Enter a new password of at least 8 characters",
  "两次输入的密码不一致。": "The passwords do not match.",
  两次输入的密码不一致: "The passwords do not match",
  "输入 6 位验证码": "Enter the 6-digit code",
  "请输入 6 位验证码": "Enter the 6-digit code",
  邀请人: "Inviter",
  "注册成功后，你可以直接登录并开始邀请好友。":
    "After registration, sign in and start inviting friends.",
  "还没有好友通过你的链接完成注册，分享邀请链接开始邀请吧。":
    "No one has registered through your link yet. Share it to get started.",
  "分享你的链接，": "Share your link,",
  "让奖励自然发生。": "and let rewards follow.",
  权限数量: "Permissions",
  菜单数量: "Menus",
  "角色关联的后台菜单资源。": "Admin menu resources associated with this role.",
  "当前角色没有关联菜单。": "This role has no associated menus.",
  "当前角色没有授予管理权限。": "This role has no admin permissions.",
  "根据当前选中的角色实时计算。": "Calculated from the selected roles.",
  "超级管理员拥有全部管理权限。": "Super administrators have all permissions.",
  "角色决定该账号能够进入哪些管理功能。":
    "Roles determine which admin features this account can access.",
  "配置角色、权限点和可访问资源，保持管理边界清晰可审计。":
    "Configure roles, permissions, and resources with clear, auditable access boundaries.",
  "集中维护可进入管理端的账号与访问状态。":
    "Manage administrator accounts and access status in one place.",
  "维护用户资料、访问状态与邀请资产。":
    "Manage user profiles, access status, and referral assets.",
  "追踪每一段邀请关系和每一笔 Credit 奖励。":
    "Track every referral relationship and Credit reward.",
  没有符合条件的邀请记录: "No matching referrals",
  "该用户还没有成功邀请。": "This user has no successful referrals.",
  "该用户还没有 Credit 流水。": "This user has no Credit transactions.",
  尚未生成: "Not generated",
  "当前账户没有查看 Referral 管理数据的权限。":
    "This account cannot view Referral admin data.",
  "当前账户没有用户管理权限。":
    "This account does not have user-management permission.",
  当前账户没有管理端访问权限: "This account does not have admin access.",
  权限资源类型: "Permission resource type",
  用户角色分配: "User role assignment",
  已分配角色: "Assigned roles",
  保存角色分配: "Save role assignment",
  角色授权已保存: "Role access saved",
  角色分配已保存: "Role assignment saved",
  管理账号角色已更新: "Administrator roles updated",
  "管理权限已移除，账号已转入普通用户":
    "Admin access removed; account converted to a standard user",
  用户资料已更新: "User profile updated",
  密码已重置: "Password reset",
  "API 同步完成": "API sync complete",
  按邮箱搜索: "Search by email",
  操作失败: "Action failed",
  数据加载失败: "Could not load data",
  登录状态验证失败: "Could not verify sign-in status",
  登录失败: "Sign in failed",
  "登录失败，请稍后重试": "Sign in failed. Please try again.",
  邮箱或密码错误: "Incorrect email or password",
  登录响应缺少用户信息: "The sign-in response is missing user information",
  获取验证码失败: "Could not get a verification code",
  "注册失败，请稍后重试": "Registration failed. Please try again.",
  请至少输入一个邮箱地址: "Enter at least one email address",
  无法载入邀请数据: "Could not load referral data",
  "邀请好友 · Referral": "Invite friends · Referral",
  "登录 · Referral": "Sign in · Referral",
  "注册 · Referral": "Register · Referral",
  "接受邀请 · Referral": "Accept invitation · Referral",
  "用户管理 · Referral Admin": "Users · Referral Admin",
  "管理账号 · Referral Admin": "Administrators · Referral Admin",
  "权限管理 · Referral Admin": "Access control · Referral Admin",
  "管理登录 · Referral": "Admin sign in · Referral",
  "分享邀请链接，获得 Credit 奖励":
    "Share invitation links and earn Credit rewards",
  "登录 Referral 邀请奖励系统": "Sign in to Referral rewards",
  "注册 Referral 邀请奖励系统": "Create a Referral rewards account",
  "通过邀请链接注册 Referral":
    "Register for Referral through an invitation link",
  邀请奖励: "Referral rewards",
  邀请奖励系统: "Referral rewards",
  管理中心: "Admin console",
  管理导航: "Administration",
  业务管理: "Operations",
  数据审计: "Data & audit",
  邀请好友: "Invite friends",
  邀请链接: "Invitation link",
  邀请记录: "Referral history",
  邮件邀请: "Email invites",
  当前账户: "Current account",
  "当前 Credit 余额": "Current Credit balance",
  退出登录: "Sign out",
  正在退出: "Signing out",
  成功邀请: "Successful referrals",
  "累计 Credit": "Total Credit",
  你的邀请链接: "Your invitation link",
  "每个链接唯一对应你的账户。": "Each link is unique to your account.",
  复制链接: "Copy link",
  已复制: "Copied",
  分享链接: "Share link",
  在邮件中打开: "Open in email",
  我邀请的用户: "People I invited",
  "好友通过你的邀请链接完成注册后会显示在这里。":
    "Friends who register through your link will appear here.",
  "正在加载邀请记录…": "Loading referral history…",
  通过邮件邀请: "Invite by email",
  "一次可以邀请多位好友。": "Invite multiple friends at once.",
  邮箱地址: "Email addresses",
  "多个邮箱请用逗号分隔。": "Separate multiple addresses with commas.",
  邀请留言: "Invitation message",
  生成邀请邮件: "Create invitation email",
  "邀请链接会自动添加到邮件末尾。":
    "Your invitation link will be appended to the email.",
  欢迎回来: "Welcome back",
  安全登录: "Secure sign in",
  "登录后查看邀请进度与 Credit 奖励。":
    "Sign in to see referral progress and Credit rewards.",
  "使用 Google 登录或注册": "Continue with Google",
  或使用邮箱登录: "Or sign in with email",
  邮箱: "Email",
  密码: "Password",
  输入登录密码: "Enter your password",
  登录: "Sign in",
  "正在登录…": "Signing in…",
  "没有账户？": "No account? ",
  创建新账户: "Create an account",
  创建你的账户: "Create your account",
  姓名: "Name",
  你的姓名: "Your name",
  确认密码: "Confirm password",
  再次输入密码: "Enter password again",
  输入邮箱验证码: "Enter email verification code",
  注册: "Register",
  "正在注册…": "Registering…",
  注册成功: "Registration complete",
  "加入 Referral": "Join Referral",
  你收到了一份邀请: "You've received an invitation",
  "邀请人获得 100 Credit": "The inviter earns 100 Credit",
  普通用户: "Users",
  管理账号: "Administrators",
  角色与权限: "Roles & permissions",
  邀请概览: "Referral overview",
  数据概览: "Dashboard",
  权限管理: "Access control",
  配置角色权限: "Configure roles",
  创建普通用户: "Create user",
  搜索姓名或邮箱: "Search name or email",
  全部状态: "All statuses",
  已启用: "Enabled",
  已停用: "Disabled",
  查询: "Search",
  用户: "User",
  角色: "Role",
  状态: "Status",
  注册时间: "Registered",
  操作: "Actions",
  启用: "Enable",
  停用: "Disable",
  编辑: "Edit",
  保存: "Save",
  关闭: "Close",
  删除: "Delete",
  新增: "New",
  全部: "All",
  仅启用: "Enabled only",
  仅停用: "Disabled only",
  加载中: "Loading",
  "正在保存…": "Saving…",
  "保存中…": "Saving…",
  "正在创建…": "Creating…",
  创建用户: "Create user",
  用户详情: "User details",
  管理账号详情: "Administrator details",
  初始密码: "Initial password",
  确认初始密码: "Confirm initial password",
  填写信息后保存更改: "Enter details and save changes",
  角色分配: "Role assignment",
  未配置角色: "No role assigned",
  权限点: "Permissions",
  菜单: "Menus",
  用户角色: "User roles",
  管理角色: "Admin role",
  权限分组: "Permission groups",
  权限点列表: "Permissions",
  角色列表: "Roles",
  菜单列表: "Menus",
  "API 列表": "APIs",
  关键词: "Keyword",
  所属模块: "Module",
  说明: "Description",
  排序: "Order",
  类型: "Type",
  访问路径: "Path",
  组件: "Component",
  请求方式: "Method",
  路由地址: "Route",
  权限名称: "Permission name",
  权限编码: "Permission code",
  角色名称: "Role name",
  角色编码: "Role code",
  菜单名称: "Menu name",
  "API 名称": "API name",
  顶级分组: "Top-level group",
  顶级菜单: "Top-level menu",
  系统: "System",
  系统角色: "System role",
  未分组: "Ungrouped",
  暂无描述: "No description",
  权限与菜单授权: "Permission and menu access",
  有效权限: "Effective permissions",
  可见菜单: "Visible menus",
  数据与审计: "Data & audit",
  全局统计: "Global statistics",
  用户总数: "Total users",
  新用户注册: "New registrations",
  "已发放 Credit": "Credit issued",
  奖励规则: "Reward policy",
  发起邀请用户: "Inviter",
  被邀请人: "Invitee",
  奖励: "Reward",
  关系: "Relationship",
  "Credit 变更流水": "Credit transactions",
  金额: "Amount",
  变更原因: "Reason",
  关联邀请: "Related referral",
  入账时间: "Credited at",
  创建时间: "Created at",
  开始日期: "Start date",
  结束日期: "End date",
  至: "to",
  获得用户: "New users",
  "正在验证管理权限…": "Verifying administrator access…",
  页面出错: "Something went wrong",
  "找不到这个管理页面。": "This admin page could not be found.",
  "无法加载管理中心，请稍后重试。":
    "The admin console could not be loaded. Please try again.",
  "确认删除这条记录？删除后无法恢复。":
    "Delete this record? This cannot be undone.",
  创建成功: "Created successfully",
  更新成功: "Updated successfully",
  删除成功: "Deleted successfully",
  账户已启用: "Account enabled",
  账户已停用: "Account disabled",
  普通用户已创建: "User created",
  显示密码: "Show password",
  隐藏密码: "Hide password",
  显示确认密码: "Show confirmation password",
  隐藏确认密码: "Hide confirmation password",
}

const patterns: Array<[RegExp, string]> = [
  [/^邀请码 · (.+)$/, "Invitation code · $1"],
  [
    /^好友通过你的链接注册后，你将获得 (\d+) Credit。$/,
    "When a friend registers through your link, you earn $1 Credit.",
  ],
  [/^共 (\d+) 个普通用户$/, "$1 users"],
  [/^共 (\d+) 个管理账号$/, "$1 administrators"],
  [/^邮箱格式不正确：(.*)$/, "Invalid email address: $1"],
  [/^邀请数据加载失败：(.*)$/, "Could not load referral data: $1"],
  [/^密码需要满足：(.*)$/, "Password requirements: $1"],
]

export function getLocale(): Locale {
  if (typeof window === "undefined") return "zh-CN"
  const saved = window.localStorage.getItem(STORAGE_KEY)
  if (saved === "zh-CN" || saved === "en-US") return saved
  return window.navigator.language.toLowerCase().startsWith("zh")
    ? "zh-CN"
    : "en-US"
}

export function useLocale(): Locale {
  return useContext(LocaleContext)?.locale ?? getLocale()
}

function translate(value: string) {
  const text = value.trim()
  const translated =
    en[text] ?? patterns.find(([pattern]) => pattern.test(text))?.[1]
  if (!translated) return value
  const result =
    en[text] ??
    patterns.reduce(
      (current, [pattern, replacement]) =>
        pattern.test(text) ? text.replace(pattern, replacement) : current,
      text
    )
  return value.replace(text, result)
}

function translateTree(root: ParentNode) {
  const walker = document.createTreeWalker(root, NodeFilter.SHOW_TEXT)
  let node: Node | null
  while ((node = walker.nextNode()))
    if (node.nodeValue) node.nodeValue = translate(node.nodeValue)
  root
    .querySelectorAll<HTMLElement>(
      "[aria-label], [title], input[placeholder], textarea[placeholder]"
    )
    .forEach((element) => {
      for (const name of ["aria-label", "title", "placeholder"]) {
        const value = element.getAttribute(name)
        if (value) element.setAttribute(name, translate(value))
      }
    })
  document.title = translate(document.title)
}

export function I18nProvider({
  children,
  defaultSwitcher = true,
}: {
  children: ReactNode
  defaultSwitcher?: boolean
}) {
  const [locale, setLocale] = useState<Locale>("zh-CN")
  useEffect(() => setLocale(getLocale()), [])
  useEffect(() => {
    document.documentElement.lang = locale
    if (locale === "zh-CN") return
    translateTree(document)
    const observer = new MutationObserver((changes) =>
      changes.forEach((change) => {
        if (change.type === "characterData" && change.target.nodeValue) {
          const translated = translate(change.target.nodeValue)
          if (translated !== change.target.nodeValue)
            change.target.nodeValue = translated
        }
        if (
          change.type === "attributes" &&
          change.target instanceof Element &&
          change.attributeName
        ) {
          const value = change.target.getAttribute(change.attributeName)
          if (value) {
            const translated = translate(value)
            if (translated !== value)
              change.target.setAttribute(change.attributeName, translated)
          }
        }
        change.addedNodes.forEach((node) => {
          if (node.nodeType === Node.TEXT_NODE && node.nodeValue)
            node.nodeValue = translate(node.nodeValue)
          else if (node instanceof HTMLElement) translateTree(node)
        })
      })
    )
    observer.observe(document.documentElement, {
      attributes: true,
      attributeFilter: ["aria-label", "title", "placeholder"],
      childList: true,
      characterData: true,
      subtree: true,
    })
    return () => observer.disconnect()
  }, [locale])
  const select = (next: Locale) => {
    window.localStorage.setItem(STORAGE_KEY, next)
    if (next === "zh-CN" && locale === "en-US") window.location.reload()
    else setLocale(next)
  }
  return (
    <LocaleContext.Provider value={{ locale, select }}>
      {children}
      {defaultSwitcher && <LanguageSwitcher />}
    </LocaleContext.Provider>
  )
}

export function LanguageSwitcher({
  className = "",
  variant = "segmented",
}: {
  className?: string
  variant?: "segmented" | "menu" | "toggle"
}) {
  const context = useContext(LocaleContext)
  if (!context) return null
  if (variant === "toggle") {
    const nextLocale = context.locale === "zh-CN" ? "en-US" : "zh-CN"
    return (
      <button
        className={`locale-toggle ${className}`.trim()}
        type="button"
        aria-label={nextLocale === "en-US" ? "Switch to English" : "切换到中文"}
        title={nextLocale === "en-US" ? "Switch to English" : "切换到中文"}
        onClick={() => context.select(nextLocale)}
      >
        {nextLocale === "en-US" ? "EN" : "中"}
      </button>
    )
  }
  if (variant === "menu") {
    return (
      <details className={`locale-menu ${className}`.trim()}>
        <summary aria-label="切换语言">
          <span aria-hidden="true">文</span>
        </summary>
        <div role="menu">
          <button
            type="button"
            role="menuitemradio"
            aria-checked={context.locale === "zh-CN"}
            className={context.locale === "zh-CN" ? "active" : ""}
            onClick={() => context.select("zh-CN")}
          >
            <span>中文</span>
            <small>简体中文</small>
          </button>
          <button
            type="button"
            role="menuitemradio"
            aria-checked={context.locale === "en-US"}
            className={context.locale === "en-US" ? "active" : ""}
            onClick={() => context.select("en-US")}
          >
            <span>English</span>
            <small>English</small>
          </button>
        </div>
      </details>
    )
  }
  return (
    <div
      className={`locale-switcher ${className}`.trim()}
      role="group"
      aria-label="Language"
    >
      <button
        type="button"
        className={context.locale === "zh-CN" ? "active" : ""}
        onClick={() => context.select("zh-CN")}
      >
        中文
      </button>
      <button
        type="button"
        className={context.locale === "en-US" ? "active" : ""}
        onClick={() => context.select("en-US")}
      >
        EN
      </button>
    </div>
  )
}
