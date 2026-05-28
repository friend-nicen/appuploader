<template>
  <a-card :bordered="false" class="card" :title="$route.meta.name">
    <template #extra>
      <a-space>
        <a-button @click="fetchData" :loading="loading">刷新数据</a-button>
        <v-button type="primary" @click="visible_add = true">添加 Device</v-button>
      </a-space>
    </template>

    <div v-if="error" class="error-msg">{{ error }}</div>

    <v-table :dataSource="items" :rowSelection="false">
      <template #bodyCell="{ column, record }">
        <template v-if="column.dataIndex && column.dataIndex[1] === 'deviceClass'">
          <a-tag color="indigo">
            {{ record.attributes?.deviceClass || 'N/A' }}
          </a-tag>
        </template>
        <template v-else-if="column.dataIndex && column.dataIndex[1] === 'status'">
          <a-tag :color="record.attributes?.status === 'ENABLED' ? 'green' : 'default'">
            {{ record.attributes?.status || 'N/A' }}
          </a-tag>
        </template>
        <template v-else-if="column.dataIndex === 'action'">
          <a-popconfirm title="确定要删除这个 Device 吗？" ok-text="确定" cancel-text="取消" @confirm="handleDelete(record.id)">
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
      name="device_add"
      title="添加 Device">
    </v-form>
  </a-card>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { message } from 'ant-design-vue';
import { ListDevices, DeleteDevice, CreateDevice } from '#/go/main/App';
import initTable from './table';
import initForm from './form';

const { table, columns } = initTable();
const { form_add, need_add } = initForm();

const items = ref([]);
const loading = ref(false);
const error = ref('');
const visible_add = ref(false);

const fetchData = async () => {
  loading.value = true;
  error.value = '';
  try {
    const jsonString = await ListDevices();
    if (!jsonString) {
      items.value = [];
      return;
    }
    const data = JSON.parse(jsonString);
    items.value = data.data || [];
  } catch (err) {
    error.value = '获取数据失败: ' + err;
    message.error('获取数据失败');
  } finally {
    loading.value = false;
  }
};

const handleAdd = async (data) => {
  try {
    await CreateDevice(data.name, data.udid);
    message.success('添加成功');
    await fetchData();
    need_add.data.name = '';
    need_add.data.udid = '';
    return true;
  } catch (err) {
    message.error('添加失败: ' + err);
    return false;
  }
};

const handleDelete = async (id) => {
  try {
    await DeleteDevice(id);
    message.success('删除成功');
    await fetchData();
  } catch (err) {
    message.error('删除失败: ' + err);
  }
};

onMounted(() => {
  fetchData();
});
</script>

<style scoped lang="scss">
@include card;
.card {
  margin: 30px;
}
.error-msg {
  background-color: #fef2f2;
  border-left: 4px solid #ef4444;
  padding: 16px;
  border-radius: 6px;
  margin-bottom: 16px;
  font-size: 14px;
  color: #b91c1c;
}
</style>
