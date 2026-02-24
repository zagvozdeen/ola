<template>
  <div class="min-h-dvh w-full flex flex-col gap-4 py-6">
    <HeaderMenu
      :title="title"
      :edit="true"
      back="products"
      @ready="onSubmit"
    />

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
      <n-form-item
        label="Название"
        path="name"
      >
        <n-input
          v-model:value="form.name"
          placeholder="Название продукта"
        />
      </n-form-item>

      <n-form-item
        label="Описание"
        path="description"
      >
        <n-input
          v-model:value="form.description"
          type="textarea"
          placeholder="Описание продукта"
          :autosize="{ minRows: 3, maxRows: 6 }"
        />
      </n-form-item>

      <n-form-item
        label="Тип"
        path="type"
      >
        <n-select
          v-model:value="form.type"
          :options="ProductTypeOptions"
          placeholder="Выберите тип"
        />
      </n-form-item>

      <n-form-item
        label="Цена от"
        path="price_from"
      >
        <n-input-number
          v-model:value="form.price_from"
          class="w-full"
          placeholder="0"
          :show-button="false"
          :min="0"
        />
      </n-form-item>

      <n-form-item
        label="Цена до (необязательно)"
        path="price_to"
      >
        <n-input-number
          v-model:value="form.price_to"
          class="w-full"
          placeholder="0"
          :show-button="false"
          :min="0"
        />
      </n-form-item>

      <n-form-item
        label="Изображение продукта"
        path="file_id"
      >
        <AppUploadFile
          v-model:value="form.file_id"
          :content="form.file_content"
        />
      </n-form-item>
    </n-form>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref, useTemplateRef } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import HeaderMenu from '@/components/HeaderMenu.vue'
import AppUploadFile from '@/components/AppUploadFile.vue'
import { useFetch } from '@/composables/useFetch'
import { useNotifications } from '@/composables/useNotifications'
import { useSender } from '@/composables/useSender'
import { type CreateProductRequest, ProductTypeOptions } from '@/types'
import { type FormInst, NForm, NFormItem, NInput, NInputNumber, NSelect, NSpin, type FormRules } from 'naive-ui'

const route = useRoute()
const router = useRouter()
const fetcher = useFetch()
const notify = useNotifications()
const sender = useSender()

const isCreating = String(route.name).endsWith('create')
const title = isCreating ? 'Создание продукта' : 'Редактирование продукта'

type ProductForm = CreateProductRequest & {
  file_content: string | null
}

const formRef = useTemplateRef<FormInst>('formRef')
const isLoading = ref(true)
const form = reactive<ProductForm>({
  name: null,
  description: null,
  price_from: null,
  price_to: null,
  type: null,
  file_id: null,
  file_content: null,
})
const rules: FormRules = {
  name: {
    required: true,
    type: 'string',
    message: 'Введите название',
    min: 1,
    max: 255,
  },
  description: {
    required: true,
    type: 'string',
    message: 'Введите описание',
    min: 1,
    max: 3000,
  },
  type: {
    required: true,
    type: 'string',
    message: 'Выберите тип',
  },
  price_from: {
    required: true,
    type: 'number',
    message: 'Введите цену "от"',
    min: 0,
  },
  price_to: {
    type: 'number',
    required: false,
    min: 0,
    validator: (_rule, value: number | null) => {
      if (value === null || value === undefined || form.price_from === null) {
        return true
      }
      if (value < form.price_from) {
        return new Error('Цена "до" должна быть не меньше цены "от"')
      }
      return true
    },
  },
  file_id: {
    required: true,
    type: 'number',
    message: 'Выберите изображение',
    min: 1,
  },
}

const onSubmit = () => {
  sender.submit(formRef.value, async () => {
    const payload: CreateProductRequest = {
      name: form.name,
      description: form.description,
      price_from: form.price_from,
      price_to: form.price_to,
      type: form.type,
      file_id: form.file_id,
    }

    if (isCreating) {
      const data = await fetcher.createProduct(payload)
      if (data.ok) {
        notify.info('Продукт создан')
        await router.push({ name: 'products' })
      }
      return
    }

    const uuid = route.params['uuid']
    if (typeof uuid !== 'string' || !uuid) {
      notify.error('Некорректный ID продукта')
      await router.push({ name: 'products' })
      return
    }

    const data = await fetcher.updateProduct(uuid, payload)
    if (data.ok) {
      notify.info('Продукт обновлён')
      await router.push({ name: 'products' })
    }
  })
}

onMounted(() => {
  if (isCreating) {
    isLoading.value = false
    return
  }

  const uuid = route.params['uuid']
  if (typeof uuid !== 'string' || !uuid) {
    notify.error('Некорректный ID продукта')
    router.push({ name: 'products' })
    isLoading.value = false
    return
  }

  fetcher
    .getProduct(uuid)
    .then(data => {
      if (data.ok) {
        form.name = data.data.name
        form.description = data.data.description
        form.price_from = data.data.price_from
        form.price_to = data.data.price_to ?? null
        form.type = data.data.type
        form.file_id = data.data.file_id
        form.file_content = data.data.file_content
      }
    })
    .finally(() => {
      isLoading.value = false
    })
})
</script>
