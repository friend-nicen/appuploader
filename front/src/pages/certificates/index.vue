<template>
  <a-card :bordered="false" class="card" :title="$route.meta.name">
    <template #extra>
      <a-space>
        <v-button type="primary" @click="visible_add = true">创建 Certificate</v-button>
        <v-pop/>
      </a-space>
    </template>

    <v-table :init="ListCertificatesWithLocal" :rowSelection="false">
      <template #bodyCell="{ column, record }">
        <template v-if="column.dataIndex && column.dataIndex[1] === 'certificateType'">
          <a-tag color="purple">
            {{ record.attributes?.certificateType || 'N/A' }}
          </a-tag>
        </template>
        <template v-else-if="column.dataIndex && column.dataIndex[1] === 'expirationDate'">
          {{
            record.attributes?.expirationDate ? new Date(record.attributes.expirationDate).toLocaleDateString() : 'N/A'
          }}
        </template>
        <template v-else-if="column.dataIndex === 'action'">
          <Action :record="record"
                  @exportCSR="handleExportCSR"
                  @exportP12="handleExportP12"
                  @delete="handleDelete" />
        </template>
      </template>
    </v-table>

    <!-- 添加证书弹窗 -->
    <v-form
        v-model:visible="visible_add"
        :dataSource="need_add"
        :form="form_add"
        :showBorder="false"
        :union="true"
        :submit="handleCreate"
        :hasTable="false"
        name="cert_add"
        title="创建 Certificate">
    </v-form>
  </a-card>
</template>

<script setup>
import {ref} from 'vue';
import load from "@/common/load";
import {
  CreateCertificate,
  ExportCSR,
  ExportLocalP12,
  ListCertificatesWithLocal,
  RevokeCertificate,
  SaveTextFile,
  SaveBase64File
} from '#/go/main/App';
import initTable from './table';
import initForm from './form';
import Action from './action.vue';

const {table, columns} = initTable();
const {form_add, need_add} = initForm();

const visible_add = ref(false);

/* 提交创建证书 */
const handleCreate = async (data) => {
  try {
    await CreateCertificate(data.name, data.type, data.password);
    load.success('证书创建成功');
    table.loadData();
    return true;
  } catch (err) {
    load.error('创建失败: ' + err);
    return false;
  }
};

/* 导出 CSR */
const handleExportCSR = async (record) => {
  try {
    const csr = await ExportCSR(record.id);
    const name = record.attributes?.name || 'certificate';
    await SaveTextFile(csr, name + '_' + record.id + '.csr', '*.csr');
    load.success('CSR 导出成功');
  } catch (err) {
    if (err.message && err.message.includes('取消了')) return;
    load.error('CSR 导出失败: ' + err);
  }
};

/* 导出 .p12 */
const handleExportP12 = async (record) => {
  try {
    const res = await ExportLocalP12(record.id);
    const name = record.attributes?.name || 'certificate';
    const type = record.attributes?.certificateType || 'CERTIFICATE';
    await SaveBase64File(res, name + '_' + type + '.p12', '*.p12');
    load.success('.p12 导出成功');
  } catch (err) {
    if (err.message && err.message.includes('取消了')) return;
    load.error('.p12 导出失败: ' + err);
  }
};

/* 删除证书 */
const handleDelete = async (record) => {
  try {
    await RevokeCertificate(record.id);
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
