import { i18n } from '@/composables/useI18n'
import type { ValidationError } from '@/types'
import type { Notify } from '@/composables/useNotifications'

export type ApiResult<T> =
  | { ok: true; data: T }
  | { ok: false; data: ValidationError }

type HttpConfig = {
  getAuthorizationHeader: (() => string | null) | null
  onUnauthorized: (() => void) | null
}

type HttpConfigInput = {
  getAuthorizationHeader?: () => string | null
  onUnauthorized?: () => void
}

type FetchJsonOptions = {
  notify?: Notify | null
  disable401?: boolean
}

const httpConfig: HttpConfig = {
  getAuthorizationHeader: null,
  onUnauthorized: null,
}

export const configureHttp = (config: HttpConfigInput) => {
  if (config.getAuthorizationHeader !== undefined) {
    httpConfig.getAuthorizationHeader = config.getAuthorizationHeader
  }

  if (config.onUnauthorized !== undefined) {
    httpConfig.onUnauthorized = config.onUnauthorized
  }
}

export const getJsonHeaders = (headers?: HeadersInit): Headers => {
  const result = new Headers(headers)
  result.set('Content-Type', 'application/json')

  return result
}

export const getAuthHeaders = (headers?: HeadersInit): Headers => {
  const result = new Headers(headers)
  const authorizationHeader = httpConfig.getAuthorizationHeader?.()

  if (authorizationHeader) {
    result.set('Authorization', authorizationHeader)
  }

  return result
}

export const getAuthJsonHeaders = (headers?: HeadersInit): Headers => {
  return getAuthHeaders(getJsonHeaders(headers))
}

const isJsonResponse = (res: Response): boolean => {
  const contentType = res.headers.get('Content-Type')

  return contentType !== null && contentType.includes('application/json')
}

export const fetchJson = async <T>(
  input: RequestInfo,
  init?: RequestInit,
  options: FetchJsonOptions = {},
): Promise<ApiResult<T>> => {
  const { notify = null, disable401 = false } = options

  let res: Response

  try {
    res = await fetch(input, init)
  } catch {
    return { ok: false, data: { message: 'Network error', errors: {} } }
  }

  if (!res.ok) {
    if (res.status === 401 && !disable401) {
      httpConfig.onUnauthorized?.()
    }

    if (res.status === 403) {
      notify?.error('У вас нет прав выполнять это действие!')

      return {
        ok: false,
        data: {
          message: 'Insufficient permissions',
          errors: {},
        },
      }
    }

    if (isJsonResponse(res)) {
      return { ok: false, data: await res.json() as ValidationError }
    }

    const rawText = (await res.text()).trim()
    const text = i18n[rawText] || rawText

    notify?.error(text)

    return {
      ok: false,
      data: {
        message: text,
        errors: {},
      },
    }
  }

  if (res.status === 204) {
    return { ok: true, data: null as T }
  }

  return { ok: true, data: await res.json() as T }
}
