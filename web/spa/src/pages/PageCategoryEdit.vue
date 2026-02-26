<template>
  <AppLayout
    :title="title"
    back="categories"
    @save="onSubmit"
  >
    <!--  <div class="min-h-dvh w-full flex flex-col gap-4 py-6">-->
    <!--    <HeaderMenu-->
    <!--      :title="title"-->
    <!--      :edit="true"-->
    <!--      back="categories"-->
    <!--      @ready="onSubmit"-->
    <!--    />-->

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
          placeholder="Название категории"
        />
      </n-form-item>
    </n-form>
    <!--  </div>-->
  </AppLayout>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref, useTemplateRef } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useFetch } from '@/composables/useFetch'
import { useNotifications } from '@/composables/useNotifications'
import { useSender } from '@/composables/useSender'
import { type FormInst, type FormRules, NForm, NFormItem, NInput, NSpin } from 'naive-ui'
import { type UpsertCategoryRequest } from '@/types'
import AppLayout from '@/components/AppLayout.vue'

const route = useRoute()
const router = useRouter()
const fetcher = useFetch()
const notify = useNotifications()
const sender = useSender()

const isCreating = String(route.name).endsWith('create')
const title = isCreating ? 'Создание категории' : 'Редактирование категории'

const formRef = useTemplateRef<FormInst>('formRef')
const isLoading = ref(true)
const form = reactive<UpsertCategoryRequest>({
  name: null,
})

const rules: FormRules = {
  name: {
    required: true,
    type: 'string',
    message: 'Введите название',
    min: 1,
    max: 255,
  },
}

const onSubmit = () => {
  sender.submit(formRef.value, async () => {
    if (isCreating) {
      const data = await fetcher.createCategory({ name: form.name })
      if (data.ok) {
        notify.info('Категория создана')
        await router.push({ name: 'categories' })
      }
      return
    }

    const uuid = route.params['uuid']
    if (typeof uuid !== 'string' || !uuid) {
      notify.error('Некорректный ID категории')
      await router.push({ name: 'categories' })
      return
    }

    const data = await fetcher.updateCategory(uuid, { name: form.name })
    if (data.ok) {
      notify.info('Категория обновлена')
      await router.push({ name: 'categories' })
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
    notify.error('Некорректный ID категории')
    router.push({ name: 'categories' })
    isLoading.value = false
    return
  }

  fetcher.getCategory(uuid)
    .then(data => {
      if (data.ok) {
        form.name = data.data.name
      }
    })
    .finally(() => {
      isLoading.value = false
    })
})
</script>
