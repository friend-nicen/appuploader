<template>
  <a-card :bordered="false" class="card" :title="$route.meta.name">
    <template #extra>
      <a-space>
        <a-button @click="testAuth" :loading="testingAuth">测试当前认证</a-button>
        <v-button type="primary" @click="visible_add = true">添加新密钥</v-button>
        <v-pop/>
      </a-space>
    </template>

    <v-table :init="GetApiKeys" :rowSelection="false">
      <template #bodyCell="{ column, record }">
        <template v-if="column.dataIndex === 'is_active'">
          <a-tag :color="record.is_active ? 'green' : 'default'">
            {{ record.is_active ? '当前使用' : '未激活' }}
          </a-tag>
        </template>
        <template v-else-if="column.dataIndex === 'action'">
          <a-space>
            <a-button v-if="!record.is_active" type="link" size="small" @click="setActive(record.id)">
              设为当前
            </a-button>
            <a-popconfirm title="确定要删除这个密钥吗？" ok-text="确定" cancel-text="取消"
                          @confirm="deleteKey(record.id)">
              <a-button type="link" danger size="small">删除</a-button>
            </a-popconfirm>
          </a-space>
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
        :init="addKey"
        :after="() => Object.keys(need_add.data).forEach(k => need_add.data[k] = '')"
        message="添加成功"
        name="key_add"
        title="添加 API 密钥">
    </v-form>
  </a-card>
</template>

<script setup>
import {AddApiKey, DeleteApiKey, GetApiKeys, SetCurrentKey, TestAuth} from '#/go/main/App';
import {ref} from 'vue';
import load from "@/common/load";
import initTable from './table';
import initForm from './form';

const {table, columns} = initTable();
const {form_add, need_add} = initForm();

const visible_add = ref(false);
const testingAuth = ref(false);

const addKey = (data) => AddApiKey(data.name, data.issuer_id, data.key_id, data.private_key);

const setActive = async (id) => {
  try {
    await SetCurrentKey(id);
    load.success('设置成功');
    table.loadData();
  } catch (err) {
    load.error('设置失败: ' + err);
  }
};

const deleteKey = async (id) => {
  try {
    await DeleteApiKey(id);
    load.success('删除成功');
    table.loadData();
  } catch (err) {
    load.error('删除失败: ' + err);
  }
};

const testAuth = async () => {
  testingAuth.value = true;
  try {
    const success = await TestAuth();
    if (success) {
      load.success('认证成功！');
    } else {
      load.error('认证失败。');
    }
  } catch (err) {
    load.error('测试认证异常: ' + err);
  } finally {
    testingAuth.value = false;
  }
};

</script>

<style scoped lang="scss">
@include card;
.card {
  margin: 30px;
}
</style>
