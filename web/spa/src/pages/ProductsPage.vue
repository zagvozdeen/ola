<template>
  <div class="min-h-dvh w-full flex flex-col gap-4 py-6 pb-22">
    <div class="flex items-center justify-between">
      <h1 class="text-lg font-bold">
        Все продукты
      </h1>
      <router-link
        v-if="canManageProducts"
        class="bg-blue-600 hover:bg-blue-700 rounded px-3 py-1 text-xs uppercase font-bold"
        :to="{ name: 'product-create' }"
      >
        Создать
      </router-link>
    </div>

    <div
      v-if="!permissionChecked"
      class="text-sm text-gray-300"
    >
      Проверяем доступ...
    </div>

    <div
      v-else-if="!canManageProducts"
      class="text-sm text-red-300"
    >
      Недостаточно прав для управления продуктами.
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
            <span class="text-xs uppercase bg-gray-600 px-2 py-0.5 rounded-full">{{ product.type }}</span>
          </div>
          <p class="text-xs text-gray-300 line-clamp-2 mt-1">
            {{ product.description }}
          </p>
          <p class="text-xs mt-1">
            от {{ product.price_from }} ₽{{ product.price_to ? ` до ${product.price_to} ₽` : '' }}
          </p>
          <div class="flex gap-2 mt-2">
            <router-link
              class="bg-gray-600 hover:bg-gray-700 rounded px-2 py-1 text-xs font-bold"
              :to="{ name: 'product-edit', params: { uuid: product.uuid } }"
            >
              Редактировать
            </router-link>
            <button
              class="bg-red-700 hover:bg-red-800 rounded px-2 py-1 text-xs font-bold disabled:opacity-50"
              :disabled="deletingUuid === product.uuid"
              @click="() => handleDeleteProduct(product)"
            >
              {{ deletingUuid === product.uuid ? 'Удаляем...' : 'Удалить' }}
            </button>
          </div>
        </div>
      </li>
    </ul>

    <div
      v-if="canManageProducts && !isLoading && products.length === 0"
      class="text-sm text-gray-300"
    >
      Список продуктов пуст.
    </div>

    <FooterMenu />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import FooterMenu from '@/components/FooterMenu.vue'
import { useFetch } from '@/composables/useFetch'
import { me } from '@/composables/useState'
import { useNotifications } from '@/composables/useNotifications'
import type { Product } from '@/types'

const fetcher = useFetch()
const notify = useNotifications()
const products = ref<Product[]>([])
const isLoading = ref(false)
const deletingUuid = ref<string | null>(null)
const permissionChecked = ref(false)
const canManageProducts = computed(() => me.value?.role === 'admin' || me.value?.role === 'moderator')

const ensureMe = async () => {
  if (me.value) {
    return
  }
  const data = await fetcher.getMe()
  if (data.ok) {
    me.value = data.data
  }
}

const loadProducts = async () => {
  isLoading.value = true
  const data = await fetcher.getProducts()
  if (data.ok) {
    products.value = data.data
  } else {
    notify.error(data.data.message)
  }
  isLoading.value = false
}

const handleDeleteProduct = async (product: Product) => {
  const confirmed = window.confirm(`Удалить продукт "${product.name}"?`)
  if (!confirmed) {
    return
  }

  deletingUuid.value = product.uuid
  const data = await fetcher.deleteProduct(product.uuid)
  deletingUuid.value = null

  if (!data.ok) {
    notify.error(data.data.message)
    return
  }

  products.value = products.value.filter(_product => _product.uuid !== product.uuid)
  notify.info('Продукт удалён')
}

onMounted(async () => {
  await ensureMe()
  permissionChecked.value = true

  if (!canManageProducts.value) {
    return
  }

  await loadProducts()
})
</script>
