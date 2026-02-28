<template>
  <div class="min-h-dvh w-full flex items-center justify-center">
    <n-form
      ref="formRef"
      class="w-full bg-black/5 dark:bg-gray-500/20 p-4 my-8 rounded-2xl"
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
import { reactive, useTemplateRef } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthState } from '@/composables/useAuthState'
import { useFetch } from '@/composables/useFetch'
import { useSender } from '@/composables/useSender'
import type { AuthLoginRequest } from '@/types'
import { type FormInst, NButton, NForm, NFormItem, NInput, type FormRules } from 'naive-ui'

const router = useRouter()
const auth = useAuthState()
const fetcher = useFetch()
const sender = useSender()
const formRef = useTemplateRef<FormInst>('formRef')
const form = reactive<AuthLoginRequest>({
  email: '',
  password: '',
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
    const data = await fetcher.login(form)

    if (data.ok) {
      auth.setToken(data.data.token)
      await auth.fetchMe()
      await router.push({ name: 'main' })
    }
  })
}
</script>
