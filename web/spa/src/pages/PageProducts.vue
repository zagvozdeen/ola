<template>
  <div class="min-h-dvh w-full flex flex-col gap-4 py-6 pb-22">
    <HeaderMenu
      title="Все продукты"
      :edit="false"
      back="settings"
      create="products.create"
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
        v-for="product in products"
        :key="product.id"
        class="bg-gray-500/20 border border-gray-500/20 p-2 rounded-xl overflow-hidden flex gap-2"
      >
        <img
          class="size-20 object-cover rounded-lg"
          :src="product.file_content"
          alt=""
        >
        <div class="flex-1 min-w-0">
          <div class="flex justify-between gap-2">
            <span class="font-bold text-sm truncate">{{ product.name }}</span>
            <span class="text-xs uppercase bg-gray-600 font-bold px-2 py-0.5 rounded-full">{{ ProductTypeTranslates[product.type] }}</span>
          </div>
          <p class="text-xs text-gray-300 line-clamp-2 mt-1">
            {{ product.description }}
          </p>
          <p class="text-xs mt-1 font-medium">
            от {{ product.price_from }} ₽{{ product.price_to ? ` до ${product.price_to} ₽` : '' }}
          </p>
          <div class="flex gap-2 mt-2">
            <router-link
              class="bg-gray-600 hover:bg-gray-700 rounded px-2 py-1 text-xs font-bold"
              :to="{ name: 'products.edit', params: { uuid: product.uuid } }"
            >
              Редактировать
            </router-link>
            <NPopconfirm
              negative-text="Отмена"
              positive-text="Удалить"
              @positive-click="() => handleDeleteProduct(product.uuid)"
            >
              <template #trigger>
                <button
                  class="bg-red-700 hover:bg-red-800 rounded px-2 py-1 text-xs font-bold disabled:opacity-50 cursor-pointer"
                  :disabled="deleting === product.uuid"
                >
                  {{ deleting === product.uuid ? 'Удаляем...' : 'Удалить' }}
                </button>
              </template>

              Вы уверены, что хотите удалить продукт?
            </NPopconfirm>
          </div>
        </div>
      </li>
    </ul>

    <div
      v-if="!isLoading && products.length === 0"
      class="text-sm text-gray-300"
    >
      Список продуктов пуст.
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useFetch } from '@/composables/useFetch'
import { isUserModerator, onUserLoaded } from '@/composables/useState'
import { useNotifications } from '@/composables/useNotifications'
import { type Product, ProductTypeTranslates } from '@/types'
import { NPopconfirm, NSpin } from 'naive-ui'
import HeaderMenu from '@/components/HeaderMenu.vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const fetcher = useFetch()
const notify = useNotifications()
const products = ref<Product[]>([])
const isLoading = ref(true)
const deleting = ref<string | null>(null)

const handleDeleteProduct = (uuid: string) => {
  deleting.value = uuid

  fetcher
    .deleteProduct(uuid)
    .then(data => {
      if (data.ok) {
        products.value = products.value.filter(product => product.uuid !== uuid)
        notify.info('Продукт удалён')
      }
    })
    .finally(() => {
      deleting.value = null
    })
}

onUserLoaded((user) => {
  if (isUserModerator(user)) {
    fetcher
      .getProducts()
      .then(data => {
        if (data.ok) {
          products.value = data.data
        }
      })
      .finally(() => {
        isLoading.value = false
      })
  } else {
    notify.error('У вас нет прав просматривать эту страницу!')
    router.push({ name: 'main' })
  }
})
</script>
