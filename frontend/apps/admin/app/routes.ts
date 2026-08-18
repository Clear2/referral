import { type RouteConfig, index, route } from "@react-router/dev/routes"

export default [
  index("routes/dashboard.tsx"),
  route("users", "routes/users.tsx"),
  route("administrators", "routes/administrators.tsx"),
  route("login", "routes/login.tsx"),
  route("permissions", "routes/permissions.tsx"),
] satisfies RouteConfig
