export type UUID = string
export type DateTime = string

export enum UserRole {
  User = 'user',
  Manager = 'manager',
  Moderator = 'moderator',
  Admin = 'admin',
}

export enum ProductType {
  Product = 'product',
  Service = 'service',
}

export enum OrderSource {
  Landing = 'landing',
  Spa = 'spa',
  Tma = 'tma',
}

export enum RequestStatus {
  Created = 'created',
  InProgress = 'in_progress',
  Reviewed = 'reviewed',
}

export const RequestStatusBgColor: Record<RequestStatus, string> = {
  [RequestStatus.Created]: 'bg-blue-700',
  [RequestStatus.InProgress]: 'bg-yellow-700',
  [RequestStatus.Reviewed]: 'bg-green-700',
}

export enum FeedbackType {
  ManagerContact = 'manager_contact',
  PartnershipOffer = 'partnership_offer',
  FeedbackRequest = 'feedback_request',
}

export const ProductTypeTranslates: Record<ProductType, string> = {
  [ProductType.Product]: 'Товар',
  [ProductType.Service]: 'Услуга',
}

export const ProductTypeBgColor: Record<ProductType, string> = {
  [ProductType.Product]: 'bg-green-700',
  [ProductType.Service]: 'bg-blue-700',
}

export const RequestStatusTranslates: Record<RequestStatus, string> = {
  [RequestStatus.Created]: 'Создана',
  [RequestStatus.InProgress]: 'В процессе',
  [RequestStatus.Reviewed]: 'Рассмотрена',
}

export const ProductTypeOptions = Object.values(ProductType).map((key) => ({
  value: key,
  label: ProductTypeTranslates[key],
}))

export const RequestStatusOptions = Object.values(RequestStatus).map((key) => ({
  value: key,
  label: RequestStatusTranslates[key],
}))

export const UserRoleTranslates: Record<UserRole, string> = {
  [UserRole.User]: 'Пользователь',
  [UserRole.Manager]: 'Менеджер',
  [UserRole.Moderator]: 'Модератор',
  [UserRole.Admin]: 'Администратор',
}

export const UserRoleOptions = Object.values(UserRole).map((key) => ({
  value: key,
  label: UserRoleTranslates[key],
}))

export type User = {
  id: number
  tid: number | null
  uuid: UUID
  first_name: string
  last_name: string | null
  username: string | null
  email: string | null
  phone: string | null
  role: UserRole
  created_at: DateTime
  updated_at: DateTime
}

export type File = {
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
  type: ProductType
  is_main: boolean
  file_content: string
  categories: Category[]
  user_id: number
  created_at: DateTime
  updated_at: DateTime
}

export type CartItem = {
  product_id: number
  product_uuid: UUID
  product_name: string
  price_from: number
  price_to?: number
  type: ProductType
  file_content?: string
  qty: number
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
  status: RequestStatus
  source: OrderSource
  name: string
  phone: string
  content: string
  user_id: number | null
  created_at: DateTime
  updated_at: DateTime
}

export type Feedback = {
  id: number
  uuid: UUID
  status: RequestStatus
  source: OrderSource
  type: FeedbackType
  name: string
  phone: string
  content: string
  user_id: number
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

export type AuthLoginRequest = {
  email: string
  password: string
}

export type AuthRegisterRequest = {
  first_name: string | null
  last_name: string | null
  email: string | null
  password: string | null
  password_confirmation: string | null
}

export type CreateFeedbackRequest = {
  name: string
  phone: string
  content: string
  type: FeedbackType
}

export type CreateOrderRequest = {
  name: string | null
  phone: string | null
  content: string | null
}

export type UpdateRequestStatusRequest = {
  status: RequestStatus | null
}

export type UpsertCategoryRequest = {
  name: string | null
}

export type UpdateUserRoleRequest = {
  role: UserRole | null
}

export type CreateProductRequest = {
  name: string | null
  description: string | null
  price_from: number | null
  price_to: number | null
  type: ProductType | null
  file_content: string | null
  category_uuids: UUID[]
}

export type ValidationError = {
  message: string
  errors: Record<string, string>
}

export type Levels = 'error' | 'warn' | 'info'

export interface Notification {
  id: number
  level: Levels
  msg: string
  date: number
}

export type PusherFunc = (level: Levels, msg: string, date: number) => void

export type Cart = {
  items: CartItem[]
}
