import './styles.css'
import App from './App.vue'
import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import PageTest from './pages/PageTest.vue'
import PageMain from './pages/PageMain.vue'

const routes = createRouter({
  history: createWebHistory('/admin/'),
  routes: [
    { path: '', name: 'main', component: PageMain },
    { path: '/test', name: 'test', component: PageTest },
  ],
})

createApp(App).use(routes).mount('#app')