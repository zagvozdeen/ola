import { type State, useState } from './useState'
import { i18n } from './useI18n'
import type {
  AuthLoginRequest,
  AuthLoginResponse,
  AuthRegisterRequest,
  AuthRegisterResponse,
  Category,
  CreateFeedbackRequest,
  CreateGuestFeedbackRequest,
  CreateGuestOrderRequest,
  CreateOrderRequest,
  Feedback,
  Order,
  Product,
  Review,
  User, ValidationError,
} from '@shared/types'
import { type Notify, useNotifications } from '@shared/composables/useNotifications'

type ApiResult<T> = { ok: true; data: T } | { ok: false, data: ValidationError }

const getJsonHeaders = (headers?: HeadersInit): Headers => {
  const result = new Headers(headers)
  result.set('Content-Type', 'application/json')
  return result
}

const getAuthHeaders = (state: State, headers?: HeadersInit): Headers => {
  const result = new Headers(headers)
  result.set('Authorization', state.getAuthorizationHeader())
  return result
}

const getAuthJsonHeaders = (state: State): Headers => {
  return getAuthHeaders(state, getJsonHeaders())
}

const fetchJson = async <T>(state: State, notify: Notify, input: RequestInfo, init?: RequestInit): Promise<ApiResult<T>> => {
  let res: Response

  try {
    res = await fetch(input, init)
  } catch {
    return {
      ok: false,
      data: {
        message: 'Network error',
        errors: {},
      },
    }
  }

  if (!res.ok) {
    if (res.status === 401) {
      if (!state.isTelegramEnv()) {
        state.unsetToken()
        location.reload()
      }
    }

    if (res.headers.get('Content-Type') === 'application/json') {
      return {
        ok: false,
        data: await res.json(),
      }
    }

    const rawText = (await res.text()).trim()
    const text=  i18n[rawText] || rawText

    notify.error(text)

    return {
      ok: false,
      data: {
        message: text,
        errors: {},
      },
    }
  }

  return { ok: true, data: await res.json() as T }
}

const login = async (state: State, notify: Notify, payload: AuthLoginRequest) => {
  return fetchJson<AuthLoginResponse>(state, notify, `${state.getApiUrl()}/api/auth/login`, {
    method: 'POST',
    headers: getJsonHeaders(),
    body: JSON.stringify(payload),
  })
}

const register = async (state: State, notify: Notify, payload: AuthRegisterRequest) => {
  return fetchJson<AuthRegisterResponse>(state, notify, `${state.getApiUrl()}/api/auth/register`, {
    method: 'POST',
    headers: getJsonHeaders(),
    body: JSON.stringify(payload),
  })
}

const getMe = async (state: State, notify: Notify) => {
  return fetchJson<User>(state, notify, `${state.getApiUrl()}/api/me`, {
    headers: getAuthHeaders(state),
  })
}

const createGuestFeedback = async (state: State, notify: Notify, payload: CreateGuestFeedbackRequest) => {
  return fetchJson<Feedback>(state, notify, `${state.getApiUrl()}/api/guest/feedback`, {
    method: 'POST',
    headers: getJsonHeaders(),
    body: JSON.stringify(payload),
  })
}

const createGuestOrder = async (state: State, notify: Notify, payload: CreateGuestOrderRequest) => {
  return fetchJson<Order>(state, notify, `${state.getApiUrl()}/api/guest/orders`, {
    method: 'POST',
    headers: getJsonHeaders(),
    body: JSON.stringify(payload),
  })
}

const getProducts = async (state: State, notify: Notify) => {
  return fetchJson<Product[]>(state, notify, `${state.getApiUrl()}/api/products`, {
    headers: getAuthHeaders(state),
  })
}

const getCategories = async (state: State, notify: Notify) => {
  return fetchJson<Category[]>(state, notify, `${state.getApiUrl()}/api/categories`, {
    headers: getAuthHeaders(state),
  })
}

const getFeedback = async (state: State, notify: Notify) => {
  return fetchJson<Feedback[]>(state, notify, `${state.getApiUrl()}/api/feedback`, {
    headers: getAuthHeaders(state),
  })
}

const createFeedback = async (state: State, notify: Notify, payload: CreateFeedbackRequest) => {
  return fetchJson<Feedback>(state, notify, `${state.getApiUrl()}/api/feedback`, {
    method: 'POST',
    headers: getAuthJsonHeaders(state),
    body: JSON.stringify(payload),
  })
}

const getReviews = async (state: State, notify: Notify) => {
  return fetchJson<Review[]>(state, notify, `${state.getApiUrl()}/api/reviews`, {
    headers: getAuthHeaders(state),
  })
}

const getOrders = async (state: State, notify: Notify) => {
  return fetchJson<Order[]>(state, notify, `${state.getApiUrl()}/api/orders`, {
    headers: getAuthHeaders(state),
  })
}

const createOrder = async (state: State, notify: Notify, payload: CreateOrderRequest) => {
  return fetchJson<Order>(state, notify, `${state.getApiUrl()}/api/orders`, {
    method: 'POST',
    headers: getAuthJsonHeaders(state),
    body: JSON.stringify(payload),
  })
}

const getUsers = async (state: State, notify: Notify) => {
  return fetchJson<User[]>(state, notify, `${state.getApiUrl()}/api/users`, {
    headers: getAuthHeaders(state),
  })
}

export const useFetch = () => {
  const state = useState()
  const notify = useNotifications()

  return {
    login: (payload: AuthLoginRequest) => login(state, notify, payload),
    register: (payload: AuthRegisterRequest) => register(state, notify, payload),
    getMe: () => getMe(state, notify),
    createGuestFeedback: (payload: CreateGuestFeedbackRequest) => createGuestFeedback(state, notify, payload),
    createGuestOrder: (payload: CreateGuestOrderRequest) => createGuestOrder(state, notify, payload),
    getProducts: () => getProducts(state, notify),
    getCategories: () => getCategories(state, notify),
    getFeedback: () => getFeedback(state, notify),
    createFeedback: (payload: CreateFeedbackRequest) => createFeedback(state, notify, payload),
    getReviews: () => getReviews(state, notify),
    getOrders: () => getOrders(state, notify),
    createOrder: (payload: CreateOrderRequest ) => createOrder(state, notify, payload),
    getUsers: () => getUsers(state, notify),
  }
}
