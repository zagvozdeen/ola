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
      class="grid grid-cols-1 gap-3"
    >
      <li
        v-for="order in orders"
        :key="order.id"
        class="bg-gray-500/20 border border-gray-500/20 p-3 rounded-2xl overflow-hidden flex flex-col gap-3"
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

        <div
          v-if="(order.items ?? []).length > 0"
          class="rounded-xl border border-gray-500/20 bg-gray-500/10 p-3"
        >
          <p class="text-[11px] font-semibold uppercase tracking-wide text-gray-300">
            Состав заказа
          </p>
          <ul class="mt-3 flex flex-col gap-2">
            <li
              v-for="item in order.items ?? []"
              :key="`${order.uuid}-${item.product_id}`"
              class="grid grid-cols-[3rem_1fr] gap-3 rounded-xl bg-gray-500/10 p-2"
            >
              <img
                v-if="item.file_content"
                :src="item.file_content"
                :alt="item.product_name"
                class="size-12 rounded-lg object-cover"
              >
              <div
                v-else
                class="size-12 rounded-lg bg-gray-500/20"
              />
              <div class="min-w-0">
                <p class="truncate text-sm font-medium">
                  {{ item.product_name }}
                </p>
                <p class="text-xs text-gray-300">
                  {{ formatPrice(item) }} x {{ item.qty }}
                </p>
              </div>
            </li>
          </ul>
        </div>

        <div
          v-if="(order.comments ?? []).length > 0"
          class="rounded-xl border border-gray-500/20 bg-gray-500/10 p-3"
        >
          <p class="text-[11px] font-semibold uppercase tracking-wide text-gray-300">
            Комментарии сотрудников
          </p>
          <ul class="mt-3 flex flex-col gap-2">
            <li
              v-for="comment in order.comments ?? []"
              :key="comment.uuid"
              class="rounded-xl bg-gray-500/10 p-3"
            >
              <div class="flex flex-wrap items-center gap-2 text-xs text-gray-300">
                <span class="font-semibold text-white">{{ formatAuthor(comment) }}</span>
                <span>{{ formatDate(comment.created_at) }}</span>
              </div>
              <p class="mt-2 whitespace-pre-line text-sm">
                {{ comment.content }}
              </p>
            </li>
          </ul>
        </div>

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
import { useFetch } from '@/composables/useFetch'
import { type Order, type OrderComment, type OrderItem, RequestStatusBgColor, RequestStatusTranslates } from '@/types'
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

const formatAuthor = (comment: OrderComment) => {
  const author = comment.author

  if (!author) {
    return 'Сотрудник'
  }

  return [author.first_name, author.last_name].filter(Boolean).join(' ') || author.username || 'Сотрудник'
}

const formatDate = (value: string) => {
  return new Date(value).toLocaleString('ru-RU')
}

const formatPrice = (item: OrderItem) => {
  if (typeof item.price_to === 'number') {
    return `${item.price_from}-${item.price_to} RUB`
  }

  return `${item.price_from} RUB`
}
</script>
