<template>
  <div class="min-h-screen bg-gradient-to-br from-white to-blue-50 p-6 flex justify-center font-sans">
    <div class="w-full max-w-4xl">
      <header class="mb-10">
        <h1 class="text-3xl sm:text-4xl font-light text-gray-800 tracking-tight">
          <span class="text-blue-600 font-medium">Glimpse</span>
        </h1>
        <p class="text-gray-500 mt-2">Find text in your screenshots instantly</p>
      </header>

      <div class="bg-white rounded-xl shadow-sm border border-gray-100 p-6 mb-8">
        <button
          @click="scan"
          class="mb-6 px-5 py-2.5 bg-blue-500 hover:bg-blue-600 text-white text-sm font-medium rounded-md inline-flex items-center gap-2 transition-colors"
          :disabled="isScanning"
        >
          <svg
            v-if="!isScanning"
            xmlns="http://www.w3.org/2000/svg"
            class="h-4 w-4"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14"
            />
          </svg>
          <svg 
            v-else 
            class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" 
            xmlns="http://www.w3.org/2000/svg" 
            fill="none" 
            viewBox="0 0 24 24"
          >
            <circle
              class="opacity-25"
              cx="12"
              cy="12"
              r="10"
              stroke="currentColor"
              stroke-width="4"
            />
            <path
              class="opacity-75"
              fill="currentColor"
              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
            />
          </svg>
          {{ isScanning ? 'Scanning...' : 'Scan Screenshot Library' }}
        </button>

        <div class="relative">
          <input
            v-model="searchQuery"
            @keyup.enter="search"
            class="w-full px-4 py-3 pr-12 text-gray-700 border border-gray-200 rounded-lg focus:ring-2 focus:ring-blue-200 focus:border-blue-400 focus:outline-none transition-all"
            placeholder="Search for text in screenshots..."
            :disabled="isSearching"
          />
          <button
            @click="search"
            class="absolute right-2 top-1/2 -translate-y-1/2 p-2 text-gray-400 hover:text-blue-500 transition-colors"
            :disabled="isSearching"
          >
            <svg
              v-if="!isSearching"
              xmlns="http://www.w3.org/2000/svg"
              class="h-5 w-5"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
              />
            </svg>
            <svg 
              v-else 
              class="animate-spin h-5 w-5 text-blue-500" 
              xmlns="http://www.w3.org/2000/svg" 
              fill="none" 
              viewBox="0 0 24 24"
            >
              <circle
                class="opacity-25"
                cx="12"
                cy="12"
                r="10"
                stroke="currentColor"
                stroke-width="4"
              />
              <path
                class="opacity-75"
                fill="currentColor"
                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
              />
            </svg>
          </button>
        </div>
      </div>

      <!-- Tabs -->
      <div class="mb-6 border-b border-gray-200">
        <ul class="flex flex-wrap -mb-px text-sm font-medium text-center">
          <li class="mr-2">
            <button 
              @click="activeTab = 'search'"
              :class="[
                'inline-block p-4 rounded-t-lg border-b-2',
                activeTab === 'search' 
                  ? 'text-blue-600 border-blue-600' 
                  : 'border-transparent hover:text-gray-600 hover:border-gray-300'
              ]"
              @disabled="isScanning"
            >
              Search Results ({{ searchResults.length }})
            </button>
          </li>
          <li>
            <button 
              @click="activeTab = 'scan'"
              :class="[
                'inline-block p-4 rounded-t-lg border-b-2',
                activeTab === 'scan' 
                  ? 'text-blue-600 border-blue-600' 
                  : 'border-transparent hover:text-gray-600 hover:border-gray-300'
              ]"
              @disabled="isSearching"
            >
              Scan Results ({{ scanResults.length }})
            </button>
          </li>
        </ul>
      </div>

      <!-- Tab content -->
      <div v-if="activeTab === 'search'">
        <div
          v-if="isSearching"
          class="mt-8 text-center text-gray-500 bg-white p-8 rounded-xl shadow-sm border border-gray-100 flex flex-col items-center justify-center h-60"
        >
          <svg
            class="animate-spin h-10 w-10 text-blue-500 mb-4"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
          >
            <circle
              class="opacity-25"
              cx="12"
              cy="12"
              r="10"
              stroke="currentColor"
              stroke-width="4"
            />
            <path
              class="opacity-75"
              fill="currentColor"
              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
            />
          </svg>
          <p class="text-lg font-medium">Searching...</p>
          <p class="text-sm">Please wait while we find your screenshots.</p>
        </div>
        <div
          v-else-if="searchResults.length > 0"
          class="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden"
        >
          <div class="px-6 py-4 border-b border-gray-100">
            <h2 class="text-lg font-medium text-gray-800">Search Results</h2>
          </div>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4 p-6">
            <div
              v-for="(screenshot, index) in searchResults"
              :key="index"
              class="rounded-lg overflow-hidden border border-gray-100 hover:shadow-md transition-shadow"
            >
              <img 
                :src="screenshot.url" 
                :alt="getFileName(screenshot.path)"
                class="w-full h-40 object-cover object-top"
              />
              <div class="p-3">
                <p class="text-sm text-gray-700 truncate">{{ getFileName(screenshot.path) }}</p>
                <p class="text-xs text-gray-400 truncate">{{ screenshot.path }}</p>
              </div>
            </div>
          </div>
        </div>
        <div
          v-else
          class="mt-8 text-center text-gray-500 bg-white p-8 rounded-xl shadow-sm border border-gray-100 h-60 flex flex-col items-center justify-center"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            class="h-12 w-12 mx-auto text-gray-300 mb-3"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="1"
              d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14"
            />
          </svg>
          <p class="font-medium">No Search Results</p>
          <p class="text-sm">Try searching for text in your screenshots.</p>
        </div>
      </div>

      <div v-else>
        <div
          v-if="isScanning"
          class="mt-8 text-center text-gray-500 bg-white p-8 rounded-xl shadow-sm border border-gray-100 flex flex-col items-center justify-center h-60"
        >
          <svg
            class="animate-spin h-10 w-10 text-blue-500 mb-4"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
          >
            <circle
              class="opacity-25"
              cx="12"
              cy="12"
              r="10"
              stroke="currentColor"
              stroke-width="4"
            />
            <path
              class="opacity-75"
              fill="currentColor"
              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
            />
          </svg>
          <p class="text-lg font-medium">Scanning Library...</p>
          <p class="text-sm">This might take a moment, please wait.</p>
        </div>
        <div
          v-else-if="scanResults.length > 0"
          class="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden"
        >
          <div class="px-6 py-4 border-b border-gray-100">
            <h2 class="text-lg font-medium text-gray-800">Scan Results</h2>
          </div>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4 p-6">
            <div
              v-for="(screenshot, index) in scanResults"
              :key="index"
              class="rounded-lg overflow-hidden border border-gray-100 hover:shadow-md transition-shadow"
            >
              <img 
                :src="screenshot.url" 
                :alt="getFileName(screenshot.path)"
                class="w-full h-40 object-cover object-top"
              />
              <div class="p-3">
                <p class="text-sm text-gray-700 truncate">{{ getFileName(screenshot.path) }}</p>
                <p class="text-xs text-gray-400 truncate">{{ screenshot.path }}</p>
              </div>
            </div>
          </div>
        </div>
        <div
          v-else
          class="mt-8 text-center text-gray-500 bg-white p-8 rounded-xl shadow-sm border border-gray-100 h-60 flex flex-col items-center justify-center"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            class="h-12 w-12 mx-auto text-gray-300 mb-3"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="1"
              d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14"
            />
          </svg>
          <p class="font-medium">No Scan Results</p>
          <p class="text-sm">Try scanning your screenshot library to see results here.</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
    import { ref } from 'vue';
    import { ScanScreenshots, SearchScreenshots } from "../../wailsjs/go/main/App.js";  
    import { EventsOn } from "../../wailsjs/runtime/runtime.js";
    import { SearchResult } from '../types.js';

    const searchQuery = ref('');
    const searchResults = ref<SearchResult[]>([]);
    const scanResults = ref<SearchResult[]>([]);
    const activeTab = ref('search');
    const isScanning = ref(false);
    const isSearching = ref(false);
    
    async function scan() {
      if (isScanning.value) return;
      scanResults.value = [];
      isScanning.value = true;
      activeTab.value = 'scan';

      try {
        await ScanScreenshots();
      } catch (e: unknown) {
        console.error("Scan error:", e);
      } finally {
        isScanning.value = false;
      }
    }

    async function search() {
      if (!searchQuery.value.trim() || isSearching.value) return;
      searchResults.value = [];
      isSearching.value = true;
      activeTab.value = 'search';
      
      try {
        await SearchScreenshots(searchQuery.value);
      } catch (e: unknown) {
        console.error("Search error:", e);
      } finally {
        isSearching.value = false;
      }
    }

    EventsOn("result:found", (entry: SearchResult) => {
      if (!entry || !entry.url) return;

      scanResults.value.push({
        path: entry.path,
        url: entry.url.startsWith('data:image') ? entry.url : `data:image/png;base64,${entry.url}`
      });
    });

    EventsOn("search:found", (entry: SearchResult) => {
      if (!entry || !entry.url) return;

      searchResults.value.push({
        path: entry.path,
        url: entry.url.startsWith('data:image') ? entry.url : `data:image/png;base64,${entry.url}`
      });
    });
    
    function getFileName(path: string): string {
      if (!path) return 'Unknown File';
      return path.split(/[\\/]/).pop() || path;
    }
</script>

<style scoped>
.h-60 {
  height: 15rem;
}
</style>