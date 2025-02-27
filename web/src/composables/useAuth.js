import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { cartService } from '@/services/cartService'
import { useToast } from '@/composables/useToast'
import apiClient from '@/utils/api'

export function useAuth() {
  const router = useRouter()
  const { showToast } = useToast()

  const isAuthenticated = computed(() => {
    return !!localStorage.getItem('token') !== null
  })

  const isAdmin = computed(() => {
    return localStorage.getItem('isAdmin') === 'true'
  })

  const userRole = computed(() => {
    return isAdmin.value ? 'admin' : 'user'
  })

  const userName = computed(() => {
    return localStorage.getItem('userName') || 'User'
  })

  // Fungsi baru untuk sinkronisasi cart
  const syncCartAfterLogin = async () => {
    try {
      // Hanya jika user sudah login dan ada cart di localStorage
      if (!localStorage.getItem('token')) return

      const localCart = JSON.parse(localStorage.getItem('cart') || '[]')
      if (!localCart.length) return

      let syncCount = 0

      // Add each item to server cart
      for (const item of localCart) {
        await cartService.addToCart(item.id || item.product_id, item.quantity)
        syncCount += item.quantity
      }

      // Clear local cart after sync
      localStorage.removeItem('cart')

      if (syncCount > 0) {
        showToast(`${syncCount} items synchronized to your account`, 'success')
      }
    } catch (error) {
      console.error('Error syncing cart after login:', error)
    }
  }

  const login = async (credentials) => {
    try {
      const response = await apiClient.post('/auth/login', credentials)

      if (response.data.status) {
        localStorage.setItem('token', response.data.data.access_token)
        localStorage.setItem('isAdmin', response.data.data.is_admin)

        // Sync cart setelah login berhasil
        await syncCartAfterLogin()

        return {
          success: true,
          data: response.data.data
        }
      } else {
        return {
          success: false,
          message: response.data.message || 'Login failed'
        }
      }
    } catch (error) {
      return {
        success: false,
        message: error.response?.data?.message || 'An unexpected error occurred'
      }
    }
  }

  const checkAuth = (requiredRole = null) => {
    // cek auth
    if (!isAuthenticated.value) {
      router.push({
        name: 'LoginAuth',
        query: { redirect: router.currentRoute.value.fullPath }
      })
      return false
    }

    // cek role
    if (requiredRole && requiredRole !== userRole.value) {
      router.push(isAdmin.value ? '/admin-dashboard' : '/user/products')
      return false
    }

    return true
  }

  const logout = () => {
    try {
      localStorage.removeItem('token')
      localStorage.removeItem('isAdmin')
      localStorage.removeItem('userName')
      router.push({ name: 'LoginAuth' })
    } catch (error) {
      console.error('Error logout:', error)
    }
  }

  // helper untuk dapat initial user
  const getUserInitials = () => {
    return userName.value
      .split(' ')
      .map((n) => n[0])
      .join('')
      .toUpperCase()
  }

  return {
    isAuthenticated,
    isAdmin,
    userRole,
    userName,
    login,
    syncCartAfterLogin,
    checkAuth,
    logout,
    getUserInitials
  }
}
