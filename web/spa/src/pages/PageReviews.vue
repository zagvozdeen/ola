<template>
  <div class="min-h-dvh w-full flex flex-col gap-4 py-6 pb-22">
    <HeaderMenu
      title="Отзывы"
      :edit="false"
      back="settings"
      create="reviews.create"
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
        v-for="review in reviews"
        :key="review.id"
        class="bg-gray-500/20 border border-gray-500/20 p-3 rounded-xl overflow-hidden flex gap-2"
      >
        <img
          class="size-16 object-cover rounded-lg"
          :src="review.file_content"
          alt=""
        >
        <div class="flex-1 min-w-0 flex flex-col gap-1">
          <span class="font-bold text-sm truncate">{{ review.name }}</span>
          <p class="text-xs text-gray-300 line-clamp-2">
            {{ review.content }}
          </p>
          <p class="text-xs text-gray-300">
            {{ formatDate(review.published_at) }}
          </p>
          <div class="flex gap-2 mt-1">
            <router-link
              class="bg-gray-600 hover:bg-gray-700 rounded px-2 py-1 text-xs font-bold"
              :to="{ name: 'reviews.edit', params: { uuid: review.uuid } }"
            >
              Редактировать
            </router-link>
            <NPopconfirm
              negative-text="Отмена"
              positive-text="Удалить"
              @positive-click="() => handleDeleteReview(review.uuid)"
            >
              <template #trigger>
                <button
                  class="bg-red-700 hover:bg-red-800 rounded px-2 py-1 text-xs font-bold disabled:opacity-50 cursor-pointer"
                  :disabled="deleting === review.uuid"
                >
                  {{ deleting === review.uuid ? 'Удаляем...' : 'Удалить' }}
                </button>
              </template>

              Вы уверены, что хотите удалить отзыв?
            </NPopconfirm>
          </div>
        </div>
      </li>
    </ul>

    <div
      v-if="!isLoading && reviews.length === 0"
      class="text-sm text-gray-300"
    >
      Список отзывов пуст.
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import HeaderMenu from '@/components/HeaderMenu.vue'
import { useFetch } from '@/composables/useFetch'
import { useNotifications } from '@/composables/useNotifications'
import { type Review } from '@/types'
import { NPopconfirm, NSpin } from 'naive-ui'

const fetcher = useFetch()
const notify = useNotifications()
const reviews = ref<Review[]>([])
const isLoading = ref(true)
const deleting = ref<string | null>(null)

const formatDate = (value: string) => {
  return new Date(value).toLocaleString('ru-RU')
}

const handleDeleteReview = (uuid: string) => {
  deleting.value = uuid

  fetcher.deleteReview(uuid)
    .then(data => {
      if (data.ok) {
        reviews.value = reviews.value.filter(review => review.uuid !== uuid)
        notify.info('Отзыв удалён')
      }
    })
    .finally(() => {
      deleting.value = null
    })
}

const initPage = async () => {
  isLoading.value = true
  const data = await fetcher.getReviews()
  if (data.ok) {
    reviews.value = data.data
  }
  isLoading.value = false
}

onMounted(() => {
  void initPage()
})
</script>
