import './styles.css'
import IMask from 'imask'

type FeedbackFormConfig = {
  endpoint: string
  successMessage: string
  submitErrorMessage: string
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

  const setConsentError = (message = ''): void => {
    const errorNode = form.querySelector<HTMLElement>('[data-error-for="consent"]')
    if (!(errorNode instanceof HTMLElement)) {
      return
    }

    errorNode.textContent = message
    errorNode.classList.toggle('invisible', message.length === 0)
  }

  const phoneMask = IMask(phoneInput, {
    mask: '+{7} (000) 000-00-00',
  })

  consentInput.addEventListener('change', () => {
    if (consentInput.checked) {
      setConsentError()
    }
  })

  form.addEventListener('submit', async (event) => {
    event.preventDefault()
    setStatus()
    setConsentError()

    const name = nameInput.value.trim()
    const content = contentInput.value.trim()
    const phone = phoneInput.value.trim()

    if (!consentInput.checked) {
      setConsentError('Нужно согласие на обработку данных')
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
        }),
      })

      if (!response.ok) {
        const errorText = (await response.text()).trim()
        setStatus(errorText || config.submitErrorMessage, 'error')
        return
      }

      setStatus(config.successMessage, 'success')
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
  initReviewsSlider()
  initReviewForms()
}

if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', initLandingPage, { once: true })
} else {
  initLandingPage()
}
