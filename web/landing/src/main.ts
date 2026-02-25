import './styles.css'
import IMask from 'imask'
import { i18n } from '../../spa/src/composables/useI18n'
import type { ValidationError } from '../../spa/src/types'

type FeedbackFormConfig = {
  endpoint: string
  successMessage: string
  submitErrorMessage: string
}

type ValidationField = 'name' | 'phone' | 'content' | 'consent'
type FormField = ValidationField

const validationFields: ValidationField[] = ['name', 'phone', 'content', 'consent']

const getValidationMessage = (field: string, tag: string): string => {
  const fieldMessage = i18n[`validation.${field}.${tag}`]
  if (fieldMessage !== undefined) {
    return fieldMessage
  }

  const tagMessage = i18n[`validation.${tag}`]
  if (tagMessage !== undefined) {
    return tagMessage
  }

  return i18n['validation.invalid'] || 'Некорректное значение'
}

const initMobileMenu = (): void => {
  const toggle = document.getElementById('menu-toggle')
  const menu = document.getElementById('mobile-menu')

  if (!(toggle instanceof HTMLButtonElement) || !(menu instanceof HTMLElement)) {
    return
  }

  const body = document.body

  const setOpen = (isOpen: boolean): void => {
    menu.classList.toggle('is-open', isOpen)
    toggle.classList.toggle('is-open', isOpen)
    body.classList.toggle('menu-open', isOpen)
    toggle.setAttribute('aria-expanded', String(isOpen))
    menu.setAttribute('aria-hidden', String(!isOpen))
  }

  toggle.addEventListener('click', () => {
    const isOpen = !menu.classList.contains('is-open')
    setOpen(isOpen)
  })

  menu.addEventListener('click', (event) => {
    const target = event.target instanceof Element ? event.target : null
    const link = target?.closest<HTMLAnchorElement>('a[href]')

    if (!link) {
      if (target === menu) {
        setOpen(false)
      }
      return
    }

    setOpen(false)

    const hash = link.getAttribute('href') || ''
    if (!hash.startsWith('#') || hash === '#') {
      return
    }

    const section = document.querySelector<HTMLElement>(hash)
    if (!section) {
      return
    }

    event.preventDefault()
    section.scrollIntoView({ behavior: 'smooth', block: 'start' })
    history.pushState(null, '', hash)
  })

  window.addEventListener('keydown', (event) => {
    if (event.key === 'Escape' && menu.classList.contains('is-open')) {
      setOpen(false)
      toggle.focus()
    }
  })
}

const initFeedbackForm = (
  form: HTMLFormElement,
  config: FeedbackFormConfig,
): void => {
  const nameInput = form.elements.namedItem('name')
  const phoneInput = form.elements.namedItem('phone')
  const contentInput = form.elements.namedItem('content')
  const consentInput = form.elements.namedItem('consent')
  const statusNode = form.querySelector<HTMLElement>('[data-form-status]')
  const submitButton = form.querySelector<HTMLButtonElement>('button[type="submit"]')

  if (
    !(nameInput instanceof HTMLInputElement) ||
    !(phoneInput instanceof HTMLInputElement) ||
    !(contentInput instanceof HTMLTextAreaElement) ||
    !(consentInput instanceof HTMLInputElement) ||
    !(statusNode instanceof HTMLElement) ||
    !(submitButton instanceof HTMLButtonElement)
  ) {
    return
  }

  const setStatus = (message = '', type: 'idle' | 'error' | 'success' = 'idle'): void => {
    statusNode.textContent = message
    statusNode.classList.remove('text-red-600', 'text-green-700')

    if (type === 'error') {
      statusNode.classList.add('text-red-600')
    }

    if (type === 'success') {
      statusNode.classList.add('text-green-700')
    }

    statusNode.classList.toggle('invisible', message.length === 0)
  }

  const setFieldError = (field: FormField, message = ''): void => {
    const errorNode = form.querySelector<HTMLElement>(`[data-error-for="${field}"]`)
    if (!(errorNode instanceof HTMLElement)) {
      return
    }

    errorNode.textContent = message
    errorNode.classList.toggle('invisible', message.length === 0)
  }

  const clearFieldErrors = (): void => {
    validationFields.forEach((field) => setFieldError(field))
  }

  const setValidationErrors = (errors: Record<string, string>): boolean => {
    let hasValidationErrors = false

    validationFields.forEach((field) => {
      const tag = errors[field]
      if (typeof tag !== 'string' || tag.length === 0) {
        setFieldError(field)
        return
      }

      setFieldError(field, getValidationMessage(field, tag))
      hasValidationErrors = true
    })

    return hasValidationErrors
  }

  const phoneMask = IMask(phoneInput, {
    mask: '+{7} (000) 000-00-00',
  })

  consentInput.addEventListener('change', () => {
    if (consentInput.checked) {
      setFieldError('consent')
    }
  })

  form.addEventListener('submit', async (event) => {
    event.preventDefault()
    setStatus()
    clearFieldErrors()

    const name = nameInput.value.trim()
    const content = contentInput.value.trim()
    const phone = phoneInput.value.trim()

    if (!consentInput.checked) {
      setFieldError('consent', i18n['form.consent_required'] || 'Нужно согласие на обработку данных')
      return
    }

    submitButton.disabled = true

    try {
      const response = await fetch(config.endpoint, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          name,
          phone,
          content,
          consent: consentInput.checked,
        }),
      })

      if (!response.ok) {
        const contentType = response.headers.get('Content-Type')?.toLowerCase() || ''
        if (contentType.includes('application/json')) {
          let error: ValidationError
          try {
            error = await response.json() as ValidationError
          } catch {
            setStatus(config.submitErrorMessage, 'error')
            return
          }

          if (setValidationErrors(error.errors)) {
            return
          }

          setStatus(config.submitErrorMessage, 'error')
          return
        }

        const errorText = (await response.text()).trim()
        setStatus(errorText || config.submitErrorMessage, 'error')
        return
      }

      setStatus(config.successMessage, 'success')
      form.reset()
      phoneMask.value = ''
    } catch {
      setStatus(i18n['form.network_error'] || 'Ошибка сети. Попробуйте позже', 'error')
    } finally {
      submitButton.disabled = false
    }
  })
}

const initReviewForms = (): void => {
  const feedbackForm = document.querySelector<HTMLFormElement>('#feedback-form')
  if (feedbackForm) {
    initFeedbackForm(feedbackForm, {
      endpoint: '/api/guest/feedback',
      successMessage: 'Спасибо! Отзыв отправлен',
      submitErrorMessage: 'Не удалось отправить отзыв',
    })
  }

  const orderForm = document.querySelector<HTMLFormElement>('#order-form')
  if (orderForm) {
    initFeedbackForm(orderForm, {
      endpoint: '/api/guest/orders',
      successMessage: 'Спасибо! Заявка отправлена',
      submitErrorMessage: 'Не удалось отправить заявку',
    })
  }
}

const initLandingPage = (): void => {
  initMobileMenu()
  initReviewForms()
}

if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', initLandingPage, { once: true })
} else {
  initLandingPage()
}
