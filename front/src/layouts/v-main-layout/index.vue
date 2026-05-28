<template>
  <a-layout class="v-main-layout">
    <VSidebar/>

    <a-layout class="v-main-layout__body">
      <VHeader :title="currentRouteName"/>

      <a-layout-content class="v-main-layout__content">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component"/>
          </transition>
        </router-view>
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script setup>
import {computed} from 'vue';
import {useRoute} from 'vue-router';
import VSidebar from '@/layouts/v-sidebar/index.vue';
import VHeader from '@/layouts/v-header/index.vue';

const route = useRoute();

const currentRouteName = computed(() => {
  return route.meta?.displayName || route.name;
});
</script>

<style lang="scss" scoped>
.v-main-layout {
  min-height: 100vh;

  &__body {
    margin-left: 256px;
    transition: all 0.2s;
    min-height: 100vh;
    display: flex;
    flex-direction: column;
  }

  &__content {
    background: transparent;
    min-height: 280px;
    flex: 1;
    overflow: auto;
  }
}
</style>
