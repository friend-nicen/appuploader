<template>

  <a-card :title="$route.meta.name" class="card" :bordered="false">

    <template #extra>
      <v-pop refreshText="查询"/>
    </template>

    <v-status :data="summary"/>

    <v-batch :show-reset="true" :data="table.condition.data" :form="form" layout="horizontal"/>

    <v-table :sticky="{offsetHeader:0}" :init="face.metrics" size="small" :row-selection="true"/>

  </a-card>

</template>

<script setup>

import initTable from './table';
import VStatus from './stats.vue';


/* 初始化表单 */
const {
  table,
  face,
  form,
  summary
} = initTable()


</script>

<style scoped lang="scss">

@include card;

.card {
  margin: 30px;
}


</style>
