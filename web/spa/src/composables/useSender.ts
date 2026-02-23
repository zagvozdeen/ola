import { type FormInst, type LoadingBarApi, useLoadingBar } from 'naive-ui'
import { type Notify, useNotifications } from './useNotifications'
// import type { ValidateError } from 'async-validator'

type cb = () => Promise<void>

const submit = async (loadingBar: LoadingBarApi, notify: Notify, formRef: FormInst | null, callback: cb): Promise<void> => {
  try {
    await formRef?.validate()
  } catch (e) {
    if (!(e instanceof Array)) {
      throw e
    }
    return
  }

  loadingBar.start()
  try {
    await callback()

    loadingBar.finish()
  } catch {
    loadingBar.error()

    // console.log(e)

    // notify.error(e as string)
    // message.error(await getMessage(e))
  }
}

export const useSender = () => {
  const loadingBar = useLoadingBar()
  const notify = useNotifications()

  return {
    submit: (formRef: FormInst | null, cb: cb) => submit(loadingBar, notify, formRef, cb),
  }
}
