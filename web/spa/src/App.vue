<template>
  <n-config-provider
    :theme="darkTheme"
    preflight-style-disabled
  >
    <n-loading-bar-provider>
      <NotificationProvider>
        <main class="max-w-md mx-auto px-4">
          <router-view />
        </main>
      </NotificationProvider>
    </n-loading-bar-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { darkTheme, NConfigProvider, NLoadingBarProvider } from 'naive-ui'
import NotificationProvider from '@/components/NotificationProvider.vue'
</script>
