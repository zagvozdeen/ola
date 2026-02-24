<template>
  <div class="min-h-dvh w-full flex flex-col gap-4 py-6 pb-22">
    <HeaderMenu
      title="Заявки обратной связи"
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
        v-for="item in feedbacks"
        :key="item.id"
        class="bg-gray-500/20 border border-gray-500/20 p-3 rounded-xl overflow-hidden flex flex-col gap-2"
      >
        <div class="flex justify-between gap-2">
          <span class="font-bold text-sm truncate">{{ item.name }}</span>
          <span class="text-xs uppercase bg-gray-600 font-bold px-2 py-0.5 rounded-full">{{ RequestStatusTranslates[item.status] }}</span>
        </div>
        <p class="text-xs text-gray-300 line-clamp-2">
          {{ item.content }}
        </p>
        <p class="text-xs text-gray-300">
          {{ item.phone }}
        </p>
        <div class="flex justify-end">
          <router-link
            class="bg-gray-600 hover:bg-gray-700 rounded px-2 py-1 text-xs font-bold"
            :to="{ name: 'feedback.edit', params: { uuid: item.uuid } }"
          >
            Редактировать
          </router-link>
        </div>
      </li>
    </ul>

    <div
      v-if="!isLoading && feedbacks.length === 0"
      class="text-sm text-gray-300"
    >
      Список заявок пуст.
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import HeaderMenu from '@/components/HeaderMenu.vue'
import { useFetch } from '@/composables/useFetch'
import { type Feedback, RequestStatusTranslates } from '@/types'
import { NSpin } from 'naive-ui'

const fetcher = useFetch()
const feedbacks = ref<Feedback[]>([])
const isLoading = ref(true)

const initPage = async () => {
  isLoading.value = true
  const data = await fetcher.getFeedback()
  if (data.ok) {
    feedbacks.value = data.data
  }
  isLoading.value = false
}

onMounted(() => {
  void initPage()
})
</script>
