<template>
  <div class="min-h-dvh w-full flex flex-col gap-4 items-center justify-center py-6">
    <ul class="flex flex-col gap-px w-full rounded-2xl border border-gray-500/30 overflow-hidden">
      <li class="w-full">
        <router-link
          class="grid grid-cols-[min-content_1fr_min-content] items-center w-full gap-2 p-2 cursor-pointer bg-gray-500/20 hover:bg-gray-500/30"
          type="button"
          :to="{ name: 'main' }"
        >
          <span class="size-6 flex items-center justify-center rounded-lg bg-gray-400">
            <i class="bi bi-list-check text-sm flex" />
          </span>
          <span class="text-left text-sm font-medium">Все карточки</span>
          <span class="text-gray-400">
            <i class="bi bi-chevron-right text-sm flex" />
          </span>
        </router-link>
      </li>
      <li class="w-full">
        <router-link
          class="grid grid-cols-[min-content_1fr_min-content] items-center w-full gap-2 p-2 cursor-pointer bg-gray-500/20 hover:bg-gray-500/30"
          type="button"
          :to="{ name: 'main' }"
        >
          <span class="size-6 flex items-center justify-center rounded-lg bg-orange-400">
            <i class="bi bi-graph-up-arrow text-sm flex" />
          </span>
          <span class="text-left text-sm font-medium">Статистика</span>
          <span class="text-gray-400">
            <i class="bi bi-chevron-right text-sm flex" />
          </span>
        </router-link>
      </li>
      <li class="w-full">
        <router-link
          class="grid grid-cols-[min-content_1fr_min-content] items-center w-full gap-2 p-2 cursor-pointer bg-gray-500/20 hover:bg-gray-500/30"
          type="button"
          :to="{ name: 'main' }"
        >
          <span class="size-6 flex items-center justify-center rounded-lg bg-purple-400">
            <i class="bi bi-bookmark-star-fill text-sm flex" />
          </span>
          <span class="text-left text-sm font-medium">Начать тест</span>
          <span class="text-gray-400">
            <i class="bi bi-chevron-right text-sm flex" />
          </span>
        </router-link>
      </li>
      <li class="w-full">
        <a
          class="grid grid-cols-[min-content_1fr_min-content] items-center w-full gap-2 p-2 cursor-pointer bg-gray-500/20 hover:bg-gray-500/30"
          type="button"
          href="/admin/"
        >
          <span class="size-6 flex items-center justify-center rounded-lg bg-purple-400">
            <i class="bi bi-bookmark-star-fill text-sm flex" />
          </span>
          <span class="text-left text-sm font-medium">Перейти в панель администратора</span>
          <span class="text-gray-400">
            <i class="bi bi-chevron-right text-sm flex" />
          </span>
        </a>
      </li>
    </ul>

    <FooterMenu />
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import FooterMenu from '@/components/FooterMenu.vue'

onMounted(() => {
})
</script>
