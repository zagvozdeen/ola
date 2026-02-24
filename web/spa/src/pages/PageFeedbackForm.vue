<template>
  <div class="min-h-dvh w-full flex flex-col gap-4 py-6">
    <HeaderMenu
      :title="title"
      :edit="true"
      back="settings"
      @ready="onSubmit"
    />

    <n-form
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
          placeholder="Введите имя"
        />
      </n-form-item>

      <n-form-item
        label="Телефон"
        path="phone"
      >
        <n-input
          v-model:value="form.phone"
          v-maska
          class="w-full"
          placeholder="+7 (___) ___-__-__"
          :input-props="phoneInputProps"
        />
      </n-form-item>

      <n-form-item
        label="Комментарий"
        path="content"
      >
        <n-input
          v-model:value="form.content"
          type="textarea"
          placeholder="Введите сообщение"
          :autosize="{ minRows: 3, maxRows: 6 }"
        />
      </n-form-item>
    </n-form>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, type InputHTMLAttributes, useTemplateRef } from 'vue'
import { useRoute } from 'vue-router'
import HeaderMenu from '@/components/HeaderMenu.vue'
import { useAuthState } from '@/composables/useAuthState'
import { useFetch } from '@/composables/useFetch'
import { useNotifications } from '@/composables/useNotifications'
import { useSender } from '@/composables/useSender'
import { type CreateFeedbackRequest, FeedbackType } from '@/types'
import { type FormInst, NForm, NFormItem, NInput, type FormRules } from 'naive-ui'
import { vMaska } from 'maska/vue'

const route = useRoute()
const auth = useAuthState()
const fetcher = useFetch()
const notify = useNotifications()
const sender = useSender()

const typeByRouteName: Record<string, FeedbackType> = {
  'settings.manager': FeedbackType.ManagerContact,
  'settings.feedback': FeedbackType.FeedbackRequest,
  'settings.partnership': FeedbackType.PartnershipOffer,
}

const titleByType: Record<FeedbackType, string> = {
  [FeedbackType.ManagerContact]: 'Связаться с менеджером',
  [FeedbackType.FeedbackRequest]: 'Оставить обратную связь',
  [FeedbackType.PartnershipOffer]: 'Предложить сотрудничество',
}

const successByType: Record<FeedbackType, string> = {
  [FeedbackType.ManagerContact]: 'Заявка отправлена менеджеру',
  [FeedbackType.FeedbackRequest]: 'Обратная связь отправлена',
  [FeedbackType.PartnershipOffer]: 'Предложение отправлено',
}

const feedbackType = computed(() => {
  return typeByRouteName[String(route.name)] ?? FeedbackType.FeedbackRequest
})

const title = computed(() => {
  return titleByType[feedbackType.value]
})

const formRef = useTemplateRef<FormInst>('formRef')
const form = reactive<CreateFeedbackRequest>({
  name: '',
  phone: '',
  content: '',
  type: FeedbackType.FeedbackRequest,
})

const phoneInputProps = {
  'data-maska': '+7 (###) ###-##-##',
} as unknown as InputHTMLAttributes

const rules: FormRules = {
  name: {
    required: true,
    type: 'string',
    message: 'Введите имя',
    min: 1,
    max: 255,
  },
  phone: {
    required: true,
    type: 'string',
    message: 'Введите телефон',
    min: 1,
    max: 255,
  },
  content: {
    required: true,
    type: 'string',
    message: 'Введите сообщение',
    min: 1,
    max: 3000,
  },
}

const onSubmit = () => {
  sender.submit(formRef.value, async () => {
    form.type = feedbackType.value

    const data = await fetcher.createFeedback({
      name: form.name,
      phone: form.phone,
      content: form.content,
      type: form.type,
    })

    if (!data.ok) {
      return
    }

    notify.info(successByType[form.type])
    form.name = ''
    form.phone = ''
    form.content = ''
  })
}

onMounted(async () => {
  await auth.ensureUserLoaded()

  if (auth.currentUser.value) {
    form.name = auth.currentUser.value.first_name
    form.phone = auth.currentUser.value.phone ?? ''
  }
})
</script>
