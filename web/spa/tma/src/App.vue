<template>
  <n-config-provider :theme="darkTheme">
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
import NotificationProvider from '@shared/components/NotificationProvider.vue'
import { useFetch } from '@shared/composables/useFetch'
import { onMounted } from 'vue'
import { me } from '@shared/composables/useState'

const fetcher = useFetch()

onMounted(() => {
  fetcher
    .getMe()
    .then(data => {
      if (data.ok) {
        me.value = data.data
      }
    })
})
</script>
