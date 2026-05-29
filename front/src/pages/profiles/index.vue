<template>
  <a-card :bordered="false" class="card" :title="$route.meta.name">
    <template #extra>
      <a-space>
        <v-button type="primary" @click="handleOpenAdd">添加 Profile</v-button>
        <v-pop/>
      </a-space>
    </template>

    <v-table :init="ListProfiles" :rowSelection="false">
      <template #bodyCell="{ column, record }">
        <template v-if="column.dataIndex && column.dataIndex[1] === 'profileType'">
          <a-tag color="green">
            {{ record.attributes?.profileType || 'N/A' }}
          </a-tag>
        </template>
        <template v-else-if="column.dataIndex && column.dataIndex[1] === 'expirationDate'">
          {{
            record.attributes?.expirationDate ? new Date(record.attributes.expirationDate).toLocaleDateString() : 'N/A'
          }}
        </template>
        <template v-else-if="column.dataIndex === 'action'">
          <Action :record="record"
                  @download="handleDownload"
                  @delete="handleDelete" />
        </template>
      </template>
    </v-table>

    <v-form
        v-model:visible="visible_add"
        :dataSource="need_add"
        :form="form_add"
        :showBorder="false"
        :union="true"
        :submit="handleCreate"
        :hasTable="false"
        name="profile_add"
        title="添加 Profile">
    </v-form>
  </a-card>
</template>

<script setup>
import {ref} from 'vue';
import load from "@/common/load";
import {
  CreateProfile,
  DeleteProfile,
  DownloadProfile,
  ListBundleIds,
  ListCertificatesWithLocal,
  ListDevices,
  ListProfiles,
  SaveBase64File
} from '#/go/main/App';
import initTable from './table';
import initForm, {BUNDLE_ID_INDEX, CERT_INDEX, DEVICE_INDEX} from './form';
import Action from './action.vue';

const {table, columns} = initTable();
const {form_add, need_add} = initForm();

const visible_add = ref(false);

/* 加载表单选项数据 */
const loadFormOptions = async () => {
  try {
    const [bundlesRes, certsRes, devicesRes] = await Promise.all([
      ListBundleIds(),
      ListCertificatesWithLocal(),
      ListDevices()
    ]);

    const bundles = JSON.parse(bundlesRes);
    form_add[BUNDLE_ID_INDEX].attr.options = (bundles.data || []).map(b => ({
      label: (b.attributes?.name || '') + ' (' + (b.attributes?.identifier || '') + ')',
      value: b.id
    }));

    const certs = JSON.parse(certsRes);
    form_add[CERT_INDEX].attr.options = (certs.data || []).map(c => ({
      label: (c.attributes?.name || '') + ' (' + (c.attributes?.certificateType || '') + ')',
      value: c.id
    }));

    const devices = JSON.parse(devicesRes);
    form_add[DEVICE_INDEX].attr.options = (devices.data || []).map(d => ({
      label: (d.attributes?.name || '') + ' (' + (d.attributes?.udid || '') + ')',
      value: d.id
    }));
  } catch (e) {
    /* 选项加载失败不影响页面使用 */
    console.warn('加载表单选项失败:', e);
  }
};

/* 打开添加弹窗时刷新选项 */
const handleOpenAdd = () => {
  visible_add.value = true;
  loadFormOptions();
};

/* 提交创建描述文件 */
const handleCreate = async (data) => {
  try {
    /* 转换证书 ID 数组为 Apple API 格式 */
    const certData = (data.certIds || []).map(id => ({
      id: id,
      type: 'certificates'
    }));
    /* 转换设备 ID 数组为 Apple API 格式 */
    const deviceData = (data.deviceIds || []).map(id => ({
      id: id,
      type: 'devices'
    }));
    await CreateProfile(
        data.name,
        data.profileType,
        data.bundleId,
        JSON.stringify(certData),
        deviceData.length > 0 ? JSON.stringify(deviceData) : ''
    );
    load.success('描述文件创建成功');
    table.loadData();
    return true;
  } catch (err) {
    load.error('创建失败: ' + err);
    return false;
  }
};

/* 下载描述文件 */
const handleDownload = async (record) => {
  try {
    const res = await DownloadProfile(record.id);
    const parsed = JSON.parse(res);
    const content = parsed.data?.attributes?.profileContent;
    if (!content) {
      load.error('未找到描述文件内容');
      return;
    }
    const name = record.attributes?.name || 'profile';
    await SaveBase64File(content, name + '_' + record.id + '.mobileprovision', '*.mobileprovision');
    load.success('描述文件下载成功');
  } catch (err) {
    if (err.message && err.message.includes('取消了')) return;
    load.error('下载失败: ' + err);
  }
};

/* 删除描述文件 */
const handleDelete = async (record) => {
  try {
    await DeleteProfile(record.id);
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
