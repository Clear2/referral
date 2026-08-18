# Keeper deployment

The deployment uses host Nginx and systemd. Backend binds only to localhost;
Admin is built locally as a static SPA and served directly by Nginx. PostgreSQL
and Redis must already be running on the server.

## First-time setup

Install Nginx, PostgreSQL, and Redis. Node.js and pnpm are only needed
on the development machine. Then create the local, gitignored production config
and replace every placeholder:

```bash
cp deploy/config.production.example.yaml deploy/config.production.yaml
# Edit deploy/config.production.yaml with production secrets and settings.
ssh-copy-id root@8.147.104.113
```

Every deployment uploads `deploy/config.production.yaml`, backs up the current
remote config as `config.yaml.backup-<release-id>`, and atomically replaces
`/opt/keep-2026/shared/config.yaml`. The real config remains gitignored. Set
`DEPLOY_CONFIG=/absolute/path/config.yaml` to deploy from another secure path.

## Deploy

```bash
./deploy/deploy.sh
```

Useful overrides include `DEPLOY_HOST`, `DEPLOY_USER`, `DEPLOY_PORT`,
`DEPLOY_DIR`, `DEPLOY_DOMAIN`, `DEPLOY_CONFIG`, `TLS_CERT`, `TLS_KEY`, and `GOARCH`. The script runs package tests (the Docker-based
integration suite remains a separate workflow), builds locally, uploads an
atomic release, migrates the database, restarts both services, checks health,
and retains the five newest releases.

Inspect service state and logs with:

```bash
ssh root@8.147.104.113 'systemctl status keeper-backend'
ssh root@8.147.104.113 'journalctl -u keeper-backend -f'
```

## Keep 线上凭证调试

用户成功请求 `/api/v1/activities/:code/user-info` 时，后端会自动保存或刷新该用户的
Keep 凭证。

线上故障时，优先使用保存的凭证直接回放 Keep `userInfo` 请求：

```bash
ssh root@8.147.104.113 \
  'cd /opt/keep-2026/current/backend &&
   CONFIG_PATH=/opt/keep-2026/shared/config.yaml \
   ./keepdebug --user-id KEEP_USER_ID --request userInfo'
```

确实需要手工调用 Keep 接口时，可显式查看保存的请求头：

```bash
ssh root@8.147.104.113 \
  'cd /opt/keep-2026/current/backend &&
   CONFIG_PATH=/opt/keep-2026/shared/config.yaml \
   ./keepdebug --user-id KEEP_USER_ID --show'
```

输出包含：

```text
Authorization: Bearer ...
x-user-id: ...
x-version-name: ...
```

`--show` 会在终端打印敏感凭证，不要复制到工单、聊天或普通日志中，也不要在
Nginx、应用日志或 Sentry 中记录 `Authorization`。一般应先用
`--request userInfo` 定位问题，它不会暴露原始 Token。调试程序位于
`/opt/keep-2026/current/backend/keepdebug`。

常见故障：

- `ent: keepaccount not found`：该用户尚未成功访问过 `/user-info`，请让用户重新打开活动页面。
- Keep 返回 401 或其他接口错误：保存的 Token 已过期或被撤销，请让用户重新打开活动页面以刷新凭证。
- 找不到 `keepdebug`：服务器尚未部署包含该工具的最新后端版本。

The checked-in deployment expects the local, gitignored certificate files
`deploy/nginx/referral.vivl.cc.pem` and `referral.vivl.cc.key`. It deploys them with
restricted permissions and serves the application at `https://referral.vivl.cc`.
