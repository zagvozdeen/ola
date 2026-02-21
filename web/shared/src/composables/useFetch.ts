import { type State, useState } from './useState'
import { i18n } from './useI18n'

type ApiResult<T> = { ok: true; data: T } | { ok: false }
type OnError = (text: string) => void

const fetchJson = async <T>(state: State, input: RequestInfo, init?: RequestInit, onError?: OnError): Promise<ApiResult<T>> => {
  const res = await fetch(input, init)

  if (!res.ok) {
    if (res.status === 401) {
      // if (!state.isTelegramEnv()) {
      //   state.unsetToken()
      //   location.reload()
      // }
    }
    const text = (await res.text()).trim()
    if (onError !== undefined) {
      onError(i18n[text] || text)
    }
    return { ok: false }
  }

  return { ok: true, data: await res.json() }
}

const login = async (state: State, username: string, password: string, onError?: OnError) => {
  return fetchJson<{ token: string }>(state, `${state.getApiUrl()}/api/auth`, {
    method: 'POST',
    body: JSON.stringify({ username, password }),
  }, onError)
}

export const useFetch = () => {
  const state = useState()

  return {
    login: (username: string, password: string) => login(state, username, password),
  }
}
