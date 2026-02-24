<template>
  <div class="min-h-dvh w-full flex flex-col gap-4 py-6">
    <HeaderMenu
      :title="title"
      :edit="true"
      back="reviews"
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
        label="Имя"
        path="name"
      >
        <n-input
          v-model:value="form.name"
          placeholder="Имя автора"
        />
      </n-form-item>

      <n-form-item
        label="Текст отзыва"
        path="content"
      >
        <n-input
          v-model:value="form.content"
          type="textarea"
          placeholder="Текст отзыва"
          :autosize="{ minRows: 3, maxRows: 6 }"
        />
      </n-form-item>

      <n-form-item
        label="Дата публикации"
        path="published_at"
      >
        <input
          v-model="publishedAtLocal"
          type="datetime-local"
          class="w-full rounded-md bg-transparent border border-gray-500/40 px-3 py-2 text-sm"
        >
      </n-form-item>

      <n-form-item
        label="Изображение"
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
import { computed, onMounted, reactive, ref, useTemplateRef } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import HeaderMenu from '@/components/HeaderMenu.vue'
import AppUploadFile from '@/components/AppUploadFile.vue'
import { useFetch } from '@/composables/useFetch'
import { useNotifications } from '@/composables/useNotifications'
import { useSender } from '@/composables/useSender'
import { type FormInst, type FormRules, NForm, NFormItem, NInput, NSpin } from 'naive-ui'
import type { UpsertReviewRequest } from '@/types'

const route = useRoute()
const router = useRouter()
const fetcher = useFetch()
const notify = useNotifications()
const sender = useSender()

const isCreating = String(route.name).endsWith('create')
const title = isCreating ? 'Создание отзыва' : 'Редактирование отзыва'

type ReviewForm = UpsertReviewRequest & {
  file_content: string | null
}

const formRef = useTemplateRef<FormInst>('formRef')
const isLoading = ref(true)
const form = reactive<ReviewForm>({
  name: null,
  content: null,
  file_id: null,
  published_at: null,
  file_content: null,
})

const publishedAtLocal = computed({
  get: () => {
    if (!form.published_at) {
      return ''
    }
    const date = new Date(form.published_at)
    const offset = date.getTimezoneOffset()
    const local = new Date(date.getTime() - offset * 60 * 1000)
    return local.toISOString().slice(0, 16)
  },
  set: (value: string) => {
    form.published_at = value ? new Date(value).toISOString() : null
  },
})

const rules: FormRules = {
  name: {
    required: true,
    type: 'string',
    message: 'Введите имя',
    min: 1,
    max: 255,
  },
  content: {
    required: true,
    type: 'string',
    message: 'Введите текст',
    min: 1,
    max: 3000,
  },
  file_id: {
    required: true,
    type: 'number',
    message: 'Выберите изображение',
    min: 1,
  },
  published_at: {
    required: true,
    type: 'string',
    message: 'Выберите дату публикации',
    min: 1,
  },
}

const onSubmit = () => {
  sender.submit(formRef.value, async () => {
    const payload: UpsertReviewRequest = {
      name: form.name,
      content: form.content,
      file_id: form.file_id,
      published_at: form.published_at,
    }

    if (isCreating) {
      const data = await fetcher.createReview(payload)
      if (data.ok) {
        notify.info('Отзыв создан')
        await router.push({ name: 'reviews' })
      }
      return
    }

    const uuid = route.params['uuid']
    if (typeof uuid !== 'string' || !uuid) {
      notify.error('Некорректный ID отзыва')
      await router.push({ name: 'reviews' })
      return
    }

    const data = await fetcher.updateReview(uuid, payload)
    if (data.ok) {
      notify.info('Отзыв обновлён')
      await router.push({ name: 'reviews' })
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
    notify.error('Некорректный ID отзыва')
    router.push({ name: 'reviews' })
    isLoading.value = false
    return
  }

  fetcher.getReview(uuid)
    .then(data => {
      if (data.ok) {
        form.name = data.data.name
        form.content = data.data.content
        form.file_id = data.data.file_id
        form.file_content = data.data.file_content ?? null
        form.published_at = data.data.published_at
      }
    })
    .finally(() => {
      isLoading.value = false
    })
})
</script>
