<template>
  <a-card :bordered="false" class="card" :title="$route.meta.name">
    <template #extra>
      <a-space>
        <v-button type="primary" @click="visible_add = true">添加 Certificate</v-button>
      </a-space>
    </template>

    <v-table :init="ListCertificates" :rowSelection="false">
      <template #bodyCell="{ column, record }">
        <template v-if="column.dataIndex && column.dataIndex[1] === 'certificateType'">
          <a-tag color="purple">
            {{ record.attributes?.certificateType || 'N/A' }}
          </a-tag>
        </template>
        <template v-else-if="column.dataIndex && column.dataIndex[1] === 'expirationDate'">
          {{ record.attributes?.expirationDate ? new Date(record.attributes.expirationDate).toLocaleDateString() : 'N/A' }}
        </template>
        <template v-else-if="column.dataIndex === 'action'">
          <a-popconfirm title="确定要撤销这个 Certificate 吗？" ok-text="确定" cancel-text="取消" @confirm="handleDelete(record.id)">
            <a-button type="link" danger>撤销</a-button>
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
      name="cert_add"
      title="添加 Certificate">
    </v-form>
  </a-card>
</template>

<script setup>
import { ref } from 'vue';
import { message } from 'ant-design-vue';
import { ListCertificates, RevokeCertificate } from '#/go/main/App';
import initTable from './table';
import initForm from './form';

const { table, columns } = initTable();
const { form_add, need_add } = initForm();

const visible_add = ref(false);

const handleAdd = async (data) => {
  try {
    message.info('创建 Certificate 暂未在后端实现');
    need_add.data.name = '';
    need_add.data.type = '';
    return true;
  } catch (err) {
    message.error('添加失败: ' + err);
    return false;
  }
};

const handleDelete = async (id) => {
  try {
    await RevokeCertificate(id);
    message.success('撤销成功');
    table.loadData();
  } catch (err) {
    message.error('撤销失败: ' + err);
  }
};

</script>

<style scoped lang="scss">
@include card;
.card {
  margin: 30px;
}
</style>
