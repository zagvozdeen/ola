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
        :key="category.slug"
        class="bg-gray-500/20 border border-gray-500/20 p-3 rounded-xl overflow-hidden flex justify-between items-center gap-2"
      >
        <div class="min-w-0">
          <span class="block font-bold text-sm truncate">{{ category.name }}</span>
          <span class="block text-[10px] text-gray-300 truncate">{{ category.slug }}</span>
        </div>
        <div class="flex gap-2 shrink-0">
          <router-link
            class="bg-gray-600 hover:bg-gray-700 rounded px-2 py-1 text-xs font-bold"
            :to="{ name: 'categories.edit', params: { slug: category.slug } }"
          >
            Редактировать
          </router-link>
          <NPopconfirm
            negative-text="Отмена"
            positive-text="Удалить"
            @positive-click="() => handleDeleteCategory(category.slug)"
          >
            <template #trigger>
              <button
                class="bg-red-700 hover:bg-red-800 rounded px-2 py-1 text-xs font-bold disabled:opacity-50 cursor-pointer"
                :disabled="deleting === category.slug"
              >
                {{ deleting === category.slug ? 'Удаляем...' : 'Удалить' }}
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

const handleDeleteCategory = (slug: string) => {
  deleting.value = slug

  fetcher.deleteCategory(slug)
    .then(data => {
      if (data.ok) {
        categories.value = categories.value.filter(category => category.slug !== slug)
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
