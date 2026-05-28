<template>
  <a-dropdown :trigger="['click']">
    <a class="action-trigger" @click.prevent>
      操作
      <DownOutlined style="font-size: 12px;"/>
    </a>
    <template #overlay>
      <a-menu>
        <a-menu-item @click="$emit('download', record)">
          下载描述文件
        </a-menu-item>
        <a-menu-item @click="handleDelete">
          删除描述文件
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

const emit = defineEmits(['download', 'delete']);

const handleDelete = () => {
  load.confirm(
      '确定要删除这个 Profile 吗？',
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
