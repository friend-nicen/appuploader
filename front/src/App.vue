<template>
  <a-style-provider :transformers="[
      px2remTransformer({
        rootValue: 15.36,
        precision: 8,
        mediaQuery: false
      })
    ]"
  >
    <a-config-provider :locale="zhCN" :theme="theme">
      <a-layout style="min-height: 100vh">
        <!-- 左侧侧边栏 (Fixed Left Sidebar) -->
        <a-layout-sider
            theme="dark"
            :style="{ overflow: 'auto', height: '100vh', position: 'fixed', left: 0, top: 0, bottom: 0 }"
            :width="256"
        >
          <div class="h-16 flex items-center justify-center border-b border-gray-800"
               style="height: 64px; display: flex; align-items: center; justify-content: center; border-bottom: 1px solid #333;">
          <span class="text-xl font-bold text-white tracking-wider"
                style="color: white; font-size: 1.25rem; font-weight: bold;">AppStore Connect</span>
          </div>

          <a-menu
              v-model:selectedKeys="selectedKeys"
              theme="dark"
              mode="inline"
              @click="handleMenuClick"
              style="margin-top: 16px; border-right: 0;"
          >
            <a-menu-item key="/">
              Dashboard
            </a-menu-item>

            <a-menu-item key="/auth">
              API Keys
            </a-menu-item>

            <a-menu-item key="/apps">
              Apps
            </a-menu-item>

            <a-menu-item key="/bundle-ids">
              Bundle IDs
            </a-menu-item>

            <a-menu-item key="/certificates">
              Certificates
            </a-menu-item>

            <a-menu-item key="/profiles">
              Profiles
            </a-menu-item>

            <a-menu-item key="/devices">
              Devices
            </a-menu-item>

            <a-menu-item key="/testflight">
              TestFlight
            </a-menu-item>
          </a-menu>

        </a-layout-sider>

        <!-- 右侧内容区 -->
        <a-layout
            :style="{ marginLeft: '256px', transition: 'all 0.2s', minHeight: '100vh', display: 'flex', flexDirection: 'column' }">
          <!-- 顶部导航栏 -->
          <a-layout-header
              style="background: #fff; padding: 0 24px; display: flex; align-items: center; justify-content: space-between; border-bottom: 1px solid #f0f0f0; height: 64px;">
            <div style="display: flex; align-items: center;">
              <h1 style="font-size: 1.25rem; font-weight: 600; color: #1f2937; margin: 0;">{{ currentRouteName }}</h1>
            </div>
            <div style="display: flex; align-items: center;">
              <KeySelector/>
            </div>
          </a-layout-header>

          <!-- 主视图区域 -->
          <a-layout-content style="margin: 24px; background: transparent; min-height: 280px; flex: 1; overflow: auto;">
            <router-view v-slot="{ Component }">
              <transition name="fade" mode="out-in">
                <component :is="Component"/>
              </transition>
            </router-view>
          </a-layout-content>
        </a-layout>
      </a-layout>
    </a-config-provider>
  </a-style-provider>
</template>

<script setup>
import {computed, ref, watch} from 'vue';
import {useRoute, useRouter} from 'vue-router';
import KeySelector from '@/components/v-key-select.vue';

/* 引入全局中文语言包 */
import {px2remTransformer} from 'ant-design-vue';
import zhCN from 'ant-design-vue/es/locale/zh_CN';
import {dynamicTheme} from "@/theme/theme";
import dayjs from 'dayjs';
import 'dayjs/locale/zh-cn';

/* 动态主题 */
const theme = dynamicTheme();

/* dayjs中文语言包 */
dayjs.locale('zh');

const route = useRoute();
const router = useRouter();

const selectedKeys = ref([route.path]);

watch(
    () => route.path,
    (newPath) => {
      selectedKeys.value = [newPath];
    },
    {immediate: true}
);

const handleMenuClick = ({key}) => {
  router.push(key);
};

const menuItems = [
  {name: 'Dashboard', path: '/'},
  {name: 'API Keys', path: '/auth'},
  {name: 'Apps', path: '/apps'},
  {name: 'Bundle IDs', path: '/bundle-ids'},
  {name: 'Certificates', path: '/certificates'},
  {name: 'Profiles', path: '/profiles'},
  {name: 'Devices', path: '/devices'},
  {name: 'TestFlight', path: '/testflight'}
];

const currentRouteName = computed(() => {
  const match = menuItems.find(item => item.path === route.path);
  return match ? match.name : route.name;
});
</script>

<style lang="scss">
@import "@/theme/antd.scss";
@import "@/theme/reset.scss";
@import "@/theme/iconfont.css";

html,
body {
  height: 100%;
}

body {
  margin: 0;
  background: #f5f5f5;
}

#app {
  height: 100%;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
