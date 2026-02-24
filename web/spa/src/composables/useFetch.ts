import {
  fetchJson,
  getAuthHeaders,
  getAuthJsonHeaders,
  getJsonHeaders,
} from '@/composables/httpCore'
import { type Notify, useNotifications } from '@/composables/useNotifications'
import type {
  AuthLoginRequest,
  AuthLoginResponse,
  CartItem,
  AuthRegisterRequest,
  AuthRegisterResponse,
  Category,
  CreateFeedbackRequest,
  CreateGuestFeedbackRequest,
  CreateGuestOrderRequest,
  CreateOrderRequest,
  CreateProductRequest,
  Feedback,
  File as UploadedFile,
  Order,
  Product,
  Review,
  UpdateProductRequest,
  User,
} from '@/types'

const login = async (notify: Notify, payload: AuthLoginRequest) => {
  return fetchJson<AuthLoginResponse>(
    '/api/auth/login',
    {
      method: 'POST',
      headers: getJsonHeaders(),
      body: JSON.stringify(payload),
    },
    {
      notify,
      disable401: true,
    },
  )
}

const register = async (notify: Notify, payload: AuthRegisterRequest) => {
  return fetchJson<AuthRegisterResponse>(
    '/api/auth/register',
    {
      method: 'POST',
      headers: getJsonHeaders(),
      body: JSON.stringify(payload),
    },
    { notify },
  )
}

const createGuestFeedback = async (notify: Notify, payload: CreateGuestFeedbackRequest) => {
  return fetchJson<Feedback>(
    '/api/guest/feedback',
    {
      method: 'POST',
      headers: getJsonHeaders(),
      body: JSON.stringify(payload),
    },
    { notify },
  )
}

const createGuestOrder = async (notify: Notify, payload: CreateGuestOrderRequest) => {
  return fetchJson<Order>(
    '/api/guest/orders',
    {
      method: 'POST',
      headers: getJsonHeaders(),
      body: JSON.stringify(payload),
    },
    { notify },
  )
}

const getProducts = async (notify: Notify) => {
  return fetchJson<Product[]>('/api/products', {
    headers: getAuthHeaders(),
  }, { notify })
}

const getProduct = async (notify: Notify, uuid: string) => {
  return fetchJson<Product>(`/api/products/${uuid}`, {
    headers: getAuthHeaders(),
  }, { notify })
}

const createProduct = async (notify: Notify, payload: CreateProductRequest) => {
  return fetchJson<Product>(
    '/api/products',
    {
      method: 'POST',
      headers: getAuthJsonHeaders(),
      body: JSON.stringify(payload),
    },
    { notify },
  )
}

const updateProduct = async (notify: Notify, uuid: string, payload: UpdateProductRequest) => {
  return fetchJson<Product>(
    `/api/products/${uuid}`,
    {
      method: 'PATCH',
      headers: getAuthJsonHeaders(),
      body: JSON.stringify(payload),
    },
    { notify },
  )
}

const deleteProduct = async (notify: Notify, uuid: string) => {
  return fetchJson<null>(
    `/api/products/${uuid}`,
    {
      method: 'DELETE',
      headers: getAuthHeaders(),
    },
    { notify },
  )
}

const uploadFile = async (notify: Notify, file: globalThis.File) => {
  const formData = new FormData()
  formData.append('file', file)

  return fetchJson<UploadedFile>(
    '/api/files',
    {
      method: 'POST',
      headers: getAuthHeaders(),
      body: formData,
    },
    { notify },
  )
}

const getCategories = async (notify: Notify) => {
  return fetchJson<Category[]>('/api/categories', {
    headers: getAuthHeaders(),
  }, { notify })
}

const getFeedback = async (notify: Notify) => {
  return fetchJson<Feedback[]>('/api/feedback', {
    headers: getAuthHeaders(),
  }, { notify })
}

const createFeedback = async (notify: Notify, payload: CreateFeedbackRequest) => {
  return fetchJson<Feedback>(
    '/api/feedback',
    {
      method: 'POST',
      headers: getAuthJsonHeaders(),
      body: JSON.stringify(payload),
    },
    { notify },
  )
}

const getReviews = async (notify: Notify) => {
  return fetchJson<Review[]>('/api/reviews', {
    headers: getAuthHeaders(),
  }, { notify })
}

const getOrders = async (notify: Notify) => {
  return fetchJson<Order[]>('/api/orders', {
    headers: getAuthHeaders(),
  }, { notify })
}

const getCart = async (notify: Notify) => {
  return fetchJson<CartItem[]>('/api/cart', {
    headers: getAuthHeaders(),
  }, { notify })
}

const upsertCartItem = async (notify: Notify, productID: number, qty: number) => {
  return fetchJson<null>(
    '/api/cart/items',
    {
      method: 'POST',
      headers: getAuthJsonHeaders(),
      body: JSON.stringify({
        product_id: productID,
        qty,
      }),
    },
    { notify },
  )
}

const deleteCartItem = async (notify: Notify, productUUID: string) => {
  return fetchJson<null>(
    `/api/cart/items/${encodeURIComponent(productUUID)}`,
    {
      method: 'DELETE',
      headers: getAuthHeaders(),
    },
    { notify },
  )
}

const createOrder = async (notify: Notify, payload: CreateOrderRequest) => {
  return fetchJson<Order>(
    '/api/orders',
    {
      method: 'POST',
      headers: getAuthJsonHeaders(),
      body: JSON.stringify(payload),
    },
    { notify },
  )
}

const createOrderFromCart = async (notify: Notify, payload: CreateOrderRequest) => {
  return fetchJson<Order>(
    '/api/orders/from-cart',
    {
      method: 'POST',
      headers: getAuthJsonHeaders(),
      body: JSON.stringify(payload),
    },
    { notify },
  )
}

const getUsers = async (notify: Notify) => {
  return fetchJson<User[]>('/api/users', {
    headers: getAuthHeaders(),
  }, { notify })
}

export const useFetch = () => {
  const notify = useNotifications()

  return {
    login: (payload: AuthLoginRequest) => login(notify, payload),
    register: (payload: AuthRegisterRequest) => register(notify, payload),
    createGuestFeedback: (payload: CreateGuestFeedbackRequest) => createGuestFeedback(notify, payload),
    createGuestOrder: (payload: CreateGuestOrderRequest) => createGuestOrder(notify, payload),
    getProducts: () => getProducts(notify),
    getProduct: (uuid: string) => getProduct(notify, uuid),
    createProduct: (payload: CreateProductRequest) => createProduct(notify, payload),
    updateProduct: (uuid: string, payload: UpdateProductRequest) => updateProduct(notify, uuid, payload),
    deleteProduct: (uuid: string) => deleteProduct(notify, uuid),
    uploadFile: (file: globalThis.File) => uploadFile(notify, file),
    getCategories: () => getCategories(notify),
    getFeedback: () => getFeedback(notify),
    createFeedback: (payload: CreateFeedbackRequest) => createFeedback(notify, payload),
    getReviews: () => getReviews(notify),
    getOrders: () => getOrders(notify),
    getCart: () => getCart(notify),
    upsertCartItem: (productID: number, qty: number) => upsertCartItem(notify, productID, qty),
    deleteCartItem: (productUUID: string) => deleteCartItem(notify, productUUID),
    createOrder: (payload: CreateOrderRequest) => createOrder(notify, payload),
    createOrderFromCart: (payload: CreateOrderRequest) => createOrderFromCart(notify, payload),
    getUsers: () => getUsers(notify),
  }
}
