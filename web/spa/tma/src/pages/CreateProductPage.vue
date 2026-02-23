<template>
  <div class="min-h-dvh w-full flex flex-col gap-4 py-6 pb-22">
    <h1 class="text-lg font-bold">
      Создать продукт
    </h1>

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
      Недостаточно прав для создания продукта.
    </div>

    <n-form
      v-else
      ref="formRef"
      class="w-full bg-gray-500/20 p-4 rounded-2xl"
      :rules="rules"
      :model="form"
      @submit.prevent="onSubmitForm"
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
          :options="typeOptions"
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

      <n-form-item label="Изображение">
        <input
          type="file"
          accept="image/*"
          class="w-full text-sm file:mr-3 file:rounded file:border-0 file:bg-blue-700 file:px-3 file:py-1 file:text-white file:cursor-pointer"
          @change="onFileChange"
        >
      </n-form-item>

      <p
        v-if="selectedFileName"
        class="text-xs text-gray-300 mb-4"
      >
        Выбран файл: {{ selectedFileName }}
      </p>

      <div class="flex gap-2">
        <router-link
          class="bg-gray-600 hover:bg-gray-700 rounded px-3 py-2 text-xs font-bold"
          :to="{ name: 'products' }"
        >
          Отмена
        </router-link>
        <n-button
          attr-type="submit"
          type="success"
          class="flex-1"
        >
          Создать продукт
        </n-button>
      </div>
    </n-form>

    <FooterMenu />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, useTemplateRef } from 'vue'
import { useRouter } from 'vue-router'
import { type FormRules, type FormInst, NButton, NForm, NFormItem, NInput, NInputNumber, NSelect } from 'naive-ui'
import FooterMenu from '@shared/components/FooterMenu.vue'
import { useFetch } from '@shared/composables/useFetch'
import { me } from '@shared/composables/useState'
import { useNotifications } from '@shared/composables/useNotifications'
import { useSender } from '@shared/composables/useSender'
import type { CreateProductRequest, ProductType } from '@shared/types'

type ProductForm = {
  name: string | null
  description: string | null
  price_from: number | null
  price_to: number | null
  type: ProductType | null
}

const router = useRouter()
const fetcher = useFetch()
const notify = useNotifications()
const sender = useSender()
const formRef = useTemplateRef<FormInst>('formRef')
const permissionChecked = ref(false)
const selectedFile = ref<globalThis.File | null>(null)
const selectedFileName = ref('')
const canManageProducts = computed(() => me.value?.role === 'admin' || me.value?.role === 'moderator')
const typeOptions = [
  { label: 'Товар', value: 'product' },
  { label: 'Услуга', value: 'service' },
]
const form = reactive<ProductForm>({
  name: null,
  description: null,
  price_from: null,
  price_to: null,
  type: 'product',
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
}

const ensureMe = async () => {
  if (me.value) {
    return
  }
  const data = await fetcher.getMe()
  if (data.ok) {
    me.value = data.data
  }
}

const onFileChange = (e: Event) => {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]

  if (!file) {
    selectedFile.value = null
    selectedFileName.value = ''
    return
  }

  selectedFile.value = file
  selectedFileName.value = file.name
}

const buildPayload = (fileID: number): CreateProductRequest | null => {
  if (form.name === null || form.description === null || form.price_from === null || form.type === null) {
    return null
  }

  return {
    name: form.name,
    description: form.description,
    price_from: form.price_from,
    price_to: form.price_to,
    type: form.type,
    file_id: fileID,
  }
}

const onSubmitForm = () => {
  sender.submit(formRef.value, async () => {
    if (!selectedFile.value) {
      notify.error('Выберите изображение продукта')
      return
    }

    const fileData = await fetcher.uploadFile(selectedFile.value)
    if (!fileData.ok) {
      notify.error(fileData.data.message)
      return
    }

    const payload = buildPayload(fileData.data.id)
    if (!payload) {
      notify.error('Форма заполнена некорректно')
      return
    }

    const data = await fetcher.createProduct(payload)
    if (!data.ok) {
      notify.error(data.data.message)
      return
    }

    notify.info('Продукт создан')
    await router.push({ name: 'products' })
  })
}

onMounted(async () => {
  await ensureMe()
  permissionChecked.value = true

  if (!canManageProducts.value) {
    return
  }
})
</script>
