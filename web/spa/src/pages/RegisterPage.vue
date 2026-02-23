<template>
  <div class="min-h-dvh w-full flex items-center justify-center">
    <n-form
      ref="formRef"
      class="w-full bg-gray-500/20 p-4 my-8 rounded-2xl"
      :rules="rules"
      :model="form"
      @submit.prevent="onSubmitForm"
    >
      <n-form-item
        :show-feedback="false"
        :show-label="false"
      >
        <h1 class="text-xl mb-4 text-center w-full">
          Создать аккаунт
        </h1>
      </n-form-item>
      <n-form-item
        label="Имя"
        path="first_name"
      >
        <n-input
          v-model:value="form.first_name"
          placeholder="Введите имя"
          autofocus
        />
      </n-form-item>
      <n-form-item
        label="Фамилия"
        path="last_name"
      >
        <n-input
          v-model:value="form.last_name"
          placeholder="Введите фамилию"
        />
      </n-form-item>
      <n-form-item
        label="Электронная почта"
        path="email"
      >
        <n-input
          v-model:value="form.email"
          placeholder="Введите электронную почту"
        />
      </n-form-item>
      <n-form-item
        label="Пароль"
        path="password"
      >
        <n-input
          v-model:value="form.password"
          type="password"
          placeholder="Введите пароль"
        />
      </n-form-item>
      <n-form-item
        label="Подтвердите пароль"
        path="password_confirmation"
      >
        <n-input
          v-model:value="form.password_confirmation"
          type="password"
          placeholder="Повторите пароль"
        />
      </n-form-item>
      <n-form-item
        :show-feedback="false"
        :show-label="false"
      >
        <n-button
          attr-type="submit"
          type="success"
          class="flex-1"
        >
          Зарегистрироваться
        </n-button>
      </n-form-item>


      <p class="text-center mt-4">
        Уже есть аккаунт? <router-link
          :to="{ name: 'login' }"
          class="underline"
        >
          Войти в аккаунт
        </router-link>
      </p>
    </n-form>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { reactive, useTemplateRef } from 'vue'
import { useFetch } from '../composables/useFetch'
import type { AuthRegisterRequest } from '../types'
import { type FormRules, type FormInst, NForm, NFormItem, NButton, NInput } from 'naive-ui'
import { useSender } from '../composables/useSender'
import { useNotifications } from '../composables/useNotifications'

const router = useRouter()
const fetcher = useFetch()
const sender = useSender()
const notify = useNotifications()
const formRef = useTemplateRef<FormInst>('formRef')
const form = reactive<AuthRegisterRequest>({
  first_name: null,
  last_name: null,
  email: null,
  password: null,
  password_confirmation: null,
})
const rules: FormRules = {
  first_name: {
    required: true,
    type: 'string',
    message: 'Введите имя',
    min: 1,
    max: 255,
  },
  last_name: {
    required: true,
    type: 'string',
    message: 'Введите фамилию',
    min: 1,
    max: 255,
  },
  email: {
    required: true,
    type: 'email',
    message: 'Введите электронную почту',
    min: 1,
    max: 256,
  },
  password: {
    required: true,
    type: 'string',
    message: 'Введите пароль от аккаунта',
    min: 8,
    max: 72,
  },
  password_confirmation: {
    required: true,
    type: 'string',
    message: 'Подтвердите пароль',
    min: 8,
    max: 72,
    validator: (_rule, value: string) => {
      if (!value) {
        return new Error('Введите подтверждение пароля')
      }
      if (value !== form.password) {
        return new Error('Пароли не совпадают')
      }
      return true
    },
  },
}

const onSubmitForm = () => {
  sender.submit(formRef.value, async () => {
    const data = await fetcher.register(form)

    if (data.ok) {
      notify.info('Аккаунт создан, теперь войдите в систему!')
      await router.push({ name: 'login' })
      return
    }

    notify.error(data.data.message)
  })
}
</script>
