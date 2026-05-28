<template>
  <a-card :bordered="false" class="card" :title="$route.meta.name">
    <template #extra>
      <v-pop/>
    </template>
    <v-table :init="ListApps" :rowSelection="false" />
  </a-card>
</template>

<script setup>
import {ListApps} from '#/go/main/App';
import initTable from './table';

const {table, columns} = initTable();

</script>

<style scoped lang="scss">
@include card;
.card {
  margin: 30px;
}
</style>
