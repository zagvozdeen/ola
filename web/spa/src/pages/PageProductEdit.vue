<template>
  <AppLayout
    :title="title"
    back="products"
    @save="onSubmit"
  >
    <!--    <div class="min-h-dvh w-full flex flex-col gap-4 py-6">-->
    <!--      <HeaderMenu-->
    <!--        :title="title"-->
    <!--        :edit="true"-->
    <!--        back="products"-->
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
        path="file_content"
      >
        <AppUploadFile
          v-model:value="form.file_content"
          :content="form.file_content"
        />
      </n-form-item>

      <n-form-item label="Категории">
        <n-checkbox-group v-model:value="form.category_uuids">
          <n-space vertical>
            <n-checkbox
              v-for="category in categories"
              :key="category.uuid"
              :value="category.uuid"
              :label="category.name"
            />
          </n-space>
        </n-checkbox-group>
      </n-form-item>
    </n-form>
    <!--    </div>-->
  </AppLayout>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref, useTemplateRef } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import AppUploadFile from '@/components/AppUploadFile.vue'
import { useFetch } from '@/composables/useFetch'
import { useNotifications } from '@/composables/useNotifications'
import { useSender } from '@/composables/useSender'
import { type Category, type CreateProductRequest, ProductTypeOptions } from '@/types'
import { NCheckbox, NCheckboxGroup, type FormInst, NForm, NFormItem, NInput, NInputNumber, NSelect, NSpace, NSpin, type FormRules } from 'naive-ui'
import AppLayout from '@/components/AppLayout.vue'

const route = useRoute()
const router = useRouter()
const fetcher = useFetch()
const notify = useNotifications()
const sender = useSender()

const isCreating = String(route.name).endsWith('create')
const title = isCreating ? 'Создание продукта' : 'Редактирование продукта'

const formRef = useTemplateRef<FormInst>('formRef')
const isLoading = ref(true)
const categories = ref<Category[]>([])
const form = reactive<CreateProductRequest>({
  name: null,
  description: null,
  price_from: null,
  price_to: null,
  type: null,
  file_content: null,
  category_uuids: [],
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
  file_content: {
    required: true,
    type: 'string',
    message: 'Выберите изображение',
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
      file_content: form.file_content,
      category_uuids: form.category_uuids,
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
  const init = async () => {
    const categoriesData = await fetcher.getCategories()
    if (categoriesData.ok) {
      categories.value = categoriesData.data
    }

    if (isCreating) {
      isLoading.value = false
      return
    }

    const uuid = route.params['uuid']
    if (typeof uuid !== 'string' || !uuid) {
      notify.error('Некорректный ID продукта')
      await router.push({ name: 'products' })
      isLoading.value = false
      return
    }

    const data = await fetcher.getProduct(uuid)
    if (data.ok) {
      form.name = data.data.name
      form.description = data.data.description
      form.price_from = data.data.price_from
      form.price_to = data.data.price_to ?? null
      form.type = data.data.type
      form.file_content = data.data.file_content
      form.category_uuids = data.data.categories.map(category => category.uuid)
    }

    isLoading.value = false
  }

  void init()
})
</script>
