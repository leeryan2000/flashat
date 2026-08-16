import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react-swc'
import tailwindcss from '@tailwindcss/vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), tailwindcss()],
  server: {
    // Mirrors nginx/flashat.conf's two /api locations for local dev,
    // where there's no nginx in front of the Vite dev server. Order
    // matters here (unlike nginx's longest-prefix match) — Vite checks
    // proxy rules in declaration order, so /api/posts must come first.
    proxy: {
      '/api/websocket': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        ws: true,
      },
      '/api/posts': {
        target: 'http://localhost:8081',
        changeOrigin: true,
      },
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
})
