<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useAuth } from '@/composables/useAuth';
import { productService } from '@/services/productService';
import UserLayout from '@/layouts/UserLayout.vue';
import { 
  ShoppingBag, 
  Package, 
  ListOrdered, 
  TrendingUp,
  LucideShoppingCart,
  UserIcon
} from 'lucide-vue-next';
import { formatCurrency } from '@/utils/formatters';

const { checkAuth } = useAuth();
const router = useRouter();

const isLoading = ref(true);
const stats = ref({
  totalProducts: 0,
  categories: [],
  topProducts: []
});

// Mock data for analytics that would come from an actual API
const revenueData = ref([
  { month: 'Jan', amount: 3200000 },
  { month: 'Feb', amount: 4100000 },
  { month: 'Mar', amount: 3800000 },
  { month: 'Apr', amount: 5200000 },
  { month: 'May', amount: 4600000 },
  { month: 'Jun', amount: 5100000 }
]);

const loadDashboardData = async () => {
  try {
    isLoading.value = true;
    
    // Fetch products
    const productsResponse = await productService.getProducts({ limit: 100 });
    if (productsResponse.data.status) {
      const products = productsResponse.data.data.products;
      
      // Count total products
      stats.value.totalProducts = products.length;
      
      // Process categories
      const categoryMap = {};
      products.forEach(product => {
        if (!categoryMap[product.category]) {
          categoryMap[product.category] = 0;
        }
        categoryMap[product.category]++;
      });
      
      stats.value.categories = Object.keys(categoryMap).map(name => ({
        name,
        count: categoryMap[name]
      }));
      
      // Get top 5 products (by price as a proxy for importance)
      stats.value.topProducts = [...products]
        .sort((a, b) => b.price - a.price)
        .slice(0, 5);
    }
  } catch (error) {
    console.error('Error loading dashboard data:', error);
  } finally {
    isLoading.value = false;
  }
};

onMounted(() => {
  checkAuth();
  loadDashboardData();
});
</script>

<template>
  <UserLayout>
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
      <!-- Header Section -->
      <div class="mb-6">
        <h1 class="text-2xl font-bold text-gray-900">ElectroShop Dashboard</h1>
        <p class="text-gray-600 mt-1">Monitor and manage your electronic store</p>
      </div>
      
      <!-- Stats Cards -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-6">
        <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-200">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-gray-500 text-sm">Total Products</p>
              <p class="text-2xl font-bold text-gray-900">{{ stats.totalProducts }}</p>
            </div>
            <div class="bg-blue-100 p-3 rounded-full">
              <Package class="w-6 h-6 text-blue-500" />
            </div>
          </div>
        </div>
        
        <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-200">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-gray-500 text-sm">Categories</p>
              <p class="text-2xl font-bold text-gray-900">{{ stats.categories.length }}</p>
            </div>
            <div class="bg-purple-100 p-3 rounded-full">
              <ListOrdered class="w-6 h-6 text-purple-500" />
            </div>
          </div>
        </div>
        
        <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-200">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-gray-500 text-sm">Revenue (Monthly)</p>
              <p class="text-2xl font-bold text-gray-900">{{ formatCurrency(5100000) }}</p>
            </div>
            <div class="bg-green-100 p-3 rounded-full">
              <TrendingUp class="w-6 h-6 text-green-500" />
            </div>
          </div>
        </div>
        
        <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-200">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-gray-500 text-sm">Customers</p>
              <p class="text-2xl font-bold text-gray-900">125</p>
            </div>
            <div class="bg-amber-100 p-3 rounded-full">
              <UserIcon class="w-6 h-6 text-amber-500" />
            </div>
          </div>
        </div>
      </div>
      
      <!-- Main Content Grid -->
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <!-- Recent Products -->
        <div class="lg:col-span-2 bg-white rounded-xl shadow-sm p-6 border border-gray-200">
          <div class="flex justify-between items-center mb-4">
            <h2 class="text-lg font-semibold text-gray-900">Recent Products</h2>
            <button 
              @click="router.push('/products')"
              class="text-sm text-indigo-600 hover:text-indigo-500"
            >
              View all
            </button>
          </div>
          
          <div class="overflow-x-auto">
            <table class="min-w-full divide-y divide-gray-200">
              <thead>
                <tr>
                  <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Product</th>
                  <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Category</th>
                  <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Price</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-200">
                <template v-if="!isLoading && stats.topProducts.length">
                  <tr v-for="product in stats.topProducts" :key="product.id" class="hover:bg-gray-50">
                    <td class="px-4 py-3 whitespace-nowrap">
                      <div class="flex items-center">
                        <img 
                          :src="product.thumbnail || '/assets/placeholder.jpg'" 
                          class="w-8 h-8 rounded-full mr-3"
                          :alt="product.name"
                        />
                        <span class="font-medium text-gray-900">{{ product.name }}</span>
                      </div>
                    </td>
                    <td class="px-4 py-3 whitespace-nowrap">
                      <span class="px-2 py-1 text-xs rounded-full bg-blue-100 text-blue-700">
                        {{ product.category }}
                      </span>
                    </td>
                    <td class="px-4 py-3 whitespace-nowrap text-right">
                      <span class="text-gray-900 font-medium">
                        {{ formatCurrency(product.price) }}
                      </span>
                    </td>
                  </tr>
                </template>
                <tr v-else-if="isLoading">
                  <td colspan="3" class="px-4 py-6 text-center text-gray-500">
                    <div class="animate-spin rounded-full h-6 w-6 border-b-2 border-indigo-600 mx-auto"></div>
                    <p class="mt-2">Loading products...</p>
                  </td>
                </tr>
                <tr v-else>
                  <td colspan="3" class="px-4 py-6 text-center text-gray-500">
                    No products available
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
        
        <!-- Category Distribution -->
        <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-200">
          <h2 class="text-lg font-semibold text-gray-900 mb-4">Category Distribution</h2>
          
          <div v-if="isLoading" class="flex justify-center py-8">
            <div class="animate-spin rounded-full h-6 w-6 border-b-2 border-indigo-600"></div>
          </div>
          
          <div v-else-if="stats.categories.length" class="space-y-4">
            <div v-for="(category, index) in stats.categories" :key="index" class="space-y-2">
              <div class="flex justify-between items-center">
                <div class="font-medium text-gray-700">{{ category.name }}</div>
                <div class="text-sm text-gray-500">{{ category.count }} products</div>
              </div>
              <div class="w-full bg-gray-200 rounded-full h-2">
                <div class="h-2 rounded-full bg-indigo-600" 
                     :style="{width: `${(category.count / stats.totalProducts) * 100}%`}">
                </div>
              </div>
            </div>
          </div>
          
          <div v-else class="py-8 text-center text-gray-500">
            No categories available
          </div>
        </div>
      </div>
      
      <!-- Revenue Chart -->
      <div class="mt-6 bg-white rounded-xl shadow-sm p-6 border border-gray-200">
        <h2 class="text-lg font-semibold text-gray-900 mb-4">Revenue Overview</h2>
        
        <div class="h-72">
          <!-- In a real app, you would use a chart library like Chart.js -->
          <div class="flex h-full items-end justify-between">
            <div v-for="(item, index) in revenueData" :key="index" class="flex-1 flex flex-col items-center">
              <div class="w-full px-2">
                <div 
                  class="w-full bg-indigo-500 rounded-t-lg" 
                  :style="{
                    height: `${(item.amount / 6000000) * 100}%`,
                    opacity: 0.7 + (index * 0.05)
                  }"
                ></div>
              </div>
              <div class="text-xs text-gray-500 mt-2">{{ item.month }}</div>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Quick Actions -->
      <div class="mt-6 grid grid-cols-1 md:grid-cols-3 gap-6">
        <div class="bg-gradient-to-r from-indigo-500 to-purple-500 rounded-xl shadow-lg p-6 text-white">
          <div class="flex items-center gap-4">
            <div class="p-3 bg-white/20 rounded-full">
              <ShoppingBag class="w-6 h-6" />
            </div>
            <div>
              <h3 class="font-semibold text-lg">Manage Products</h3>
              <p class="text-white/80 text-sm">View and edit your product listings</p>
            </div>
          </div>
          <button 
            @click="router.push('/products')"
            class="mt-4 w-full py-2 bg-white/20 hover:bg-white/30 rounded-lg text-sm font-medium"
          >
            View Products
          </button>
        </div>
        
        <div class="bg-gradient-to-r from-emerald-500 to-teal-500 rounded-xl shadow-lg p-6 text-white">
          <div class="flex items-center gap-4">
            <div class="p-3 bg-white/20 rounded-full">
              <LucideShoppingCart class="w-6 h-6" />
            </div>
            <div>
              <h3 class="font-semibold text-lg">Manage Orders</h3>
              <p class="text-white/80 text-sm">View and process customer orders</p>
            </div>
          </div>
          <button 
            class="mt-4 w-full py-2 bg-white/20 hover:bg-white/30 rounded-lg text-sm font-medium"
          >
            View Orders
          </button>
        </div>
        
        <div class="bg-gradient-to-r from-amber-500 to-orange-500 rounded-xl shadow-lg p-6 text-white">
          <div class="flex items-center gap-4">
            <div class="p-3 bg-white/20 rounded-full">
              <UserIcon class="w-6 h-6" />
            </div>
            <div>
              <h3 class="font-semibold text-lg">Customer Management</h3>
              <p class="text-white/80 text-sm">View and manage customer accounts</p>
            </div>
          </div>
          <button 
            class="mt-4 w-full py-2 bg-white/20 hover:bg-white/30 rounded-lg text-sm font-medium"
          >
            View Customers
          </button>
        </div>
      </div>
    </div>
  </UserLayout>
</template>