import { reactive, ref } from 'vue'
import type { Cart, User } from '../types'

const state = {
  tma: window.Telegram?.WebApp?.initData || null,
  token: localStorage.getItem('token'),
  apiUrl: import.meta.env['VITE_API_URL'],
}

export const me = ref<User>()
export const cart = reactive<Cart>({
  product_ids: [],
})

export type State = ReturnType<typeof useState>

export const useState = () => {
  return {
    isTelegramEnv: () => state.tma !== null,
    isLoggedIn: () => state.token !== null,
    getAuthorizationHeader: () => state.tma !== null ? `tma ${state.tma}` : `Bearer ${state.token}`,
    getApiUrl: () => state.apiUrl,
    setToken: (token: string) => {
      localStorage.setItem('token', token)
      state.token = token
    },
    unsetToken: () => {
      localStorage.removeItem('token')
      state.token = null
    },
  }
}
