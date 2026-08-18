# Referral Frontend

前端使用 pnpm workspace，由两个可独立开发和部署的 React Router SPA，以及两个共享包组成。

```text
apps/
  web/       用户邀请工作台、登录、邀请注册
  admin/     Referral 管理中心、Credit 审计、RBAC
packages/
  api/       Cookie 会话刷新、统一请求与错误处理
  ui/        两端共享的视觉令牌和基础样式
```

## 本地开发

先启动后端 `:8999`，然后分别启动两个前端：

```bash
pnpm install
pnpm dev:web      # http://localhost:5173
pnpm dev:admin    # http://localhost:5174/admin
```

如后端不在 `http://localhost:8999`，启动时设置 `API_PROXY_TARGET`。
管理端跳转用户端默认使用 `http://localhost:5173`，可通过 `VITE_WEB_URL` 修改。

## 验证

```bash
pnpm typecheck    # 分别检查 web 和 admin
pnpm build        # 分别生成 apps/*/build/client
pnpm test:e2e     # Playwright 验证受邀注册成功与错误路径
```

## 容器构建

Docker 镜像会组合两个独立产物，并按生产路径提供服务：

```bash
docker build -t referral-frontend .
```

生产环境将两个独立构建组合在同一个站点：用户端部署在 `/`，管理端部署在 `/admin`。两者共享后端 `/api`，但不共享路由树和构建产物。
