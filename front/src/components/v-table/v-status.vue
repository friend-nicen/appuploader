<template>
  <template v-if="status === 1">
    <a-tag color="green">正常</a-tag>
  </template>
  <template v-if="status === 2">
    <a-tag color="red">已离职</a-tag>
  </template>
</template>

<script setup>
defineProps(['status'])
</script>