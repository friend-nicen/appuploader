<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center">
      <h2 class="text-2xl font-bold text-gray-800">Apps</h2>
      <div class="space-x-2">
        <a-button type="primary" @click="showAddModal">添加 App</a-button>
        <a-button @click="fetchData" :loading="loading">刷新数据</a-button>
      </div>
    </div>

    <div v-if="error" class="bg-red-50 border-l-4 border-red-500 p-4 rounded-md">
      <p class="text-sm text-red-700">{{ error }}</p>
    </div>

    <div class="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden p-4">
      <a-table 
        :dataSource="items" 
        :columns="columns" 
        rowKey="id" 
        :loading="loading"
        :pagination="{ pageSize: 10 }"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'action'">
            <a-popconfirm
              title="确定要删除这个 App 吗？"
              ok-text="确定"
              cancel-text="取消"
              @confirm="handleDelete(record.id)"
            >
              <a-button type="link" danger>删除</a-button>
            </a-popconfirm>
          </template>
        </template>
      </a-table>
    </div>

    <!-- Add Modal Stub -->
    <a-modal
      v-model:open="isModalVisible"
      title="添加 App"
      @ok="handleAdd"
      :confirmLoading="isSubmitting"
      okText="提交"
      cancelText="取消"
    >
      <a-form layout="vertical">
        <a-form-item label="名称">
          <a-input v-model:value="form.name" placeholder="请输入名称" />
        </a-form-item>
        <a-form-item label="Bundle ID">
          <a-input v-model:value="form.bundleId" placeholder="请输入 Bundle ID" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { ListApps } from '../../wailsjs/go/main/App'

const items = ref([])
const loading = ref(false)
const error = ref('')

const isModalVisible = ref(false)
const isSubmitting = ref(false)
const form = ref({ name: '', bundleId: '' })

const columns = [
  {
    title: 'Name',
    dataIndex: ['attributes', 'name'],
    key: 'name',
  },
  {
    title: 'Bundle ID',
    dataIndex: ['attributes', 'bundleId'],
    key: 'bundleId',
  },
  {
    title: 'SKU',
    dataIndex: ['attributes', 'sku'],
    key: 'sku',
  },
  {
    title: 'Primary Locale',
    dataIndex: ['attributes', 'primaryLocale'],
    key: 'primaryLocale',
  },
  {
    title: '操作',
    key: 'action',
    width: 100,
  }
]

const fetchData = async () => {
  loading.value = true
  error.value = ''
  
  try {
    const jsonString = await ListApps()
    if (!jsonString) {
      items.value = []
      return
    }
    const data = JSON.parse(jsonString)
    items.value = data.data || []
  } catch (err) {
    error.value = '获取数据失败: ' + err
    message.error('获取数据失败')
  } finally {
    loading.value = false
  }
}

const showAddModal = () => {
  form.value = { name: '', bundleId: '' }
  isModalVisible.value = true
}

const handleAdd = async () => {
  isSubmitting.value = true
  try {
    // 这是一个 stub，因为后端还没有实现 CreateApp
    message.info('添加 App 功能暂未在后端实现')
    isModalVisible.value = false
  } catch (err) {
    message.error('添加失败: ' + err)
  } finally {
    isSubmitting.value = false
  }
}

const handleDelete = async (id) => {
  try {
    // 这是一个 stub，因为后端还没有实现 DeleteApp
    message.info('删除 App 功能暂未在后端实现')
  } catch (err) {
    message.error('删除失败: ' + err)
  }
}

onMounted(() => {
  fetchData()
})
</script>
