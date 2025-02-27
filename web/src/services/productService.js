import apiClient from '@/utils/api'

export const productService = {
  async getProducts(params = {}) {
    return apiClient.get('/product', { params })
  },

  async getProductById(id) {
    return apiClient.get(`/product/${id}`)
  },

  async createProduct(data) {
    return apiClient.post(`/product-management`, data)
  },

  async updateProduct(id, data) {
    return apiClient.put(`/product-management/${id}`, data)
  },

  async deleteProduct(id) {
    return apiClient.delete(`/product-management/${id}`)
  },

  async getCategories() {
    return apiClient.get('/product/categories')
  }
}
