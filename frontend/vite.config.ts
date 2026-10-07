import react from '@vitejs/plugin-react'
import { defineConfig } from 'vitest/config'

export default defineConfig({
  plugins: [react()],
  server: {
    // 5173 is taken by ddev-router on this machine, so use a different port.
    host: '127.0.0.1',
    port: 5180,
    strictPort: true,
    // The Go API runs on :8080; proxying keeps the browser on one origin,
    // so there is no CORS setup in development.
    proxy: {
      '/api': 'http://localhost:8080',
    },
  },
  test: {
    environment: 'jsdom',
    setupFiles: ['./src/test/setup.ts'],
  },
})
