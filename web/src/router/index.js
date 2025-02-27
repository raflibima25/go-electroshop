import { createRouter, createWebHistory } from 'vue-router'

// Gunakan lazy loading untuk semua komponen
const Login = () => import('@/views/LoginAuth.vue')
const Register = () => import('@/views/RegisterAuth.vue')
const DashboardAdmin = () => import('@/views/DashboardAdmin.vue')
const ProductList = () => import('@/views/ProductList.vue')
const PageNotFound = () => import('@/views/PageNotFound.vue')
const LandingPage = () => import('@/views/LandingPage.vue')
const UserProductList = () => import('@/views/UserProductList.vue')
const ShoppingCart = () => import('@/views/ShoppingCart.vue')
const ChatAssistant = () => import('@/views/ChatAssistant.vue')

const routes = [
  // public routes
  {
    path: '/',
    name: 'LandingPage',
    component: LandingPage,
    meta: { requiresAuth: false, title: 'ElectroShop - Home' }
  },
  {
    path: '/login',
    name: 'LoginAuth',
    component: Login,
    meta: {
      requiresAuth: false,
      title: 'Login'
    }
  },
  {
    path: '/register',
    name: 'Register',
    component: Register,
    meta: {
      requiresAuth: false,
      title: 'Register'
    }
  },
  {
    path: '/auth/google/callback',
    name: 'GoogleCallback',
    component: () => import('@/views/GoogleCallback.vue'),
    meta: {
      requiresAuth: false
    }
  },

  // Admin routes (role admin)
  {
    path: '/admin-dashboard',
    name: 'DashboardAdmin',
    component: DashboardAdmin,
    meta: {
      requiresAuth: true,
      role: 'admin',
      title: 'Admin Dashboard',
      layout: 'UserLayout'
    }
  },
  {
    path: '/products',
    name: 'ProductList',
    component: ProductList,
    meta: {
      requiresAuth: true,
      role: 'admin',
      title: 'Products Management',
      layout: 'UserLayout'
    }
  },

  // User routes (role user)
  {
    path: '/user/products',
    name: 'UserProductList',
    component: UserProductList,
    meta: {
      requiresAuth: true,
      role: 'user',
      title: 'Browse Products',
      layout: 'UserLayout'
    }
  },
  {
    path: '/user/cart',
    name: 'ShoppingCart',
    component: ShoppingCart,
    meta: {
      requiresAuth: true,
      role: 'user',
      title: 'Shopping Cart',
      layout: 'UserLayout'
    }
  },
  {
    path: '/user/chat',
    name: 'ChatAssistant',
    component: ChatAssistant,
    meta: {
      requiresAuth: true,
      role: 'user',
      title: 'Chat Assistant',
      layout: 'UserLayout'
    }
  },

  // 404 route - always at the end
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: PageNotFound,
    meta: { requiresAuth: false }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else {
      return { top: 0 }
    }
  }
})

// Fungsi helper untuk validasi
const isAuthenticated = () => !!localStorage.getItem('token')
const isAdmin = () => localStorage.getItem('isAdmin') === 'true'
const getUserRole = () => (isAdmin() ? 'admin' : 'user')

// Navigation guard
router.beforeEach(async (to, from, next) => {
  // Update title
  document.title = `${to.meta.title || 'ElectroShop'}`

  // Jika route memerlukan auth dan user tidak terautentikasi
  if (to.meta.requiresAuth && !isAuthenticated()) {
    next({
      name: 'LoginAuth',
      query: { redirect: to.fullPath }
    })
    return
  }

  // Jika user sudah login
  if (isAuthenticated()) {
    const userRole = getUserRole()

    // Mencegah akses ke login/register
    if (to.path === '/login' || to.path === '/register') {
      if (userRole === 'admin') {
        next('/admin-dashboard')
      } else {
        next('/user/products')
      }
      return
    }

    // Validasi role-based access - jika spesifik role diperlukan
    if (to.meta.role && to.meta.role !== userRole) {
      if (userRole === 'admin') {
        next('/admin-dashboard')
      } else {
        next('/user/products')
      }
      return
    }

    // Redirect user berdasarkan role jika mencoba akses halaman 404
    if (to.name === 'NotFound') {
      if (userRole === 'admin') {
        // Keep this behavior since you may want to preserve the 404
        // but ensure the "go back home" button redirects correctly
      } else {
        // Keep this behavior
      }
    }
  }

  next()
})

export default router
