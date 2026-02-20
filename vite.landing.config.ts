// vite.landing.config.ts
import { defineConfig } from 'vite'
import { resolve } from 'node:path'
import tailwindcss from '@tailwindcss/vite'

export default defineConfig({
  root: resolve(__dirname, 'web/landing'),
  base: '/',
  appType: 'mpa',
  plugins: [
    tailwindcss(),
  ],
  build: {
    outDir: resolve(__dirname, 'dist/landing'),
    emptyOutDir: true,
    manifest: true,
  },
})
