<template>
  <AppLayout
    title="Категории"
    back="settings"
    no-save
  >
    <!--  <div class="min-h-dvh w-full flex flex-col gap-4 py-6 pb-22">-->
    <!--    <HeaderMenu-->
    <!--      title="Категории"-->
    <!--      :edit="false"-->
    <!--      back="settings"-->
    <!--      create="categories.create"-->
    <!--    />-->

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
        v-for="category in categories"
        :key="category.id"
        class="bg-gray-500/20 border border-gray-500/20 p-3 rounded-xl overflow-hidden flex justify-between items-center gap-2"
      >
        <span class="font-bold text-sm truncate">{{ category.name }}</span>
        <div class="flex gap-2 shrink-0">
          <router-link
            class="bg-gray-600 hover:bg-gray-700 rounded px-2 py-1 text-xs font-bold"
            :to="{ name: 'categories.edit', params: { uuid: category.uuid } }"
          >
            Редактировать
          </router-link>
          <NPopconfirm
            negative-text="Отмена"
            positive-text="Удалить"
            @positive-click="() => handleDeleteCategory(category.uuid)"
          >
            <template #trigger>
              <button
                class="bg-red-700 hover:bg-red-800 rounded px-2 py-1 text-xs font-bold disabled:opacity-50 cursor-pointer"
                :disabled="deleting === category.uuid"
              >
                {{ deleting === category.uuid ? 'Удаляем...' : 'Удалить' }}
              </button>
            </template>

            Вы уверены, что хотите удалить категорию?
          </NPopconfirm>
        </div>
      </li>
    </ul>

    <div
      v-if="!isLoading && categories.length === 0"
      class="text-sm text-gray-300"
    >
      Список категорий пуст.
    </div>
    <!--  </div>-->
  </AppLayout>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import HeaderMenu from '@/components/HeaderMenu.vue'
import { useFetch } from '@/composables/useFetch'
import { useNotifications } from '@/composables/useNotifications'
import { type Category } from '@/types'
import { NPopconfirm, NSpin } from 'naive-ui'
import AppLayout from '@/components/AppLayout.vue'

const fetcher = useFetch()
const notify = useNotifications()
const categories = ref<Category[]>([])
const isLoading = ref(true)
const deleting = ref<string | null>(null)

const handleDeleteCategory = (uuid: string) => {
  deleting.value = uuid

  fetcher.deleteCategory(uuid)
    .then(data => {
      if (data.ok) {
        categories.value = categories.value.filter(category => category.uuid !== uuid)
        notify.info('Категория удалена')
      }
    })
    .finally(() => {
      deleting.value = null
    })
}

const initPage = async () => {
  isLoading.value = true
  const data = await fetcher.getCategories()
  if (data.ok) {
    categories.value = data.data
  }
  isLoading.value = false
}

onMounted(() => {
  void initPage()
})
</script>
