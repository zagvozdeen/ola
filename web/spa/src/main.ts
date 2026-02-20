import './index.css'
import App from './App.vue'
import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import PageTest from './pages/PageTest.vue'

const routes = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/test', component: PageTest },
  ],
})

createApp(App).use(routes).mount('#app')