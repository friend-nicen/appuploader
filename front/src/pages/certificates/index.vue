<template>
  <a-card :bordered="false" class="card" :title="$route.meta.name">
    <template #extra>
      <a-space>
        <v-button type="primary" @click="visible_add = true">创建 Certificate</v-button>
        <v-pop/>
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
          {{
            record.attributes?.expirationDate ? new Date(record.attributes.expirationDate).toLocaleDateString() : 'N/A'
          }}
        </template>
        <template v-else-if="column.dataIndex === 'action'">
          <a-space>
            <a-button type="link" size="small" @click="handleDownload(record)">下载 .cer</a-button>
            <a-popconfirm title="确定要撤销这个 Certificate 吗？" ok-text="确定" cancel-text="取消"
                          @confirm="handleDelete(record.id)">
              <a-button type="link" danger size="small">撤销</a-button>
            </a-popconfirm>
          </a-space>
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

    <!-- 创建成功后的下载弹窗 -->
    <a-modal
        v-model:open="visible_download"
        title="证书创建成功"
        :footer="null"
        destroy-on-close>
      <p>名称: <strong>{{ createdCert?.name }}</strong></p>
      <p>证书 ID: <code>{{ createdCert?.certificateId }}</code></p>
      <p>显示名称: <strong>{{ createdCert?.displayName }}</strong></p>
      <div style="margin-top: 24px; display: flex; gap: 12px; justify-content: center;">
        <a-button type="primary" @click="downloadCer">下载 .cer</a-button>
        <a-button @click="promptP12Password">下载 .p12 (含私钥)</a-button>
      </div>
    </a-modal>
  </a-card>
</template>

<script setup>
import {ref} from 'vue';
import load from "@/common/load";
import {CreateCertificate, DownloadCertificate, ExportP12, ListCertificates, RevokeCertificate} from '#/go/main/App';
import initTable from './table';
import initForm from './form';

const {table, columns} = initTable();
const {form_add, need_add} = initForm();

const visible_add = ref(false);
const visible_download = ref(false);
const createdCert = ref(null);

/* 提交创建证书 */
const handleCreate = async (data) => {
  try {
    const res = await CreateCertificate(data.name, data.type);
    const parsed = JSON.parse(res);
    /* 保存创建结果，用于后续下载 */
    createdCert.value = {
      certPem: parsed.certPem,
      keyPem: parsed.keyPem,
      name: parsed.name,
      certificateId: parsed.certificateId,
      displayName: parsed.displayName,
    };
    load.success('证书创建成功');
    /* 关闭添加弹窗后弹出下载弹窗 */
    setTimeout(() => {
      visible_download.value = true;
    }, 300);
    return true;
  } catch (err) {
    load.error('创建失败: ' + err);
    return false;
  }
};

/* 下载 .cer 文件（公钥证书） */
const downloadCer = () => {
  const cert = createdCert.value;
  if (!cert || !cert.certPem) {
    load.error('证书数据丢失');
    return;
  }
  /* certPem 是 base64 编码的 DER 证书，需要解码 */
  try {
    const binary = atob(cert.certPem);
    const bytes = new Uint8Array(binary.length);
    for (let i = 0; i < binary.length; i++) {
      bytes[i] = binary.charCodeAt(i);
    }
    const blob = new Blob([bytes], { type: 'application/x-x509-ca-cert' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = (cert.name || 'certificate') + '_' + (cert.certificateId || '') + '.cer';
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  } catch (e) {
    load.error('下载 .cer 失败: ' + e);
  }
};

/* 提示用户输入 .p12 密码 */
const promptP12Password = () => {
  const password = prompt('请设置 .p12 文件的导出密码（至少 6 位）:', '');
  if (!password || password.length < 6) {
    load.warning('密码至少 6 位');
    return;
  }
  downloadP12(password);
};

/* 下载 .p12 文件（证书 + 私钥，密码保护） */
const downloadP12 = async (password) => {
  const cert = createdCert.value;
  if (!cert || !cert.certPem || !cert.keyPem) {
    load.error('证书或私钥数据丢失');
    return;
  }
  try {
    const res = await ExportP12(cert.certPem, cert.keyPem, password);
    const binary = atob(res);
    const bytes = new Uint8Array(binary.length);
    for (let i = 0; i < binary.length; i++) {
      bytes[i] = binary.charCodeAt(i);
    }
    const blob = new Blob([bytes], { type: 'application/x-pkcs12' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = (cert.name || 'certificate') + '_' + (cert.certificateId || '') + '.p12';
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
    load.success('.p12 下载完成');
  } catch (err) {
    load.error('生成 .p12 失败: ' + err);
  }
};

/* 从列表下载已有证书的 .cer */
const handleDownload = async (record) => {
  try {
    const res = await DownloadCertificate(record.id);
    const parsed = JSON.parse(res);
    const content = parsed.data?.attributes?.certificateContent;
    if (!content) {
      load.error('未找到证书内容');
      return;
    }
    const binary = atob(content);
    const bytes = new Uint8Array(binary.length);
    for (let i = 0; i < binary.length; i++) {
      bytes[i] = binary.charCodeAt(i);
    }
    const blob = new Blob([bytes], { type: 'application/x-x509-ca-cert' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    const type = record.attributes?.certificateType || 'CERTIFICATE';
    a.download = (record.attributes?.name || 'certificate') + '_' + type + '.cer';
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  } catch (err) {
    load.error('下载失败: ' + err);
  }
};

const handleDelete = async (id) => {
  try {
    await RevokeCertificate(id);
    load.success('撤销成功');
    table.loadData();
  } catch (err) {
    load.error('撤销失败: ' + err);
  }
};

</script>

<style scoped lang="scss">
@include card;
.card {
  margin: 30px;
}
</style>