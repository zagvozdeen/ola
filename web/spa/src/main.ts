import '@/styles.css'
import App from '@/App.vue'
import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import MainPage from '@/pages/MainPage.vue'
import LoginPage from '@/pages/LoginPage.vue'
import PageRegister from '@/pages/PageRegister.vue'
import { me, useState } from '@/composables/useState'
import CartPage from '@/pages/CartPage.vue'
import SettingsPage from '@/pages/SettingsPage.vue'
import PageProducts from '@/pages/PageProducts.vue'
import PageProductEdit from '@/pages/PageProductEdit.vue'
import { getMe } from '@/composables/useFetch'

const router = createRouter({
  history: createWebHistory('/spa/'),
  routes: [
    { path: '/', name: 'main', component: MainPage },
    { path: '/login', name: 'login', component: LoginPage },
    { path: '/register', name: 'register', component: PageRegister },
    { path: '/cart', name: 'cart', component: CartPage },
    { path: '/settings', name: 'settings', component: SettingsPage },
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
  getMe(state)
    .then(data => {
      if (data.ok) {
        me.value = data.data
      }
    })
}

// const meta = document.createElement('meta')
// meta.name = 'naive-ui-style'
// document.head.appendChild(meta)

createApp(App).use(router).mount('#app')
