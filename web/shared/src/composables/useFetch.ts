import { type State, useState } from './useState'
import { i18n } from './useI18n'
import type {
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
  Service,
  User, ValidationError,
} from '../types'
import { type Notify, useNotifications } from '@shared/composables/useNotifications'

type ApiResult<T> = { ok: true; data: T } | { ok: false, data: ValidationError }
type OnError = (text: string) => void

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
    // if (onError !== undefined) {
    //   onError('Network error')
    // }
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
      // if (!state.isTelegramEnv()) {
      //   state.unsetToken()
      //   location.reload()
      // }
    }
    if (res.headers.get('Content-Type') === 'application/json') {
      return {
        ok: false,
        data: await res.json(),
      }
    }
    const text = (await res.text()).trim()

    notify.error(i18n[text] || text)
    // if (onError !== undefined) {
    //   onError()
    // }
    return {
      ok: false,
      data: {
        message: i18n[text] || text,
        errors: {},
      },
    }
  }

  return { ok: true, data: await res.json() as T }
}

const login = async (state: State, notify: Notify, payload: object) => {
  return fetchJson<AuthLoginResponse>(state, notify, `${state.getApiUrl()}/api/auth/login`, {
    method: 'POST',
    headers: getJsonHeaders(),
    body: JSON.stringify(payload),
  })
}

const register = async (state: State, payload: AuthRegisterRequest, onError?: OnError) => {
  return fetchJson<AuthRegisterResponse>(state, `${state.getApiUrl()}/api/auth/register`, {
    method: 'POST',
    headers: getJsonHeaders(),
    body: JSON.stringify(payload),
  }, onError)
}

const createGuestFeedback = async (state: State, payload: CreateGuestFeedbackRequest, onError?: OnError) => {
  return fetchJson<Feedback>(state, `${state.getApiUrl()}/api/guest/feedback`, {
    method: 'POST',
    headers: getJsonHeaders(),
    body: JSON.stringify(payload),
  }, onError)
}

const createGuestOrder = async (state: State, payload: CreateGuestOrderRequest, onError?: OnError) => {
  return fetchJson<Order>(state, `${state.getApiUrl()}/api/guest/orders`, {
    method: 'POST',
    headers: getJsonHeaders(),
    body: JSON.stringify(payload),
  }, onError)
}

const getMe = async (state: State, onError?: OnError) => {
  return fetchJson<User>(state, `${state.getApiUrl()}/api/me`, {
    headers: getAuthHeaders(state),
  }, onError)
}

const getProducts = async (state: State, onError?: OnError) => {
  return fetchJson<Product[]>(state, `${state.getApiUrl()}/api/products`, {
    headers: getAuthHeaders(state),
  }, onError)
}

const getServices = async (state: State, onError?: OnError) => {
  return fetchJson<Service[]>(state, `${state.getApiUrl()}/api/services`, {
    headers: getAuthHeaders(state),
  }, onError)
}

const getCategories = async (state: State, onError?: OnError) => {
  return fetchJson<Category[]>(state, `${state.getApiUrl()}/api/categories`, {
    headers: getAuthHeaders(state),
  }, onError)
}

const getFeedback = async (state: State, onError?: OnError) => {
  return fetchJson<Feedback[]>(state, `${state.getApiUrl()}/api/feedback`, {
    headers: getAuthHeaders(state),
  }, onError)
}

const createFeedback = async (state: State, payload: CreateFeedbackRequest, onError?: OnError) => {
  return fetchJson<Feedback>(state, `${state.getApiUrl()}/api/feedback`, {
    method: 'POST',
    headers: getAuthJsonHeaders(state),
    body: JSON.stringify(payload),
  }, onError)
}

const getReviews = async (state: State, onError?: OnError) => {
  return fetchJson<Review[]>(state, `${state.getApiUrl()}/api/reviews`, {
    headers: getAuthHeaders(state),
  }, onError)
}

const getOrders = async (state: State, onError?: OnError) => {
  return fetchJson<Order[]>(state, `${state.getApiUrl()}/api/orders`, {
    headers: getAuthHeaders(state),
  }, onError)
}

const createOrder = async (state: State, payload: CreateOrderRequest, onError?: OnError) => {
  return fetchJson<Order>(state, `${state.getApiUrl()}/api/orders`, {
    method: 'POST',
    headers: getAuthJsonHeaders(state),
    body: JSON.stringify(payload),
  }, onError)
}

const getUsers = async (state: State, onError?: OnError) => {
  return fetchJson<User[]>(state, `${state.getApiUrl()}/api/users`, {
    headers: getAuthHeaders(state),
  }, onError)
}

export const useFetch = () => {
  const state = useState()
  const notify = useNotifications()

  return {
    login: (payload: object) => login(state, notify, payload),
    register: (payload: AuthRegisterRequest, onError?: OnError) => register(state, payload, onError),
    createGuestFeedback: (payload: CreateGuestFeedbackRequest, onError?: OnError) => createGuestFeedback(state, payload, onError),
    createGuestOrder: (payload: CreateGuestOrderRequest, onError?: OnError) => createGuestOrder(state, payload, onError),
    getMe: (onError?: OnError) => getMe(state, onError),
    getProducts: (onError?: OnError) => getProducts(state, onError),
    getServices: (onError?: OnError) => getServices(state, onError),
    getCategories: (onError?: OnError) => getCategories(state, onError),
    getFeedback: (onError?: OnError) => getFeedback(state, onError),
    createFeedback: (payload: CreateFeedbackRequest, onError?: OnError) => createFeedback(state, payload, onError),
    getReviews: (onError?: OnError) => getReviews(state, onError),
    getOrders: (onError?: OnError) => getOrders(state, onError),
    createOrder: (payload: CreateOrderRequest, onError?: OnError) => createOrder(state, payload, onError),
    getUsers: (onError?: OnError) => getUsers(state, onError),
  }
}
