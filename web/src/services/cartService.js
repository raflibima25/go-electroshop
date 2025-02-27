import apiClient from '@/utils/api'

export const cartService = {
  async getUserCart() {
    return apiClient.get('/cart')
  },

  async addToCart(productId, quantity) {
    return apiClient.post('/cart', { product_id: productId, quantity })
  },

  async updateCartItem(itemId, quantity) {
    return apiClient.put(`/cart/${itemId}`, { quantity })
  },

  async removeCartItem(itemId) {
    return apiClient.delete(`/cart/${itemId}`)
  },

  async clearCart() {
    return apiClient.delete('/cart')
  }
}
