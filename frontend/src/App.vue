<script setup lang="ts">
import Navbar from "./components/NavBar.vue";
</script>

<template>
  <div class="min-h-screen bg-gray-50 flex flex-col">
    <Navbar />
    <main class="container mx-auto px-4 py-6 flex-grow">
      <router-view />
    </main>
    <footer class="bg-white shadow-inner py-4 mt-auto">
      <div class="container mx-auto px-4 text-center text-gray-500">
        &copy; 2025 Stock Analyzer - All rights reserved
      </div>
    </footer>
  </div>
</template>
