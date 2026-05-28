<template>
  <a-layout-header
    class="v-header"
  >
    <div class="v-header__left">
      <h1 class="v-header__title">{{ title }}</h1>
    </div>
    <div class="v-header__right">
      <a-space>
        <v-button @click="handleExport">导出数据</v-button>
        <v-button @click="triggerImport">导入数据</v-button>
        <KeySelector />
      </a-space>
      <input ref="fileInput" type="file" accept=".json" style="display: none" @change="handleImport">
    </div>
  </a-layout-header>
</template>

<script setup>
import {ref} from 'vue';
import load from "@/common/load";
import {ExportLocalData, ImportLocalData} from '#/go/main/App';
import KeySelector from '@/components/v-key-select.vue';

defineProps({
  title: {
    type: String,
    default: ''
  }
});

const fileInput = ref(null);

/* 触发文件选择 */
const triggerImport = () => {
  fileInput.value?.click();
};

/* 导出数据 */
const handleExport = async () => {
  try {
    const json = await ExportLocalData();
    const blob = new Blob([json], { type: 'application/json;charset=utf-8' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'appstore_data_' + new Date().toISOString().slice(0, 10) + '.json';
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
    load.success('数据导出成功');
  } catch (err) {
    load.error('导出失败: ' + err);
  }
};

/* 导入数据 */
const handleImport = async (e) => {
  const file = e.target?.files?.[0];
  if (!file) return;
  try {
    const text = await file.text();
    await ImportLocalData(text);
    load.success('数据导入成功，请刷新页面');
  } catch (err) {
    load.error('导入失败: ' + err);
  }
  /* 重置 input，允许重复选择同一文件 */
  e.target.value = '';
};
</script>

<style lang="scss" scoped>
.v-header {
  background: #fff;
  padding: 0 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid #f0f0f0;
  height: 64px;

  &__left {
    display: flex;
    align-items: center;
  }

  &__title {
    font-size: 1.25rem;
    font-weight: 600;
    color: #1f2937;
    margin: 0;
  }

  &__right {
    display: flex;
    align-items: center;
  }
}
</style>
