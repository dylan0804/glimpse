<template>
  <input
    v-model="hey"
    class="hey"
  />
  <button @click="init">init</button>
</template>

<script lang="ts" setup>
    import { ref } from 'vue';
    import { ScanScreenshots } from "../../wailsjs/go/main/App.js";

    const hey = ref('')
    
    async function init() {
      try {
        const result = await ScanScreenshots()
        console.log(result)
      } catch (e) {
        console.log(e)
      }
    }

</script>

<style scoped>
</style>