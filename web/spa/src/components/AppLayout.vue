<template>
  <div
    class="min-h-dvh w-full pb-22"
    style="padding-top: calc(var(--tg-content-safe-area-inset-top) + var(--tg-safe-area-inset-top) + calc(var(--spacing) * 6));"
  >
    <div
      class="max-w-md w-full pb-6 px-4 top-0 left-1/2 -translate-x-1/2 z-10 bg-linear-to-t from-transparent to-gray-950"
      :class="{ 'fixed': back || authState.isTelegramEnv.value, 'absolute': !(back || authState.isTelegramEnv.value) }"
      style="padding-top: calc(var(--tg-safe-area-inset-top) + calc(var(--spacing) * 2))"
    >
      <div
        class="grid items-center min-h-8"
        :class="{ 'grid-cols-[80px_1fr_80px] gap-2': back || authState.isTelegramEnv.value, 'grid-cols-[min-content_1fr_min-content]': !(back || authState.isTelegramEnv.value) }"
      >
        <router-link
          v-if="!authState.isTelegramEnv.value && back"
          class="inline-flex items-center justify-center size-8 text-sm font-semibold transition bg-gray-500/20 backdrop-blur-lg border border-gray-500/20 rounded-full shadow-lg hover:bg-gray-500/25 cursor-pointer"
          type="button"
          :to="{ name: back }"
        >
          <i class="bi bi-chevron-left flex justify-center items-center" />
        </router-link>
        <div v-else />
        <h1
          class="font-bold select-none"
          :class="{ 'text-sm text-center': back || authState.isTelegramEnv.value, 'text-lg text-left': !(back || authState.isTelegramEnv.value) }"
        >
          {{ title }}
        </h1>
        <div />
      </div>
    </div>

    <main>
      <slot />
    </main>

    <div
      v-if="!authState.isTelegramEnv.value && back && !noSave"
      class="fixed max-w-md w-full pb-6 px-4 bottom-0 left-1/2 -translate-x-1/2 z-10"
    >
      <button
        class="w-full bg-blue-600 focus:bg-blue-700 hover:bg-blue-700 rounded-xl px-2 py-2 cursor-pointer"
        type="button"
        @click="save"
      >
        Сохранить
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAuthState } from '@/composables/useAuthState'
import { onBeforeUnmount, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const authState = useAuthState()

const { title, back = undefined, noSave = false } = defineProps<{
  title: string
  back?: string
  noSave?: boolean
}>()

const emit = defineEmits<{
  (e: 'save'): void
}>()

const save = () => emit('save')
const previous = () => router.push({ name: back })

onMounted(() => {
  if (authState.isTelegramEnv.value) {
    if (back && !noSave) {
      window.Telegram.WebApp.SecondaryButton.setText('Сохранить')
      window.Telegram.WebApp.SecondaryButton.show()
      window.Telegram.WebApp.SecondaryButton.enable()
      window.Telegram.WebApp.SecondaryButton.color = window.Telegram.WebApp.themeParams.button_color || 'blue'
      window.Telegram.WebApp.SecondaryButton.textColor = window.Telegram.WebApp.themeParams.text_color || 'white'
      window.Telegram.WebApp.SecondaryButton.onClick(save)
    }

    if (back) {
      window.Telegram.WebApp.BackButton.show()
      window.Telegram.WebApp.BackButton.onClick(previous)
    }
  }
})

onBeforeUnmount(() => {
  if (authState.isTelegramEnv.value) {
    if (!noSave) {
      window.Telegram.WebApp.SecondaryButton.offClick(save)
      window.Telegram.WebApp.SecondaryButton.hide()
    }

    if (back) {
      window.Telegram.WebApp.BackButton.offClick(previous)
      window.Telegram.WebApp.BackButton.hide()
    }
  }
})
</script>
