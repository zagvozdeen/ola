import { computed, reactive, ref } from 'vue'
import { type Cart, type User, UserRole } from '@/types'

const state = {
  tma: window.Telegram?.WebApp?.initData || null,
  token: localStorage.getItem('token'),
  apiUrl: import.meta.env['VITE_API_URL'],
}

const me = ref<User>()
export const cart = reactive<Cart>({
  product_ids: [],
})
export const isUserModerator = (user: User) => user.role === UserRole.Moderator || user.role === UserRole.Admin
export const isUserAdmin = (user: User) => user.role === UserRole.Admin
export const isModerator = computed(() => me.value ? isUserModerator(me.value) : false)
export const isAdmin = computed(() =>  me.value ? isUserAdmin(me.value) : false)

// const subscribers: Array<> = []
// watch(me, (user) => {
//   if (user) {
//     for (const hook of subscribers) {
//       hook(user)
//     }
//   }
// })
let _hook: ((user: User) => void) | null  = null
export const onUserLoaded = (hook: (user: User) => void) => {
  if (me.value) {
    hook(me.value)
  } else {
    _hook = hook
  }
}

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
    setMe: (user: User) => {
      me.value = user
      if (_hook) {
        _hook(user)
      }
      // for (const hook of subscribers) {
      //   hook(user)
      // }
    },
  }
}
