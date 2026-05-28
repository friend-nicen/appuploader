<template>
  <a-layout-sider
    theme="dark"
    class="v-sidebar"
    :width="256"
  >
    <div class="v-sidebar__logo">
      <span class="v-sidebar__logo-text">AppStore Connect</span>
    </div>

    <a-menu
      :selectedKeys="selectedKeys"
      theme="dark"
      mode="inline"
      @click="handleMenuClick"
      class="v-sidebar__menu"
    >
      <a-menu-item v-for="item in data" :key="item.path">
        {{ item.name }}
      </a-menu-item>
    </a-menu>
  </a-layout-sider>
</template>

<script setup>
import { computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import router from '@/router';

const route = useRoute();
const routerPush = useRouter();

const data = computed(() => {
  return router.getRoutes()
    .filter(route => route.meta?.displayName)
    .map(route => ({
      name: route.meta.displayName,
      path: route.path
    }));
});

const selectedKeys = computed(() => [route.path]);

const handleMenuClick = ({ key }) => {
  routerPush.push(key);
};
</script>

<style lang="scss" scoped>
.v-sidebar {
  overflow: auto;
  height: 100vh;
  position: fixed;
  left: 0;
  top: 0;
  bottom: 0;

  &__logo {
    height: 64px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-bottom: 1px solid #333;
  }

  &__logo-text {
    color: white;
    font-size: 1.25rem;
    font-weight: bold;
  }

  &__menu {
    margin-top: 16px;
    border-right: 0;
  }
}
</style>
