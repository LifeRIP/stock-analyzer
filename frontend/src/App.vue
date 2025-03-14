<template>
  <div class="min-h-screen bg-gray-50">
    <Navbar />
    <main class="container mx-auto px-4 py-6">
      <router-view />
    </main>
    <footer class="bg-white shadow-inner py-4 mt-8">
      <div class="container mx-auto px-4 text-center text-gray-500">
        &copy; 2025 Stock Analyzer - All rights reserved
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import Navbar from './components/Navbar.vue'
</script>

