<template>
  <n-config-provider
    :theme="isDarkTheme ? darkTheme : lightTheme"
    preflight-style-disabled
  >
    <n-loading-bar-provider>
      <NotificationProvider>
        <main class="max-w-md mx-auto min-h-dvh px-4 bg-white text-slate-900 dark:bg-transparent dark:text-white">
          <router-view />
        </main>
      </NotificationProvider>
    </n-loading-bar-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue'
import { darkTheme, lightTheme, NConfigProvider, NLoadingBarProvider } from 'naive-ui'
import NotificationProvider from '@/components/NotificationProvider.vue'
import { useAuthState } from '@/composables/useAuthState'

const authState = useAuthState()
const isDarkTheme = ref(true)
let mediaQuery: MediaQueryList | null = null

const syncTheme = () => {
  if (!authState.isTelegramEnv.value) {
    document.documentElement.classList.add('dark')
    isDarkTheme.value = true
    return
  }

  const telegramColorScheme = window.Telegram?.WebApp?.colorScheme

  if (telegramColorScheme === 'dark') {
    document.documentElement.classList.add('dark')
    isDarkTheme.value = true
    return
  }

  if (telegramColorScheme === 'light') {
    document.documentElement.classList.remove('dark')
    isDarkTheme.value = false
    return
  }

  isDarkTheme.value = window.matchMedia('(prefers-color-scheme: dark)').matches
}

onMounted(() => {
  syncTheme()

  mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
  mediaQuery.addEventListener?.('change', syncTheme)
  window.Telegram?.WebApp?.onEvent('themeChanged', syncTheme)
})

onBeforeUnmount(() => {
  mediaQuery?.removeEventListener?.('change', syncTheme)
  window.Telegram?.WebApp?.offEvent('themeChanged', syncTheme)
})
</script>
