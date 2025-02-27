<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useToast } from '@/composables/useToast';
import UserLayout from '@/layouts/UserLayout.vue';
import { cartService } from '@/services/cartService';
import { useAuth } from '@/composables/useAuth';
import { formatCurrency } from '@/utils/formatters';
import { Button } from '@/components/ui/button';
import { 
  Trash2Icon, 
  MinusIcon, 
  PlusIcon, 
  ShoppingBagIcon,
  AlertTriangleIcon
} from 'lucide-vue-next';

const router = useRouter();
const { showToast } = useToast();
const { isAuthenticated } = useAuth();
const cart = ref([]);
const isLoading = ref(false);
const isProcessing = ref(false);

const cartTotal = computed(() => {
  return cart.value.reduce((total, item) => {
    const price = item.product?.price || item.price || 0;
    return total + (price * item.quantity);
  }, 0);
});

const cartItemCount = computed(() => {
  return cart.value.reduce((total, item) => total + item.quantity, 0);
});

const loadCart = async () => {
  isLoading.value = true;
  try {
    if (isAuthenticated.value) {
      // Load dari API
      const response = await cartService.getUserCart();
      if (response.data.status) {
        cart.value = response.data.data.items;
      }
    } else {
      // Load dari localStorage
      loadCartFromLocalStorage();
    }
  } catch (error) {
    showToast('Error loading cart', 'error');
    console.error('Error loading cart:', error);
    // Fallback ke localStorage jika API error
    loadCartFromLocalStorage();
  } finally {
    isLoading.value = false;
  }
};

const loadCartFromLocalStorage = () => {
  const savedCart = localStorage.getItem('cart');
  if (savedCart) {
    try {
      cart.value = JSON.parse(savedCart);
    } catch (e) {
      console.error('Error parsing cart from localStorage', e);
      cart.value = [];
    }
  }
};

const saveCart = () => {
  localStorage.setItem('cart', JSON.stringify(cart.value));
};

// Fungsi untuk menambah jumlah item
const increaseQuantity = async (item) => {
  if (isAuthenticated.value) {
    try {
      await cartService.updateCartItem(item.id, item.quantity + 1);
      await loadCart(); // Reload cart dari API
    } catch (error) {
      showToast('Error updating cart', 'error');
    }
  } else {
    item.quantity++;
    saveCartToLocalStorage();
  }
};

// Fungsi untuk mengurangi jumlah item
const decreaseQuantity = async (item) => {
  if (item.quantity <= 1) return;
  
  if (isAuthenticated.value) {
    try {
      await cartService.updateCartItem(item.id, item.quantity - 1);
      await loadCart(); // Reload cart dari API
    } catch (error) {
      showToast('Error updating cart', 'error');
    }
  } else {
    item.quantity--;
    saveCartToLocalStorage();
  }
};

// Fungsi untuk menghapus item
const removeItem = async (item, index) => {
  if (isAuthenticated.value) {
    try {
      // Pastikan item.id ada dan bukan undefined
      if (!item.id) {
        showToast('Cannot remove item: Invalid item ID', 'error');
        return;
      }
      
      await cartService.removeCartItem(item.id);
      await loadCart(); // Reload cart dari API
      showToast('Item removed from cart', 'success');
    } catch (error) {
      console.error('Error removing item:', error);
      showToast('Error removing item from cart', 'error');
    }
  } else {
    cart.value.splice(index, 1);
    saveCartToLocalStorage();
    showToast('Item removed from cart', 'success');
  }
};

// Fungsi untuk mengosongkan cart
const clearCart = async () => {
  if (isAuthenticated.value) {
    try {
      await cartService.clearCart();
      await loadCart(); // Reload cart dari API
      showToast('Cart cleared', 'success');
    } catch (error) {
      showToast('Error clearing cart', 'error');
    }
  } else {
    cart.value = [];
    saveCartToLocalStorage();
    showToast('Cart cleared', 'success');
  }
};

// Simpan cart ke localStorage untuk user yang tidak login
const saveCartToLocalStorage = () => {
  localStorage.setItem('cart', JSON.stringify(cart.value));
};

// Sinkronisasi cart localStorage ke server setelah login
const syncCartWithServer = async () => {
  if (!isAuthenticated.value) return;
  
  // Ambil cart dari localStorage
  const localCart = JSON.parse(localStorage.getItem('cart') || '[]');
  if (!localCart.length) return;
  
  // Tambahkan semua item ke cart server
  for (const item of localCart) {
    try {
      await cartService.addToCart(item.id, item.quantity);
    } catch (error) {
      console.error('Error syncing cart item to server:', error);
    }
  }
  
  // Kosongkan localStorage cart setelah sync
  localStorage.removeItem('cart');
  
  // Reload cart dari server
  await loadCart();
};

// Proses checkout
const checkout = async () => {
  if (cart.value.length === 0) {
    showToast('Your cart is empty', 'error');
    return;
  }
  
  isProcessing.value = true;
  
  try {
    // Dummy checkout process
    await new Promise(resolve => setTimeout(resolve, 1500));
    
    if (isAuthenticated.value) {
      await cartService.clearCart();
    } else {
      localStorage.removeItem('cart');
    }
    
    showToast('Order placed successfully', 'success');
    cart.value = [];
    
    // Redirect ke success page
    setTimeout(() => {
      router.push('/user/products');
    }, 1000);
  } catch (error) {
    showToast('Error processing your order', 'error');
  } finally {
    isProcessing.value = false;
  }
};

onMounted(() => {
  loadCart();
  
  // Sinkronisasi cart jika user baru login
  if (isAuthenticated.value) {
    syncCartWithServer();
  }
});
</script>

<template>
  <UserLayout>
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
      <div class="mb-6">
        <h1 class="text-2xl font-bold text-gray-900">Shopping Cart</h1>
        <p class="text-gray-500 mt-1">Review and manage your selected items</p>
      </div>
      
      <div v-if="isLoading" class="flex justify-center py-20">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600"></div>
      </div>
      
      <div v-else-if="cart.length === 0" class="bg-white rounded-xl shadow-sm border border-gray-200 p-8 text-center">
        <div class="mx-auto w-20 h-20 bg-gray-100 rounded-full flex items-center justify-center mb-4">
          <ShoppingBagIcon class="w-10 h-10 text-gray-400" />
        </div>
        <h3 class="text-lg font-medium text-gray-900">Your cart is empty</h3>
        <p class="mt-2 text-gray-500">Start shopping to add products to your cart</p>
        <Button 
          @click="router.push('/user/products')"
          class="mt-6 bg-indigo-600 text-white hover:bg-indigo-500"
        >
          Browse Products
        </Button>
      </div>
      
      <div v-else class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <!-- Cart Items -->
        <div class="lg:col-span-2 space-y-4">
          <div class="bg-white rounded-xl shadow-sm border border-gray-200 p-4">
            <div class="flex justify-between items-center mb-4">
              <h2 class="text-lg font-semibold text-gray-900">Cart Items ({{ cartItemCount }})</h2>
              <Button 
                @click="clearCart"
                variant="outline"
                class="text-red-600 border-red-200 hover:bg-red-50"
              >
                Clear Cart
              </Button>
            </div>
            
            <div class="divide-y divide-gray-200">
              <div 
                v-for="(item, index) in cart" 
                :key="index"
                class="py-4 flex items-center"
              >
                <div class="h-20 w-20 flex-shrink-0 overflow-hidden rounded-md border border-gray-200">
                  <img 
                    :src="item.product?.thumbnail || item.thumbnail || '/assets/placeholder.jpg'" 
                    :alt="item.product?.name || item.name || 'Product'"
                    class="h-full w-full object-cover object-center"
                  />
                </div>
                
                <div class="ml-4 flex-1">
                  <h3 class="text-base font-medium text-gray-900">{{ item.product?.name || item.name || 'Unknown Product' }}</h3>
                  <p class="mt-1 text-sm text-gray-500">{{ item.product?.category || item.category || 'No Category' }}</p>
                  <p class="mt-1 text-sm font-medium text-gray-900">{{ formatCurrency(item.product?.price || item.price || 0) }}</p>
                </div>
                
                <div class="flex items-center">
                  <div class="flex items-center border border-gray-300 rounded-md">
                    <button 
                      @click="decreaseQuantity(item)"
                      class="p-2 hover:bg-gray-100"
                    >
                      <MinusIcon class="w-4 h-4" />
                    </button>
                    <span class="px-4">{{ item.quantity }}</span>
                    <button 
                      @click="increaseQuantity(item)"
                      class="p-2 hover:bg-gray-100"
                    >
                      <PlusIcon class="w-4 h-4" />
                    </button>
                  </div>
                  
                  <button 
                    @click="removeItem(item, index)"
                    class="ml-4 text-red-500 hover:text-red-700"
                  >
                    <Trash2Icon class="w-5 h-5" />
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <!-- Order Summary -->
        <div>
          <div class="bg-white rounded-xl shadow-sm border border-gray-200 p-4">
            <h2 class="text-lg font-semibold text-gray-900 mb-4">Order Summary</h2>
            
            <div class="space-y-3">
              <div class="flex justify-between text-base text-gray-600">
                <p>Subtotal</p>
                <p>{{ formatCurrency(cartTotal) }}</p>
              </div>
              
              <div class="flex justify-between text-base text-gray-600">
                <p>Shipping</p>
                <p>Free</p>
              </div>
              
              <div class="flex justify-between text-base text-gray-600">
                <p>Tax</p>
                <p>{{ formatCurrency(cartTotal * 0.1) }}</p>
              </div>
              
              <div class="border-t border-gray-200 pt-3 flex justify-between text-base font-semibold text-gray-900">
                <p>Total</p>
                <p>{{ formatCurrency(cartTotal + (cartTotal * 0.1)) }}</p>
              </div>
            </div>
            
            <Button 
              @click="checkout"
              :disabled="isProcessing"
              class="mt-6 w-full bg-indigo-600 text-white hover:bg-indigo-500"
            >
              <span v-if="isProcessing">Processing...</span>
              <span v-else>Checkout</span>
            </Button>
            
            <div class="mt-6 flex items-center text-sm text-gray-500">
              <AlertTriangleIcon class="w-4 h-4 mr-2 text-yellow-500" />
              <p>This is a demo. No actual purchases will be made.</p>
            </div>
          </div>
          
          <div class="mt-4">
            <Button 
              @click="router.push('/user/products')"
              variant="outline"
              class="w-full"
            >
              Continue Shopping
            </Button>
          </div>
        </div>
      </div>
    </div>
  </UserLayout>
</template>