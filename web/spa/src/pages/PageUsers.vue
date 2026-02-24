<template>
  <div class="min-h-dvh w-full flex flex-col gap-4 py-6 pb-22">
    <HeaderMenu
      title="Пользователи"
      :edit="false"
      back="settings"
    />

    <div
      v-if="isLoading"
      class="flex justify-center my-4"
    >
      <n-spin size="small" />
    </div>

    <ul
      v-else
      class="grid grid-cols-1 gap-2"
    >
      <li
        v-for="user in users"
        :key="user.id"
        class="bg-gray-500/20 border border-gray-500/20 p-3 rounded-xl overflow-hidden flex flex-col gap-2"
      >
        <div class="flex justify-between gap-2">
          <span class="font-bold text-sm truncate">{{ user.first_name }} {{ user.last_name || '' }}</span>
          <span class="text-xs uppercase bg-gray-600 font-bold px-2 py-0.5 rounded-full">{{ UserRoleTranslates[user.role] }}</span>
        </div>
        <p class="text-xs text-gray-300 truncate">
          {{ user.email || user.username || 'Без e-mail' }}
        </p>
        <div class="flex justify-end">
          <router-link
            class="bg-gray-600 hover:bg-gray-700 rounded px-2 py-1 text-xs font-bold"
            :to="{ name: 'users.edit', params: { uuid: user.uuid } }"
          >
            Редактировать
          </router-link>
        </div>
      </li>
    </ul>

    <div
      v-if="!isLoading && users.length === 0"
      class="text-sm text-gray-300"
    >
      Список пользователей пуст.
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import HeaderMenu from '@/components/HeaderMenu.vue'
import { useFetch } from '@/composables/useFetch'
import { type User, UserRoleTranslates } from '@/types'
import { NSpin } from 'naive-ui'

const fetcher = useFetch()
const users = ref<User[]>([])
const isLoading = ref(true)

const initPage = async () => {
  isLoading.value = true
  const data = await fetcher.getUsers()
  if (data.ok) {
    users.value = data.data
  }
  isLoading.value = false
}

onMounted(() => {
  void initPage()
})
</script>
