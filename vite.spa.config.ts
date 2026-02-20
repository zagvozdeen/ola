import { resolve } from 'node:path'
import vue from '@vitejs/plugin-vue'
import { defineConfig } from 'vite'
import tailwindcss from '@tailwindcss/vite'

export default defineConfig({
  root: resolve(__dirname, 'web/spa'),
  base: '/',
  plugins: [
    vue(),
    tailwindcss(),
  ],
  build: {
    outDir: resolve(__dirname, 'dist/spa'),
    emptyOutDir: true,
  },
})
