<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center">
      <h2 class="text-2xl font-bold text-gray-800">TestFlight</h2>
    </div>

    <div class="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
      <a-card title="Upload Build to TestFlight" :bordered="false">
        <div class="flex flex-col space-y-4 max-w-md">
          <a-button @click="handleSelectFile">选择 IPA 文件</a-button>
          
          <div v-if="filePath" class="p-3 bg-gray-50 border border-gray-200 rounded text-sm text-gray-700 break-all">
            已选择文件: {{ filePath }}
          </div>

          <a-button 
            type="primary" 
            :disabled="!filePath" 
            :loading="uploading"
            @click="handleUpload"
          >
            上传
          </a-button>
        </div>
      </a-card>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { message } from 'ant-design-vue'
import { SelectFile, UploadIPA } from '../../wailsjs/go/main/App'

const filePath = ref('')
const uploading = ref(false)

const handleSelectFile = async () => {
  try {
    const path = await SelectFile()
    if (path) {
      filePath.value = path
      message.success('已选择文件')
    }
  } catch (err) {
    message.error('选择文件失败: ' + err)
  }
}

const handleUpload = async () => {
  if (!filePath.value) return
  
  uploading.value = true
  try {
    await UploadIPA(filePath.value)
    message.success('上传成功')
  } catch (err) {
    message.error('上传失败: ' + err)
  } finally {
    uploading.value = false
  }
}
</script>
