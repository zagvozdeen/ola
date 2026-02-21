<template>
  <div class="min-h-dvh w-full flex items-center justify-center">
    <form @submit.prevent="onSubmitForm">
      <input
        v-model="form.username"
        type="text"
        placeholder="Username"
      >
      <input
        v-model="form.password"
        type="password"
        placeholder="Password"
      >
      <button type="submit">
        Login
      </button>
    </form>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { reactive } from 'vue'
import { useState } from '@shared/composables/useState'
import { useFetch } from '@shared/composables/useFetch'

const router = useRouter()
const state = useState()
const fetcher = useFetch()
const form = reactive({
  username: '',
  password: '',
})

const onSubmitForm = () => {
  fetcher
    .login(form.username, form.password)
    .then(data => {
      if (data.ok) {
        state.setToken(data.data.token)
        router.push({ name: 'main' })
      }
    })
}
</script>
