export const i18n = {
  'user answer status must be null': 'Вы уже ответили на этот вопрос, если хотите ответить на вопрос повторно, то начните новый тест',
  'test session is not active': 'Этот тест устарел и закрыт, начните новый тест',
  'tma user not found: no rows in result set': 'Чтобы использовать мини-приложение, необходимо зарегистрироваться: введите команду /start в боте',
  'form.network_error': 'Ошибка сети. Попробуйте позже',
  'form.consent_required': 'Нужно согласие на обработку данных',
  'validation.invalid': 'Некорректное значение',
  'validation.required': 'Поле обязательно для заполнения',
  'validation.max': 'Значение слишком длинное',
  'validation.min': 'Значение слишком короткое',
  'validation.ru_phone': 'Введите телефон в формате +7 (999) 999-99-99',
  'validation.name.required': 'Введите имя',
  'validation.phone.required': 'Введите телефон',
  'validation.content.required': 'Введите текст',
  'validation.consent.required': 'Подтвердите согласие на обработку данных',
} as Readonly<Record<string, string>>

const cases = [2, 0, 1, 1, 1, 2]

export const pluralize = (count: number, words: string[]): string => {
  return words[ (count % 100 > 4 && count % 100 < 20) ? 2 : cases[Math.min(count % 10, 5)] as number ] || ''
}
