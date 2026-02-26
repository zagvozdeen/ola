import {
  fetchJson,
  getAuthHeaders,
  getAuthJsonHeaders,
  getJsonHeaders,
} from '@/composables/httpCore'
import { type Notify, useNotifications } from '@/composables/useNotifications'
import type {
  AuthLoginRequest,
  CartItem,
  AuthRegisterRequest,
  Category,
  CreateFeedbackRequest,
  CreateOrderRequest,
  CreateProductRequest,
  Feedback,
  File as UploadedFile,
  Order,
  Product,
  // Review,
  UpdateRequestStatusRequest,
  UpdateUserRoleRequest,
  UpsertCategoryRequest,
  User,
} from '@/types'

const login = async (notify: Notify, payload: AuthLoginRequest) => {
  return fetchJson<User>('/api/auth/login', {
    method: 'POST',
    headers: getJsonHeaders(),
    body: JSON.stringify(payload),
  }, {
    notify,
    disable401: true,
  })
}

const register = async (notify: Notify, payload: AuthRegisterRequest) => {
  return fetchJson<User>('/api/auth/register', {
    method: 'POST',
    headers: getJsonHeaders(),
    body: JSON.stringify(payload),
  }, {
    notify,
  })
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

const updateProduct = async (notify: Notify, uuid: string, payload: CreateProductRequest) => {
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

const getCategory = async (notify: Notify, uuid: string) => {
  return fetchJson<Category>(`/api/categories/${uuid}`, {
    headers: getAuthHeaders(),
  }, { notify })
}

const createCategory = async (notify: Notify, payload: UpsertCategoryRequest) => {
  return fetchJson<Category>(
    '/api/categories',
    {
      method: 'POST',
      headers: getAuthJsonHeaders(),
      body: JSON.stringify(payload),
    },
    { notify },
  )
}

const updateCategory = async (notify: Notify, uuid: string, payload: UpsertCategoryRequest) => {
  return fetchJson<Category>(
    `/api/categories/${uuid}`,
    {
      method: 'PATCH',
      headers: getAuthJsonHeaders(),
      body: JSON.stringify(payload),
    },
    { notify },
  )
}

const deleteCategory = async (notify: Notify, uuid: string) => {
  return fetchJson<null>(
    `/api/categories/${uuid}`,
    {
      method: 'DELETE',
      headers: getAuthHeaders(),
    },
    { notify },
  )
}

const getFeedback = async (notify: Notify) => {
  return fetchJson<Feedback[]>('/api/feedback', {
    headers: getAuthHeaders(),
  }, { notify })
}

const getFeedbackItem = async (notify: Notify, uuid: string) => {
  return fetchJson<Feedback>(`/api/feedback/${uuid}`, {
    headers: getAuthHeaders(),
  }, { notify })
}

const updateFeedbackStatus = async (notify: Notify, uuid: string, payload: UpdateRequestStatusRequest) => {
  return fetchJson<Feedback>(
    `/api/feedback/${uuid}/status`,
    {
      method: 'PATCH',
      headers: getAuthJsonHeaders(),
      body: JSON.stringify(payload),
    },
    { notify },
  )
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

// const getReviews = async (notify: Notify) => {
//   return fetchJson<Review[]>('/api/reviews', {
//     headers: getAuthHeaders(),
//   }, { notify })
// }
//
// const getReview = async (notify: Notify, uuid: string) => {
//   return fetchJson<Review>(`/api/reviews/${uuid}`, {
//     headers: getAuthHeaders(),
//   }, { notify })
// }
//
// const createReview = async (notify: Notify, payload: UpsertReviewRequest) => {
//   return fetchJson<Review>(
//     '/api/reviews',
//     {
//       method: 'POST',
//       headers: getAuthJsonHeaders(),
//       body: JSON.stringify(payload),
//     },
//     { notify },
//   )
// }
//
// const updateReview = async (notify: Notify, uuid: string, payload: UpsertReviewRequest) => {
//   return fetchJson<Review>(
//     `/api/reviews/${uuid}`,
//     {
//       method: 'PATCH',
//       headers: getAuthJsonHeaders(),
//       body: JSON.stringify(payload),
//     },
//     { notify },
//   )
// }
//
// const deleteReview = async (notify: Notify, uuid: string) => {
//   return fetchJson<null>(
//     `/api/reviews/${uuid}`,
//     {
//       method: 'DELETE',
//       headers: getAuthHeaders(),
//     },
//     { notify },
//   )
// }

const getOrders = async (notify: Notify) => {
  return fetchJson<Order[]>('/api/orders', {
    headers: getAuthHeaders(),
  }, { notify })
}

const getOrder = async (notify: Notify, uuid: string) => {
  return fetchJson<Order>(`/api/orders/${uuid}`, {
    headers: getAuthHeaders(),
  }, { notify })
}

const updateOrderStatus = async (notify: Notify, uuid: string, payload: UpdateRequestStatusRequest) => {
  return fetchJson<Order>(
    `/api/orders/${uuid}/status`,
    {
      method: 'PATCH',
      headers: getAuthJsonHeaders(),
      body: JSON.stringify(payload),
    },
    { notify },
  )
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

// const createOrder = async (notify: Notify, payload: CreateOrderRequest) => {
//   return fetchJson<Order>(
//     '/api/orders',
//     {
//       method: 'POST',
//       headers: getAuthJsonHeaders(),
//       body: JSON.stringify(payload),
//     },
//     { notify },
//   )
// }

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

const getUser = async (notify: Notify, uuid: string) => {
  return fetchJson<User>(`/api/users/${uuid}`, {
    headers: getAuthHeaders(),
  }, { notify })
}

const updateUserRole = async (notify: Notify, uuid: string, payload: UpdateUserRoleRequest) => {
  return fetchJson<User>(
    `/api/users/${uuid}/role`,
    {
      method: 'PATCH',
      headers: getAuthJsonHeaders(),
      body: JSON.stringify(payload),
    },
    { notify },
  )
}

export const useFetch = () => {
  const notify = useNotifications()

  return {
    login: (payload: AuthLoginRequest) => login(notify, payload),
    register: (payload: AuthRegisterRequest) => register(notify, payload),
    getProducts: () => getProducts(notify),
    getProduct: (uuid: string) => getProduct(notify, uuid),
    createProduct: (payload: CreateProductRequest) => createProduct(notify, payload),
    updateProduct: (uuid: string, payload: CreateProductRequest) => updateProduct(notify, uuid, payload),
    deleteProduct: (uuid: string) => deleteProduct(notify, uuid),
    uploadFile: (file: globalThis.File) => uploadFile(notify, file),
    getCategories: () => getCategories(notify),
    getCategory: (uuid: string) => getCategory(notify, uuid),
    createCategory: (payload: UpsertCategoryRequest) => createCategory(notify, payload),
    updateCategory: (uuid: string, payload: UpsertCategoryRequest) => updateCategory(notify, uuid, payload),
    deleteCategory: (uuid: string) => deleteCategory(notify, uuid),
    getFeedback: () => getFeedback(notify),
    getFeedbackItem: (uuid: string) => getFeedbackItem(notify, uuid),
    updateFeedbackStatus: (uuid: string, payload: UpdateRequestStatusRequest) => updateFeedbackStatus(notify, uuid, payload),
    createFeedback: (payload: CreateFeedbackRequest) => createFeedback(notify, payload),
    // getReviews: () => getReviews(notify),
    // getReview: (uuid: string) => getReview(notify, uuid),
    // createReview: (payload: UpsertReviewRequest) => createReview(notify, payload),
    // updateReview: (uuid: string, payload: UpsertReviewRequest) => updateReview(notify, uuid, payload),
    // deleteReview: (uuid: string) => deleteReview(notify, uuid),
    getOrders: () => getOrders(notify),
    getOrder: (uuid: string) => getOrder(notify, uuid),
    updateOrderStatus: (uuid: string, payload: UpdateRequestStatusRequest) => updateOrderStatus(notify, uuid, payload),
    getCart: () => getCart(notify),
    upsertCartItem: (productID: number, qty: number) => upsertCartItem(notify, productID, qty),
    deleteCartItem: (productUUID: string) => deleteCartItem(notify, productUUID),
    // createOrder: (payload: CreateOrderRequest) => createOrder(notify, payload),
    createOrderFromCart: (payload: CreateOrderRequest) => createOrderFromCart(notify, payload),
    getUsers: () => getUsers(notify),
    getUser: (uuid: string) => getUser(notify, uuid),
    updateUserRole: (uuid: string, payload: UpdateUserRoleRequest) => updateUserRole(notify, uuid, payload),
  }
}
