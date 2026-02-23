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
          Войти в аккаунт
        </h1>
      </n-form-item>
      <n-form-item
        label="Введите электронную почту"
        path="email"
      >
        <n-input
          v-model:value="form.email"
          placeholder="Электронная почта"
          autofocus
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
        :show-feedback="false"
        :show-label="false"
      >
        <n-button
          attr-type="submit"
          type="success"
          class="flex-1"
        >
          Войти
        </n-button>
      </n-form-item>


      <p class="text-center mt-4">
        Ещё нет аккаунта? <router-link
          :to="{ name: 'register' }"
          class="underline"
        >
          Зарегистрироваться
        </router-link>
      </p>
    </n-form>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { reactive, useTemplateRef } from 'vue'
import { useState } from '@shared/composables/useState'
import { useFetch } from '@shared/composables/useFetch'
import { type FormRules, type FormInst, NForm, NFormItem, NButton, NInput } from 'naive-ui'
import { useSender } from '@shared/composables/useSender'
import { useNotifications } from '@shared/composables/useNotifications'

const router = useRouter()
const state = useState()
const fetcher = useFetch()
const sender = useSender()
const notify = useNotifications()
const formRef = useTemplateRef<FormInst>('formRef')
const form = reactive({
  email: null as string | null,
  password: null as string | null,
})
const rules: FormRules = {
  email: {
    required: true,
    type: 'email',
    message: 'Введите ваш электронную почту',
    min: 1,
    max: 255,
  },
  password: {
    required: true,
    type: 'string',
    message: 'Введите пароль от аккаунта',
    min: 1,
    max: 255,
  },
}

const onSubmitForm = () => {
  sender.submit(formRef.value, async () => {
    fetcher
      .login(form)
      .then(data => {
        if (data.ok) {
          state.setToken(data.data.token)
          router.push({ name: 'main' })
        }
        // } else {
        //   notify.error(data.data.message)
        // }
      })
  })
}
</script>
