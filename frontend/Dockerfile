FROM node:20-alpine AS build-env
WORKDIR /app
RUN corepack enable
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
COPY apps/web/package.json apps/web/package.json
COPY apps/admin/package.json apps/admin/package.json
COPY packages/api/package.json packages/api/package.json
COPY packages/ui/package.json packages/ui/package.json
RUN pnpm install --frozen-lockfile
COPY . .
RUN pnpm build

FROM nginx:1.29-alpine
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build-env /app/apps/web/build/client /usr/share/nginx/html
COPY --from=build-env /app/apps/admin/build/client /usr/share/nginx/html/admin
