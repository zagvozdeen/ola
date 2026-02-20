import './index.css'
import App from './App.vue'
import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router';

const routes = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/test', component: App },
  ],
})

createApp(App).use(routes).mount('#app')