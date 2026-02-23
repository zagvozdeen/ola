<template>
  <div class="fixed max-w-md top-0 pt-6 pb-2 left-1/2 -translate-x-1/2 w-full px-4 z-10 bg-linear-to-t from-transparent to-gray-950">
    <div class="grid grid-cols-[80px_1fr_80px] items-center gap-2">
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
      <h1 class="font-bold text-sm text-center select-none">
        {{ title }}
      </h1>
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
    </div>
  </div>
  <div class="h-10 w-full" />
</template>

<script setup lang="ts">
defineProps<{
  title: string
  edit: boolean
  back: string
  create?: string
}>()

const emit = defineEmits<{
  (e: 'ready'): void
}>()
</script>

