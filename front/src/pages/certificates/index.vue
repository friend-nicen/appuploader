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
                  @exportP12OpenSSL="handleExportP12OpenSSL"
                  @exportPEM="handleExportPEM"
                  @changePwd="handleOpenChangePwd"
                  @delete="handleDelete"/>
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
        name="cert_add"
        title="创建 Certificate">
    </v-form>

    <!-- 修改密码弹窗 -->
    <v-form
        v-if="currentRecord"
        v-model:visible="visible_change_pwd"
        :dataSource="need_change_pwd"
        :form="form_change_pwd"
        :showBorder="false"
        :union="true"
        :init="changePwd"
        name="cert_change_pwd"
        title="修改证书密码">
    </v-form>
  </a-card>
</template>

<script setup>
import {ref} from 'vue';
import load from "@/common/load";
import {store} from '@/common';
import {
  CreateCertificate,
  ExportCSR,
  ExportLocalP12,
  ExportP12WithOpenSSL,
  ExportPEM,
  ListCertificatesWithLocal,
  RevokeCertificate,
  SaveBase64File,
  SaveTextFile,
  UpdateCertificatePassword,
  VerifyP12Password
} from '#/go/main/App';
import initTable from './table';
import initForm from './form';
import Action from './action.vue';

const {table, columns} = initTable();
const {form_add, need_add} = initForm();

const visible_add = ref(false);
const visible_change_pwd = ref(false);
const currentRecord = ref(null);

/* 修改密码表单配置 */
const form_change_pwd = [
  {
    key: 'password',
    type: 'input-password',
    label: '新密码',
    attr: {
      required: true,
      placeholder: '请输入新的证书密码（至少6位）',
      rules: [
        {min: 6, message: '密码至少 6 位'}
      ]
    }
  }
];

const need_change_pwd = store({password: ''});

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

/* 打开修改密码弹窗 */
const handleOpenChangePwd = (record) => {
  currentRecord.value = record;
  need_change_pwd.data.password = '';
  visible_change_pwd.value = true;
};

/* 修改密码（init 函数） */
const changePwd = (data) => UpdateCertificatePassword(currentRecord.value.id, data.password);

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

/* 导出 .p12（Go 库方式） */
const handleExportP12 = async (record) => {
  try {
    const result = await VerifyP12Password(record.id);
    if (result.includes('失败')) {
      load.info('密码校验: ' + result);
      return;
    }
    const name = record.attributes?.name || 'certificate';
    const type = record.attributes?.certificateType || 'CERTIFICATE';

    const res = await ExportLocalP12(record.id);
    await SaveBase64File(res, name + '_' + type + '.p12', '*.p12');
    load.success('.p12 导出成功');
  } catch (err) {
    if (err.message && err.message.includes('取消了')) return;
    load.error('.p12 导出失败: ' + err);
  }
};

/* 导出 .p12（OpenSSL 方式） */
const handleExportP12OpenSSL = async (record) => {
  try {
    const name = record.attributes?.name || 'certificate';
    const type = record.attributes?.certificateType || 'CERTIFICATE';
    const res = await ExportP12WithOpenSSL(record.id);
    await SaveBase64File(res, name + '_' + type + '.p12', '*.p12');
    load.success('.p12 导出成功（OpenSSL 方式）');
  } catch (err) {
    if (err.message && err.message.includes('取消了')) return;
    load.error('.p12 导出失败: ' + err);
  }
};

/* 导出 PEM（使用系统对话框选择保存位置） */
const handleExportPEM = async (record) => {
  try {
    const res = await ExportPEM(record.id);
    const parsed = JSON.parse(res);
    const name = record.attributes?.name || 'certificate';

    const keyPath = await SaveTextFile(parsed.keyPem, name + '_key.pem', '*.pem');
    if (!keyPath) return;

    const certPath = await SaveTextFile(parsed.certPem, name + '_cert.pem', '*.pem');
    if (!certPath) return;

    load.info('PEM 导出成功\n私钥: ' + keyPath + '\n证书: ' + certPath + '\n\n可使用 OpenSSL 打包:\nopenssl pkcs12 -export -in ' + certPath + ' -inkey ' + keyPath + ' -out output.p12');
  } catch (err) {
    if (err.message && err.message.includes('取消了')) return;
    load.error('PEM 导出失败: ' + err);
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
