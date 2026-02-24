<template>
  <div class="min-h-dvh w-full flex flex-col gap-4 items-center justify-center py-6">
    <ul class="flex flex-col gap-px w-full rounded-2xl border border-gray-500/30 overflow-hidden">
      <li
        v-if="me && isUserModerator(me)"
        class="w-full"
      >
        <router-link
          class="grid grid-cols-[min-content_1fr_min-content] items-center w-full gap-2 p-2 cursor-pointer bg-gray-500/20 hover:bg-gray-500/30"
          type="button"
          :to="{ name: 'products' }"
        >
          <span class="size-6 flex items-center justify-center rounded-lg bg-blue-500">
            <i class="bi bi-box-seam text-sm flex" />
          </span>
          <span class="text-left text-sm font-medium">Управление продуктами</span>
          <span class="text-gray-400">
            <i class="bi bi-chevron-right text-sm flex" />
          </span>
        </router-link>
      </li>
    </ul>

    <FooterMenu />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import FooterMenu from '@/components/FooterMenu.vue'
import { isUserModerator, useAuthState } from '@/composables/useAuthState'

const auth = useAuthState()
const me = computed(() => auth.currentUser.value)

onMounted(() => {
  void auth.ensureUserLoaded()
})
</script>
