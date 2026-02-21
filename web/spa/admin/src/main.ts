import './styles.css'
import App from './App.vue'
import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import MainPage from './pages/MainPage.vue'

const routes = createRouter({
  history: createWebHistory('/admin/'),
  routes: [
    { path: '', name: 'main', component: MainPage },
  ],
})

createApp(App).use(routes).mount('#app')