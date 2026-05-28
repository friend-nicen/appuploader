<template>
  <a-card :bordered="false" class="card" :title="$route.meta.name">
    <template #extra>
      <a-space>
        <v-button type="primary" @click="visible_add = true">添加 Bundle ID</v-button>
        <v-pop/>
      </a-space>
    </template>

    <v-table :init="ListBundleIds" :rowSelection="false">
      <template #bodyCell="{ column, record }">
        <template v-if="column.dataIndex && column.dataIndex[1] === 'platform'">
          <a-tag color="blue">
            {{ record.attributes?.platform || 'N/A' }}
          </a-tag>
        </template>
        <template v-else-if="column.dataIndex === 'action'">
          <a-popconfirm title="确定要删除这个 Bundle ID 吗？" ok-text="确定" cancel-text="取消"
                        @confirm="handleDelete(record.id)">
            <a-button type="link" danger>删除</a-button>
          </a-popconfirm>
        </template>
      </template>
    </v-table>

    <v-form
        v-model:visible="visible_add"
        :dataSource="need_add"
        :form="form_add"
        :showBorder="false"
        :union="true"
        :init="addBundleId"
        :after="() => Object.keys(need_add.data).forEach(k => need_add.data[k] = '')"
        message="添加成功"
        name="bundleid_add"
        title="添加 Bundle ID">
    </v-form>
  </a-card>
</template>

<script setup>
import {ref} from 'vue';
import load from "@/common/load";
import {CreateBundleId, DeleteBundleId, ListBundleIds} from '#/go/main/App';
import initTable from './table';
import initForm from './form';

const {table, columns} = initTable();
const {form_add, need_add} = initForm();

const visible_add = ref(false);

const addBundleId = (data) => CreateBundleId(data.name, data.identifier, data.platform);

const handleDelete = async (id) => {
  try {
    await DeleteBundleId(id);
    load.success('删除成功');
    table.loadData();
  } catch (err) {
    load.error('删除失败: ' + err);
  }
};

</script>

<style scoped lang="scss">
@include card;
.card {
  margin: 30px;
}
</style>
