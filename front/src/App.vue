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
      <VMainLayout />
    </a-config-provider>
  </a-style-provider>
</template>

<script setup>
import VMainLayout from '@/layouts/v-main-layout/index.vue';

import { px2remTransformer } from 'ant-design-vue';
import zhCN from 'ant-design-vue/es/locale/zh_CN';
import { dynamicTheme } from "@/theme/theme";
import dayjs from 'dayjs';
import 'dayjs/locale/zh-cn';

const theme = dynamicTheme();

dayjs.locale('zh');
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
