import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import Dashboard from '../views/Dashboard.vue'
import IncidentContext from '../components/IncidentContext.vue'
import Login from '../views/Login.vue'
import LogExplorer from '../views/LogExplorer.vue'
import TalonRules from '../views/TalonRules.vue'
import { useAuthStore } from '../store/auth'

const routes: Array<RouteRecordRaw> = [
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { requiresAuth: false }
  },
  {
    path: '/',
    name: 'Dashboard',
    component: Dashboard,
    meta: { requiresAuth: true }
  },
  {
    path: '/logs',
    name: 'LogExplorer',
    component: LogExplorer,
    meta: { requiresAuth: true }
  },
  {
    path: '/rules',
    name: 'TalonRules',
    component: TalonRules,
    meta: { requiresAuth: true }
  },
  {
    path: '/alert/:id',
    name: 'IncidentContext',
    component: IncidentContext,
    props: true,
    meta: { requiresAuth: true }
  }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

// Navigation Guard pentru protejarea rutelor
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth)

  if (requiresAuth && !authStore.isAuthenticated) {
    // Îl trimitem la login dacă vrea o pagină privată și nu are token
    next('/login')
  } else if (to.path === '/login' && authStore.isAuthenticated) {
    // Dacă e deja logat și merge pe /login, îl trimitem pe dashboard
    next('/')
  } else {
    next()
  }
})

export default router
