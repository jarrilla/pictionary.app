import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView
    },
    {
      path: '/donate',
      name: 'donate',
      // Lazy-loaded component
      component: () => import('../views/DonateView.vue')
    }
  ]
})

export default router 