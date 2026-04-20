<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center">
      <h2 class="text-2xl font-bold text-gray-800">Provisioning Profiles</h2>
      <div class="space-x-2">
        <a-button type="primary" @click="showAddModal">添加 Profile</a-button>
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
          <template v-if="column.key === 'profileType'">
            <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
              {{ record.attributes?.profileType || 'N/A' }}
            </span>
          </template>
          <template v-if="column.key === 'expirationDate'">
            {{ record.attributes?.expirationDate ? new Date(record.attributes.expirationDate).toLocaleDateString() : 'N/A' }}
          </template>
          <template v-if="column.key === 'action'">
            <a-popconfirm
              title="确定要删除这个 Profile 吗？"
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
      title="添加 Profile"
      @ok="handleAdd"
      :confirmLoading="isSubmitting"
      okText="提交"
      cancelText="取消"
    >
      <a-form layout="vertical">
        <a-form-item label="名称">
          <a-input v-model:value="form.name" placeholder="请输入名称" />
        </a-form-item>
        <a-form-item label="Profile Type">
          <a-input v-model:value="form.profileType" placeholder="请输入 Profile Type (e.g. IOS_APP_DEVELOPMENT)" />
        </a-form-item>
        <a-form-item label="Bundle ID">
          <a-input v-model:value="form.bundleId" placeholder="请输入 Bundle ID 的 ID" />
        </a-form-item>
        <a-form-item label="Certificate ID">
          <a-input v-model:value="form.certId" placeholder="请输入 Certificate 的 ID" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { ListProfiles, DeleteProfile, CreateProfile } from '../../wailsjs/go/main/App'

const items = ref([])
const loading = ref(false)
const error = ref('')

const isModalVisible = ref(false)
const isSubmitting = ref(false)
const form = ref({ name: '', profileType: '', bundleId: '', certId: '' })

const columns = [
  {
    title: '名称',
    dataIndex: ['attributes', 'name'],
    key: 'name',
  },
  {
    title: '类型',
    key: 'profileType',
  },
  {
    title: 'Platform',
    dataIndex: ['attributes', 'platform'],
    key: 'platform',
  },
  {
    title: '状态',
    dataIndex: ['attributes', 'profileState'],
    key: 'profileState',
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
    const jsonString = await ListProfiles()
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
  form.value = { name: '', profileType: '', bundleId: '', certId: '' }
  isModalVisible.value = true
}

const handleAdd = async () => {
  isSubmitting.value = true
  try {
    await CreateProfile(form.value.name, form.value.profileType, form.value.bundleId, form.value.certId)
    message.success('添加成功')
    isModalVisible.value = false
    fetchData()
  } catch (err) {
    message.error('添加失败: ' + err)
  } finally {
    isSubmitting.value = false
  }
}

const handleDelete = async (id) => {
  try {
    await DeleteProfile(id)
    message.success('删除成功')
    fetchData()
  } catch (err) {
    message.error('删除失败: ' + err)
  }
}

onMounted(() => {
  fetchData()
})
</script>
