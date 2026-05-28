<template>
  <div style="display: flex; align-items: center; gap: 8px;">
    <a-select
      v-model:value="selectedKeyId"
      :options="keyOptions"
      style="width: 200px"
      placeholder="选择 API Key"
      @change="handleKeyChange"
      :loading="loading"
    />
    <a-button type="text" size="small" @click="fetchKeys" :loading="loading" title="刷新">
      <v-icon name="sync" />
    </a-button>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import { GetApiKeys, SetCurrentKey } from '#/go/main/App';

const keys = ref([]);
const selectedKeyId = ref(null);
const loading = ref(false);

const keyOptions = computed(() => {
  return keys.value.map(k => ({
    value: k.id,
    label: k.name
  }));
});

const fetchKeys = async () => {
  loading.value = true;
  try {
    const data = await GetApiKeys();
    keys.value = data || [];
    
    /* 找出当前处于激活状态的密钥 */
    const activeKey = keys.value.find(k => k.is_active);
    if (activeKey) {
      selectedKeyId.value = activeKey.id;
    } else {
      selectedKeyId.value = null;
    }
  } catch (err) {
    console.error('获取密钥列表失败:', err);
  } finally {
    loading.value = false;
  }
};

const handleKeyChange = async (id) => {
  try {
    await SetCurrentKey(id);
    /* 切换密钥后刷新整个页面，以刷新主视图状态 */
    window.location.reload();
  } catch (err) {
    console.error('设置当前密钥失败:', err);
  }
};

onMounted(() => {
  fetchKeys();
});
</script>
