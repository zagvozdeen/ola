import { inject } from 'vue'
import type { PusherFunc } from '@shared/types'

export interface Notify {
  info: (n: string) => void,
  warn: (n: string) => void,
  error: (n: string) => void,
}

export const useNotifications = (): Notify => {
  const pusher = inject('notifications') as PusherFunc

  return {
    info: (n: string) => pusher('info', n, Date.now()),
    warn: (n: string) => pusher('warn', n, Date.now()),
    error: (n: string) => pusher('error', n, Date.now()),
  }
}
