<template>
  <a-dropdown :trigger="['click']">
    <a class="action-trigger" @click.prevent>
      操作
      <DownOutlined style="font-size: 12px;"/>
    </a>
    <template #overlay>
      <a-menu>
        <a-menu-item @click="$emit('exportCSR', record)">
          导出 CSR
        </a-menu-item>
        <a-sub-menu key="exportP12" title="导出 .p12">
          <a-menu-item @click="$emit('exportP12', record)">
            默认方式
          </a-menu-item>
          <a-menu-item @click="$emit('exportP12OpenSSL', record)">
            OpenSSL 方式
          </a-menu-item>
        </a-sub-menu>
        <a-menu-item @click="$emit('exportPEM', record)">
          导出 PEM
        </a-menu-item>
        <a-menu-item @click="$emit('changePwd', record)">
          修改密码
        </a-menu-item>
        <a-menu-item @click="handleDelete">
          删除证书
        </a-menu-item>
      </a-menu>
    </template>
  </a-dropdown>
</template>

<script setup>
import {DownOutlined} from '@ant-design/icons-vue';
import load from "@/common/load";

const props = defineProps({
  record: {type: Object, required: true}
});

const emit = defineEmits(['exportCSR', 'exportP12', 'exportP12OpenSSL', 'exportPEM', 'changePwd', 'delete']);

const handleDelete = () => {
  load.confirm(
      '确定要删除这个 Certificate 吗？',
      () => emit('delete', props.record),
      null,
      {okText: '确定', cancelText: '取消'}
  );
};
</script>

<style scoped lang="scss">
.action-trigger {
  color: #1677ff;
  cursor: pointer;
  user-select: none;
}
</style>
