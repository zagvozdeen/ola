import './styles.css'

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

if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', initReviewsSlider, { once: true })
} else {
  initReviewsSlider()
}
