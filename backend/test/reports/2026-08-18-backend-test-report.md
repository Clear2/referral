# Referral 后端版本测试报告

## 1. 版本信息

| 项目 | 内容 |
| --- | --- |
| 测试版本 | Local RC · 2026-08-18 |
| 测试时间 | 2026-08-18（Asia/Shanghai） |
| 测试范围 | `backend/` Go 服务、认证、Referral、Credit、RBAC、用户管理及 HTTP API |
| 运行平台 | macOS arm64 |
| Go | 1.25.5（项目声明 1.25.3） |
| k6 | 2.0.0 |
| 服务地址 | `http://localhost:8999` |
| Git 标识 | 当前目录未检测到 Git 元数据，无法记录 commit SHA |

## 2. 测试结论

**结论：READY WITH CAVEATS**

Go 单元测试、竞态检测和本地 k6 API 测试均通过。k6 共完成 99 项断言，成功率 100%，HTTP 请求失败率 0%，响应时间满足当前阈值。

唯一未完成项为 PostgreSQL 容器集成测试：本机没有可用的 `docker` 命令，因此 `make test-integration` 未能启动独立测试数据库。这不影响已完成测试的结果，但发布前应在具备 Docker 的 CI 或开发环境补跑。

## 3. 执行结果

| 检查 | 命令 | 结果 |
| --- | --- | --- |
| 全量 Go 测试 | `go test ./...` | 通过 |
| 竞态与覆盖率测试 | `make test` | 通过 |
| k6 API 冒烟/轻负载 | `VUS=1 DURATION=5s k6 run --quiet test/k6/backend-smoke.js` | 通过 |
| PostgreSQL 集成测试 | `make test-integration` | 未执行：本机缺少 Docker |

`make test` 在 macOS 链接阶段输出了 `LC_DYSYMTAB` 警告，但测试进程正常退出且所有测试通过；该警告属于本地链接器兼容性提示，不是测试失败。

## 4. Go 测试覆盖情况

已通过的主要模块：

- `internal/modules/auth`：安全跳转、OAuth 用户创建、邮箱验证、管理员引导、Cookie 身份验证。
- `internal/modules/rbac`：超级管理员权限、资源授权、空关系数组、菜单授权、系统角色保护、API 同步。
- `internal/modules/referral/client`：邀请码幂等、邀请注册、重复邮箱冲突、事务性奖励、邀请数据面板。
- `internal/modules/user/admin`：普通用户与管理账号分离、自锁定保护、分页、密码哈希、密码确认。
- `internal/modules/user/client`：普通用户创建。
- `internal/router`：本地与生产前端来源的 CORS 预检。
- `pkg/passwordpolicy`：长度、大小写、数字、特殊字符和空白字符校验。

本次 `make test` 报告的主要包覆盖率：

| 包 | 覆盖率 |
| --- | ---: |
| `internal/modules/auth` | 23.1% |
| `internal/modules/rbac` | 28.2% |
| `internal/modules/referral/client` | 40.0% |
| `internal/modules/user/admin` | 29.3% |
| `internal/modules/user/client` | 13.0% |
| `internal/router` | 17.1% |

部分基础设施和组合包没有独立测试文件，覆盖率为 0%。本报告不将整体覆盖率作为发布阻断阈值。

## 5. k6 测试范围

k6 脚本位置：[`backend-smoke.js`](../k6/backend-smoke.js)

已验证接口：

| 领域 | 接口 |
| --- | --- |
| 健康检查 | `GET /api/v1/healthz` |
| 认证 | `POST /api/v1/auth/login-with-account` |
| 普通用户 | `GET /api/v1/users/me` |
| Referral | `POST /api/v1/users/:id/referral-code` |
| Referral | `GET /api/v1/users/:id/referral-dashboard` |
| 管理会话 | `GET /api/v1/admin/session` |
| RBAC | `GET /api/v1/access/me` |
| 管理统计 | `GET /api/v1/admin/referral-stats` |
| 邀请审计 | `GET /api/v1/admin/referrals` |
| Credit 审计 | `GET /api/v1/admin/credit-transactions` |
| 普通用户管理 | `GET /api/v1/admin/users?account_type=customer` |
| 管理账号管理 | `GET /api/v1/admin/users?account_type=admin` |
| RBAC 配置 | `GET /api/v1/admin/rbac` |

k6 最终指标：

| 指标 | 结果 | 阈值 | 状态 |
| --- | ---: | ---: | --- |
| 断言 | 99 / 99 | `rate > 99%` | 通过 |
| HTTP 请求 | 42 | — | 完成 |
| HTTP 失败率 | 0.00% | `< 1%` | 通过 |
| 总体 p95 | 168.71 ms | `< 750 ms` | 通过 |
| 总体 p99 | 186.00 ms | `< 1500 ms` | 通过 |
| 管理接口 p95 | 180.40 ms | `< 750 ms` | 通过 |
| 用户接口 p95 | 82.51 ms | `< 750 ms` | 通过 |
| 最大响应时间 | 188.20 ms | — | 记录 |

## 6. 数据与安全说明

- k6 脚本不保存账号密码，凭据仅通过环境变量注入。
- 测试不会批量创建用户、邀请关系或 Credit 流水。
- 邀请码接口会返回账号已有邀请码，按当前实现可安全重复调用。
- 本次短测使用引导管理员账号同时执行用户自身接口；脚本支持通过 `USER_EMAIL` 和 `USER_PASSWORD` 单独指定普通用户。
- RBAC 响应已验证 `roles`、`permissions`、`menu_ids` 和 `menus` 始终为数组，不再返回 `null`。

## 7. 发布前建议

1. 在具备 Docker 的环境执行 `make test-integration`，验证真实 PostgreSQL 迁移与邀请奖励事务。
2. 使用独立普通用户凭据再执行一次 k6，隔离普通用户与管理员身份场景。
3. CI 中保留 `go test ./...`、`make test` 和短时 k6 冒烟检查。
4. 若进行容量评估，将 `VUS` 和 `DURATION` 分阶段提升，并为测试数据库准备可回收的独立数据集。

## 8. 复现命令

```bash
cd /Users/xx/Work/own2/referral/backend

go test ./...
make test

ADMIN_EMAIL='admin@example.com' \
ADMIN_PASSWORD='replace-me' \
USER_EMAIL='user@example.com' \
USER_PASSWORD='replace-me' \
VUS=5 DURATION=1m \
make k6-smoke

make test-integration
```
