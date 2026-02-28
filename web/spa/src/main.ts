import '@/styles.css'
import { createApp } from 'vue'
import { createRouter, createWebHistory, START_LOCATION } from 'vue-router'
import App from '@/App.vue'
import { configureHttp } from '@/composables/httpCore'
import { isUserAdmin, isUserModerator, isUserOrderManager, useAuthState } from '@/composables/useAuthState'
import PageCart from '@/pages/PageCart.vue'
import PageCategories from '@/pages/PageCategories.vue'
import PageCategoryEdit from '@/pages/PageCategoryEdit.vue'
import PageFeedbackForm from '@/pages/PageFeedbackForm.vue'
import PageFeedbackEdit from '@/pages/PageFeedbackEdit.vue'
import PageFeedbacks from '@/pages/PageFeedbacks.vue'
import PageLogin from '@/pages/PageLogin.vue'
import PageMain from '@/pages/PageMain.vue'
import PageOrderEdit from '@/pages/PageOrderEdit.vue'
import PageOrders from '@/pages/PageOrders.vue'
import PageProductEdit from '@/pages/PageProductEdit.vue'
import PageProducts from '@/pages/PageProducts.vue'
import PageRegister from '@/pages/PageRegister.vue'
import PageSettings from '@/pages/PageSettings.vue'
import PageUserEdit from '@/pages/PageUserEdit.vue'
import PageUsers from '@/pages/PageUsers.vue'
import { Buffer } from 'buffer'
import base64url from 'base64url'

window.Buffer = Buffer

const router = createRouter({
  history: createWebHistory('/spa/'),
  routes: [
    { path: '/', name: 'main', component: PageMain },
    { path: '/login', name: 'login', component: PageLogin },
    { path: '/register', name: 'register', component: PageRegister },
    { path: '/cart', name: 'cart', component: PageCart },
    { path: '/settings', name: 'settings', component: PageSettings },
    { path: '/settings/manager', name: 'settings.manager', component: PageFeedbackForm },
    { path: '/settings/feedback', name: 'settings.feedback', component: PageFeedbackForm },
    { path: '/settings/partnership', name: 'settings.partnership', component: PageFeedbackForm },
    { path: '/products', name: 'products', component: PageProducts, meta: { requiresModerator: true } },
    { path: '/products/create', name: 'products.create', component: PageProductEdit, meta: { requiresModerator: true } },
    { path: '/products/:uuid/edit', name: 'products.edit', component: PageProductEdit, meta: { requiresModerator: true } },
    { path: '/feedback', name: 'feedback', component: PageFeedbacks, meta: { requiresModerator: true } },
    { path: '/feedback/:uuid/edit', name: 'feedback.edit', component: PageFeedbackEdit, meta: { requiresModerator: true } },
    { path: '/orders', name: 'orders', component: PageOrders, meta: { requiresOrderManager: true } },
    { path: '/orders/:uuid/edit', name: 'orders.edit', component: PageOrderEdit, meta: { requiresOrderManager: true } },
    { path: '/categories', name: 'categories', component: PageCategories, meta: { requiresModerator: true } },
    { path: '/categories/create', name: 'categories.create', component: PageCategoryEdit, meta: { requiresModerator: true } },
    { path: '/categories/:uuid/edit', name: 'categories.edit', component: PageCategoryEdit, meta: { requiresModerator: true } },
    { path: '/users', name: 'users', component: PageUsers, meta: { requiresAdmin: true } },
    { path: '/users/:uuid/edit', name: 'users.edit', component: PageUserEdit, meta: { requiresAdmin: true } },
  ],
})

const auth = useAuthState()
auth.initAuthSource()

configureHttp({
  getAuthorizationHeader: () => auth.authorizationHeader.value,
  onUnauthorized: () => {
    if (!auth.isTelegramEnv.value) {
      auth.unsetToken()
      location.reload()
      return
    }

    auth.clearMe()
  },
})

let isFirst = true
router.beforeEach(async (to, from) => {
  if (from === START_LOCATION && isFirst) {
    isFirst = false
    if (auth.authSource.value.mode === 'telegram') {
      if (auth.authSource.value.startParam) {
        try {
          const param = base64url.decode(auth.authSource.value.startParam)
          const [model, uuid] = param.split(':', 2)
          switch (model) {
          case 'order':
            return { name: 'orders.edit', params: { uuid } }
          case 'feedback':
            return { name: 'feedback.edit', params: { uuid } }
          }
        } catch (e) {
          console.error(e)
        }
      }
    }
  }

  const isAuthPage = to.name === 'login' || to.name === 'register'
  const requiresOrderManager = Boolean(to.meta['requiresOrderManager'])
  const requiresModerator = Boolean(to.meta['requiresModerator'])
  const requiresAdmin = Boolean(to.meta['requiresAdmin'])

  if (!isAuthPage && !auth.hasAuthCredentials.value) {
    return { name: 'login' }
  }

  if (isAuthPage && auth.hasAuthCredentials.value) {
    return { name: 'main' }
  }

  if (requiresOrderManager) {
    await auth.ensureUserLoaded()

    const user = auth.currentUser.value

    if (!user || !isUserOrderManager(user)) {
      return { name: 'main' }
    }
  }

  if (requiresModerator) {
    await auth.ensureUserLoaded()

    const user = auth.currentUser.value

    if (!user || !isUserModerator(user)) {
      return { name: 'main' }
    }
  }

  if (requiresAdmin) {
    await auth.ensureUserLoaded()

    const user = auth.currentUser.value

    if (!user || !isUserAdmin(user)) {
      return { name: 'main' }
    }
  }

  return true
})

if (auth.hasAuthCredentials.value) {
  void auth.fetchMe()
}

createApp(App).use(router).mount('#app')
