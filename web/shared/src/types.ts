export type UUID = string
export type DateTime = string

export type UserRole = 'user' | 'moderator' | 'admin'

export type User = {
  id: number
  tid?: number
  uuid: UUID
  first_name: string
  last_name?: string
  username?: string
  email?: string
  role: UserRole
  created_at: DateTime
  updated_at: DateTime
}

export type File = {
  id: number
  uuid: UUID
  content: string
  size: number
  mime_type: string
  origin_name: string
  user_id: number
  created_at: DateTime
}

export type Product = {
  id: number
  uuid: UUID
  name: string
  description: string
  price_from: number
  price_to?: number
  file_id: number
  file_content?: string
  user_id: number
  created_at: DateTime
  updated_at: DateTime
}

export type Service = {
  id: number
  uuid: UUID
  name: string
  description: string
  price_from: number
  price_to?: number
  file_id: number
  file_content?: string
  user_id: number
  created_at: DateTime
  updated_at: DateTime
}

export type Review = {
  id: number
  uuid: UUID
  name: string
  content: string
  file_id: number
  file_content?: string
  user_id: number
  published_at: DateTime
  created_at: DateTime
  updated_at: DateTime
}

export type Order = {
  id: number
  uuid: UUID
  name: string
  phone: string
  content: string
  user_id?: number
  created_at: DateTime
  updated_at: DateTime
}

export type Feedback = {
  id: number
  uuid: UUID
  name: string
  phone: string
  content: string
  user_id?: number
  created_at: DateTime
  updated_at: DateTime
}

export type Category = {
  id: number
  uuid: UUID
  name: string
  created_at: DateTime
  updated_at: DateTime
}

export type CategoryProduct = {
  category_id: number
  product_id: number
}

export type CategoryService = {
  category_id: number
  service_id: number
}

export type AuthLoginRequest = {
  username: string
  password: string
}

export type AuthLoginResponse = {
  token: string
}

export type AuthRegisterRequest = {
  first_name: string | null
  last_name: string | null
  email: string | null
  password: string | null
  password_confirmation: string | null
}

export type AuthRegisterResponse = User

export type CreateFeedbackRequest = {
  name: string
  phone: string
  content: string
}

export type CreateOrderRequest = {
  name: string
  phone: string
  content: string
}

export type CreateGuestFeedbackRequest = {
  name: string
  phone: string
  content: string
  consent: boolean
}

export type CreateGuestOrderRequest = {
  name: string
  phone: string
  content: string
  consent: boolean
}

export type ValidationError = {
  message: string
  errors: Record<string, string>
}

export type Levels = 'error'| 'warn' | 'info'

export interface Notification {
  id: number
  level: Levels
  msg: string
  date: number
}

export type PusherFunc = (level: Levels, msg: string, date: number) => void