import '@/styles.css'
import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import App from '@/App.vue'
import { configureHttp } from '@/composables/httpCore'
import { isUserModerator, useAuthState } from '@/composables/useAuthState'
import PageCart from '@/pages/PageCart.vue'
import PageLogin from '@/pages/PageLogin.vue'
import PageMain from '@/pages/PageMain.vue'
import PageProductEdit from '@/pages/PageProductEdit.vue'
import PageProducts from '@/pages/PageProducts.vue'
import PageRegister from '@/pages/PageRegister.vue'
import PageSettings from '@/pages/PageSettings.vue'

const router = createRouter({
  history: createWebHistory('/spa/'),
  routes: [
    { path: '/', name: 'main', component: PageMain },
    { path: '/login', name: 'login', component: PageLogin },
    { path: '/register', name: 'register', component: PageRegister },
    { path: '/cart', name: 'cart', component: PageCart },
    { path: '/settings', name: 'settings', component: PageSettings },
    { path: '/products', name: 'products', component: PageProducts, meta: { requiresModerator: true } },
    { path: '/products/create', name: 'products.create', component: PageProductEdit, meta: { requiresModerator: true } },
    { path: '/products/:uuid/edit', name: 'products.edit', component: PageProductEdit, meta: { requiresModerator: true } },
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

router.beforeEach(async (to) => {
  const isAuthPage = to.name === 'login' || to.name === 'register'
  const requiresModerator = Boolean(to.meta['requiresModerator'])

  if (!isAuthPage && !auth.hasAuthCredentials.value) {
    return { name: 'login' }
  }

  if (isAuthPage && auth.hasAuthCredentials.value) {
    return { name: 'main' }
  }

  if (requiresModerator) {
    await auth.ensureUserLoaded()

    const user = auth.currentUser.value

    if (!user || !isUserModerator(user)) {
      return { name: 'main' }
    }
  }

  return true
})

if (auth.hasAuthCredentials.value) {
  void auth.fetchMe()
}

createApp(App).use(router).mount('#app')
