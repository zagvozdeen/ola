<template>
  <AppLayout
    title="Заказы"
    back="settings"
    no-save
  >
    <!--  <div class="min-h-dvh w-full flex flex-col gap-4 py-6 pb-22">-->
    <!--    <HeaderMenu-->
    <!--      title="Заказы"-->
    <!--      :edit="false"-->
    <!--      back="settings"-->
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
        v-for="order in orders"
        :key="order.id"
        class="bg-gray-500/20 border border-gray-500/20 p-3 rounded-xl overflow-hidden flex flex-col gap-2"
      >
        <div class="flex justify-between gap-2">
          <span class="font-bold text-sm truncate">{{ order.name }}</span>
          <span
            class="text-xs uppercase font-bold px-2 py-0.5 rounded-full"
            :class="{
              [RequestStatusBgColor[order.status]]:true,
            }"
          >{{ RequestStatusTranslates[order.status] }}</span>
        </div>
        <p class="text-xs text-gray-300 line-clamp-2">
          {{ order.content }}
        </p>
        <p class="text-xs text-gray-300">
          {{ order.phone }}
        </p>
        <div class="flex justify-end">
          <router-link
            class="bg-gray-600 hover:bg-gray-700 rounded px-2 py-1 text-xs font-bold"
            :to="{ name: 'orders.edit', params: { uuid: order.uuid } }"
          >
            Редактировать
          </router-link>
        </div>
      </li>
    </ul>

    <div
      v-if="!isLoading && orders.length === 0"
      class="text-sm text-gray-300"
    >
      Список заказов пуст.
    </div>
    <!--  </div>-->
  </AppLayout>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import HeaderMenu from '@/components/HeaderMenu.vue'
import { useFetch } from '@/composables/useFetch'
import { type Order, RequestStatusBgColor, RequestStatusTranslates } from '@/types'
import { NSpin } from 'naive-ui'
import AppLayout from '@/components/AppLayout.vue'

const fetcher = useFetch()
const orders = ref<Order[]>([])
const isLoading = ref(true)

const initPage = async () => {
  isLoading.value = true
  const data = await fetcher.getOrders()
  if (data.ok) {
    orders.value = data.data
  }
  isLoading.value = false
}

onMounted(() => {
  void initPage()
})
</script>
