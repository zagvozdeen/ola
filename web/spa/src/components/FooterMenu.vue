<template>
  <div class="fixed flex flex-col gap-2 w-full max-w-md px-4 bottom-4 left-1/2 -translate-x-1/2">
    <div class="grid grid-cols-3 grid-cвols-[1fr_min-content_1fr_min-content_1fr] gap-1 p-1 bg-black/5 dark:bg-gray-500/20 backdrop-blur-lg border border-black/10 dark:border-gray-500/20 rounded-full shadow-lg">
      <router-link
        class="flex flex-col rounded-full py-1 px-3 transition hover:bg-black/10 dark:hover:bg-gray-500/25 cursor-pointer text-xs font-bold text-center"
        :class="{
          'bg-black/10 dark:bg-gray-500/25': route.name === 'main',
        }"
        :to="{ name: 'main' }"
        type="button"
      >
        <i class="bi bi-balloon-fill text-sm" />
        <span>Ассортимент</span>
      </router-link>
      <router-link
        class="flex flex-col rounded-full py-1 px-3 transition hover:bg-black/10 dark:hover:bg-gray-500/25 cursor-pointer text-xs font-bold text-center relative"
        :class="{
          'bg-black/10 dark:bg-gray-500/25': route.name === 'cart',
        }"
        :to="{ name: 'cart' }"
        type="button"
      >
        <i class="bi bi-basket2-fill text-sm" />
        <small
          v-show="items"
          class="absolute left-17 top-0.5 text-xs bg-green-700 rounded-full font-bold px-1"
        >{{ items }}</small>
        <span>Корзина</span>
      </router-link>
      <router-link
        class="flex flex-col rounded-full py-1 px-3 transition hover:bg-black/10 dark:hover:bg-gray-500/25 cursor-pointer text-xs font-bold text-center"
        :class="{
          'bg-black/10 dark:bg-gray-500/25': route.name === 'settings',
        }"
        :to="{ name: 'settings' }"
        type="button"
      >
        <i class="bi bi-gear-fill text-sm" />
        <span>Настройки</span>
      </router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { cart } from '@/composables/useAuthState'

const route = useRoute()
const items = computed(() => cart.items.reduce((sum, item) => sum + item.qty, 0))
</script>
