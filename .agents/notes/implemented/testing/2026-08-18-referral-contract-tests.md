# Agent Note: Referral registration contract tests

Status: implemented

## Problem

受邀注册曾要求密码和演示验证码，超出了题目规定的姓名、邮箱最小契约；邀请事务也只在 SQLite 上验证，缺少真实 PostgreSQL 并发证据。

## Decision

受邀注册只接收邀请码、姓名和邮箱，成功后仍建立浏览器会话。普通独立注册继续使用密码和验证码。测试分为 Go 公共 HTTP/服务契约、隔离 PostgreSQL 集成测试和 Playwright 浏览器契约测试。

## Alternatives considered

没有复用独立注册的密码与验证码，因为这会改变面试题明确要求的核心路径。没有用 SQLite 代替最终集成测试，因为事务、约束和并发更新必须由 PostgreSQL 验证。

## Consequences

受邀用户初始没有密码，但可使用创建后的浏览器会话；后续可增加单独的设置密码能力。集成测试需要 Docker，浏览器测试首次运行需要安装 Playwright Chromium。

## Verification

- `go test ./internal/modules/referral/client`
- `make test-integration`
- `pnpm typecheck && pnpm test:e2e`
