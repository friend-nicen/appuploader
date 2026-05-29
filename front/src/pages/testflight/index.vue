<template>
  <a-card :bordered="false" class="card" :title="$route.meta.name">
    <div class="bg-white rounded-xl shadow-sm border border-gray-200 p-6 max-w-2xl mx-auto mt-6">
      <h3 class="text-lg font-medium mb-4">上传 IPA 到 TestFlight</h3>
      <div class="flex flex-col space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">选择 App</label>
          <a-select
            v-model:value="selectedApp"
            placeholder="请选择 App"
            style="width: 100%"
            :options="appOptions"
            :loading="loadingApps"
          />
        </div>

        <a-button @click="handleSelectFile" block :disabled="!selectedApp">选择 IPA 文件</a-button>

        <div v-if="filePath" class="p-3 bg-gray-50 border border-gray-200 rounded text-sm text-gray-700 break-all">
          已选择文件: {{ filePath }}
        </div>

        <div v-if="fileInfo" class="flex flex-col space-y-3 p-4 bg-blue-50 border border-blue-200 rounded">
          <div class="text-sm font-medium text-blue-800 mb-1">从 IPA 读取的信息</div>
          <div>
            <label class="block text-xs text-gray-600 mb-1">版本号</label>
            <a-input v-model:value="fileInfo.versionString" placeholder="未读取到版本号" />
          </div>
          <div>
            <label class="block text-xs text-gray-600 mb-1">构建号</label>
            <a-input v-model:value="fileInfo.buildNumber" placeholder="未读取到构建号" />
          </div>
          <div>
            <label class="block text-xs text-gray-600 mb-1">平台</label>
            <a-select
              v-model:value="fileInfo.platform"
              style="width: 100%"
              :options="[
                { label: 'iOS', value: 'IOS' },
                { label: 'macOS', value: 'MAC_OS' },
                { label: 'tvOS', value: 'TV_OS' },
                { label: 'visionOS', value: 'VISION_OS' }
              ]"
            />
          </div>
        </div>

        <a-button
          type="primary"
          :disabled="!filePath || !selectedApp"
          :loading="uploading"
          @click="handleUpload"
          block
        >
          {{ uploading ? '上传中...' : '上传' }}
        </a-button>
      </div>
    </div>
  </a-card>
</template>

<script setup>
import { ref, onMounted, reactive } from 'vue';
import load from "@/common/load";
import { GetIPAInfo, SelectFile, UploadIPA, ListApps } from '#/go/main/App';

const filePath = ref('');
const uploading = ref(false);
const selectedApp = ref(undefined);
const appOptions = ref([]);
const loadingApps = ref(false);
const fileInfo = ref(null);

onMounted(async () => {
  await loadApps();
});

const loadApps = async () => {
  loadingApps.value = true;
  try {
    const res = await ListApps();
    const parsed = JSON.parse(res);
    appOptions.value = (parsed.data || []).map(app => ({
      label: (app.attributes?.name || '') + ' (' + (app.attributes?.bundleId || '') + ')',
      value: app.id
    }));
  } catch (err) {
    load.error('加载 App 列表失败: ' + err);
  } finally {
    loadingApps.value = false;
  }
};

const handleSelectFile = async () => {
  try {
    const path = await SelectFile();
    if (path) {
      filePath.value = path;
      /* 自动读取 IPA 信息 */
      try {
        const res = await GetIPAInfo(path);
        const info = JSON.parse(res);
        fileInfo.value = reactive({
          versionString: info.versionString || '',
          buildNumber: info.buildNumber || '',
          platform: info.platform || ''
        });
        if (!info.versionString || !info.buildNumber) {
          load.info('未能自动读取版本号/构建号，请手动填写');
        }
      } catch (e) {
        fileInfo.value = reactive({
          versionString: '',
          buildNumber: '',
          platform: 'IOS'
        });
        load.info('无法自动读取 IPA 信息，请手动填写');
      }
    }
  } catch (err) {
    if (err.message && err.message.includes('取消了')) return;
    load.error('选择文件失败: ' + err);
  }
};

const handleUpload = async () => {
  if (!filePath.value || !selectedApp.value) return;

  uploading.value = true;
  try {
    const versionString = fileInfo.value?.versionString || '';
    const buildNumber = fileInfo.value?.buildNumber || '';
    const platform = fileInfo.value?.platform || 'IOS';

    if (!versionString || !buildNumber) {
      load.error('请填写版本号和构建号');
      uploading.value = false;
      return;
    }

    await UploadIPA(selectedApp.value, filePath.value, versionString, buildNumber, platform);
    load.success('上传成功，Apple 正在处理中，请稍后在 TestFlight 查看');
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
.space-y-3 > :not([hidden]) ~ :not([hidden]) {
  margin-top: 12px;
}
.p-4 {
  padding: 16px;
}
.bg-blue-50 {
  background-color: #eff6ff;
}
.text-blue-800 {
  color: #1e40af;
}
.text-xs {
  font-size: 0.75rem;
  line-height: 1rem;
}
.text-gray-600 {
  color: #4b5563;
}
</style>
