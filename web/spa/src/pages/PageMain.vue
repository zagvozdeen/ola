<template>
  <AppLayout title="Товары и услуги">
    <ul class="grid grid-cols-2 gap-2">
      <li
        v-for="product in products"
        :key="product.id"
        class="bg-black/5 dark:bg-gray-500/20 border border-black/10 dark:border-gray-500/20 p-2 rounded-xl overflow-hidden flex flex-col"
      >
        <img
          class="h-40 w-full object-cover rounded-xl"
          :src="product.file_content"
          alt=""
        >
        <div class="my-2 flex flex-col gap-1 h-full">
          <span class="font-bold text-sm">{{ product.name }}</span>
          <div class="mt-auto">
            <span class="bg-blue-500/20 pl-1 pr-2 py-1 text-xs font-bold rounded-full inline-flex items-center gap-1">
              <span class="bg-blue-500 size-4 rounded-full text-center">₽</span>
              <span>от {{ product.price_from }} ₽{{ product.price_to ? ` до ${product.price_to} ₽` : '' }}</span>
            </span>
          </div>
        </div>
        <div class="flex justify-center">
          <button
            v-if="!isInCart(product.id)"
            class="w-full bg-grape-500 text-black rounded hover:bg-grape-600 px-4 py-1 mt-auto cursor-pointer text-xs uppercase font-bold text-center disabled:opacity-50"
            :disabled="isSubmitting(product.id)"
            @click="() => handleAddProductButton(product.id)"
          >
            Добавить
          </button>
          <button
            v-else
            class="w-full bg-gray-700 text-white dark:bg-gray-600 rounded hover:bg-gray-800 dark:hover:bg-gray-700 px-4 py-1 mt-auto cursor-pointer text-xs uppercase font-bold text-center disabled:opacity-50"
            :disabled="isSubmitting(product.id)"
            @click="() => handleRemoveProductButton(product.id, product.uuid)"
          >
            Убрать из корзины
          </button>
        </div>
      </li>
    </ul>

    <FooterMenu />
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import FooterMenu from '@/components/FooterMenu.vue'
import { cart } from '@/composables/useAuthState'
import { useFetch } from '@/composables/useFetch'
import { useNotifications } from '@/composables/useNotifications'
import type { Product } from '@/types'
import AppLayout from '@/components/AppLayout.vue'

const fetcher = useFetch()
const notify = useNotifications()
const products = ref<Product[]>([])
const submitting = ref<number | null>(null)

const cartProductIDs = computed(() => {
  return new Set(cart.items.map(item => item.product_id))
})

const isSubmitting = (productID: number) => {
  return submitting.value === productID
}

const isInCart = (productID: number) => {
  return cartProductIDs.value.has(productID)
}

const refreshCart = async () => {
  const data = await fetcher.getCart()

  if (data.ok) {
    cart.items = data.data
  }
}

const handleAddProductButton = async (productID: number) => {
  submitting.value = productID

  try {
    const data = await fetcher.upsertCartItem(productID, 1)

    if (!data.ok) {
      return
    }

    await refreshCart()
    notify.info('Товар добавлен в корзину')
  } finally {
    submitting.value = null
  }
}

const handleRemoveProductButton = async (productID: number, productUUID: string) => {
  submitting.value = productID

  try {
    const data = await fetcher.deleteCartItem(productUUID)

    if (!data.ok) {
      return
    }

    await refreshCart()
    notify.info('Товар убран из корзины')
  } finally {
    submitting.value = null
  }
}

onMounted(async () => {
  const [productsData] = await Promise.all([
    fetcher.getProducts(),
    refreshCart(),
  ])

  if (productsData.ok) {
    products.value = productsData.data
  }
})
</script>
