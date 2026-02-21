<template>
  <n-config-provider :theme="darkTheme">
    <n-loading-bar-provider>
      <main class="max-w-md mx-auto px-4">
        <router-view />
      </main>
    </n-loading-bar-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { darkTheme, NConfigProvider, NLoadingBarProvider } from 'naive-ui'
</script>
