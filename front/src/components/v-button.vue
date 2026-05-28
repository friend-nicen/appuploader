<template>
  <!--
    v-button：对 a-button 的轻量封装
    用途：
    - 兼容你提供的组件库中对 <v-button> 的依赖
    - 保持原有使用方式（type/loading/size 等透传）
  -->
  <a-button v-bind="$attrs">
    <slot />
  </a-button>
</template>

<script setup>
/* 该组件仅用于透传属性与插槽，不额外定义 props，避免限制 antd 的按钮能力 */
</script>

