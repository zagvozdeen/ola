import './styles.css'
import { i18n } from '../../spa/src/composables/useI18n'
import type { ValidationError } from '../../spa/src/types'
import { MaskInput } from 'maska'

new MaskInput('[data-maska]')

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
  const popupToggle = document.getElementById('mobile-menu-toggle')
  const menu = document.getElementById('mobile-menu')

  if (!(toggle instanceof HTMLButtonElement) || !(menu instanceof HTMLElement)) {
    return
  }

  const toggles = [toggle]
  if (popupToggle instanceof HTMLButtonElement) {
    toggles.push(popupToggle)
  }

  const setOpen = (isOpen: boolean): void => {
    menu.classList.toggle('is-open', isOpen)
    toggles.forEach((button) => {
      button.classList.toggle('is-open', isOpen)
      button.setAttribute('aria-expanded', String(isOpen))
    })
    document.body.classList.toggle('menu-open', isOpen)
    menu.setAttribute('aria-hidden', String(!isOpen))
  }

  toggles.forEach((button) => {
    button.addEventListener('click', () => {
      const isOpen = !menu.classList.contains('is-open')
      setOpen(isOpen)
    })
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

    // const hash = link.getAttribute('href') || ''
    // if (!hash.startsWith('#') || hash === '#') {
    //   return
    // }
    //
    // const section = document.querySelector<HTMLElement>(hash)
    // if (!section) {
    //   return
    // }
    //
    // event.preventDefault()
    // section.scrollIntoView({ behavior: 'smooth', block: 'start' })
    // history.pushState(null, '', hash)
  })

  window.addEventListener('keydown', (event) => {
    if (event.key === 'Escape' && menu.classList.contains('is-open')) {
      setOpen(false)
      toggle.focus()
    }
  })
}

const initReviewForms = (): void => {
  const form = document.querySelector<HTMLFormElement>('#order-form')
  if (form) {
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

    consentInput.addEventListener('change', () => {
      if (consentInput.checked) {
        setFieldError('consent')
      }
    })

    form.addEventListener('submit', async (event) => {
      event.preventDefault()
      setStatus()
      clearFieldErrors()

      // if (!consentInput.checked) {
      //   setFieldError('consent', i18n['form.consent_required'] || 'Нужно согласие на обработку данных')
      //   return
      // }

      submitButton.disabled = true

      try {
        const response = await fetch('/api/guest/orders', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            name: nameInput.value.trim(),
            phone: contentInput.value.trim(),
            content: phoneInput.value.trim(),
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
              setStatus('Не удалось отправить заявку', 'error')
              return
            }

            if (setValidationErrors(error.errors)) {
              return
            }

            setStatus('Не удалось отправить заявку', 'error')
            return
          }

          const errorText = (await response.text()).trim()
          setStatus(errorText || 'Не удалось отправить заявку', 'error')
          return
        }

        setStatus('Спасибо, заявка отправлена!', 'success')
        form.reset()
      } catch {
        setStatus(i18n['form.network_error'] || 'Ошибка сети. Попробуйте позже', 'error')
      } finally {
        submitButton.disabled = false
      }
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
