import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'

const isMobile = () => {
  return /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent)
    || window.innerWidth <= 768
}

console.log(isMobile())

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => isMobile() 
        ? import('../views/MobileHomeView.vue')
        : Promise.resolve(HomeView)
    },
    {
      path: '/tip',
      name: 'tip',
      component: () => import('../views/TipView.vue')
    }
  ]
})

export default router 