<script setup>
import { ref, onMounted, onUnmounted, watch, computed } from 'vue';
import { 
  SearchIcon,
  FilterIcon,
  ChevronLeftIcon,
  ChevronRightIcon,
  ChevronDownIcon,
  ShoppingCartIcon
} from 'lucide-vue-next';

import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { useToast } from '@/composables/useToast';
import { productService } from '@/services/productService';
import UserLayout from '@/layouts/UserLayout.vue';
import { formatCurrency } from '@/utils/formatters';
import { useRouter } from 'vue-router';
import { cartService } from '@/services/cartService';
import { useAuth } from '@/composables/useAuth';

const router = useRouter();
const { showToast } = useToast();
const products = ref([]);
const categories = ref([]);
const isLoading = ref(false);
const searchQuery = ref('');
const isFilterExpanded = ref(false);
const cart = ref([]);
const { isAuthenticated } = useAuth();

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

const cartCount = ref(0);

const fetchCartCount = async () => {
  if (isAuthenticated.value) {
    try {
      const response = await cartService.getUserCart();
      if (response.data.status) {
        cartCount.value = response.data.data.total_items;
      }
    } catch (error) {
      console.error('Error fetching cart count:', error);
    }
  } else {
    // Ambil dari localStorage
    try {
      const localCart = JSON.parse(localStorage.getItem('cart') || '[]');
      cartCount.value = localCart.reduce((total, item) => total + item.quantity, 0);
    } catch (e) {
      cartCount.value = 0;
    }
  }
};

const addToCart = async (product) => {
  try {
    if (isAuthenticated.value) {
      // Tambahkan ke cart server
      await cartService.addToCart(product.id, 1);
      showToast(`${product.name} added to cart`, 'success');
      // Update cart count langsung setelah berhasil menambahkan
      await fetchCartCount();
    } else {
      // Tambahkan ke localStorage cart
      const existingItem = cart.value.find(item => item.id === product.id);
      
      if (existingItem) {
        existingItem.quantity += 1;
      } else {
        cart.value.push({
          ...product,
          quantity: 1
        });
      }
      
      localStorage.setItem('cart', JSON.stringify(cart.value));
      showToast(`${product.name} added to cart`, 'success');
      // Update cart count langsung untuk localStorage
      cartCount.value = cart.value.reduce((total, item) => total + item.quantity, 0);
    }
  } catch (error) {
    showToast('Error adding product to cart', 'error');
    console.error('Error adding to cart:', error);
  }
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

watch(isAuthenticated, () => {
  fetchCartCount();
}, { immediate: true });

// Watch cart for changes and update cartCount
watch(cart, () => {
  if (!isAuthenticated.value) {
    cartCount.value = cart.value.reduce((total, item) => total + item.quantity, 0);
  }
}, { deep: true });

onMounted(() => {
  fetchCartCount();
  
  // Set interval untuk update jumlah cart secara periodik
  const interval = setInterval(fetchCartCount, 30000); // Setiap 30 detik
  
  // Clear interval ketika komponen di-unmount
  onUnmounted(() => {
    clearInterval(interval);
  });
});

onMounted(() => {
  // Load cart from localStorage
  const savedCart = localStorage.getItem('cart');
  if (savedCart) {
    try {
      cart.value = JSON.parse(savedCart);
    } catch (e) {
      console.error('Error parsing cart from localStorage', e);
      cart.value = [];
    }
  }
  
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
              Browse Products
            </h1>
            <p class="text-gray-500 mt-1">Find the latest electronics for your needs</p>
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
              @click="router.push('/user/cart')"
              class="bg-indigo-600 text-white hover:bg-indigo-700 relative"
            >
              <ShoppingCartIcon class="w-5 h-5" />
              <span 
                v-if="cartCount > 0" 
                class="absolute -top-2 -right-2 bg-red-500 text-white text-xs rounded-full w-5 h-5 flex items-center justify-center"
              >
                {{ cartCount }}
              </span>
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

      <!-- Products Grid -->
      <div>
        <div v-if="isLoading" class="flex justify-center py-20">
          <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600"></div>
        </div>
        
        <div v-else-if="filteredProducts.length === 0" class="py-20 text-center text-gray-500">
          <div class="mx-auto w-24 h-24 rounded-full bg-gray-100 flex items-center justify-center mb-4">
            <SearchIcon class="w-12 h-12 text-gray-400" />
          </div>
          <h3 class="text-lg font-medium text-gray-900">No products found</h3>
          <p class="mt-2">Try adjusting your search or filter criteria</p>
        </div>
        
        <div v-else class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
          <div 
            v-for="product in filteredProducts" 
            :key="product.id" 
            class="bg-white overflow-hidden rounded-lg shadow-md border border-gray-200 hover:shadow-lg transition-shadow"
          >
            <div class="h-48 overflow-hidden">
              <img 
                :src="product.thumbnail || '/assets/placeholder.jpg'" 
                :alt="product.name"
                class="w-full h-full object-cover object-center"
              />
            </div>
            
            <div class="p-4">
              <div class="flex justify-between items-start">
                <div>
                  <span class="inline-block px-2 py-1 text-xs font-semibold rounded-full bg-blue-100 text-blue-800">
                    {{ product.category }}
                  </span>
                  <h3 class="mt-2 text-lg font-semibold text-gray-900 line-clamp-2">{{ product.name }}</h3>
                </div>
              </div>
              
              <div class="mt-3 flex items-center justify-between">
                <span class="text-xl font-bold text-gray-900">{{ formatCurrency(product.price) }}</span>
              </div>
              <div class="mt-4 flex justify-end">
                <Button 
                  @click="addToCart(product)"
                  class="bg-indigo-600 text-white hover:bg-indigo-500"
                >
                  <ShoppingCartIcon class="w-4 h-4 mr-1" />
                  Add to Cart
                </Button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Pagination -->
      <div v-if="pagination.total_page > 1" class="px-6 py-4 flex items-center justify-between bg-white rounded-xl shadow-sm border border-gray-200">
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
  </UserLayout>
</template>