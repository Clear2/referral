import { type RouteConfig, index, route } from "@react-router/dev/routes"

export default [
  index("routes/home.tsx"),
  route("login", "routes/login.tsx"),
  route("register", "routes/signup.tsx"),
  route("ref/:code", "routes/register.tsx"),
] satisfies RouteConfig
