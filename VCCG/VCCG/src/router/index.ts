import { createRouter, createWebHistory } from 'vue-router'
import Home from '@/views/Home.vue'
import Login from '@/views/Login.vue'
import Register from '@/views/Register.vue'
import MemeGenerator from '@/views/MemeGenerator.vue'
import { useUserStore } from '@/stores/user'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: Home
    },
    {
      path: '/login',
      name: 'login',
      component: Login
    },
    {
      path: '/register',
      name: 'register',
      component: Register
    },
    {
      path: '/generate',
      name: 'generate',
      component: MemeGenerator,
      beforeEnter: (to, from, next) => {
        const userStore = useUserStore()
        if (userStore.isAuthenticated) {
          next()
        } else {
          next('/login')
        }
      }
    }
  ]
})

export default router