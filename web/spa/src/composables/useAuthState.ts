import { computed, reactive, ref } from 'vue'
import { fetchJson, getAuthHeaders } from '@/composables/httpCore'
import { type Cart, type User, UserRole } from '@/types'

type AuthMode = 'guest' | 'bearer' | 'telegram'

type AuthSource =
  | { mode: 'guest' }
  | { mode: 'bearer'; token: string }
  | { mode: 'telegram'; initData: string }

const TOKEN_KEY = 'token'

const currentUser = ref<User | null>(null)
const isUserLoaded = ref(false)
const isLoadingUser = ref(false)

export const cart = reactive<Cart>({
  items: [],
})

const authSource = ref<AuthSource>({ mode: 'guest' })

let loadMePromise: Promise<void> | null = null

const readTelegramInitData = (): string | null => {
  return window.Telegram?.WebApp?.initData || null
}

const readStoredToken = (): string | null => {
  return localStorage.getItem(TOKEN_KEY)
}

const writeStoredToken = (token: string | null) => {
  if (token) {
    localStorage.setItem(TOKEN_KEY, token)
    return
  }

  localStorage.removeItem(TOKEN_KEY)
}

const detectAuthSource = (): AuthSource => {
  const tmaInitData = readTelegramInitData()
  if (tmaInitData) {
    window.Telegram.WebApp.isVerticalSwipesEnabled = false

    return { mode: 'telegram', initData: tmaInitData }
  }

  const token = readStoredToken()

  if (token) {
    return { mode: 'bearer', token }
  }

  return { mode: 'guest' }
}

const authHeaderFromSource = (source: AuthSource): string | null => {
  if (source.mode === 'telegram') {
    return `tma ${source.initData}`
  }

  if (source.mode === 'bearer') {
    return `Bearer ${source.token}`
  }

  return null
}

const resetUserLoadingState = (isLoaded: boolean) => {
  currentUser.value = null
  cart.items = []
  isUserLoaded.value = isLoaded
  isLoadingUser.value = false
  loadMePromise = null
}

export const isUserModerator = (user: User): boolean => {
  return user.role === UserRole.Moderator || user.role === UserRole.Admin
}

export const isUserAdmin = (user: User): boolean => {
  return user.role === UserRole.Admin
}

export const useAuthState = () => {
  const authMode = computed<AuthMode>(() => authSource.value.mode)

  const isTelegramEnv = computed(() => authSource.value.mode === 'telegram')
  const isBearerEnv = computed(() => authSource.value.mode === 'bearer')

  const hasAuthCredentials = computed(() => authSource.value.mode !== 'guest')
  const isAuthenticated = computed(() => currentUser.value !== null)

  const authorizationHeader = computed(() => authHeaderFromSource(authSource.value))

  const isModerator = computed(() => {
    return currentUser.value ? isUserModerator(currentUser.value) : false
  })

  const isAdmin = computed(() => {
    return currentUser.value ? isUserAdmin(currentUser.value) : false
  })

  const initAuthSource = () => {
    authSource.value = detectAuthSource()
  }

  const refreshAuthSource = () => {
    const prevMode = authSource.value.mode
    const prevHeader = authHeaderFromSource(authSource.value)
    const nextSource = detectAuthSource()
    const nextHeader = authHeaderFromSource(nextSource)

    authSource.value = nextSource

    if (prevMode !== nextSource.mode || prevHeader !== nextHeader) {
      resetUserLoadingState(false)
    }
  }

  const setToken = (token: string) => {
    writeStoredToken(token)
    authSource.value = { mode: 'bearer', token }
    resetUserLoadingState(false)
  }

  const unsetToken = () => {
    writeStoredToken(null)

    const tmaInitData = readTelegramInitData()
    authSource.value = tmaInitData
      ? { mode: 'telegram', initData: tmaInitData }
      : { mode: 'guest' }

    resetUserLoadingState(false)
  }

  const setMe = (user: User) => {
    currentUser.value = user
    isUserLoaded.value = true
  }

  const clearMe = () => {
    currentUser.value = null
    isUserLoaded.value = true
  }

  const logout = () => {
    if (authSource.value.mode === 'bearer') {
      writeStoredToken(null)
    }

    const tmaInitData = readTelegramInitData()
    authSource.value = tmaInitData
      ? { mode: 'telegram', initData: tmaInitData }
      : { mode: 'guest' }

    resetUserLoadingState(true)
  }

  const fetchMe = async () => {
    if (!hasAuthCredentials.value) {
      currentUser.value = null
      isUserLoaded.value = true
      return
    }

    if (isLoadingUser.value) {
      return loadMePromise ?? Promise.resolve()
    }

    isLoadingUser.value = true

    loadMePromise = (async () => {
      try {
        const res = await fetchJson<User>(
          '/api/me',
          {
            headers: getAuthHeaders(),
          },
          { notify: null },
        )

        if (res.ok) {
          currentUser.value = res.data
          return
        }

        currentUser.value = null

        if (authSource.value.mode === 'bearer') {
          unsetToken()
        }
      } finally {
        isLoadingUser.value = false
        isUserLoaded.value = true
      }
    })()

    return loadMePromise
  }

  const ensureUserLoaded = async () => {
    if (isUserLoaded.value) {
      return
    }

    await fetchMe()
  }

  return {
    currentUser,
    isUserLoaded,
    isLoadingUser,

    authSource,
    authMode,
    isTelegramEnv,
    isBearerEnv,
    hasAuthCredentials,
    isAuthenticated,
    authorizationHeader,

    isModerator,
    isAdmin,

    initAuthSource,
    refreshAuthSource,
    setToken,
    unsetToken,
    setMe,
    clearMe,
    logout,
    fetchMe,
    ensureUserLoaded,
  }
}
