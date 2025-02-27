<template>
  <div id="app" class="font-poppins ">
    <router-view></router-view>
    <ToastContainer />
    <CartSyncManager />
  </div>
</template>

<script setup>
import ToastContainer from "@/components/ToastContainer.vue";
import CartSyncManager from "@/components/CartSyncManager.vue";
</script>