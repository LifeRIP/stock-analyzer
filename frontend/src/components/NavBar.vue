<script setup lang="ts">
import { ref } from "vue";
import { useStockStore } from "../stores/stockStore";
import { Bars3Icon, XMarkIcon } from "@heroicons/vue/24/outline";

const stockStore = useStockStore();
const isMobileMenuOpen = ref(false);
const isLoading = ref(false);

const toggleMobileMenu = () => {
  isMobileMenuOpen.value = !isMobileMenuOpen.value;
};

const syncStocks = async () => {
  isLoading.value = true;
  await stockStore.syncStocks();
  isLoading.value = false;
};

const syncAndCloseMenu = async () => {
  await syncStocks();
  isMobileMenuOpen.value = false;
};
</script>

<template>
  <nav class="bg-white shadow-md">
    <div class="container mx-auto px-4">
      <div class="flex justify-between items-center h-16">
        <div class="flex items-center">
          <router-link to="/" class="flex items-center">
            <span class="text-xl font-bold text-green-600">StockAnalyzer</span>
          </router-link>
        </div>

        <div class="hidden md:block">
          <div class="flex items-center space-x-4">
            <router-link to="/" class="px-3 py-2 rounded-md text-sm font-medium hover:bg-gray-100">
              Stocks
            </router-link>
            <router-link
              to="/recommendations"
              class="px-3 py-2 rounded-md text-sm font-medium hover:bg-gray-100"
            >
              Recommendations
            </router-link>
            <button
              @click="syncStocks"
              class="px-3 py-2 rounded-md text-sm font-medium bg-green-600 text-white hover:bg-green-700"
              :class="{ 'cursor-progress': isLoading, 'cursor-pointer': !isLoading }"
              :disabled="isLoading"
            >
              <span v-if="isLoading">Syncing...</span>
              <span v-else>Sync Stocks</span>
            </button>
          </div>
        </div>

        <div class="md:hidden">
          <button
            @click="toggleMobileMenu"
            class="text-gray-500 hover:text-gray-700 focus:outline-none"
          >
            <Bars3Icon v-if="!isMobileMenuOpen" class="h-6 w-6" />
            <XMarkIcon v-else class="h-6 w-6" />
          </button>
        </div>
      </div>
    </div>

    <div v-if="isMobileMenuOpen" class="md:hidden">
      <div class="px-2 pt-2 pb-3 space-y-1 sm:px-3">
        <router-link
          to="/"
          class="block px-3 py-2 rounded-md text-base font-medium hover:bg-gray-100"
          @click="isMobileMenuOpen = false"
        >
          Stocks
        </router-link>
        <router-link
          to="/recommendations"
          class="block px-3 py-2 rounded-md text-base font-medium hover:bg-gray-100"
          @click="isMobileMenuOpen = false"
        >
          Recommendations
        </router-link>
        <button
          @click="syncAndCloseMenu"
          class="w-full text-left px-3 py-2 rounded-md text-base font-medium bg-green-600 text-white hover:bg-green-700"
          :disabled="isLoading"
        >
          <span v-if="isLoading">Syncing...</span>
          <span v-else>Sync Stocks</span>
        </button>
      </div>
    </div>
  </nav>
</template>
