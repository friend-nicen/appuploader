<template>
  <a-layout-header
    class="v-header"
  >
    <div class="v-header__left">
      <h1 class="v-header__title">{{ title }}</h1>
    </div>
    <div class="v-header__right">
      <KeySelector />
    </div>
  </a-layout-header>
</template>

<script setup>
import KeySelector from '@/components/v-key-select.vue';

defineProps({
  title: {
    type: String,
    default: ''
  }
});
</script>

<style lang="scss" scoped>
.v-header {
  background: #fff;
  padding: 0 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid #f0f0f0;
  height: 64px;

  &__left {
    display: flex;
    align-items: center;
  }

  &__title {
    font-size: 1.25rem;
    font-weight: 600;
    color: #1f2937;
    margin: 0;
  }

  &__right {
    display: flex;
    align-items: center;
  }
}
</style>
