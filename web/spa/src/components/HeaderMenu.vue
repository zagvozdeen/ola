<template>
  <div
    class="fixed max-w-md pfft-6 pb-2 top-0 left-1/2 -translate-x-1/2 w-full px-4 z-10 bg-linear-to-t from-transparent to-gray-950"
    style="padding-top: calc(var(--tg-safe-area-inset-top, 0.5rem) + 1rem)"
  >
    <div class="grid grid-cols-[80px_1fr_80px] items-center gap-2">
      <template v-if="!authState.isTelegramEnv">
        <router-link
          v-if="edit"
          class="flex items-center justify-center py-1 px-4 text-sm font-semibold transition bg-gray-500/20 backdrop-blur-lg border border-gray-500/20 rounded-full shadow-lg hover:bg-gray-500/25 cursor-pointer"
          type="button"
          :to="{ name: back }"
        >
          Отмена
        </router-link>
        <div v-else>
          <router-link
            class="inline-flex items-center justify-center size-8 text-sm font-semibold transition bg-gray-500/20 backdrop-blur-lg border border-gray-500/20 rounded-full shadow-lg hover:bg-gray-500/25 cursor-pointer"
            type="button"
            :to="{ name: back }"
          >
            <i class="bi bi-chevron-left flex justify-center items-center" />
          </router-link>
        </div>
      </template>
      <div v-else />
      <h1 class="font-bold text-sm text-center select-none">
        {{ title }}
      </h1>
      <div v-if="authState.isTelegramEnv" />
      <template v-else>
        <button
          v-if="edit"
          class="flex items-center justify-center py-1 px-4 text-sm font-semibold transition bg-gray-500/20 backdrop-blur-lg border border-gray-500/20 rounded-full shadow-lg hover:bg-gray-500/25 cursor-pointer"
          type="button"
          @click="emit('ready')"
        >
          Готово
        </button>
        <router-link
          v-else-if="!edit && create"
          class="flex items-center justify-center py-1 px-4 text-sm font-semibold transition bg-gray-500/20 backdrop-blur-lg border border-gray-500/20 rounded-full shadow-lg hover:bg-gray-500/25 cursor-pointer"
          type="button"
          :to="{ name: create }"
        >
          Создать
        </router-link>
      </template>
    </div>
  </div>
  <div class="h-10 w-full" />
</template>

<script setup lang="ts">
import { onBeforeUnmount, onMounted } from 'vue'
import { useAuthState } from '@/composables/useAuthState'
import { useRouter } from 'vue-router'

const { title, edit, back, create = undefined } = defineProps<{
  title: string
  edit: boolean
  back: string
  create?: string
}>()

const emit = defineEmits<{
  (e: 'ready'): void
}>()

const router = useRouter()
const authState = useAuthState()

// style="top: calc(var(--tg-content-safe-area-inset-top, 0px) + var(--tg-safe-area-inset-top, 0px))"

const ready = () => emit('ready')

onMounted(() => {
  if (authState.isTelegramEnv.value) {
    if (edit) {
      window.Telegram.WebApp.SecondaryButton.setText('Готово')
      window.Telegram.WebApp.SecondaryButton.show()
      window.Telegram.WebApp.SecondaryButton.enable()
      window.Telegram.WebApp.SecondaryButton.color = window.Telegram.WebApp.themeParams.button_color || 'blue'
      window.Telegram.WebApp.SecondaryButton.textColor = window.Telegram.WebApp.themeParams.text_color || 'white'
      window.Telegram.WebApp.SecondaryButton.onClick(ready)
    }

    window.Telegram.WebApp.BackButton.show()
    window.Telegram.WebApp.BackButton.onClick(() => {
      router.push({ name: back })
    })
  }
})

onBeforeUnmount(() => {
  if (authState.isTelegramEnv.value) {
    if (edit) {
      window.Telegram.WebApp.SecondaryButton.hide()
      window.Telegram.WebApp.SecondaryButton.offClick(ready)
    }

    window.Telegram.WebApp.BackButton.hide()
  }
})
</script>
