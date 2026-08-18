# Referral System Backend

一个基于 Go、Gin、Fx、Ent 和 PostgreSQL 的简化邀请返利系统。用户可以生成唯一邀请码；新用户通过邀请码注册后，邀请人获得 100 Credit。

## 功能

- 认证用户查询个人资料，管理端创建和管理用户
- 为用户幂等生成唯一邀请码和邀请链接
- 通过邀请码注册新用户
- 在单个数据库事务中创建邀请关系、发放 100 Credit、记录 Credit 流水
- Referral Dashboard：余额、成功邀请人数、邀请历史、Credit 流水
- 数据库唯一约束防止重复邮箱、重复受邀和重复奖励
- Swagger API 文档、结构化错误响应、OpenTelemetry 可观测性

## 本地运行

要求 Go 1.25+、Docker 和 Docker Compose。

```bash
cp config.example.yaml config.yaml
make compose-up
go run ./cmd/migrate up
go run ./cmd/app
```

默认 HTTP 地址为 `http://localhost:8999`，Swagger 地址为 `http://localhost:8999/swagger/index.html`。

如果已有旧版 `config.yaml`，请用精简后的 `config.example.yaml` 重新创建；真实配置和密码不应提交。

## API

| Method | Path | 说明 |
| --- | --- | --- |
| GET | `/api/v1/users/me` | 查询当前登录用户 |
| GET | `/api/v1/auth/google/login` | 跳转到 Google OAuth 授权页 |
| GET | `/api/v1/auth/google/callback` | Google 回调，注册或登录用户并签发 Token |
| POST | `/api/v1/users/:id/referral-code` | 生成或返回邀请码和链接 |
| POST | `/api/v1/referrals/register` | 被邀请人注册并触发奖励 |
| GET | `/api/v1/users/:id/referral-dashboard` | 邀请统计、历史和 Credit 流水 |
| GET | `/api/v1/admin/referrals` | 管理端查看全部邀请奖励（分页、筛选） |
| GET | `/api/v1/admin/credit-transactions` | 管理端查看全部 Credit 流水 |
| GET | `/api/v1/admin/referral-stats` | 管理端查看全局邀请统计 |
| GET | `/api/v1/admin/users` | 管理端分页查询用户 |
| GET | `/api/v1/admin/users/:id` | 管理端查看用户详情 |
| PUT | `/api/v1/admin/users/:id` | 管理端修改用户资料 |
| PUT | `/api/v1/admin/users/:id/status` | 管理端启用或禁用用户 |
| PUT | `/api/v1/admin/users/:id/roles` | 管理端分配用户角色 |

管理端接口通过签名后的 HttpOnly 登录 Cookie 识别用户，用户必须拥有 `referral:read` 或超级管理员权限。列表接口支持 `page`、`page_size`、`user_id`、`email`、`created_at_from` 和 `created_at_to` 查询参数。

Google 登录需要在 `config.yaml` 中填写 `google.client_id`、`google.client_secret` 和 Google Console 中登记过的 `google.redirect_url`。回调成功后会返回 Bearer Token，同时写入 HttpOnly 登录 Cookie；首次登录会按 Google 已验证邮箱自动创建用户。

受邀注册按题目要求无需密码或验证码，仅提交姓名和邮箱：

```json
{
  "code": "ABC12345",
  "name": "Bob",
  "email": "bob@example.com"
}
```

## 设计

代码保持单向依赖：`Controller -> Service -> Repository`。`user` 模块处理基础用户能力，`referral` 模块拥有邀请与奖励用例。

数据模型包括：

- `users`：用户、唯一邀请码、当前 Credit 余额
- `referrals`：邀请人、被邀请人、奖励值；被邀请人唯一
- `credit_transactions`：不可变 Credit 流水；每条邀请最多一条奖励流水

业务模块在领域内部区分用户端与管理端：`internal/modules/referral/{client,admin}` 和
`internal/modules/user/{client,admin}`。顶层 `module.go` 只负责 Fx 聚合装配；管理端查询、
分页、筛选和 RBAC 中间件不会进入用户端用例。

余额是便于读取的聚合值，流水是审计依据。接受邀请时，用户、邀请记录、余额增量和流水在同一个 PostgreSQL 事务内提交。邀请码使用密码学安全随机数生成，并依靠数据库唯一索引处理极低概率的碰撞。

## 验证

```bash
go test ./...
make test-integration
go vet ./...
make swag-v1
```

## 取舍与后续工作

接受邀请按题目要求无需认证且只需姓名和邮箱；成功后会为新用户建立浏览器会话。用户资料、邀请面板和管理端接口使用 Cookie/JWT 与 RBAC 保护。当前奖励固定为 100 Credit，后续可以抽象奖励规则、活动期限、邀请码停用、分页、幂等键与余额对账任务。
