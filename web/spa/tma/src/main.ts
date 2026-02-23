import './styles.css'
import App from './App.vue'
import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import MainPage from './pages/MainPage.vue'
import LoginPage from './pages/LoginPage.vue'
import RegisterPage from './pages/RegisterPage.vue'
import { useState } from '@shared/composables/useState'
import CartPage from './pages/CartPage.vue'
import SettingsPage from './pages/SettingsPage.vue'

const router = createRouter({
  history: createWebHistory('/tma/'),
  routes: [
    { path: '/', name: 'main', component: MainPage },
    { path: '/login', name: 'login', component: LoginPage },
    { path: '/register', name: 'register', component: RegisterPage },
    { path: '/cart', name: 'cart', component: CartPage },
    { path: '/settings', name: 'settings', component: SettingsPage },
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

// const meta = document.createElement('meta')
// meta.name = 'naive-ui-style'
// document.head.appendChild(meta)

createApp(App).use(router).mount('#app')
