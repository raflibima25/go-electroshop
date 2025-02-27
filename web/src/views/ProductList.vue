<script setup>
import { ref, onMounted, computed } from 'vue';
import { 
  PlusIcon, 
  SearchIcon,
  FilterIcon,
  ChevronLeftIcon,
  ChevronRightIcon,
  MoreHorizontalIcon,
  ChevronDownIcon
} from 'lucide-vue-next';

import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { useToast } from '@/composables/useToast';
import { productService } from '@/services/productService';
import ProductFormModal from '@/components/ProductFormModal.vue';
import UserLayout from '@/layouts/UserLayout.vue';
import { formatCurrency } from '@/utils/formatters';

import {
  DropdownMenu,
  DropdownMenuTrigger,
  DropdownMenuContent,
  DropdownMenuItem,
} from "@/components/ui/dropdown-menu";

const { showToast } = useToast();
const products = ref([]);
const categories = ref([]);
const isLoading = ref(false);
const showProductFormModal = ref(false);
const showDeleteDialog = ref(false);
const selectedProduct = ref(null);
const searchQuery = ref('');
const isFilterExpanded = ref(false);

const pagination = ref({
  current_page: 1,
  total_page: 1,
  total_items: 0,
  item_per_page: 10
});

const filter = ref({
  category: '',
  search: '',
  min_price: '',
  max_price: '',
  page: 1,
  limit: 10
});

const fetchProducts = async () => {
  try {
    isLoading.value = true;
    const response = await productService.getProducts(filter.value);
    if (response.data.status) {
      products.value = response.data.data.products;
      pagination.value = response.data.data.pagination;
    } else {
      showToast(response.data.message || 'Error fetching products', 'error');
    }
  } catch (error) {
    showToast(error.response?.data?.message || 'Error fetching products', 'error');
  } finally {
    isLoading.value = false;
  }
};

const fetchCategories = async () => {
  try {
    const response = await productService.getCategories();
    if (response.data.status) {
      categories.value = response.data.data.categories;
    }
  } catch (error) {
    console.error('Error fetching categories:', error);
    // Use default categories as fallback
    categories.value = ['Iphone', 'Samsung', 'Xiaomi'];
  }
};

const openAddModal = () => {
  selectedProduct.value = null;
  showProductFormModal.value = true;
};

const openEditModal = (product) => {
  selectedProduct.value = product;
  showProductFormModal.value = true;
};

const closeProductModal = () => {
  showProductFormModal.value = false;
  selectedProduct.value = null;
};

const handleProductAdded = async (productData) => {
  try {
    let response;
    
    if (selectedProduct.value) {
      // Update existing product
      response = await productService.updateProduct(selectedProduct.value.id, productData);
    } else {
      // Create new product
      response = await productService.createProduct(productData);
    }
    
    if (response.data.status) {
      showToast(response.data.message || 'Product saved successfully', 'success');
      showProductFormModal.value = false;
      selectedProduct.value = null;
      await fetchProducts();
    } else {
      showToast(response.data.message || 'Error saving product', 'error');
    }
  } catch (error) {
    showToast(error.response?.data?.message || 'Error saving product', 'error');
  }
};

const confirmDelete = (product) => {
  selectedProduct.value = product;
  showDeleteDialog.value = true;
};

const handleDelete = async () => {
  try {
    const response = await productService.deleteProduct(selectedProduct.value.id);
    if (response.data.status) {
      showToast(response.data.message || 'Product deleted successfully');
      await fetchProducts();
    } else {
      showToast(response.data.message || 'Error deleting product', 'error');
    }
  } catch (error) {
    showToast(error.response?.data?.message || 'Error deleting product', 'error');
  } finally {
    showDeleteDialog.value = false;
    selectedProduct.value = null;
  }
};

const changePage = (page) => {
  if (page >= 1 && page <= pagination.value.total_page) {
    filter.value.page = page;
    fetchProducts();
  }
};

const applyFilters = () => {
  filter.value.page = 1; // Reset to first page when applying filters
  fetchProducts();
};

const resetFilters = () => {
  filter.value = {
    category: '',
    search: '',
    min_price: '',
    max_price: '',
    page: 1,
    limit: 10
  };
  fetchProducts();
};

const filteredProducts = computed(() => {
  if (!searchQuery.value) return products.value;
  
  const query = searchQuery.value.toLowerCase();
  return products.value.filter(product => 
    product.name?.toLowerCase().includes(query) ||
    product.category?.toLowerCase().includes(query) ||
    product.price.toString().includes(query)
  );
});

onMounted(() => {
  fetchProducts();
  fetchCategories();
});
</script>

<template>
  <UserLayout>
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6 space-y-6">
      <!-- Header Section -->
      <div class="bg-white/30 border border-gray-200/50 rounded-2xl p-6 shadow-sm">
        <div class="flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
          <div>
            <h1 class="text-2xl font-bold text-black bg-clip-text">
              Products
            </h1>
            <p class="text-gray-500 mt-1">Manage your product inventory</p>
          </div>
          
          <div class="flex flex-wrap items-center gap-3">
            <div class="relative">
              <SearchIcon class="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-4 h-4" />
              <Input 
                v-model="searchQuery"
                type="text"
                placeholder="Search products..."
                class="pl-10 w-64 outline-none focus-visible:ring-0 focus-visible:ring-indigo-500 focus-visible:border-indigo-500"
              />
            </div>
            
            <Button 
              @click="openAddModal"
              class="bg-gradient-to-r bg-indigo-600 text-white hover:bg-indigo-500 shadow-lg shadow-indigo-500/20"
            >
              <PlusIcon class="w-4 h-4 mr-2" />
              Add Product
            </Button>
          </div>
        </div>
      </div>

      <!-- Filter Section -->
      <div class="bg-white rounded-xl shadow-sm border border-gray-200">
        <div 
          class="p-4 flex justify-between items-center cursor-pointer"
          @click="isFilterExpanded = !isFilterExpanded"
        >
          <div class="flex items-center gap-2">
            <FilterIcon class="w-4 h-4 text-gray-500" />
            <span class="font-medium">Filters</span>
          </div>
          <ChevronDownIcon
            :class="`w-4 h-4 text-gray-500 transform transition-transform duration-200 ${isFilterExpanded ? 'rotate-180' : ''}`"
          />
        </div>
        
        <div 
          v-show="isFilterExpanded"
          class="p-4 border-t border-gray-100 grid grid-cols-1 md:grid-cols-4 gap-4"
        >
          <div class="space-y-2">
            <label class="text-sm font-medium text-gray-700">Category</label>
            <select 
              v-model="filter.category"
              class="w-full rounded-lg border border-gray-200 p-2 focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 transition-shadow"
            >
              <option value="">All Categories</option>
              <option 
                v-for="category in categories" 
                :key="category" 
                :value="category"
              >
                {{ category }}
              </option>
            </select>
          </div>
          
          <div class="space-y-2">
            <label class="text-sm font-medium text-gray-700">Min Price</label>
            <Input 
              v-model="filter.min_price"
              type="number" 
              placeholder="Min price"
              min="0"
            />
          </div>
          
          <div class="space-y-2">
            <label class="text-sm font-medium text-gray-700">Max Price</label>
            <Input 
              v-model="filter.max_price"
              type="number" 
              placeholder="Max price"
              min="0"
            />
          </div>
          
          <div class="flex items-end gap-2">
            <Button 
              @click="applyFilters"
              class="bg-indigo-600 text-white hover:bg-indigo-500"
            >
              Apply Filter
            </Button>
            <Button 
              @click="resetFilters"
              variant="outline"
            >
              Reset
            </Button>
          </div>
        </div>
      </div>

      <!-- Products Table -->
      <div class="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
        <div class="overflow-x-auto">
          <table class="w-full">
            <thead class="bg-gray-50/50">
              <tr>
                <th class="px-6 py-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Thumbnail</th>
                <th class="px-6 py-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Product</th>
                <th class="px-6 py-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Category</th>
                <th class="px-6 py-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Price</th>
                <th class="px-6 py-4 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100">
              <template v-if="!isLoading && filteredProducts.length">
                <tr 
                  v-for="product in filteredProducts" 
                  :key="product.id"
                  class="hover:bg-gray-50/50 transition-colors"
                >
                  <td class="px-6 py-4 whitespace-nowrap">
                    <img 
                      :src="product.thumbnail || '/assets/placeholder.jpg'" 
                      :alt="product.name"
                      class="h-10 w-10 rounded-full object-cover"
                    />
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                    {{ product.name }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap">
                    <span class="px-3 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-700">
                      {{ product.category }}
                    </span>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                    {{ formatCurrency(product.price) }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-right text-sm">
                    <DropdownMenu>
                      <DropdownMenuTrigger asChild>
                        <Button variant="ghost" size="sm">
                          <MoreHorizontalIcon class="w-4 h-4" />
                        </Button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent align="end">
                        <DropdownMenuItem @click="openEditModal(product)">
                          Edit
                        </DropdownMenuItem>
                        <DropdownMenuItem 
                          @click="confirmDelete(product)"
                          class="text-red-600"
                        >
                          Delete
                        </DropdownMenuItem>
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </td>
                </tr>
              </template>
              <tr v-else-if="isLoading">
                <td colspan="5">
                  <div class="flex items-center justify-center py-8">
                    <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-600"></div>
                  </div>
                </td>
              </tr>
              <tr v-else>
                <td colspan="5">
                  <div class="flex flex-col items-center justify-center py-8 text-gray-500">
                    <div class="rounded-full bg-gray-100 p-3 mb-2">
                      <SearchIcon class="w-6 h-6" />
                    </div>
                    <p>No products found</p>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Pagination -->
        <div class="px-6 py-4 flex items-center justify-between border-t border-gray-100">
          <p class="text-sm text-gray-700">
            Showing 
            <span class="font-medium">{{ ((pagination.current_page - 1) * pagination.item_per_page) + 1 }}</span>
            to
            <span class="font-medium">{{ Math.min(pagination.current_page * pagination.item_per_page, pagination.total_items) }}</span>
            of
            <span class="font-medium">{{ pagination.total_items }}</span>
            results
          </p>
          <div class="flex items-center gap-2">
            <Button
              variant="outline"
              size="sm"
              @click="changePage(pagination.current_page - 1)"
              :disabled="pagination.current_page === 1"
              class="flex items-center gap-1"
            >
              <ChevronLeftIcon class="w-4 h-4" />
              Previous
            </Button>
            <Button
              variant="outline"
              size="sm"
              @click="changePage(pagination.current_page + 1)"
              :disabled="pagination.current_page === pagination.total_page"
              class="flex items-center gap-1"
            >
              Next
              <ChevronRightIcon class="w-4 h-4" />
            </Button>
          </div>
        </div>
      </div>

      <!-- Product Form Dialog -->
      <div v-if="showProductFormModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div class="bg-white rounded-lg shadow-xl max-w-md w-full max-h-[90vh] overflow-y-auto">
          <div class="flex justify-between items-center p-4 border-b">
            <h2 class="text-lg font-semibold text-gray-900">
              {{ selectedProduct ? 'Edit Product' : 'Add New Product' }}
            </h2>
            <button @click="closeProductModal" class="text-gray-500 hover:text-gray-700">
              <XIcon class="w-5 h-5" />
            </button>
          </div>
          
          <ProductFormModal 
            :product="selectedProduct"
            :categories="categories"
            @close="closeProductModal" 
            @submit="handleProductAdded" 
          />
        </div>
      </div>

      <!-- Delete Confirmation Dialog -->
      <div v-if="showDeleteDialog" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div class="bg-white rounded-lg shadow-xl p-6 max-w-md w-full">
          <h3 class="text-lg font-semibold mb-2">Delete Product</h3>
          <p class="text-gray-600 mb-4">
            Are you sure you want to delete this product? This action cannot be undone.
          </p>
          <div class="flex justify-end gap-3">
            <Button 
              variant="outline" 
              @click="showDeleteDialog = false"
            >
              Cancel
            </Button>
            <Button 
              class="bg-red-600 text-white hover:bg-red-500"
              @click="handleDelete"
            >
              Delete
            </Button>
          </div>
        </div>
      </div>
    </div>
  </UserLayout>
</template>