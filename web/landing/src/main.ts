import './styles.css'
import IMask from 'imask'

type FeedbackFormConfig = {
  endpoint: string
  successMessage: string
  submitErrorMessage: string
  messageRequiredText: string
}

const initReviewsSlider = (): void => {
  const slider = document.getElementById('reviews-slider')
  const track = document.getElementById('reviews-track')
  const prev = document.getElementById('reviews-prev')
  const next = document.getElementById('reviews-next')

  if (
    !(slider instanceof HTMLElement) ||
    !(track instanceof HTMLElement) ||
    !(prev instanceof HTMLButtonElement) ||
    !(next instanceof HTMLButtonElement)
  ) {
    return
  }

  const slides: HTMLElement[] = Array.from(
    track.querySelectorAll<HTMLElement>('.review-slide'),
  )
  if (slides.length === 0) {
    return
  }

  let index = 0

  const getGap = (): number => {
    const styles = window.getComputedStyle(track)
    const rawGap = styles.columnGap || styles.gap || '0'
    const parsedGap = Number.parseFloat(rawGap)
    return Number.isNaN(parsedGap) ? 0 : parsedGap
  }

  const getStep = (): number => {
    const firstSlide = slides[0]
    if (!firstSlide) {
      return 0
    }
    return firstSlide.getBoundingClientRect().width + getGap()
  }

  const render = (): void => {
    const maxIndex = slides.length - 1
    if (index < 0) index = 0
    if (index > maxIndex) index = maxIndex

    track.style.transform = `translate3d(${-index * getStep()}px, 0, 0)`

    prev.disabled = index === 0
    next.disabled = index === maxIndex
    prev.classList.toggle('opacity-40', prev.disabled)
    next.classList.toggle('opacity-40', next.disabled)
  }

  prev.addEventListener('click', () => {
    index -= 1
    render()
  })

  next.addEventListener('click', () => {
    index += 1
    render()
  })

  window.addEventListener('resize', render)
  render()
}

const initFeedbackForm = (
  form: HTMLFormElement,
  config: FeedbackFormConfig,
): void => {
  const nameInput = form.elements.namedItem('name')
  const phoneInput = form.elements.namedItem('phone')
  const textInput = form.elements.namedItem('message')
  const consentInput = form.elements.namedItem('consent')
  const statusNode = form.querySelector<HTMLElement>('[data-form-status]')
  const submitButton = form.querySelector<HTMLButtonElement>('button[type="submit"]')

  if (
    !(nameInput instanceof HTMLInputElement) ||
    !(phoneInput instanceof HTMLInputElement) ||
    !(textInput instanceof HTMLTextAreaElement) ||
    !(consentInput instanceof HTMLInputElement) ||
    !(statusNode instanceof HTMLElement) ||
    !(submitButton instanceof HTMLButtonElement)
  ) {
    return
  }

  const setError = (field: string, message = ''): void => {
    const errorNode = form.querySelector<HTMLElement>(`[data-error-for="${field}"]`)
    if (!(errorNode instanceof HTMLElement)) {
      return
    }

    errorNode.textContent = message
    errorNode.classList.toggle('invisible', message.length === 0)
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

  const clearErrors = (): void => {
    setError('name')
    setError('phone')
    setError('message')
    setError('consent')
  }

  const phoneMask = IMask(phoneInput, {
    mask: '+{7} (000) 000-00-00',
  })

  form.addEventListener('submit', async (event) => {
    event.preventDefault()
    clearErrors()
    setStatus()

    const name = nameInput.value.trim()
    const message = textInput.value.trim()
    const phone = phoneInput.value.trim()
    const phoneDigits = phone.replace(/\D/g, '')

    let hasError = false

    if (!name) {
      setError('name', 'Укажите ваше имя')
      hasError = true
    }

    if (phoneDigits.length !== 11 || !phoneDigits.startsWith('7')) {
      setError('phone', 'Введите номер в формате +7 (999) 123-45-67')
      hasError = true
    }

    if (!message) {
      setError('message', config.messageRequiredText)
      hasError = true
    }

    if (!consentInput.checked) {
      setError('consent', 'Нужно согласие на обработку данных')
      hasError = true
    }

    if (hasError) {
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
          message,
          consent: consentInput.checked,
        }),
      })

      if (!response.ok) {
        const errorText = (await response.text()).trim()
        setStatus(errorText || config.submitErrorMessage, 'error')
        return
      }

      const result = (await response.json()) as { message?: string }
      setStatus(result.message ?? config.successMessage, 'success')
      form.reset()
      phoneMask.value = ''
    } catch {
      setStatus('Ошибка сети. Попробуйте позже', 'error')
    } finally {
      submitButton.disabled = false
    }
  })
}

const initReviewForms = (): void => {
  const forms = Array.from(document.querySelectorAll<HTMLFormElement>('form#review-form'))

  const reviewForm = forms[0]
  if (reviewForm) {
    initFeedbackForm(reviewForm, {
      endpoint: '/api/reviews/submit-placeholder',
      successMessage: 'Спасибо! Отзыв отправлен',
      submitErrorMessage: 'Не удалось отправить отзыв',
      messageRequiredText: 'Напишите отзыв',
    })
  }

  const requestForm = forms[1]
  if (requestForm) {
    initFeedbackForm(requestForm, {
      endpoint: '/api/requests/submit-placeholder',
      successMessage: 'Спасибо! Заявка отправлена',
      submitErrorMessage: 'Не удалось отправить заявку',
      messageRequiredText: 'Укажите, что вас интересует',
    })
  }
}

const initLandingPage = (): void => {
  initReviewsSlider()
  initReviewForms()
}

if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', initLandingPage, { once: true })
} else {
  initLandingPage()
}
