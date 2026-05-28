<template>
  <a-card :bordered="false" class="card" :title="$route.meta.name">
    <div class="bg-white rounded-xl shadow-sm border border-gray-200 p-6 max-w-2xl mx-auto mt-6">
      <h3 class="text-lg font-medium mb-4">Upload Build to TestFlight</h3>
      <div class="flex flex-col space-y-4">
        <a-button @click="handleSelectFile" block>选择 IPA 文件</a-button>
        
        <div v-if="filePath" class="p-3 bg-gray-50 border border-gray-200 rounded text-sm text-gray-700 break-all">
          已选择文件: {{ filePath }}
        </div>

        <a-button 
          type="primary" 
          :disabled="!filePath" 
          :loading="uploading"
          @click="handleUpload"
          block
        >
          上传
        </a-button>
      </div>
    </div>
  </a-card>
</template>

<script setup>
import { ref } from 'vue';
import load from "@/common/load";
import { SelectFile, UploadIPA } from '#/go/main/App';

const filePath = ref('');
const uploading = ref(false);

const handleSelectFile = async () => {
  try {
    const path = await SelectFile();
    if (path) {
      filePath.value = path;
      load.success('已选择文件');
    }
  } catch (err) {
    load.error('选择文件失败: ' + err);
  }
};

const handleUpload = async () => {
  if (!filePath.value) return;
  
  uploading.value = true;
  try {
    await UploadIPA(filePath.value);
    load.success('上传成功');
  } catch (err) {
    load.error('上传失败: ' + err);
  } finally {
    uploading.value = false;
  }
};
</script>

<style scoped lang="scss">
@include card;
.card {
  margin: 30px;
}
.bg-white {
  background-color: #fff;
  border-radius: 12px;
  border: 1px solid #e5e7eb;
  padding: 24px;
}
.max-w-2xl {
  max-width: 42rem;
}
.mx-auto {
  margin-left: auto;
  margin-right: auto;
}
.mt-6 {
  margin-top: 24px;
}
.text-lg {
  font-size: 1.125rem;
  line-height: 1.75rem;
}
.font-medium {
  font-weight: 500;
}
.mb-4 {
  margin-bottom: 16px;
}
.flex {
  display: flex;
}
.flex-col {
  flex-direction: column;
}
.space-y-4 > :not([hidden]) ~ :not([hidden]) {
  margin-top: 16px;
}
.p-3 {
  padding: 12px;
}
.bg-gray-50 {
  background-color: #f9fafb;
}
.text-sm {
  font-size: 0.875rem;
  line-height: 1.25rem;
}
.text-gray-700 {
  color: #374151;
}
.break-all {
  word-break: break-all;
}
.rounded {
  border-radius: 0.25rem;
}
</style>
