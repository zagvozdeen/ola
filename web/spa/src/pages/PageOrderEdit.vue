<template>
  <AppLayout
    title="Редактирование заказа"
    back="orders"
    @save="onSubmit"
  >
    <!--    <div class="min-h-dvh w-full flex flex-col gap-4 py-6">-->
    <!--      <HeaderMenu-->
    <!--        title="Редактирование заказа"-->
    <!--        :edit="true"-->
    <!--        back="orders"-->
    <!--        @ready="onSubmit"-->
    <!--      />-->

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
    </n-form>
    <!--  </div>-->
    <!--    </div>-->
  </AppLayout>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref, useTemplateRef } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import HeaderMenu from '@/components/HeaderMenu.vue'
import { useFetch } from '@/composables/useFetch'
import { useNotifications } from '@/composables/useNotifications'
import { useSender } from '@/composables/useSender'
import { type Order, RequestStatus, RequestStatusOptions, type UpdateRequestStatusRequest } from '@/types'
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
const form = reactive<UpdateRequestStatusRequest>({
  status: RequestStatus.Created,
})

const rules: FormRules = {
  status: {
    required: true,
    type: 'string',
    message: 'Выберите статус',
  },
}

const onSubmit = () => {
  sender.submit(formRef.value, async () => {
    const uuid = route.params['uuid']
    if (typeof uuid !== 'string' || !uuid) {
      notify.error('Некорректный ID заказа')
      await router.push({ name: 'orders' })
      return
    }

    const data = await fetcher.updateOrderStatus(uuid, { status: form.status })
    if (data.ok) {
      notify.info('Статус заказа обновлён')
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
      }
    })
    .finally(() => {
      isLoading.value = false
    })
})
</script>
