<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center">
      <h2 class="text-2xl font-bold text-gray-800">Certificates</h2>
      <div class="space-x-2">
        <a-button type="primary" @click="showAddModal">添加 Certificate</a-button>
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
          <template v-if="column.key === 'type'">
            <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-purple-100 text-purple-800">
              {{ record.attributes?.certificateType || 'N/A' }}
            </span>
          </template>
          <template v-if="column.key === 'expirationDate'">
            {{ record.attributes?.expirationDate ? new Date(record.attributes.expirationDate).toLocaleDateString() : 'N/A' }}
          </template>
          <template v-if="column.key === 'action'">
            <a-popconfirm
              title="确定要撤销这个 Certificate 吗？"
              ok-text="确定"
              cancel-text="取消"
              @confirm="handleDelete(record.id)"
            >
              <a-button type="link" danger>撤销</a-button>
            </a-popconfirm>
          </template>
        </template>
      </a-table>
    </div>

    <!-- Add Modal Stub -->
    <a-modal
      v-model:open="isModalVisible"
      title="添加 Certificate"
      @ok="handleAdd"
      :confirmLoading="isSubmitting"
      okText="提交"
      cancelText="取消"
    >
      <a-form layout="vertical">
        <a-form-item label="名称">
          <a-input v-model:value="form.name" placeholder="请输入名称" />
        </a-form-item>
        <a-form-item label="类型">
          <a-input v-model:value="form.type" placeholder="请输入类型 (e.g. IOS_DEVELOPMENT)" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { ListCertificates, RevokeCertificate } from '../../wailsjs/go/main/App'

const items = ref([])
const loading = ref(false)
const error = ref('')

const isModalVisible = ref(false)
const isSubmitting = ref(false)
const form = ref({ name: '', type: '' })

const columns = [
  {
    title: '名称',
    dataIndex: ['attributes', 'name'],
    key: 'name',
  },
  {
    title: '类型',
    key: 'type',
  },
  {
    title: '序列号',
    dataIndex: ['attributes', 'serialNumber'],
    key: 'serialNumber',
  },
  {
    title: '过期时间',
    key: 'expirationDate',
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
    const jsonString = await ListCertificates()
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
  form.value = { name: '', type: '' }
  isModalVisible.value = true
}

const handleAdd = async () => {
  isSubmitting.value = true
  try {
    // Stub
    message.info('创建 Certificate 暂未在后端实现')
    isModalVisible.value = false
  } catch (err) {
    message.error('添加失败: ' + err)
  } finally {
    isSubmitting.value = false
  }
}

const handleDelete = async (id) => {
  try {
    await RevokeCertificate(id)
    message.success('撤销成功')
    fetchData()
  } catch (err) {
    message.error('撤销失败: ' + err)
  }
}

onMounted(() => {
  fetchData()
})
</script>
