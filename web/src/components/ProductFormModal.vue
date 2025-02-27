<script setup>
import { ref, watch } from 'vue';
import { useToast } from '@/composables/useToast';
import { Money3Component } from 'v-money3';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { XIcon } from 'lucide-vue-next';

const props = defineProps({
  show: {
    type: Boolean,
    default: true
  },
  product: {
    type: Object,
    default: null
  },
  categories: {
    type: Array,
    default: () => []
  }
});

const emit = defineEmits(['close', 'submit']);
const { showToast } = useToast();

const form = ref({
  name: '',
  category: '',
  price: '',
  thumbnail: '',
  image_link: ''
});

const isSubmitting = ref(false);
const error = ref('');

// Money configuration
const moneyConfig = {
  decimal: ',',
  thousands: '.',
  prefix: 'Rp ',
  suffix: '',
  precision: 0,
  masked: false,
  disableNegative: true,
  disabled: false,
  min: 0,
  max: null,
  allowBlank: false,
  minimumNumberOfCharacters: 0,
  modelModifiers: { lazy: true },
  debounce: 0,
};

// Watch for product changes (edit mode)
watch(() => props.product, (newProduct) => {
  if (newProduct) {
    form.value = {
      name: newProduct.name || '',
      category: newProduct.category || '',
      price: newProduct.price ? `Rp ${newProduct.price}` : '',
      thumbnail: newProduct.thumbnail || '',
      image_link: newProduct.image_link || ''
    };
  } else {
    // Reset form for new product
    form.value = {
      name: '',
      category: '',
      price: '',
      thumbnail: '',
      image_link: ''
    };
  }
}, { immediate: true });

const handleSubmit = async () => {
  try {
    if (!form.value.name) {
      error.value = 'Product name is required';
      return;
    }
    
    if (!form.value.category) {
      error.value = 'Category is required';
      return;
    }
    
    // Process price from money format
    const rawPrice = form.value.price?.toString().replace('Rp ', '').replace(/\./g, '');
    const price = Number(rawPrice);
    
    if (!price || price <= 0) {
      error.value = 'Price must be greater than 0';
      return;
    }
    
    isSubmitting.value = true;
    error.value = '';
    
    const productData = {
      name: form.value.name,
      category: form.value.category,
      price: price,
      thumbnail: form.value.thumbnail,
      image_link: form.value.image_link
    };
    
    emit('submit', productData);
  } catch (err) {
    error.value = err.message || 'Error saving product';
    showToast(error.value, 'error');
  } finally {
    isSubmitting.value = false;
  }
};

const closeModal = () => {
  if (!isSubmitting.value) {
    emit('close');
  }
};
</script>

<template>
  <div class="p-4">
    <form @submit.prevent="handleSubmit" class="space-y-4">
      <!-- Product Name -->
      <div class="space-y-2">
        <label class="text-sm font-medium text-gray-700">Product Name</label>
        <Input
          v-model="form.name"
          type="text"
          required
          placeholder="Enter product name"
          class="w-full"
        />
      </div>
      
      <!-- Category -->
      <div class="space-y-2">
        <label class="text-sm font-medium text-gray-700">Category</label>
        <select
          v-model="form.category"
          required
          class="w-full rounded-md border border-gray-300 p-2 outline-none focus:ring-1 focus:border-indigo-600 focus:ring-indigo-600"
        >
          <option value="" disabled>Select category</option>
          <option
            v-for="category in categories"
            :key="category"
            :value="category"
          >
            {{ category }}
          </option>
          <!-- Default categories if API fails -->
          <option v-if="categories.length === 0" value="Iphone">Iphone</option>
          <option v-if="categories.length === 0" value="Samsung">Samsung</option>
          <option v-if="categories.length === 0" value="Xiaomi">Xiaomi</option>
        </select>
      </div>
      
      <!-- Price -->
      <div class="space-y-2">
        <label class="text-sm font-medium text-gray-700">Price</label>
        <Money3Component
          v-model="form.price"
          class="w-full rounded-md border border-gray-300 p-2 outline-none focus:ring-1 focus:border-indigo-600 focus:ring-indigo-600"
          v-bind="moneyConfig"
          required
        />
      </div>
      
      <!-- Thumbnail -->
      <div class="space-y-2">
        <label class="text-sm font-medium text-gray-700">Thumbnail URL</label>
        <Input
          v-model="form.thumbnail"
          type="text"
          placeholder="Enter thumbnail URL"
          class="w-full"
        />
      </div>
      
      <!-- Image Link -->
      <div class="space-y-2">
        <label class="text-sm font-medium text-gray-700">Image Link</label>
        <Input
          v-model="form.image_link"
          type="text"
          placeholder="Enter image URL"
          class="w-full"
        />
      </div>
      
      <!-- Error message -->
      <div v-if="error" class="p-3 bg-red-100 text-red-700 rounded-md">
        {{ error }}
      </div>
      
      <!-- Buttons -->
      <div class="flex justify-end gap-3 pt-4">
        <Button
          type="button"
          @click="closeModal"
          variant="outline"
          class="hover:bg-gray-50"
        >
          Cancel
        </Button>
        <Button
          type="submit"
          :disabled="isSubmitting"
          class="bg-indigo-600 text-white hover:bg-indigo-500 disabled:opacity-50"
        >
          {{ isSubmitting ? 'Saving...' : (product ? 'Update' : 'Add') }}
        </Button>
      </div>
    </form>
  </div>
</template>