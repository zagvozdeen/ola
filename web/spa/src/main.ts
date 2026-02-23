import '@/styles.css'
import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import { useState } from '@/composables/useState'
import { getMe } from '@/composables/useFetch'
import App from '@/App.vue'
import MainPage from '@/pages/MainPage.vue'
import PageLogin from '@/pages/PageLogin.vue'
import PageRegister from '@/pages/PageRegister.vue'
import PageCart from '@/pages/PageCart.vue'
import PageSettings from '@/pages/PageSettings.vue'
import PageProducts from '@/pages/PageProducts.vue'
import PageProductEdit from '@/pages/PageProductEdit.vue'

const router = createRouter({
  history: createWebHistory('/spa/'),
  routes: [
    { path: '/', name: 'main', component: MainPage },
    { path: '/login', name: 'login', component: PageLogin },
    { path: '/register', name: 'register', component: PageRegister },
    { path: '/cart', name: 'cart', component: PageCart },
    { path: '/settings', name: 'settings', component: PageSettings },
    { path: '/products', name: 'products', component: PageProducts },
    { path: '/products/create', name: 'products.create', component: PageProductEdit },
    { path: '/products/:uuid/edit', name: 'products.edit', component: PageProductEdit },
  ],
})

const state = useState()

router.beforeEach((to, _, next) => {
  if (state.isTelegramEnv()) {
    next()
  } else if ((to.name !== 'login' && to.name !== 'register') && !state.isLoggedIn()) {
    next({ name: 'login' })
  } else if ((to.name === 'login' || to.name === 'register') && state.isLoggedIn()) {
    next({ name: 'main' })
  } else {
    next()
  }
})

if (state.isLoggedIn() || state.isTelegramEnv()) {
  getMe(state).then(data => {
    if (data.ok) {
      state.setMe(data.data)
    }
  })
}

createApp(App).use(router).mount('#app')
