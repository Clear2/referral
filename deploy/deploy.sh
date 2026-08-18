#!/usr/bin/env bash
set -Eeuo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"

DEPLOY_HOST="${DEPLOY_HOST:-8.147.104.113}"
DEPLOY_USER="${DEPLOY_USER:-root}"
DEPLOY_PORT="${DEPLOY_PORT:-22}"
DEPLOY_DIR="${DEPLOY_DIR:-/opt/keep-2026}"
DEPLOY_DOMAIN="${DEPLOY_DOMAIN:-referral.vivl.cc}"
GOARCH="${GOARCH:-amd64}"
DEPLOY_CONFIG="${DEPLOY_CONFIG:-${SCRIPT_DIR}/config.production.yaml}"
TLS_CERT="${TLS_CERT:-${SCRIPT_DIR}/nginx/${DEPLOY_DOMAIN}.pem}"
TLS_KEY="${TLS_KEY:-${SCRIPT_DIR}/nginx/${DEPLOY_DOMAIN}.key}"
RELEASE_ID="$(date +%Y%m%d%H%M%S)"
REMOTE="${DEPLOY_USER}@${DEPLOY_HOST}"
WORK_DIR="$(mktemp -d "${TMPDIR:-/tmp}/keep-deploy.XXXXXX")"

cleanup() { rm -rf "${WORK_DIR}"; }
trap cleanup EXIT

for command in go pnpm rsync ssh; do
  command -v "${command}" >/dev/null || { echo "缺少本地命令: ${command}" >&2; exit 1; }
done
GO_BINARY="$(command -v go)"
GO_ROOT="$(env -u GOROOT "${GO_BINARY}" env GOROOT)"
for tls_file in "${TLS_CERT}" "${TLS_KEY}"; do
	[[ -f "${tls_file}" ]] || { echo "缺少 TLS 文件: ${tls_file}" >&2; exit 1; }
done
[[ -f "${DEPLOY_CONFIG}" ]] || {
	echo "缺少生产配置: ${DEPLOY_CONFIG}" >&2
	echo "请复制 deploy/config.production.example.yaml 为 deploy/config.production.yaml，并填写真实配置" >&2
	exit 1
}
grep -Fq "https://${DEPLOY_DOMAIN}" "${DEPLOY_CONFIG}" || {
	echo "生产配置 ${DEPLOY_CONFIG} 尚未包含新域名 https://${DEPLOY_DOMAIN}" >&2
	echo "请更新 allowed_origins 和 OAuth redirect_url 后再部署" >&2
	exit 1
}

echo "[1/6] 检查 Frontend"
(
  cd "${ROOT_DIR}/frontend"
  pnpm install --frozen-lockfile
  pnpm typecheck
  pnpm build
)

echo "[2/6] 检查并编译 Backend (linux/${GOARCH})"
(
  cd "${ROOT_DIR}/backend"
  export GOROOT="${GO_ROOT}"
  export GOTOOLCHAIN=local
  "${GO_BINARY}" list ./... | rg -v '/integration-test$' | xargs "${GO_BINARY}" test
  mkdir -p "${WORK_DIR}/backend/ent"
  CGO_ENABLED=0 GOOS=linux GOARCH="${GOARCH}" "${GO_BINARY}" build -trimpath -ldflags="-s -w" -o "${WORK_DIR}/backend/keeper" ./cmd/app
  CGO_ENABLED=0 GOOS=linux GOARCH="${GOARCH}" "${GO_BINARY}" build -trimpath -ldflags="-s -w" -o "${WORK_DIR}/backend/migrate" ./cmd/migrate
  cp -R ent/migrate "${WORK_DIR}/backend/ent/"
)

mkdir -p "${WORK_DIR}/web/admin"
cp -R "${ROOT_DIR}/frontend/apps/web/build/client/." "${WORK_DIR}/web/"
cp -R "${ROOT_DIR}/frontend/apps/admin/build/client/." "${WORK_DIR}/web/admin/"
mkdir -p "${WORK_DIR}/systemd"
sed "s|/opt/keep-2026|${DEPLOY_DIR}|g" "${SCRIPT_DIR}/systemd/keeper-backend.service" > "${WORK_DIR}/systemd/keeper-backend.service"
sed \
  -e "s|__DEPLOY_DIR__|${DEPLOY_DIR}|g" \
  -e "s|__DEPLOY_DOMAIN__|${DEPLOY_DOMAIN}|g" \
  -e "s|__TLS_CERT__|$(basename "${TLS_CERT}")|g" \
  -e "s|__TLS_KEY__|$(basename "${TLS_KEY}")|g" \
  "${SCRIPT_DIR}/nginx/referral.conf" > "${WORK_DIR}/referral.nginx.conf"

SSH=(ssh -p "${DEPLOY_PORT}" -o ServerAliveInterval=30 -o ServerAliveCountMax=4)
RSYNC_SSH="ssh -p ${DEPLOY_PORT} -o ServerAliveInterval=30 -o ServerAliveCountMax=4"
REMOTE_RELEASE="${DEPLOY_DIR}/releases/${RELEASE_ID}"

echo "[3/6] 检查服务器环境"
"${SSH[@]}" "${REMOTE}" bash -s -- "${DEPLOY_DIR}" <<'REMOTE_CHECK'
set -Eeuo pipefail
deploy_dir="$1"
for command in nginx systemctl runuser curl rsync; do
  command -v "$command" >/dev/null || { echo "服务器缺少命令: $command" >&2; exit 1; }
done
id keeper >/dev/null 2>&1 || useradd --system --home-dir "$deploy_dir" --shell /usr/sbin/nologin keeper
mkdir -p "${deploy_dir}/releases" "${deploy_dir}/shared/logs" "${deploy_dir}/shared/home" "${deploy_dir}/shared/tls"
chown root:keeper "${deploy_dir}/shared"
chmod 0750 "${deploy_dir}/shared"
chown -R keeper:keeper "${deploy_dir}/shared/logs" "${deploy_dir}/shared/home"
REMOTE_CHECK

echo "[4/6] 上传 release ${RELEASE_ID}"
"${SSH[@]}" "${REMOTE}" "mkdir -p '${REMOTE_RELEASE}'"
rsync -az --delete -e "${RSYNC_SSH}" "${WORK_DIR}/backend/" "${REMOTE}:${REMOTE_RELEASE}/backend/"
rsync -az --delete -e "${RSYNC_SSH}" "${WORK_DIR}/web/" "${REMOTE}:${REMOTE_RELEASE}/web/"
rsync -az -e "${RSYNC_SSH}" "${WORK_DIR}/systemd/" "${REMOTE}:${DEPLOY_DIR}/shared/systemd/"
rsync -az -e "${RSYNC_SSH}" "${WORK_DIR}/referral.nginx.conf" "${REMOTE}:${DEPLOY_DIR}/shared/referral.nginx.conf"
rsync -az -e "${RSYNC_SSH}" "${DEPLOY_CONFIG}" "${REMOTE}:${DEPLOY_DIR}/shared/config.yaml.next"
rsync -az -e "${RSYNC_SSH}" "${TLS_CERT}" "${TLS_KEY}" "${REMOTE}:${DEPLOY_DIR}/shared/tls/"

echo "[5/6] 安装依赖、迁移并切换版本"
"${SSH[@]}" "${REMOTE}" bash -s -- "${DEPLOY_DIR}" "${RELEASE_ID}" "${DEPLOY_DOMAIN}" "$(basename "${TLS_CERT}")" "$(basename "${TLS_KEY}")" <<'REMOTE_DEPLOY'
set -Eeuo pipefail
deploy_dir="$1"
release_id="$2"
deploy_domain="$3"
tls_cert_name="$4"
tls_key_name="$5"
release_dir="${deploy_dir}/releases/${release_id}"
previous="$(readlink -f "${deploy_dir}/current" 2>/dev/null || true)"

chown -R keeper:keeper "$release_dir"
config_file="${deploy_dir}/shared/config.yaml"
next_config="${config_file}.next"
config_backup="${config_file}.backup-${release_id}"
restore_config() {
  if test -f "$config_backup"; then
    cp -p "$config_backup" "$config_file"
    chown keeper:keeper "$config_file"
    chmod 0600 "$config_file"
    systemctl restart keeper-backend >/dev/null 2>&1 || true
    echo "发布失败，生产配置已恢复为 ${config_backup}" >&2
  fi
}
trap restore_config ERR
test -s "$next_config"
if test -f "$config_file"; then
  cp -p "$config_file" "$config_backup"
fi
chown keeper:keeper "$next_config"
chmod 0600 "$next_config"
mv -f "$next_config" "$config_file"
chown -R root:root "${deploy_dir}/shared/tls"
chmod 0700 "${deploy_dir}/shared/tls"
chmod 0644 "${deploy_dir}/shared/tls/${tls_cert_name}"
chmod 0600 "${deploy_dir}/shared/tls/${tls_key_name}"
runuser -u keeper -- env CONFIG_PATH="${deploy_dir}/shared/config.yaml" \
  bash -c "cd '${release_dir}/backend' && ./migrate up"

ln -sfn "$release_dir" "${deploy_dir}/current.next"
mv -Tf "${deploy_dir}/current.next" "${deploy_dir}/current"

install -m 0644 "${deploy_dir}/shared/systemd/keeper-backend.service" /etc/systemd/system/keeper-backend.service
install -m 0644 "${deploy_dir}/shared/referral.nginx.conf" /etc/nginx/conf.d/referral.conf
rm -f /etc/nginx/conf.d/keeper.conf
if systemctl cat keeper-admin.service >/dev/null 2>&1; then
  systemctl disable --now keeper-admin.service || true
  rm -f /etc/systemd/system/keeper-admin.service
fi
systemctl daemon-reload
nginx -t
systemctl enable keeper-backend >/dev/null
systemctl restart keeper-backend
systemctl reload nginx

echo "已启用站点: https://${deploy_domain}/"

healthy=false
for _ in $(seq 1 20); do
  if curl --fail --silent --max-time 2 http://127.0.0.1:8999/api/v1/healthz >/dev/null; then
    healthy=true
    break
  fi
  sleep 1
done

if [[ "$healthy" != true ]]; then
  journalctl -u keeper-backend -n 40 --no-pager >&2 || true
  if [[ -n "$previous" && -d "$previous" ]]; then
    ln -sfn "$previous" "${deploy_dir}/current.next"
    mv -Tf "${deploy_dir}/current.next" "${deploy_dir}/current"
    systemctl restart keeper-backend
    echo "健康检查失败，应用文件已回滚到 ${previous}（数据库迁移不会自动回滚）" >&2
  fi
  exit 1
fi

trap - ERR

find "${deploy_dir}/releases" -mindepth 1 -maxdepth 1 -type d -printf '%T@ %p\n' \
  | sort -nr | tail -n +6 | cut -d' ' -f2- | xargs -r rm -rf
REMOTE_DEPLOY

echo "[6/6] 部署完成: https://${DEPLOY_DOMAIN}/"
echo "查看日志: ssh -p ${DEPLOY_PORT} ${REMOTE} 'journalctl -u keeper-backend -f'"
