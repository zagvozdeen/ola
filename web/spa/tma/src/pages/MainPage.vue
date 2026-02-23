<template>
  <div class="min-h-dvh w-full flex flex-col gap-4 items-center justify-center py-6">
    <ul class="grid grid-cols-2 gap-4 mt-8 mb-16">
      <li
        v-for="product in products"
        :key="product.id"
        class="bg-gray-500/20 border border-gray-500/20 rounded-xl overflow-hidden flex flex-col"
      >
        <img
          class="h-50 w-full object-cover"
          :src="product.file_content"
          alt=""
        >
        <div class="p-2 flex flex-col gap-1 h-full">
          <span class="text-xs font-semibold">{{ product.name }}</span>
          <span class="font-bold text-xs mt-auto">от {{ product.price_from }} ₽</span>
        </div>
        <button class="w-full bg-green-700 hover:bg-green-800 px-2 py-1 mt-auto cursor-pointer text-xs uppercase font-bold text-center">
          Добавить
        </button>
      </li>
    </ul>

    <FooterMenu />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import FooterMenu from '@shared/components/FooterMenu.vue'
import { useFetch } from '@shared/composables/useFetch'
import type { Product, Service } from '@shared/types'

const fetcher = useFetch()
const products = ref<Product[]>([])
const services = ref<Service[]>([])

onMounted(() => {
  fetcher
    .getProducts()
    .then(data => {
      if (data.ok) {
        products.value = data.data
      }
    })

  fetcher
    .getServices()
    .then(data => {
      if (data.ok) {
        services.value = data.data
      }
    })
})
</script>
