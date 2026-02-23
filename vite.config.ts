import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'node:path'
import tailwindcss from '@tailwindcss/vite'

export default defineConfig({
  // root: resolve(__dirname, 'web'),
  appType: 'mpa',
  plugins: [
    vue(),
    tailwindcss(),
  ],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'web/spa/src'),
    },
  },
  // envDir: resolve(__dirname),
  build: {
    outDir: resolve(__dirname, 'dist'),
    emptyOutDir: true,
    manifest: true,
    copyPublicDir: true,
    rollupOptions: {
      input: {
        landing: resolve(__dirname, 'web/landing/index.html'),
        spa: resolve(__dirname, 'web/spa/index.html'),
      },
    },
  },
})
