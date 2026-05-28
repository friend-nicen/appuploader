/**
 * 注册全局方法/全局变量
 *
 * 规范：
 * - 尽量把 main.js 保持为“装配入口”，全局注入统一放这里
 */

import api from '@/service/api'

export default function register____global(app) {

  /* app 实例注入到 provide，便于组件内部 inject 使用 */
  app.provide('appInstance', app)

  /* 全局配置对象（兼容你提供的组件库：例如 v-form.vue 中使用 $setting.isMobile） */
  app.config.globalProperties.$setting = {
    isMobile: /Android|iPhone|iPad|iPod/i.test(navigator.userAgent),
  }

  /* 全局接口对象 */
  app.config.globalProperties.$api = api

}

