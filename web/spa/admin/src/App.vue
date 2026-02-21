<template>
  ADMIN

  <router-link :to="{ name: 'main' }">
    MAIN PAGE
  </router-link>
  <router-link :to="{ name: 'test' }">
    TEST PAGE
  </router-link>

  <router-view />
</template>

<script lang="ts" setup>

</script>
