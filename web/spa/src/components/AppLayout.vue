<template>
  <div class="min-h-dvh w-full pb-6 pt-20">
    <div
      class="fixed max-w-md w-full pb-6 px-4 top-0 left-1/2 -translate-x-1/2 z-10 bg-linear-to-t from-transparent to-gray-950"
      style="padding-top: calc(var(--tg-safe-area-inset-top, 0.5rem) + 1rem)"
    >
      <div class="grid grid-cols-[80px_1fr_80px] items-center gap-2 min-h-8">
        <router-link
          v-if="!authState.isTelegramEnv.value"
          class="inline-flex items-center justify-center size-8 text-sm font-semibold transition bg-gray-500/20 backdrop-blur-lg border border-gray-500/20 rounded-full shadow-lg hover:bg-gray-500/25 cursor-pointer"
          type="button"
          :to="{ name: back }"
        >
          <i class="bi bi-chevron-left flex justify-center items-center" />
        </router-link>
        <div v-else />
        <h1 class="font-bold text-sm text-center select-none">
          {{ title }}
        </h1>
        <div />
      </div>
    </div>

    <main>
      <slot />
    </main>

    <div
      v-if="!authState.isTelegramEnv.value"
      class="fixed max-w-md w-full pb-6 px-4 bottom-0 left-1/2 -translate-x-1/2 z-10"
    >
      <button
        class="w-full bg-blue-600 focus:bg-blue-700 hover:bg-blue-700 rounded-xl px-2 py-2 cursor-pointer"
        type="button"
        @click="() => emit('save')"
      >
        Сохранить
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">

import { useAuthState } from '@/composables/useAuthState'

const authState = useAuthState()

defineProps<{
  title: string
  back: string
}>()

const emit = defineEmits<{
  (e: 'save'): void
}>()
</script>
