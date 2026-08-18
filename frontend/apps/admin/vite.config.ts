import { reactRouter } from "@react-router/dev/vite"
import tailwindcss from "@tailwindcss/vite"
import { defineConfig } from "vite"

export default defineConfig({
  base: "/admin/",
  optimizeDeps: { exclude: ["lucide-react"], include: ["react-router > cookie"] },
  plugins: [tailwindcss(), reactRouter()],
  ssr: { noExternal: ["@referral/i18n"] },
  server: { proxy: { "/api": { target: process.env.API_PROXY_TARGET ?? "http://localhost:8999", changeOrigin: true } } },
})
