import './styles.css'
import App from './App.vue'
import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import MainPage from './pages/MainPage.vue'
import LoginPage from './pages/LoginPage.vue'
import RegisterPage from './pages/RegisterPage.vue'
import { useState } from '@shared/composables/useState'

const router = createRouter({
  history: createWebHistory('/tma/'),
  routes: [
    { path: '/', name: 'main', component: MainPage },
    { path: '/login', name: 'login', component: LoginPage },
    { path: '/register', name: 'register', component: RegisterPage },
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

createApp(App).use(router).mount('#app')
