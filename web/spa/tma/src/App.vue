<template>
  <h1 class="block">
    TMA
  </h1>

  <ul class="flex flex-col gap-2">
    <li>
      <router-link :to="{ name: 'main' }">
        Main
      </router-link>
    </li>
    <li>
      <router-link :to="{ name: 'test' }">
        Test
      </router-link>
    </li>
  </ul>

  <router-view />
</template>

<script lang="ts" setup>

</script>
