<script setup>
import { watch, onMounted } from 'vue';
import { useAuth } from '@/composables/useAuth';
import { cartService } from '@/services/cartService';
import { useToast } from '@/composables/useToast';

const { isAuthenticated } = useAuth();
const { showToast } = useToast();

const syncCartWithServer = async () => {
  if (!isAuthenticated.value) return;
  
  try {
    const localCart = JSON.parse(localStorage.getItem('cart') || '[]');
    if (!localCart.length) return;
    
    let syncCount = 0;
    
    // Add all items to server cart
    for (const item of localCart) {
      await cartService.addToCart(item.id, item.quantity);
      syncCount += item.quantity;
    }
    
    // Clear local cart
    localStorage.removeItem('cart');
    
    if (syncCount > 0) {
      showToast(`${syncCount} items synchronized to your account`, 'success');
    }
  } catch (error) {
    console.error('Error syncing cart:', error);
  }
};

// Watch auth state changes
watch(isAuthenticated, (newValue, oldValue) => {
  // Only sync when transitioning from logged out to logged in
  if (newValue === true && oldValue === false) {
    syncCartWithServer();
  }
}, { immediate: false });

onMounted(() => {
  // Initial sync if already logged in
  if (isAuthenticated.value) {
    syncCartWithServer();
  }
});
</script>

<template>
  <div style="display: none;"></div>
</template>