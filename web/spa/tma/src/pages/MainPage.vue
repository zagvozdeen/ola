<template>
  <div class="min-h-dvh w-full flex flex-col gap-4 items-center justify-center py-6 pb-22">
    <ul class="grid grid-cols-2 gap-2">
      <li
        v-for="product in products"
        :key="product.id"
        class="bg-gray-500/20 border border-gray-500/20 p-2 rounded-xl overflow-hidden flex flex-col"
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
            v-if="!cart.product_ids.includes(product.id)"
            class="w-full bg-green-700 rounded hover:bg-green-800 px-4 py-1 mt-auto cursor-pointer text-xs uppercase font-bold text-center"
            @click="() => handleAddProductButton(product.id)"
          >
            Добавить
          </button>
          <button
            v-else
            class="w-full bg-gray-600 rounded hover:bg-gray-700 px-4 py-1 mt-auto cursor-pointer text-xs uppercase font-bold text-center"
            @click="() => handleRemoveProductButton(product.id)"
          >
            Убрать из корзины
          </button>
        </div>
      </li>
      <li
        v-for="service in services"
        :key="service.id"
        class="bg-gray-500/20 border border-gray-500/20 p-2 rounded-xl overflow-hidden flex flex-col"
      >
        <img
          class="h-40 w-full object-cover rounded-xl"
          :src="service.file_content"
          alt=""
        >
        <div class="my-2 flex flex-col gap-1 h-full">
          <span class="font-bold text-sm">{{ service.name }}</span>
          <div class="mt-auto">
            <span class="bg-blue-500/20 pl-1 pr-2 py-1 text-xs font-bold rounded-full inline-flex items-center gap-1">
              <span class="bg-blue-500 size-4 rounded-full text-center">₽</span>
              <span>от {{ service.price_from }} ₽{{ service.price_to ? ` до ${service.price_to} ₽` : '' }}</span>
            </span>
          </div>
        </div>
        <div class="flex justify-center">
          <button
            v-if="!cart.product_ids.includes(service.id)"
            class="w-full bg-green-700 rounded hover:bg-green-800 px-4 py-1 mt-auto cursor-pointer text-xs uppercase font-bold text-center"
            @click="() => handleAddProductButton(service.id)"
          >
            Добавить
          </button>
          <button
            v-else
            class="w-full bg-gray-600 rounded hover:bg-gray-700 px-4 py-1 mt-auto cursor-pointer text-xs uppercase font-bold text-center"
            @click="() => handleRemoveProductButton(service.id)"
          >
            Убрать из корзины
          </button>
        </div>
      </li>
    </ul>

    <FooterMenu />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import FooterMenu from '@shared/components/FooterMenu.vue'
import { useFetch } from '@shared/composables/useFetch'
import type { Product } from '@shared/types'
import { cart } from '@shared/composables/useState'

const fetcher = useFetch()
const allProducts = ref<Product[]>([])
const products = computed(() => allProducts.value.filter(product => product.type === 'product'))
const services = computed(() => allProducts.value.filter(product => product.type === 'service'))

const handleAddProductButton = (id: number) => {
  cart.product_ids.push(id)
}

const handleRemoveProductButton = (id: number) => {
  cart.product_ids = cart.product_ids.filter(_id => _id !== id)
}

onMounted(() => {
  fetcher
    .getProducts()
    .then(data => {
      if (data.ok) {
        allProducts.value = data.data
      }
    })
})
</script>
