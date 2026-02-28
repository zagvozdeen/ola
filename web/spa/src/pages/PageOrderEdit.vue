<template>
  <AppLayout
    title="Редактирование заказа"
    back="orders"
    @save="onSubmit"
  >
    <div
      v-if="isLoading"
      class="flex justify-center my-4"
    >
      <n-spin size="small" />
    </div>

    <n-form
      v-else
      ref="formRef"
      class="w-full bg-gray-500/20 p-4 rounded-2xl"
      :rules="rules"
      :model="form"
      @submit.prevent="onSubmit"
    >
      <div
        v-if="orderItems.length > 0"
        class="mb-4 rounded-2xl border border-gray-500/20 bg-gray-500/10 p-3"
      >
        <p class="text-xs font-semibold uppercase tracking-wide text-gray-300">
          Состав заказа
        </p>
        <ul class="mt-3 flex flex-col gap-2">
          <li
            v-for="item in orderItems"
            :key="`${item.order_id}-${item.product_id}`"
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

      <n-form-item label="Имя">
        <n-input
          :value="order?.name ?? ''"
          readonly
        />
      </n-form-item>

      <n-form-item label="Телефон">
        <n-input
          :value="order?.phone ?? ''"
          readonly
        />
      </n-form-item>

      <n-form-item label="Комментарий">
        <n-input
          :value="order?.content ?? ''"
          type="textarea"
          readonly
          :autosize="{ minRows: 3, maxRows: 6 }"
        />
      </n-form-item>

      <n-form-item
        label="Статус"
        path="status"
      >
        <n-select
          v-model:value="form.status"
          :options="RequestStatusOptions"
          placeholder="Выберите статус"
        />
      </n-form-item>

      <n-form-item
        label="Комментарий сотрудника"
        path="comment"
      >
        <n-input
          v-model:value="form.comment"
          type="textarea"
          placeholder="Напишите внутренний комментарий"
          :autosize="{ minRows: 3, maxRows: 6 }"
        />
      </n-form-item>

      <div
        v-if="orderComments.length > 0"
        class="rounded-2xl border border-gray-500/20 bg-gray-500/10 p-3"
      >
        <p class="text-xs font-semibold uppercase tracking-wide text-gray-300">
          История комментариев
        </p>
        <ul class="mt-3 flex flex-col gap-3">
          <li
            v-for="comment in orderComments"
            :key="comment.uuid"
            class="rounded-xl border border-gray-500/20 bg-gray-500/10 p-3"
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
    </n-form>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, useTemplateRef } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useFetch } from '@/composables/useFetch'
import { useNotifications } from '@/composables/useNotifications'
import { useSender } from '@/composables/useSender'
import { type Order, type OrderComment, type OrderItem, RequestStatus, RequestStatusOptions, type UpdateOrderStatusRequest } from '@/types'
import { type FormInst, type FormRules, NForm, NFormItem, NInput, NSelect, NSpin } from 'naive-ui'
import AppLayout from '@/components/AppLayout.vue'

const route = useRoute()
const router = useRouter()
const fetcher = useFetch()
const notify = useNotifications()
const sender = useSender()

const formRef = useTemplateRef<FormInst>('formRef')
const isLoading = ref(true)
const order = ref<Order | null>(null)
const form = reactive<UpdateOrderStatusRequest>({
  status: RequestStatus.Created,
  comment: '',
})

const orderItems = computed<OrderItem[]>(() => order.value?.items ?? [])
const orderComments = computed<OrderComment[]>(() => order.value?.comments ?? [])

const rules: FormRules = {
  status: {
    required: true,
    type: 'string',
    message: 'Выберите статус',
  },
  comment: {
    required: true,
    message: 'Добавьте комментарий сотрудника',
  },
}

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

const onSubmit = () => {
  sender.submit(formRef.value, async () => {
    const uuid = route.params['uuid']
    if (typeof uuid !== 'string' || !uuid) {
      notify.error('Некорректный ID заказа')
      await router.push({ name: 'orders' })
      return
    }

    const data = await fetcher.updateOrderStatus(uuid, {
      status: form.status,
      comment: form.comment,
    })
    if (data.ok) {
      notify.info('Заказ обновлён')
      await router.push({ name: 'orders' })
    }
  })
}

onMounted(() => {
  const uuid = route.params['uuid']
  if (typeof uuid !== 'string' || !uuid) {
    notify.error('Некорректный ID заказа')
    router.push({ name: 'orders' })
    isLoading.value = false
    return
  }

  fetcher.getOrder(uuid)
    .then(data => {
      if (data.ok) {
        order.value = data.data
        form.status = data.data.status
        form.comment = ''
      }
    })
    .finally(() => {
      isLoading.value = false
    })
})
</script>
