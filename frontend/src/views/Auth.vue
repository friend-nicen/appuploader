<template>
  <div class="max-w-5xl mx-auto p-4 space-y-6">
    <!-- 密钥列表 -->
    <a-card title="已保存的密钥" :bordered="false" class="shadow-sm">
      <template #extra>
        <a-space>
          <a-button @click="testAuth" :loading="testingAuth">
            测试当前认证
          </a-button>
          <a-button type="primary" @click="openAddModal">
            添加新密钥
          </a-button>
        </a-space>
      </template>
      
      <a-table 
        :columns="columns" 
        :data-source="keys" 
        :row-key="record => record.id"
        :loading="loadingKeys"
        :pagination="false"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'is_active'">
            <a-tag :color="record.is_active ? 'green' : 'default'">
              {{ record.is_active ? '当前使用' : '未激活' }}
            </a-tag>
          </template>
          <template v-else-if="column.key === 'action'">
            <a-space>
              <a-button 
                v-if="!record.is_active" 
                type="link" 
                size="small" 
                @click="setActive(record.id)"
              >
                设为当前
              </a-button>
              <a-button 
                type="link" 
                size="small" 
                @click="editKey(record)"
              >
                编辑
              </a-button>
              <a-popconfirm
                title="确定要删除这个密钥吗？"
                ok-text="确定"
                cancel-text="取消"
                @confirm="deleteKey(record.id)"
              >
                <a-button type="link" danger size="small">
                  删除
                </a-button>
              </a-popconfirm>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-card>

    <!-- 弹窗：添加/编辑 -->
    <a-modal
      v-model:open="isModalVisible"
      :title="editingId ? '编辑 API 密钥' : '添加新的 API 密钥'"
      @ok="handleSaveKey"
      @cancel="cancelEdit"
      :confirmLoading="loading"
      okText="保存"
      cancelText="取消"
      destroyOnClose
    >
      <a-form :model="form" :rules="rules" ref="formRef" layout="vertical" class="mt-4">
        <a-form-item label="名称 (Name)" name="name">
          <a-input v-model:value="form.name" placeholder="e.g. My CI Key" />
        </a-form-item>
        <a-form-item label="Issuer ID" name="issuerId">
          <a-input v-model:value="form.issuerId" placeholder="e.g. 69a6de7e-..." />
        </a-form-item>
        <a-form-item label="Key ID" name="keyId">
          <a-input v-model:value="form.keyId" placeholder="e.g. 2X9R4H..." />
        </a-form-item>
        <a-form-item label="Private Key (私钥内容 .p8)" name="privateKey">
          <a-textarea 
            v-model:value="form.privateKey" 
            :rows="4" 
            :placeholder="editingId ? '(留空表示不修改)' : '-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----'"
            class="font-mono text-sm"
          />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { GetApiKeys, AddApiKey, UpdateApiKey, SetCurrentKey, DeleteApiKey, TestAuth } from '../../wailsjs/go/main/App'

const keys = ref([])
const loading = ref(false)
const loadingKeys = ref(false)
const testingAuth = ref(false)
const editingId = ref(null)
const isModalVisible = ref(false)
const formRef = ref(null)

const form = reactive({
  name: '',
  issuerId: '',
  keyId: '',
  privateKey: ''
})

const rules = {
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
  issuerId: [{ required: true, message: '请输入 Issuer ID', trigger: 'blur' }],
  keyId: [{ required: true, message: '请输入 Key ID', trigger: 'blur' }],
  privateKey: [{ 
    validator: async (rule, value) => {
      if (!editingId.value && !value) {
        return Promise.reject('请输入私钥内容')
      }
      return Promise.resolve()
    }, 
    trigger: 'blur' 
  }]
}

const columns = [
  { title: '名称', dataIndex: 'name', key: 'name' },
  { title: 'Key ID', dataIndex: 'key_id', key: 'key_id', class: 'font-mono' },
  { title: 'Issuer ID', dataIndex: 'issuer_id', key: 'issuer_id', class: 'font-mono' },
  { title: '状态', key: 'is_active', width: 100 },
  { title: '操作', key: 'action', width: 200, align: 'right' }
]

// 获取密钥列表
const fetchKeys = async () => {
  loadingKeys.value = true
  try {
    const data = await GetApiKeys()
    keys.value = data || []
  } catch (err) {
    message.error('获取密钥失败: ' + err)
  } finally {
    loadingKeys.value = false
  }
}

// 打开添加弹窗
const openAddModal = () => {
  editingId.value = null
  Object.assign(form, { name: '', issuerId: '', keyId: '', privateKey: '' })
  isModalVisible.value = true
}

// 打开编辑弹窗
const editKey = (key) => {
  editingId.value = key.id
  Object.assign(form, {
    name: key.name,
    issuerId: key.issuer_id,
    keyId: key.key_id,
    privateKey: '' // 留空，后端会判断为空则不更新私钥
  })
  isModalVisible.value = true
}

// 取消编辑/添加
const cancelEdit = () => {
  isModalVisible.value = false
  editingId.value = null
  Object.assign(form, { name: '', issuerId: '', keyId: '', privateKey: '' })
  if (formRef.value) {
    formRef.value.clearValidate()
  }
}

// 保存密钥 (添加或更新)
const handleSaveKey = async () => {
  try {
    await formRef.value.validate()
  } catch (error) {
    return
  }

  loading.value = true
  try {
    if (editingId.value) {
      await UpdateApiKey(editingId.value, form.name, form.issuerId, form.keyId, form.privateKey)
      message.success('更新成功')
    } else {
      await AddApiKey(form.name, form.issuerId, form.keyId, form.privateKey)
      message.success('添加成功')
    }
    cancelEdit()
    await fetchKeys()
  } catch (err) {
    message.error('保存失败: ' + err)
  } finally {
    loading.value = false
  }
}

// 设置为当前活跃密钥
const setActive = async (id) => {
  try {
    await SetCurrentKey(id)
    message.success('设置成功')
    await fetchKeys()
  } catch (err) {
    message.error('设置失败: ' + err)
  }
}

// 删除密钥
const deleteKey = async (id) => {
  try {
    await DeleteApiKey(id)
    message.success('删除成功')
    await fetchKeys()
  } catch (err) {
    message.error('删除失败: ' + err)
  }
}

// 测试认证
const testAuth = async () => {
  testingAuth.value = true
  try {
    const success = await TestAuth()
    if (success) {
      message.success('认证成功！')
    } else {
      message.error('认证失败。')
    }
  } catch (err) {
    message.error('测试认证异常: ' + err)
  } finally {
    testingAuth.value = false
  }
}

onMounted(() => {
  fetchKeys()
})
</script>
