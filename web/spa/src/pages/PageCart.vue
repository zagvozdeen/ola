<template>
  <div class="min-h-dvh w-full flex flex-col gap-4 py-6 pb-22">
    <h1 class="text-lg font-bold">
      Корзина
    </h1>

    <div
      v-if="isLoading"
      class="flex justify-center my-4"
    >
      <n-spin size="small" />
    </div>

    <template v-else>
      <div
        v-if="cart.items.length === 0"
        class="bg-gray-500/20 border border-gray-500/30 rounded-2xl p-4"
      >
        <p class="text-sm text-gray-300">
          В корзине пока нет товаров.
        </p>
        <router-link
          class="inline-flex mt-3 bg-gray-600 hover:bg-gray-700 rounded px-3 py-1.5 text-xs font-bold"
          :to="{ name: 'main' }"
        >
          Перейти к ассортименту
        </router-link>
      </div>

      <template v-else>
        <ul class="grid grid-cols-1 gap-2">
          <li
            v-for="item in cart.items"
            :key="item.product_id"
            class="bg-gray-500/20 border border-gray-500/20 p-2 rounded-xl overflow-hidden flex gap-2"
          >
            <img
              class="size-20 object-cover rounded-lg"
              :src="item.file_content"
              alt=""
            >

            <div class="flex-1 min-w-0">
              <div class="flex justify-between gap-2">
                <span class="font-bold text-sm truncate">{{ item.product_name }}</span>
                <span class="text-xs uppercase bg-gray-600 font-bold px-2 py-0.5 rounded-full">{{ ProductTypeTranslates[item.type] }}</span>
              </div>

              <p class="text-xs mt-1 font-medium">
                от {{ item.price_from }} ₽{{ item.price_to ? ` до ${item.price_to} ₽` : '' }}
              </p>

              <div class="flex items-center gap-2 mt-2">
                <button
                  class="bg-gray-600 hover:bg-gray-700 rounded px-2 py-1 text-xs font-bold disabled:opacity-50 cursor-pointer"
                  :disabled="isUpdating(item.product_id)"
                  @click="() => handleDecrementQty(item.product_id, item.product_uuid, item.qty)"
                >
                  -
                </button>
                <span class="text-xs min-w-6 text-center">{{ item.qty }}</span>
                <button
                  class="bg-gray-600 hover:bg-gray-700 rounded px-2 py-1 text-xs font-bold disabled:opacity-50 cursor-pointer"
                  :disabled="isUpdating(item.product_id)"
                  @click="() => handleIncrementQty(item.product_id, item.qty)"
                >
                  +
                </button>
                <button
                  class="bg-red-700 hover:bg-red-800 rounded px-2 py-1 text-xs font-bold disabled:opacity-50 cursor-pointer ml-auto"
                  :disabled="isUpdating(item.product_id)"
                  @click="() => handleRemoveItem(item.product_id, item.product_uuid)"
                >
                  Удалить
                </button>
              </div>
            </div>
          </li>
        </ul>

        <div class="bg-gray-500/20 border border-gray-500/20 rounded-2xl p-4">
          <p class="text-sm">
            Позиций: <b>{{ totalItemsQty }}</b>
          </p>
          <p class="text-sm mt-1">
            Сумма: <b>от {{ totalPriceFrom }} ₽{{ totalPriceTo !== null ? ` до ${totalPriceTo} ₽` : '' }}</b>
          </p>
        </div>

        <n-form
          ref="formRef"
          class="w-full bg-gray-500/20 p-4 rounded-2xl"
          :rules="rules"
          :model="form"
          @submit.prevent="onSubmitOrder"
        >
          <n-form-item
            label="Имя"
            path="name"
          >
            <n-input
              v-model:value="form.name"
              placeholder="Введите имя"
            />
          </n-form-item>

          <n-form-item
            label="Телефон"
            path="phone"
          >
            <n-input
              v-model:value="form.phone"
              v-maska
              class="w-full"
              placeholder="+7 (___) ___-__-__"
              :input-props="{
                'data-maska': '+7 (###) ###-##-##',
              }"
            />
          </n-form-item>

          <n-form-item
            label="Комментарий"
            path="content"
          >
            <n-input
              v-model:value="form.content"
              type="textarea"
              placeholder="Комментарий к заказу"
              :autosize="{ minRows: 3, maxRows: 6 }"
            />
          </n-form-item>

          <n-button
            attr-type="submit"
            class="w-full"
            type="success"
            :disabled="isOrdering"
          >
            Оформить заказ
          </n-button>
        </n-form>
      </template>
    </template>

    <FooterMenu />
  </div>
</template>

<script setup lang="ts">
// import IMask from 'imask'
import {
  computed,
  nextTick,
  onBeforeUnmount,
  onMounted,
  reactive,
  ref,
  useTemplateRef,
  watch,
} from 'vue'
import FooterMenu from '@/components/FooterMenu.vue'
import { cart, useAuthState } from '@/composables/useAuthState'
import { useFetch } from '@/composables/useFetch'
import { useNotifications } from '@/composables/useNotifications'
import { useSender } from '@/composables/useSender'
import { type CreateOrderRequest, ProductTypeTranslates } from '@/types'
import { type FormInst, NButton, NForm, NFormItem, NInput, NSpin, type FormRules } from 'naive-ui'
import { vMaska } from 'maska/vue'

const auth = useAuthState()
const fetcher = useFetch()
const notify = useNotifications()
const sender = useSender()

const formRef = useTemplateRef<FormInst>('formRef')
// const phoneInputWrapperRef = ref<HTMLDivElement | null>(null)
const isLoading = ref(true)
const updatingProductID = ref<number | null>(null)
const isOrdering = ref(false)
// let phoneMask: ReturnType<typeof IMask> | null = null

const form = reactive<CreateOrderRequest>({
  name: '',
  phone: '',
  content: '',
})

const rules: FormRules = {
  name: {
    required: true,
    type: 'string',
    message: 'Введите имя',
    min: 1,
    max: 255,
  },
  phone: {
    required: true,
    type: 'string',
    message: 'Введите телефон',
    min: 1,
    max: 255,
  },
  content: {
    required: true,
    type: 'string',
    message: 'Введите комментарий',
    min: 1,
    max: 3000,
  },
}

const totalItemsQty = computed(() => {
  return cart.items.reduce((sum, item) => sum + item.qty, 0)
})

const totalPriceFrom = computed(() => {
  return cart.items.reduce((sum, item) => sum + item.price_from * item.qty, 0)
})

const totalPriceTo = computed(() => {
  if (cart.items.some(item => item.price_to === undefined)) {
    return null
  }

  return cart.items.reduce((sum, item) => sum + (item.price_to ?? 0) * item.qty, 0)
})

const isUpdating = (productID: number) => {
  return updatingProductID.value === productID
}

const refreshCart = async () => {
  const data = await fetcher.getCart()

  if (data.ok) {
    cart.items = data.data
  }
}

const handleSetQty = async (productID: number, qty: number) => {
  updatingProductID.value = productID

  try {
    const data = await fetcher.upsertCartItem(productID, qty)

    if (!data.ok) {
      return
    }

    await refreshCart()
  } finally {
    updatingProductID.value = null
  }
}

const handleIncrementQty = async (productID: number, qty: number) => {
  await handleSetQty(productID, qty + 1)
}

const handleDecrementQty = async (productID: number, productUUID: string, qty: number) => {
  if (qty <= 1) {
    await handleRemoveItem(productID, productUUID)
    return
  }

  await handleSetQty(productID, qty - 1)
}

const handleRemoveItem = async (productID: number, productUUID: string) => {
  updatingProductID.value = productID

  try {
    const data = await fetcher.deleteCartItem(productUUID)

    if (!data.ok) {
      return
    }

    await refreshCart()
  } finally {
    updatingProductID.value = null
  }
}

const onSubmitOrder = () => {
  sender.submit(formRef.value, async () => {
    isOrdering.value = true

    try {
      const data = await fetcher.createOrderFromCart(form)

      if (!data.ok) {
        return
      }

      notify.info('Заказ оформлен')
      form.content = ''
      await refreshCart()
    } finally {
      isOrdering.value = false
    }
  })
}

// const initPhoneMask = () => {
//   if (phoneMask || !phoneInputWrapperRef.value) {
//     return
//   }
//
//   const phoneInput = phoneInputWrapperRef.value.querySelector('input')
//   if (!(phoneInput instanceof HTMLInputElement)) {
//     return
//   }
//
//   phoneMask = IMask(phoneInput, {
//     mask: '+{7} (000) 000-00-00',
//   })
// }

// const destroyPhoneMask = () => {
//   phoneMask?.destroy()
//   phoneMask = null
// }

// watch(
//   () => !isLoading.value && cart.items.length > 0,
//   async (isOrderFormVisible) => {
//     if (isOrderFormVisible) {
//       await nextTick()
//       initPhoneMask()
//       return
//     }
//
//     destroyPhoneMask()
//   },
// )

// onBeforeUnmount(() => {
//   destroyPhoneMask()
// })

onMounted(async () => {
  await auth.ensureUserLoaded()

  if (auth.currentUser.value) {
    form.name = auth.currentUser.value.first_name
  }

  await refreshCart()
  isLoading.value = false
})
</script>
