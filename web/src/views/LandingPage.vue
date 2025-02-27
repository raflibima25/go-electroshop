<script setup>
import { ref } from 'vue'
import { 
  ShoppingCartIcon,
  ArchiveBoxIcon, // Pengganti Package
  CreditCardIcon, // Pengganti CreditCard
  TruckIcon, // Pengganti Truck
  DevicePhoneMobileIcon, // Pengganti Smartphone
  ComputerDesktopIcon, // Pengganti Laptop
  SpeakerWaveIcon, // Pengganti Headphones
  PhotoIcon, // Pengganti Camera
  ClockIcon, // Pengganti Watch
  Bars3Icon,
  XMarkIcon 
} from '@heroicons/vue/24/outline'
// logo electroshop
import logoelectroshop from '../assets/logo-electroshop.webp'

const mobileMenuOpen = ref(false)

// Navigation data
const navigation = [
  { name: 'Products', href: '#' },
  { name: 'Categories', href: '#' },
  { name: 'Deals', href: '#' },
  { name: 'Support', href: '#' },
]

// Features data
const features = [
  {
    name: 'Wide Selection',
    description: 'Browse through thousands of the latest electronic devices and accessories from premium brands.',
    icon: ArchiveBoxIcon,
  },
  {
    name: 'Secure Payments',
    description: 'Shop with confidence using our secure payment processing systems with multiple payment options.',
    icon: CreditCardIcon,
  },
  {
    name: 'Fast Shipping',
    description: 'Get your electronic products delivered quickly and safely with our reliable shipping partners.',
    icon: TruckIcon,
  },
  {
    name: 'Technical Support',
    description: 'Our technical experts are ready to assist you with any questions about your electronic purchases.',
    icon: SpeakerWaveIcon,
  },
]

// Product categories
const categories = [
  {
    name: 'Smartphones',
    icon: DevicePhoneMobileIcon,
    description: 'Latest models from Apple, Samsung, Xiaomi and more'
  },
  {
    name: 'Laptops',
    icon: ComputerDesktopIcon,
    description: 'Performance laptops for work and gaming'
  },
  {
    name: 'Audio Devices',
    icon: SpeakerWaveIcon,
    description: 'Premium headphones, earbuds and speakers'
  },
  {
    name: 'Cameras',
    icon: PhotoIcon,
    description: 'DSLR, mirrorless and action cameras'
  },
  {
    name: 'Smartwatches',
    icon: ClockIcon,
    description: 'Fitness trackers and smartwatches'
  }
]

// Featured products
const featuredProducts = [
  {
    name: 'iPhone 13 Pro',
    price: 12000000,
    category: 'Smartphone',
    image: 'https://placehold.co/300x300/e2e8f0/1e293b?text=iPhone+13+Pro'
  },
  {
    name: 'Samsung Galaxy Z Flip',
    price: 20000000,
    category: 'Smartphone',
    image: 'https://placehold.co/300x300/e2e8f0/1e293b?text=Samsung+Z+Flip'
  },
  {
    name: 'Xiaomi Redmi Note 11 Pro',
    price: 3200000,
    category: 'Smartphone',
    image: 'https://placehold.co/300x300/e2e8f0/1e293b?text=Xiaomi+Redmi'
  }
]

// Footer navigation
const footerNavigation = {
  products: [
    { name: 'Smartphones', href: '#' },
    { name: 'Laptops', href: '#' },
    { name: 'Audio', href: '#' },
    { name: 'Accessories', href: '#' },
  ],
  support: [
    { name: 'Pricing', href: '#' },
    { name: 'Documentation', href: '#' },
    { name: 'Guides', href: '#' },
    { name: 'Warranty', href: '#' },
  ],
  company: [
    { name: 'About', href: '#' },
    { name: 'Blog', href: '#' },
    { name: 'Jobs', href: '#' },
    { name: 'Press', href: '#' },
    { name: 'Partners', href: '#' },
  ],
  legal: [
    { name: 'Claim', href: '#' },
    { name: 'Privacy', href: '#' },
    { name: 'Terms', href: '#' },
  ],
}

// FAQ data
const faqs = [
  {
    question: "Do you offer international shipping?",
    answer: "Yes, we ship to most countries worldwide. Shipping rates and delivery times vary based on your location."
  },
  {
    question: "What is your return policy?",
    answer: "We offer a 30-day return policy for most products. Items must be returned in original packaging with all accessories."
  },
  {
    question: "Are all products covered by warranty?",
    answer: "Yes, all our products come with the manufacturer's warranty. Additional extended warranty options are available for purchase."
  },
  {
    question: "Do you offer technical support for products?",
    answer: "Yes, our technical support team is available to help you with setup, troubleshooting, and usage questions."
  },
  {
    question: "How can I track my order?",
    answer: "Once your order ships, you'll receive a tracking number via email that you can use to monitor your delivery status."
  }
];

// Get the current year
const currentYear = new Date().getFullYear();

// Format price as currency
const formatPrice = (price) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    maximumFractionDigits: 0
  }).format(price);
}
</script>

<template>
  <div class="bg-white">
    <!-- Header -->
    <header class="absolute inset-x-0 top-0 z-50">
      <nav class="flex items-center justify-between p-6 lg:px-8" aria-label="Global">
        <!-- Logo -->
        <div class="flex lg:flex-1">
          <a href="#" class="-m-1.5 p-1.5">
            <span class="sr-only">ElectroShop</span>
            <img class="h-10 w-auto" :src="logoelectroshop" alt="ElectroShop Logo" />
          </a>
        </div>

        <!-- Mobile menu button -->
        <div class="flex lg:hidden">
          <button 
            type="button" 
            class="-m-2.5 inline-flex items-center justify-center rounded-md p-2.5 text-gray-700"
            @click="mobileMenuOpen = true"
          >
            <span class="sr-only">Open main menu</span>
            <Bars3Icon class="h-6 w-6" aria-hidden="true" />
          </button>
        </div>

        <!-- Desktop navigation -->
        <div class="hidden lg:flex lg:gap-x-12">
          <a v-for="item in navigation" 
             :key="item.name" 
             :href="item.href" 
             class="text-sm font-semibold leading-6 text-gray-900"
          >
            {{ item.name }}
          </a>
        </div>

        <!-- Desktop login button -->
        <div class="hidden lg:flex lg:flex-1 lg:justify-end">
          <a href="/login" class="text-sm font-semibold leading-6 text-gray-900">
            Login <span aria-hidden="true">→</span>
          </a>
        </div>
      </nav>

      <!-- Mobile menu -->
      <div v-if="mobileMenuOpen" class="lg:hidden">
        <div class="fixed inset-0 z-50"></div>
        <div class="fixed inset-y-0 right-0 z-50 w-full overflow-y-auto bg-white px-6 py-6 sm:max-w-sm sm:ring-1 sm:ring-gray-900/10">
          <div class="flex items-center justify-between">
            <a href="#" class="-m-1.5 p-1.5">
              <span class="sr-only">ElectroShop</span>
              <img class="h-10 w-auto" :src="logoelectroshop" alt="ElectroShop Logo" />
            </a>
            <button 
              type="button" 
              class="-m-2.5 rounded-md p-2.5 text-gray-700"
              @click="mobileMenuOpen = false"
            >
              <span class="sr-only">Close menu</span>
              <XMarkIcon class="h-6 w-6" aria-hidden="true" />
            </button>
          </div>
          <div class="mt-6 flow-root">
            <div class="-my-6 divide-y divide-gray-500/10">
              <div class="space-y-2 py-6">
                <a v-for="item in navigation" 
                   :key="item.name" 
                   :href="item.href" 
                   class="-mx-3 block rounded-lg px-3 py-2 text-base font-semibold leading-7 text-gray-900 hover:bg-gray-50"
                >
                  {{ item.name }}
                </a>
              </div>
              <div class="py-6">
                <a href="/login" class="-mx-3 block rounded-lg px-3 py-2.5 text-base font-semibold leading-7 text-gray-900 hover:bg-gray-50">
                  Login
                </a>
              </div>
            </div>
          </div>
        </div>
      </div>
    </header>

    <!-- Hero section -->
    <div class="relative isolate pt-14">
      <div class="absolute inset-x-0 -top-40 -z-10 transform-gpu overflow-hidden blur-3xl sm:-top-80">
        <div 
          class="relative left-[calc(50%-11rem)] aspect-[1155/678] w-[36.125rem] -translate-x-1/2 rotate-[30deg] bg-gradient-to-tr from-[#ff80b5] to-[#9089fc] opacity-30 sm:left-[calc(50%-30rem)] sm:w-[72.1875rem]" 
          style="clip-path: polygon(74.1% 44.1%, 100% 61.6%, 97.5% 26.9%, 85.5% 0.1%, 80.7% 2%, 72.5% 32.5%, 60.2% 62.4%, 52.4% 68.1%, 47.5% 58.3%, 45.2% 34.5%, 27.5% 76.7%, 0.1% 64.9%, 17.9% 100%, 27.6% 76.8%, 76.1% 97.7%, 74.1% 44.1%)"
        />
      </div>
      <div class="py-24 sm:py-32 lg:pb-40">
        <div class="mx-auto max-w-7xl px-6 lg:px-8">
          <div class="mx-auto max-w-2xl text-center">
            <h1 class="text-4xl font-bold tracking-tight text-gray-900 sm:text-6xl">
              Premium Electronics at Your Fingertips
            </h1>
            <p class="mt-6 text-lg leading-8 text-gray-600">
              Discover the latest smartphones, laptops, and electronic accessories with exceptional quality and competitive prices.
            </p>
            <div class="mt-10 flex items-center justify-center gap-x-6">
              <a href="/login" class="rounded-md bg-indigo-600 px-3.5 py-2.5 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600">
                Get Started
              </a>
              <a href="#" class="text-sm font-semibold leading-6 text-gray-900">
                Learn more <span aria-hidden="true">→</span>
              </a>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Features section -->
    <div class="mx-auto max-w-7xl px-6 sm:mt-24 lg:px-8">
      <div class="mx-auto max-w-2xl text-center">
        <h2 class="text-base font-semibold leading-7 text-indigo-600">Why Choose Us</h2>
        <p class="mt-2 text-3xl font-bold tracking-tight text-gray-900 sm:text-4xl">
          Everything you need for your electronics shopping
        </p>
        <p class="mt-6 text-lg leading-8 text-gray-600">
          We offer a seamless shopping experience with the best selection of electronic products and exceptional customer service.
        </p>
      </div>
      <div class="mx-auto mt-16 max-w-2xl sm:mt-20 lg:mt-24 lg:max-w-none">
        <dl class="grid max-w-xl grid-cols-1 gap-x-8 gap-y-16 lg:max-w-none lg:grid-cols-4">
          <div v-for="feature in features" :key="feature.name" class="flex flex-col">
            <dt class="flex items-center gap-x-3 text-base font-semibold leading-7 text-gray-900">
              <component :is="feature.icon" class="h-5 w-5 flex-none text-indigo-600" aria-hidden="true" />
              {{ feature.name }}
            </dt>
            <dd class="mt-4 flex flex-auto flex-col text-base leading-7 text-gray-600">
              <p class="flex-auto">{{ feature.description }}</p>
            </dd>
          </div>
        </dl>
      </div>
    </div>

    <!-- Product Categories section -->
    <div class="mx-auto max-w-7xl px-6 py-24 sm:py-32">
      <div class="mx-auto max-w-2xl text-center">
        <h2 class="text-base font-semibold leading-7 text-indigo-600">Explore</h2>
        <p class="mt-2 text-3xl font-bold tracking-tight text-gray-900 sm:text-4xl">
          Product Categories
        </p>
        <p class="mt-6 text-lg leading-8 text-gray-600">
          Browse our wide range of electronic products across various categories
        </p>
      </div>

      <div class="mx-auto mt-16 grid max-w-2xl grid-cols-1 gap-x-8 gap-y-16 sm:grid-cols-2 lg:mx-0 lg:max-w-none lg:grid-cols-5">
        <div v-for="category in categories" :key="category.name" class="group relative">
          <div class="aspect-square overflow-hidden rounded-lg bg-gray-100 flex items-center justify-center">
            <component :is="category.icon" class="h-20 w-20 text-indigo-600 group-hover:text-indigo-500 transition-colors" />
          </div>
          <h3 class="mt-6 text-lg font-semibold leading-8 tracking-tight text-gray-900">{{ category.name }}</h3>
          <p class="text-base leading-7 text-gray-600">{{ category.description }}</p>
          <a href="#" class="mt-4 text-sm font-medium text-indigo-600 group-hover:text-indigo-500">
            View products <span aria-hidden="true">→</span>
          </a>
        </div>
      </div>
    </div>

    <!-- Featured Products -->
    <div class="bg-gray-50 py-24 sm:py-32">
      <div class="mx-auto max-w-7xl px-6 lg:px-8">
        <div class="mx-auto max-w-2xl text-center">
          <h2 class="text-base font-semibold leading-7 text-indigo-600">Featured Products</h2>
          <p class="mt-2 text-3xl font-bold tracking-tight text-gray-900 sm:text-4xl">Our Best Sellers</p>
          <p class="mt-6 text-lg leading-8 text-gray-600">
            Check out our most popular electronic devices loved by our customers
          </p>
        </div>

        <div class="mx-auto mt-16 grid max-w-2xl grid-cols-1 gap-x-8 gap-y-16 sm:grid-cols-2 lg:mx-0 lg:max-w-none lg:grid-cols-3">
          <div v-for="product in featuredProducts" :key="product.name" class="flex flex-col overflow-hidden rounded-lg shadow-lg bg-white">
            <div class="flex-shrink-0">
              <img class="h-96 w-full object-cover" :src="product.image" :alt="product.name">
            </div>
            <div class="flex flex-1 flex-col justify-between p-6">
              <div class="flex-1">
                <p class="text-sm font-medium text-indigo-600">{{ product.category }}</p>
                <h3 class="mt-2 text-xl font-semibold text-gray-900">{{ product.name }}</h3>
                <p class="mt-4 text-2xl font-bold text-gray-900">{{ formatPrice(product.price) }}</p>
              </div>
              <div class="mt-6">
                <button type="button" class="w-full rounded-md bg-indigo-600 px-3.5 py-2.5 text-center text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600">
                  View Details
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- FAQ section -->
    <div class="mx-auto max-w-2xl divide-y divide-gray-900/10 px-6 pb-8 sm:pb-24 sm:pt-12 lg:max-w-7xl lg:px-8 lg:pb-32">
      <h2 class="text-2xl font-bold leading-10 tracking-tight text-gray-900">Frequently asked questions</h2>
      <dl class="mt-10 space-y-8 divide-y divide-gray-900/10">
        <div v-for="faq in faqs" :key="faq.question" class="pt-8 lg:grid lg:grid-cols-12 lg:gap-8">
          <dt class="text-base font-semibold leading-7 text-gray-900 lg:col-span-5">{{ faq.question }}</dt>
          <dd class="mt-4 lg:col-span-7 lg:mt-0">
            <p class="text-base leading-7 text-gray-600">{{ faq.answer }}</p>
          </dd>
        </div>
      </dl>
    </div>

    <!-- CTA section -->
    <div class="relative isolate px-6 py-32 sm:mt-24 sm:py-40 lg:px-8">
      <div class="absolute inset-x-0 top-0 -z-10 transform-gpu overflow-hidden blur-3xl sm:top-[-20rem]">
        <div 
          class="relative left-1/2 -z-10 aspect-[1155/678] w-[36.125rem] -translate-x-1/2 bg-gradient-to-tr from-[#ff80b5] to-[#9089fc] opacity-30 sm:w-[72.1875rem]" 
          style="clip-path: polygon(74.1% 44.1%, 100% 61.6%, 97.5% 26.9%, 85.5% 0.1%, 80.7% 2%, 72.5% 32.5%, 60.2% 62.4%, 52.4% 68.1%, 47.5% 58.3%, 45.2% 34.5%, 27.5% 76.7%, 0.1% 64.9%, 17.9% 100%, 27.6% 76.8%, 76.1% 97.7%, 74.1% 44.1%)"
        />
      </div>
      <div class="mx-auto max-w-2xl text-center">
        <h2 class="text-3xl font-bold tracking-tight text-gray-900 sm:text-4xl">
          Access the Admin Dashboard.<br>Manage Your Electronic Store.
        </h2>
        <p class="mx-auto mt-6 max-w-xl text-lg leading-8 text-gray-600">
          Log in to the admin portal to manage products, view analytics, and control all aspects of your electronic store.
        </p>
        <div class="mt-10 flex items-center justify-center gap-x-6">
          <a href="/login" class="rounded-md bg-indigo-600 px-3.5 py-2.5 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600">
            Login
          </a>
          <a href="#" class="text-sm font-semibold leading-6 text-gray-900">
            Learn more <span aria-hidden="true">→</span>
          </a>
        </div>
      </div>
    </div>

    <!-- Footer -->
    <footer class="bg-white" aria-labelledby="footer-heading">
      <h2 id="footer-heading" class="sr-only">Footer</h2>
      <div class="mx-auto max-w-7xl px-6 pb-8 pt-16 sm:pt-24 lg:px-8 lg:pt-32">
        <div class="xl:grid xl:grid-cols-3 xl:gap-8">
          <div class="space-y-8">
            <img class="h-7" :src="logoelectroshop" alt="ElectroShop" />
            <p class="text-sm leading-6 text-gray-600">
              Your premier destination for quality electronics and tech accessories.
            </p>
          </div>
          <div class="mt-16 grid grid-cols-2 gap-8 xl:col-span-2 xl:mt-0">
            <div class="md:grid md:grid-cols-2 md:gap-8">
              <div>
                <h3 class="text-sm font-semibold leading-6 text-gray-900">Products</h3>
                <ul role="list" class="mt-6 space-y-4">
                  <li v-for="item in footerNavigation.products" :key="item.name">
                    <a :href="item.href" class="text-sm leading-6 text-gray-600 hover:text-gray-900">
                      {{ item.name }}
                    </a>
                  </li>
                </ul>
              </div>
              <div class="mt-10 md:mt-0">
                <h3 class="text-sm font-semibold leading-6 text-gray-900">Support</h3>
                <ul role="list" class="mt-6 space-y-4">
                  <li v-for="item in footerNavigation.support" :key="item.name">
                    <a :href="item.href" class="text-sm leading-6 text-gray-600 hover:text-gray-900">
                      {{ item.name }}
                    </a>
                  </li>
                </ul>
              </div>
            </div>
            <div class="md:grid md:grid-cols-2 md:gap-8">
              <div>
                <h3 class="text-sm font-semibold leading-6 text-gray-900">Company</h3>
                <ul role="list" class="mt-6 space-y-4">
                  <li v-for="item in footerNavigation.company" :key="item.name">
                    <a :href="item.href" class="text-sm leading-6 text-gray-600 hover:text-gray-900">
                      {{ item.name }}
                    </a>
                  </li>
                </ul>
              </div>
              <div class="mt-10 md:mt-0">
                <h3 class="text-sm font-semibold leading-6 text-gray-900">Legal</h3>
                <ul role="list" class="mt-6 space-y-4">
                  <li v-for="item in footerNavigation.legal" :key="item.name">
                    <a :href="item.href" class="text-sm leading-6 text-gray-600 hover:text-gray-900">
                      {{ item.name }}
                    </a>
                  </li>
                </ul>
              </div>
            </div>
          </div>
        </div>
        <div class="mt-16 border-t border-gray-900/10 pt-8 sm:mt-20 lg:mt-24">
          <p class="text-xs leading-5 text-gray-500">
            &copy; <span>{{ currentYear }}</span> ElectroShop, Inc. All rights reserved.
          </p>
        </div>
      </div>
    </footer>
  </div>
</template>