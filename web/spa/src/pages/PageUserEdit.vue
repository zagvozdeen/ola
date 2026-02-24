<template>
  <div class="min-h-dvh w-full flex flex-col gap-4 py-6">
    <HeaderMenu
      title="Редактирование пользователя"
      :edit="true"
      back="users"
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
      <n-form-item label="Имя">
        <n-input
          :value="targetUser?.first_name ?? ''"
          readonly
        />
      </n-form-item>

      <n-form-item label="Email">
        <n-input
          :value="targetUser?.email ?? ''"
          readonly
        />
      </n-form-item>

      <n-form-item
        label="Роль"
        path="role"
      >
        <n-select
          v-model:value="form.role"
          :options="UserRoleOptions"
          placeholder="Выберите роль"
        />
      </n-form-item>
    </n-form>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref, useTemplateRef } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import HeaderMenu from '@/components/HeaderMenu.vue'
import { useFetch } from '@/composables/useFetch'
import { useNotifications } from '@/composables/useNotifications'
import { useSender } from '@/composables/useSender'
import { type FormInst, type FormRules, NForm, NFormItem, NInput, NSelect, NSpin } from 'naive-ui'
import { type UpdateUserRoleRequest, type User, UserRole, UserRoleOptions } from '@/types'

const route = useRoute()
const router = useRouter()
const fetcher = useFetch()
const notify = useNotifications()
const sender = useSender()

const formRef = useTemplateRef<FormInst>('formRef')
const isLoading = ref(true)
const targetUser = ref<User | null>(null)
const form = reactive<UpdateUserRoleRequest>({
  role: UserRole.User,
})

const rules: FormRules = {
  role: {
    required: true,
    type: 'string',
    message: 'Выберите роль',
  },
}

const onSubmit = () => {
  sender.submit(formRef.value, async () => {
    const uuid = route.params['uuid']
    if (typeof uuid !== 'string' || !uuid) {
      notify.error('Некорректный ID пользователя')
      await router.push({ name: 'users' })
      return
    }

    const data = await fetcher.updateUserRole(uuid, { role: form.role })
    if (data.ok) {
      notify.info('Роль пользователя обновлена')
      await router.push({ name: 'users' })
    }
  })
}

onMounted(() => {
  const uuid = route.params['uuid']
  if (typeof uuid !== 'string' || !uuid) {
    notify.error('Некорректный ID пользователя')
    router.push({ name: 'users' })
    isLoading.value = false
    return
  }

  fetcher.getUser(uuid)
    .then(data => {
      if (data.ok) {
        targetUser.value = data.data
        form.role = data.data.role
      }
    })
    .finally(() => {
      isLoading.value = false
    })
})
</script>
