<template>
  <a-card :bordered="false" class="card" :title="$route.meta.name">
    <template #extra>
      <a-space>
        <v-button type="primary" @click="visible_add = true">添加 App</v-button>
      </a-space>
    </template>

    <v-table :init="ListApps" :rowSelection="false">
      <template #bodyCell="{ column, record }">
        <template v-if="column.dataIndex === 'action'">
          <a-popconfirm title="确定要删除这个 App 吗？" ok-text="确定" cancel-text="取消" @confirm="handleDelete(record.id)">
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
      :centered="true"
      :union="true"
      :submit="handleAdd"
      message="添加成功！"
      name="app_add"
      title="添加 App">
    </v-form>
  </a-card>
</template>

<script setup>
import { ref } from 'vue';
import { message } from 'ant-design-vue';
import { ListApps } from '#/go/main/App';
import initTable from './table';
import initForm from './form';

const { table, columns } = initTable();
const { form_add, need_add } = initForm();

const visible_add = ref(false);

const handleAdd = async (data) => {
  try {
    message.info('添加 App 功能暂未在后端实现');
    need_add.data.name = '';
    need_add.data.bundleId = '';
    return true;
  } catch (err) {
    message.error('添加失败: ' + err);
    return false;
  }
};

const handleDelete = async (id) => {
  try {
    message.info('删除 App 功能暂未在后端实现');
  } catch (err) {
    message.error('删除失败: ' + err);
  }
};

</script>

<style scoped lang="scss">
@include card;
.card {
  margin: 30px;
}
</style>
